#!/usr/bin/env node
// Run every annotated lint fixture under `src/cases` and report every
// mismatch in one pass. The standard e2e runner halts on first failure,
// which makes triaging fixture drift slow when many fixtures have shifted.
import fs from "node:fs";
import path from "node:path";
import { pathToFileURL } from "node:url";

const here = path.dirname(new URL(import.meta.url).pathname);
const pkgRoot = path.join(here, "..");
process.chdir(pkgRoot);

// Use the package's loader chain to pull in the workspace TestLint helper.
const { TestLint } = await import("@ttsc/testing");
const { assertLintCase } = await import(
  pathToFileURL(path.join(pkgRoot, "src", "helpers", "assertLintCase.ts"))
);

const casesRoot = path.join(pkgRoot, "src", "cases");

function walk(dir) {
  const out = [];
  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    const file = path.join(dir, entry.name);
    if (entry.isDirectory()) out.push(...walk(file));
    else if (entry.isFile()) out.push(file);
  }
  return out.sort();
}

const cases = walk(casesRoot)
  .filter((f) => f.endsWith(".ts"))
  .map((f) => path.relative(casesRoot, f).replaceAll(path.sep, "/"))
  .filter((f) => {
    const src = fs.readFileSync(path.join(casesRoot, f), "utf8");
    return TestLint.parseExpectations(src).length !== 0;
  });

let pass = 0;
const failures = [];
for (const file of cases) {
  try {
    assertLintCase(file);
    pass += 1;
  } catch (err) {
    failures.push({ file, message: err?.message || String(err) });
  }
}
console.log(`PASS ${pass} / FAIL ${failures.length} / TOTAL ${cases.length}`);
for (const f of failures) {
  console.log(`--- ${f.file} ---`);
  console.log(f.message.split("\n").slice(0, 12).join("\n"));
  console.log("");
}
process.exit(failures.length > 0 ? 1 : 0);
