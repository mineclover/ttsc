import type { ExpressionStatement } from "../../ast";
import { createStringLiteral } from "../literals/createStringLiteral";
import { createExpressionStatement } from "./createExpressionStatement";

/**
 * Create the `"use strict";` prologue statement.
 *
 * This is a convenience wrapper that takes no inputs: it builds the `"use
 * strict"` string literal and wraps it in an {@link ExpressionStatement}. Place
 * it first in a function or module body to opt that scope into strict mode.
 *
 * The result is always:
 *
 * ```ts
 * "use strict";
 * ```
 *
 * @author Jeongho Nam - https://github.com/samchon
 * @returns The created {@link ExpressionStatement}.
 */
export const createUseStrictPrologue = (): ExpressionStatement =>
  createExpressionStatement(createStringLiteral("use strict"));
