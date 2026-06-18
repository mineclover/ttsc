import type { MissingDeclaration } from "../../ast";
import { make } from "../internal/make";

/**
 * Create a {@link MissingDeclaration}: a placeholder for a missing node.
 *
 * The TypeScript parser produces this node to fill a hole left by a syntax
 * error so the tree stays well formed. It carries no name, modifiers, or body
 * and is not something you would hand-author. The printer emits nothing for it,
 * so the rendered output is empty:
 *
 * @author Jeongho Nam - https://github.com/samchon
 * @returns The created {@link MissingDeclaration}.
 */
export const createMissingDeclaration = (): MissingDeclaration =>
  make("MissingDeclaration", {});
