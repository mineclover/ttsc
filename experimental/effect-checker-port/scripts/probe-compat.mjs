#!/usr/bin/env node
import { spawnSync } from 'node:child_process';
import { readFileSync, writeFileSync, readdirSync, existsSync, statSync } from 'node:fs';
import { createHash } from 'node:crypto';
import { fileURLToPath } from 'node:url';
import path from 'node:path';

const ALLOWED = new Set(['root', 'ttsc', 'effect-tsgo', 'ttsc-config', 'effect-config', 'runtime-test', 'expect-rule', 'out', 'timeout-ms']);
function parseArgs(argv) {
  const result = { root: process.cwd(), ttsc: 'ttsc', 'effect-tsgo': 'effect-tsgo', 'ttsc-config': 'tsconfig.ttsc.json', 'effect-config': 'tsconfig.effect.json', 'timeout-ms': '120000' };
  for (let i = 0; i < argv.length; i += 2) {
    const key = argv[i]?.replace(/^--/, '');
    if (!argv[i]?.startsWith('--') || !ALLOWED.has(key) || !argv[i + 1] || argv[i + 1].startsWith('--')) {
      throw new Error('Invalid arguments near ' + (argv[i] ?? '(end)') + '; expected --key value');
    }
    result[key] = argv[i + 1];
  }
  const timeoutMs = Number(result['timeout-ms']);
  if (!Number.isInteger(timeoutMs) || timeoutMs < 1 || timeoutMs > 600000) throw new Error('timeout-ms must be 1..600000');
  result.root = path.resolve(result.root);
  result['timeout-ms'] = timeoutMs;
  return result;
}

function execute(command, args, cwd, timeoutMs) {
  const started = performance.now();
  const p = spawnSync(command, args, { cwd, encoding: 'utf8', timeout: timeoutMs, maxBuffer: 2 * 1024 * 1024, windowsHide: true, shell: false });
  const errorType = p.error?.code === 'ENOENT' ? 'unavailable'
    : p.error?.code === 'ETIMEDOUT' ? 'timeout'
    : p.error ? 'spawn-failed'
    : p.signal ? 'signaled' : 'exited';
  return { status: errorType, exitCode: p.status, signal: p.signal, elapsedMs: Math.round(performance.now() - started), stdout: p.stdout ?? '', stderr: p.stderr ?? '', ...(p.error ? { error: p.error.message } : {}) };
}

function fingerprintInputs(cwd, configs) {
  const files = [...configs, 'package.json', 'package-lock.json'];
  const srcDir = path.join(cwd, 'src');
  if (existsSync(srcDir)) {
    const walk = (dir) => {
      for (const entry of readdirSync(dir, { withFileTypes: true }).sort((a, b) => a.name.localeCompare(b.name))) {
        const file = path.join(dir, entry.name);
        if (entry.isDirectory()) walk(file);
        else if (entry.isFile() && /\.[cm]?tsx?$/.test(entry.name)) files.push(path.relative(cwd, file));
      }
    };
    walk(srcDir);
  }
  const digest = createHash('sha256');
  const entries = [];
  for (const file of [...new Set(files)].sort()) {
    const full = path.resolve(cwd, file);
    if (path.relative(cwd, full).startsWith('..')) throw new Error('Input outside workspace: ' + file);
    const bytes = existsSync(full) && statSync(full).isFile() ? readFileSync(full) : null;
    const hash = bytes ? createHash('sha256').update(bytes).digest('hex') : null;
    entries.push({ file, hash });
    digest.update(file + '\0' + (hash ?? 'MISSING') + '\0');
  }
  return { digest: digest.digest('hex'), entries };
}

function normalizeEffect(raw) {
  if (!raw || typeof raw !== 'object' || !Array.isArray(raw.diagnostics) || !raw.summary || typeof raw.summary !== 'object') throw new Error('Effect diagnostic output missing diagnostics[] or summary');
  const s = raw.summary;
  for (const k of ['filesChecked', 'filesWithoutPlugin', 'totalFiles', 'errors', 'warnings', 'messages']) {
    if (!Number.isInteger(s[k]) || s[k] < 0) throw new Error('Invalid Effect summary.' + k);
  }
  const diagnostics = raw.diagnostics.map((d, i) => {
    if (typeof d.name !== 'string' || typeof d.message !== 'string' || typeof d.file !== 'string' || !['error', 'warning', 'message'].includes(d.severity) || !Number.isInteger(d.code)) {
      throw new Error('Invalid Effect diagnostics[' + i + ']');
    }
    return { source: 'effect-tsgo', rule: d.name, code: d.code, severity: d.severity, file: d.file, start: d.start, length: d.length, line: d.line, column: d.column, message: d.message };
  });
  return { summary: s, diagnostics, files: raw.files ?? [], coverage: s.filesChecked === 0 ? (s.filesWithoutPlugin > 0 ? 'plugin-disabled-or-unselected' : 'no-checked-files') : (s.filesWithoutPlugin > 0 ? 'partial' : 'observed') };
}

export function runProbe(options) {
  const cwd = options.root;
  const inputBefore = fingerprintInputs(cwd, [options['ttsc-config'], options['effect-config']]);
  const ttsc = execute(options.ttsc, ['check', '-p', options['ttsc-config']], cwd, options['timeout-ms']);
  const effect = execute(options['effect-tsgo'], ['diagnostics', '--project', options['effect-config'], '--format', 'json'], cwd, options['timeout-ms']);
  let analysis;
  if (effect.status === 'exited') {
    try { analysis = { status: 'parsed', ...normalizeEffect(JSON.parse(effect.stdout)) }; }
    catch (e) { analysis = { status: 'invalid-output', error: String(e) }; }
  } else {
    analysis = { status: 'not-observed' };
  }
  const runtime = options['runtime-test'] ? execute(process.execPath, ['--test', options['runtime-test']], cwd, options['timeout-ms']) : null;
  const inputAfter = fingerprintInputs(cwd, [options['ttsc-config'], options['effect-config']]);
  const inputsStable = inputBefore.digest === inputAfter.digest;
  const expectedRule = options['expect-rule'];
  const expectation = expectedRule ? { rule: expectedRule, observed: analysis.status === 'parsed' && analysis.diagnostics.some(d => d.rule === expectedRule) } : null;
  const complete = inputsStable && ttsc.status === 'exited' && effect.status === 'exited' && analysis.status === 'parsed' && analysis.coverage === 'observed' && (!runtime || runtime.status === 'exited') && (!expectation || expectation.observed);
  // Never equate clean diagnostics, nonzero lint exit, an expected error, or successful execution with each other.
  return {
    schema: 'ttsc-effect-checker-probe/v1',
    generatedAt: new Date().toISOString(),
    scope: { cwd, ttscConfig: options['ttsc-config'], effectConfig: options['effect-config'], runtimeRequested: Boolean(runtime), inputBefore, inputAfter, inputsStable },
    observations: {
      ttsc: { ...ttsc, stdout: ttsc.stdout.slice(-131072), stderr: ttsc.stderr.slice(-131072) },
      effectTsgo: { ...effect, stdout: effect.stdout.slice(-131072), stderr: effect.stderr.slice(-131072), analysis },
      runtime: runtime && { ...runtime, stdout: runtime.stdout.slice(-131072), stderr: runtime.stderr.slice(-131072) }
    },
    expectation,
    observedCompletely: complete,
    // A fully observed probe is not a compatibility pass until all cross-engine expectations are independently assessed.
    verdict: complete ? 'observed-needs-compatibility-review' : 'incomplete-or-unavailable'
  };
}

function main() {
  try {
    const opts = parseArgs(process.argv.slice(2));
    const result = runProbe(opts);
    const json = JSON.stringify(result, null, 2) + '\n';
    if (opts.out) writeFileSync(path.resolve(opts.root, opts.out), json);
    else process.stdout.write(json);
    if (!result.observedCompletely) process.exitCode = 1;
  } catch (err) {
    process.stderr.write('effect-checker-probe: ' + err.message + '\n');
    process.exitCode = 2;
  }
}

if (process.argv[1] && path.resolve(process.argv[1]) === path.resolve(fileURLToPath(import.meta.url))) main();
