import assert from 'node:assert/strict';
import { test } from 'node:test';
import { mkdtempSync, writeFileSync, chmodSync, rmSync } from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import { runProbe } from '../scripts/probe-compat.mjs';

function fixtures(kind = 'valid') {
  const root = mkdtempSync(path.join(os.tmpdir(), 'ttsc-effect-probe-'));
  const ttsc = path.join(root, 'ttsc');
  const effect = path.join(root, 'effect-tsgo');
  writeFileSync(ttsc, '#!/usr/bin/env node\nprocess.stdout.write("TS check complete\\n")\n');
  writeFileSync(effect, kind === 'malformed' ? '#!/usr/bin/env node\nprocess.stdout.write("not json")\n' : '#!/usr/bin/env node\nprocess.stdout.write(JSON.stringify({diagnostics: [{file:"src/invalid.ts", start:1, length:1, line:1, column:1, code:377001, name:"floatingEffect", severity:"error", message:"Unrun effect"}], summary:{filesChecked:2,filesWithoutPlugin:0,totalFiles:2,errors:1,warnings:0,messages:0}}))\nprocess.exitCode=1\n');
  chmodSync(ttsc, 0o755); chmodSync(effect, 0o755);
  return { root, ttsc, 'effect-tsgo': effect, 'ttsc-config': 'ttsc.json', 'effect-config': 'effect.json', 'timeout-ms': 3000 };
}

test('reports both compiler observations and expected diagnostic without claiming compatibility', () => {
  const o = fixtures();
  try {
    const result = runProbe({ ...o, 'expect-rule': 'floatingEffect' });
    assert.equal(result.observedCompletely, true);
    assert.equal(result.verdict, 'observed-needs-compatibility-review');
    assert.equal(result.observations.effectTsgo.analysis.diagnostics[0].source, 'effect-tsgo');
    assert.equal(result.observations.effectTsgo.exitCode, 1);
  } finally { rmSync(o.root, { recursive: true, force: true }); }
});

test('malformed diagnostics are incomplete rather than empty/pass', () => {
  const o = fixtures('malformed');
  try {
    const result = runProbe(o);
    assert.equal(result.observedCompletely, false);
    assert.equal(result.observations.effectTsgo.analysis.status, 'invalid-output');
  } finally { rmSync(o.root, { recursive: true, force: true }); }
});

test('missing tools are unavailable, never success', () => {
  const o = fixtures();
  try {
    const result = runProbe({ ...o, ttsc: path.join(o.root, 'does-not-exist') });
    assert.equal(result.observations.ttsc.status, 'unavailable');
    assert.equal(result.observedCompletely, false);
  } finally { rmSync(o.root, { recursive: true, force: true }); }
});

test('unselected Effect plugin does not produce a false clean finding', () => {
  const o = fixtures();
  try {
    writeFileSync(o['effect-tsgo'], '#!/usr/bin/env node\nprocess.stdout.write(JSON.stringify({diagnostics:[],summary:{filesChecked:0,filesWithoutPlugin:3,totalFiles:3,errors:0,warnings:0,messages:0}}))\n');
    const result = runProbe(o);
    assert.equal(result.observations.effectTsgo.analysis.coverage, 'plugin-disabled-or-unselected');
    assert.equal(result.observedCompletely, false);
  } finally { rmSync(o.root, { recursive: true, force: true }); }
});
