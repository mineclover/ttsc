/**
 * A placeholder for a declaration that failed to parse. It emits nothing.
 *
 * Built by {@link factory.createMissingDeclaration}.
 *
 * @author Jeongho Nam - https://github.com/samchon
 */
export interface MissingDeclaration {
  /** Discriminant tag; always `"MissingDeclaration"`. */
  kind: "MissingDeclaration";
}
