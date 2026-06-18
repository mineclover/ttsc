import type { RedirectedSourceFile, SourceFile } from "../../ast";
import { make } from "../internal/make";

/**
 * Create a {@link RedirectedSourceFile}: a thin wrapper that emits another
 * source file's content in its place.
 *
 * This is a simplified model. The legacy compiler builds the node from a
 * `redirectInfo` record carrying both the redirect target and the unredirected
 * file plus path bookkeeping. Here it accepts and stores just the redirect
 * target source file, and the printer emits exactly that target's content with
 * nothing added.
 *
 * Given a target source file that re-exports `a` from `"./a"`, printing the
 * redirected file yields the target's own output:
 *
 * ```ts
 * export { a } from "./a";
 * ```
 *
 * @author Jeongho Nam - https://github.com/samchon
 * @param redirectTarget The source file whose content is emitted.
 * @returns The created {@link RedirectedSourceFile}.
 */
export const createRedirectedSourceFile = (
  redirectTarget: SourceFile,
): RedirectedSourceFile => make("RedirectedSourceFile", { redirectTarget });
