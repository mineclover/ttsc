import type {
  Expression,
  ObjectLiteralExpression,
  PropertyAssignment,
} from "../../ast";
import { createObjectLiteralExpression } from "./createObjectLiteralExpression";
import { createPropertyAssignment } from "./createPropertyAssignment";

/**
 * The attributes accepted by {@link createPropertyDescriptor}.
 *
 * @author Jeongho Nam - https://github.com/samchon
 */
export interface PropertyDescriptorAttributes {
  /** The `enumerable` flag. */
  enumerable?: Expression;

  /** The `configurable` flag. */
  configurable?: Expression;

  /** The `writable` flag (data descriptor). */
  writable?: Expression;

  /** The `value` (data descriptor). */
  value?: Expression;

  /** The `get` accessor (accessor descriptor). */
  get?: Expression;

  /** The `set` accessor (accessor descriptor). */
  set?: Expression;
}

const tryAddPropertyAssignment = (
  properties: PropertyAssignment[],
  propertyName: string,
  expression: Expression | undefined,
): boolean => {
  if (expression) {
    properties.push(createPropertyAssignment(propertyName, expression));
    return true;
  }
  return false;
};

/**
 * Create a property descriptor object literal suitable for
 * `Object.defineProperty`.
 *
 * Each attribute present on `attributes` becomes a property assignment, added
 * in a fixed order: `enumerable`, `configurable`, `writable`, `value`, `get`,
 * `set`. Absent attributes are skipped. `singleLine` controls layout; the
 * underlying object literal is multi-line by default, so when `singleLine` is
 * falsy each property prints on its own line with a trailing comma.
 *
 * With `attributes` of `{ enumerable: true, value: x }` and the default layout,
 * the printer emits:
 *
 * ```ts
 * {
 *   enumerable: true,
 *   value: x,
 * }
 * ```
 *
 * @author Jeongho Nam - https://github.com/samchon
 * @param attributes The descriptor attributes.
 * @param singleLine When `true`, print on a single line.
 * @returns The created {@link ObjectLiteralExpression}.
 */
export const createPropertyDescriptor = (
  attributes: PropertyDescriptorAttributes,
  singleLine?: boolean,
): ObjectLiteralExpression => {
  const properties: PropertyAssignment[] = [];
  tryAddPropertyAssignment(properties, "enumerable", attributes.enumerable);
  tryAddPropertyAssignment(properties, "configurable", attributes.configurable);
  tryAddPropertyAssignment(properties, "writable", attributes.writable);
  tryAddPropertyAssignment(properties, "value", attributes.value);
  tryAddPropertyAssignment(properties, "get", attributes.get);
  tryAddPropertyAssignment(properties, "set", attributes.set);
  return createObjectLiteralExpression(properties, !singleLine);
};
