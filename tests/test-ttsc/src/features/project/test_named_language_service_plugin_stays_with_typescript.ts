import assert from "node:assert/strict";
import fs from "node:fs";
import path from "node:path";

import { readProjectConfig } from "../../../../../packages/ttsc/src/compiler/internal/project/readProjectConfig";
import { ProjectPluginEntries } from "../../../../../packages/ttsc/src/plugin/internal/load/ProjectPluginEntries";
import { hasProjectPluginEntries } from "../../../../../packages/ttsc/src/plugin/internal/load/hasProjectPluginEntries";
import { TestProject } from "../../../../utils/src/TestProject";

/**
 * A TypeScript language-service plugin belongs to the upstream compiler, not
 * to ttsc's native transform descriptor selection.
 *
 * @evidence contracts/testing.md#behavioral-verification Tests the actual config reader and production plugin selector against mixed name/transform entries, including explicit entries.
 * @evidence contracts/testing.md#independent-expectations A valid name-only entry is excluded from native loading but preserved in tsconfig; a transform entry stays, and malformed transforms must not become silently ignored plugins.
 * @evidence contracts/testing.md#distinguishing-cases Includes name-only, transform-only, mixed and malformed entries, plus name-only presence and explicit selection.
 * @evidence contracts/testing.md#execution-ownership Owns a temporary config; no compiler spawn, native plugin builds or network package downloads.
 */
export function test_named_language_service_plugin_stays_with_typescript(): void {
  const root = TestProject.physicalPath(
    TestProject.tmpdir("ttsc-named-typescript-plugin-"),
  );
  try {
    const config = path.join(root, "tsconfig.json");
    const language = { name: "@effect/language-service" };
    const native = { transform: "./local-plugin.cjs" };
    const malformed = { name: "other/language-plugin", transform: null };
    const unknown = { metadata: "missing transform and name" };
    fs.writeFileSync(
      config,
      JSON.stringify({
        compilerOptions: {
          plugins: [language, native, malformed, unknown],
        },
      }),
      "utf8",
    );
    const project = readProjectConfig({ cwd: root, tsconfig: config });
    assert.deepEqual(project.compilerOptions.plugins, [
      language, native, malformed, unknown,
    ]);
    const selected = ProjectPluginEntries.resolvePluginEntries(project);
    assert.deepEqual(selected.map((entry) => entry.config), [
      native, malformed, unknown,
    ]);
    assert.equal(hasProjectPluginEntries(project), true);
    assert.deepEqual(
      ProjectPluginEntries.resolvePluginEntries(project, [language, native]).map(
        (entry) => entry.config,
      ),
      [native],
    );

    fs.writeFileSync(
      config,
      JSON.stringify({ compilerOptions: { plugins: [language] } }),
      "utf8",
    );
    const languageOnly = readProjectConfig({ cwd: root, tsconfig: config });
    assert.equal(hasProjectPluginEntries(languageOnly), false);
    assert.deepEqual(languageOnly.compilerOptions.plugins, [language]);
  } finally {
    fs.rmSync(root, { recursive: true, force: true });
  }
}
