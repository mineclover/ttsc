#!/usr/bin/env node
import { readFileSync, existsSync, readdirSync } from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const HERE = path.dirname(fileURLToPath(import.meta.url));
const DEFAULT_ROOT = path.resolve(HERE, '../../..');
const REQUIRED = ['002', '003', '004', '005', '023', '033'];
const EFFECT_TYPESCRIPT_GO_HEAD = '2bd066d87f5bafd315be9f40889d0a60b9e58e0b';
const TTSC_PINNED_TYPESCRIPT_GO = '7fc57c005063';

export function preflight(root = DEFAULT_ROOT, patchDirectory = path.resolve(HERE, '../upstream-patches')) {
  const goModPath = path.join(root, 'packages/ttsc/go.mod');
  const mod = existsSync(goModPath) ? readFileSync(goModPath, 'utf8') : null;
  const staged = existsSync(patchDirectory) ? readdirSync(patchDirectory).filter(f => f.endsWith('.patch')) : [];
  const missing = REQUIRED.filter(id => !staged.some(f => f.startsWith(id + '-')));
  const moduleRef = mod?.match(/github\.com\/microsoft\/typescript-go\s+(v[^\s]+)/)?.[1] ?? null;
  const sourceHead = moduleRef?.split('-').at(-1) ?? null;
  const blockers = [];
  if (!mod) blockers.push('ttsc-go-mod-missing');
  if (missing.length) blockers.push('checker-source-patches-missing');
  if (!sourceHead || !EFFECT_TYPESCRIPT_GO_HEAD.startsWith(sourceHead)) blockers.push('typescript-go-baseline-mismatch');
  // A green preflight proves prerequisites only. It never proves a patch was applied or that a Go build passes.
  return {
    schema: 'ttsc-effect-checker-preflight/v1',
    target: path.resolve(root),
    stagedPatchCount: staged.length,
    missingCheckerPatches: missing,
    compiler: { ttscGoModule: moduleRef, actualPrefix: sourceHead, ttscPinnedPrefix: TTSC_PINNED_TYPESCRIPT_GO, effectRequiredHead: EFFECT_TYPESCRIPT_GO_HEAD },
    blockers,
    status: blockers.length ? 'blocked' : 'preflight-only'
  };
}

if (process.argv[1] && path.resolve(process.argv[1]) === fileURLToPath(import.meta.url)) {
  const index = process.argv.indexOf('--root');
  const root = index >= 0 && process.argv[index + 1] ? path.resolve(process.argv[index + 1]) : DEFAULT_ROOT;
  const report = preflight(root);
  process.stdout.write(JSON.stringify(report, null, 2) + '\n');
  if (report.blockers.length) process.exitCode = 1;
}
