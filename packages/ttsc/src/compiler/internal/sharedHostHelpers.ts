import path from "node:path";

import type { ITtscLoadedNativePlugin } from "../../structures/internal/ITtscLoadedNativePlugin";

/**
 * Resolve the caller-declared plugin config anchor to an absolute path, or
 * `undefined` when the caller did not declare one.
 *
 * The anchor exists for embedders that compile through a generated tsconfig
 * outside the project (the bundler adapters' alias overlay): they set
 * `pluginConfigDir` to the real project directory and every native plugin spawn
 * forwards it as `TTSC_PLUGIN_CONFIG_DIR`, so config-file discovery walks the
 * project instead of the wrapper's temp-dir ancestry. Callers that point at a
 * user-authored tsconfig (even a wrapper outside the project) leave it unset,
 * keeping discovery anchored at the tsconfig's own directory.
 */
export function resolvePluginConfigDir(options: {
  cwd?: string;
  pluginConfigDir?: string;
}): string | undefined {
  if (options.pluginConfigDir === undefined || options.pluginConfigDir === "") {
    return undefined;
  }
  return path.resolve(options.cwd ?? process.cwd(), options.pluginConfigDir);
}

/**
 * Reports whether the given transform source is linked into another compiler
 * host instead of owning the process itself.
 */
export function isLinkedTransform(plugin: ITtscLoadedNativePlugin): boolean {
  return plugin.stage === "transform" && plugin.kind === "linked";
}

/**
 * Verifies that all transform plugins in `plugins` either resolve to the same
 * native binary (the common case) after linked sources are removed from the
 * compiler-owner set.
 *
 * Two callers exist with subtly different error wording: the build path
 * (`runBuild.ts`) reports "multiple compiler native backends cannot share one
 * emit pass" while the source-to-source path (`transformProjectInMemory.ts`)
 * reports "cannot share one source-to-source pass". The `pass` argument selects
 * the appropriate phrase so the error message remains diagnostic-grade instead
 * of generic.
 */
export function assertSharedHostCompatibility(
  plugins: readonly ITtscLoadedNativePlugin[],
  pass: "emit" | "source-to-source",
): void {
  const binaries = [...new Set(plugins.map((plugin) => plugin.binary))];
  if (binaries.length <= 1) {
    return;
  }
  const ownerBinaries = [
    ...new Set(
      plugins
        .filter((plugin) => !isLinkedTransform(plugin))
        .map((plugin) => plugin.binary),
    ),
  ];
  if (ownerBinaries.length <= 1) {
    return;
  }
  const phrase =
    pass === "emit"
      ? "multiple compiler native backends cannot share one emit pass"
      : "multiple transform native backends cannot share one source-to-source pass";
  throw new Error(
    "ttsc: " +
      phrase +
      "; compose transform libraries through one aggregate native host",
  );
}

/**
 * Picks the native binary that must own the compiler pass. Linked transform
 * sources ride inside a host that uses driver.LoadProgram, so an executable
 * transform wins when one is present.
 */
export function selectSharedHostPlugin(
  plugins: readonly ITtscLoadedNativePlugin[],
): ITtscLoadedNativePlugin {
  return plugins.find((plugin) => !isLinkedTransform(plugin)) ?? plugins[0]!;
}

/**
 * Return every plugin whose transform source is linked into another host binary
 * rather than owning the process. The host binary passes these via
 * `TTSC_LINKED_PLUGINS_JSON` so their Go code runs inside the same process.
 */
export function linkedTransformPlugins(
  plugins: readonly ITtscLoadedNativePlugin[],
): ITtscLoadedNativePlugin[] {
  return plugins.filter(isLinkedTransform);
}
