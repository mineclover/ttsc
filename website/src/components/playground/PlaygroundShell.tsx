"use client";

import {
  createSandboxRequire,
  loadTypiaRuntimePack,
  PlaygroundShell as PackagedPlaygroundShell,
} from "@ttsc/playground";

import { PLAYGROUND_EXAMPLES } from "../../compiler/PlaygroundExamples";
import typiaTypes from "../../compiler/typia-types.json";

const PLAYGROUND_DEFAULT_SCRIPT = PLAYGROUND_EXAMPLES[0]?.source ?? "";

const TYPIA_RUNTIME_PACK_URL = "/compiler/typia-runtime-pack.json";

/**
 * Runtime files installed by `installPlaygroundDependencies` are mounted
 * by the worker itself; this wrapper-side require sandbox handles only the
 * typia transform's runtime references via the prebuilt typia runtime pack.
 */
const executeBundle = async (
  code: string,
  sandbox: { console: Record<string, (...args: unknown[]) => void> },
): Promise<void> => {
  const runtimePack = await loadTypiaRuntimePack(TYPIA_RUNTIME_PACK_URL);
  const sandboxRequire = createSandboxRequire(runtimePack, {
    console: sandbox.console,
  });
  const moduleObj: { exports: Record<string, unknown> } = { exports: {} };
  const wrapped = `(function(require, module, exports, console) {\n${code}\n})`;
  const factory = new Function("return " + wrapped)() as (
    req: (s: string) => unknown,
    mod: typeof moduleObj,
    exp: typeof moduleObj.exports,
    c: typeof sandbox.console,
  ) => void;
  factory(sandboxRequire, moduleObj, moduleObj.exports, sandbox.console);
};

export default function PlaygroundShell() {
  return (
    <PackagedPlaygroundShell
      workerUrl="/compiler/index.js"
      defaultScript={PLAYGROUND_DEFAULT_SCRIPT}
      examples={PLAYGROUND_EXAMPLES}
      exampleGroupLabels={{
        typia: "typia",
        lint: "@ttsc/lint",
        mixed: "mixed",
      }}
      staticEditorLibs={typiaTypes as Record<string, string>}
      executeBundle={executeBundle}
      brand={
        <a
          href="/"
          className="font-mono text-sm font-bold text-white hover:text-blue-400 transition-colors"
        >
          ttsc
        </a>
      }
      resultCaption={(options) =>
        options.typia
          ? "dist/playground.js"
          : "dist/playground.js · typia disabled"
      }
    />
  );
}
