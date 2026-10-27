import assert from 'node:assert/strict';
import { test } from 'node:test';
import { mkdtempSync, mkdirSync, writeFileSync, rmSync } from 'node:fs';
import path from 'node:path';
import os from 'node:os';
import { preflight } from '../scripts/preflight.mjs';

const IDS = ['002-checker-checker.patch','003-checker-exports.patch','004-checker-relater.patch','005-checker-types.patch','023-checker-effect-links.patch','033-checker-expression-checked-callback.patch'];
function f(content){
 const root=mkdtempSync(path.join(os.tmpdir(),'ttsc-effect-preflight-'));
 const patches=path.join(root,'patches');
 mkdirSync(path.join(root,'packages','ttsc'),{recursive:true}); mkdirSync(patches);
 writeFileSync(path.join(root,'packages','ttsc','go.mod'),content);
 for(const id of IDS) writeFileSync(path.join(patches,id),'diff --git a/x b/x\n');
 return {root,patches};
}

test('rejects the pinned TTSC June baseline even when patches are staged',()=>{
 const o=f('module x\nrequire github.com/microsoft/typescript-go v0.0.0-20260610182825-7fc57c005063\n');
 try { const r=preflight(o.root,o.patches); assert.equal(r.status,'blocked'); assert.deepEqual(r.blockers,['typescript-go-baseline-mismatch']); }
 finally {rmSync(o.root,{recursive:true,force:true});}
});
test('correct upstream head is only a preflight, not proof of successful port',()=>{
 const o=f('module x\nrequire github.com/microsoft/typescript-go v0.0.0-20260708042240-2bd066d87f5b\n');
 try { const r=preflight(o.root,o.patches); assert.equal(r.status,'preflight-only'); assert.deepEqual(r.blockers,[]); }
 finally {rmSync(o.root,{recursive:true,force:true});}
});
test('missing required Checker patch fails closed',()=>{
 const o=f('module x\nrequire github.com/microsoft/typescript-go v0.0.0-20260708042240-2bd066d87f5b\n');
 try {rmSync(path.join(o.patches,IDS[0]));const r=preflight(o.root,o.patches);assert.equal(r.status,'blocked'); assert.deepEqual(r.missingCheckerPatches,['002']);}
 finally {rmSync(o.root,{recursive:true,force:true});}
});
