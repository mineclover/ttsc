import ttscFactory from "@ttsc/factory";
import fs from "node:fs";
import path from "node:path";
import ts from "typescript";

/**
 * Prove that `@ttsc/factory` mirrors **every** non-deprecated `create*` member
 * of the legacy `ts.factory`, against a real (legacy) `typescript` install.
 *
 * The check is deliberately exhaustive: if even a single non-deprecated factory
 * function is missing, it throws an error listing the whole gap. Two classes of
 * `ts.factory` member are excluded from the requirement:
 *
 * 1. **Deprecated** members — detected automatically by parsing the real
 *    `typescript.d.ts` and treating a name as deprecated only when _every_ one
 *    of its overloads carries an `@deprecated` JSDoc tag.
 * 2. **Out-of-scope** members — JSDoc-node and JSX-node builders, plus emitter /
 *    compiler-internal synthetic-node and convenience-call helpers, none of
 *    which belong to a source-code-generation factory. These are an explicit,
 *    curated allowlist ({@link OUT_OF_SCOPE}) so each omission is a reviewed
 *    decision rather than an accident.
 *
 * The allowlist is also kept honest: a name that `@ttsc/factory` has since
 * implemented (or that is no longer a real `ts.factory` member) is reported as
 * stale and must be removed.
 *
 * @author Jeongho Nam - https://github.com/samchon
 */

/**
 * JSDoc-node builders (`createJSDoc*`) — `@ttsc/factory` models comments, not
 * JSDoc AST.
 */
const isJsDoc = (name: string): boolean => name.startsWith("createJSDoc");

/** JSX-node builders (`createJsx*`) — JSX is out of scope. */
const isJsx = (name: string): boolean => name.startsWith("createJsx");

/**
 * Emitter / compiler-internal node builders and convenience-call helpers that a
 * source-construction factory does not need. Curated allowlist.
 */
const OUT_OF_SCOPE: ReadonlySet<string> = new Set([
  // synthetic / emit-internal nodes
  "createBundle",
  "createRedirectedSourceFile",
  "createSyntaxList",
  "createSyntheticExpression",
  "createSyntheticReferenceExpression",
  "createNotEmittedStatement",
  "createNotEmittedTypeElement",
  "createPartiallyEmittedExpression",
  "createMissingDeclaration",
  "createLiteralLikeNode",
  "createTemplateLiteralLikeNode",
  "createStringLiteralFromNode",
  "createModifiersFromModifierFlags",
  "createUseStrictPrologue",
  // transformer-allocated names
  "createTempVariable",
  "createLoopVariable",
  "createUniqueName",
  "createUniquePrivateName",
  // convenience call / emit helpers (built on the primitive factory)
  "createArrayConcatCall",
  "createArraySliceCall",
  "createAssignmentTargetWrapper",
  "createCallBinding",
  "createFunctionApplyCall",
  "createFunctionBindCall",
  "createFunctionCallCall",
  "createGlobalMethodCall",
  "createMethodCall",
  "createObjectDefinePropertyCall",
  "createObjectGetOwnPropertyDescriptorCall",
  "createPropertyDescriptor",
  "createReflectGetCall",
  "createReflectSetCall",
  "createTypeCheck",
  "createIsNotTypeCheck",
  // modern import-attribute syntax (not yet modelled)
  "createImportAttribute",
  "createImportAttributes",
]);

const isOutOfScope = (name: string): boolean =>
  isJsDoc(name) || isJsx(name) || OUT_OF_SCOPE.has(name);

/** All `create*` members present on the real `ts.factory` at runtime. */
const realFactoryNames = (): string[] =>
  Object.keys(ts.factory).filter((key) => key.startsWith("create"));

/**
 * Names whose _every_ `NodeFactory` overload is `@deprecated`, parsed from the
 * real `typescript.d.ts`.
 */
const deprecatedFactoryNames = (): ReadonlySet<string> => {
  // `getDefaultLibFilePath` points inside the typescript `lib/` directory,
  // which also holds the bundled `typescript.d.ts` — resolved without relying
  // on `import.meta` / `require`.
  const dts: string = path.join(
    path.dirname(ts.getDefaultLibFilePath({})),
    "typescript.d.ts",
  );
  const source: ts.SourceFile = ts.createSourceFile(
    "typescript.d.ts",
    fs.readFileSync(dts, "utf8"),
    ts.ScriptTarget.Latest,
    true,
  );

  const tally = new Map<string, { total: number; deprecated: number }>();
  const visit = (node: ts.Node): void => {
    if (ts.isInterfaceDeclaration(node) && node.name.text === "NodeFactory")
      for (const member of node.members) {
        if (member.name === undefined || !ts.isIdentifier(member.name))
          continue;
        const name: string = member.name.text;
        if (!name.startsWith("create")) continue;
        const deprecated: boolean = ts
          .getJSDocTags(member)
          .some((tag) => tag.tagName.text === "deprecated");
        const counter = tally.get(name) ?? { total: 0, deprecated: 0 };
        counter.total += 1;
        if (deprecated) counter.deprecated += 1;
        tally.set(name, counter);
      }
    ts.forEachChild(node, visit);
  };
  visit(source);

  const names = new Set<string>();
  for (const [name, counter] of tally)
    if (counter.total > 0 && counter.total === counter.deprecated)
      names.add(name);
  return names;
};

/** `create*` members implemented by `@ttsc/factory`. */
const ttscFactoryNames = (): ReadonlySet<string> =>
  new Set(Object.keys(ttscFactory).filter((key) => key.startsWith("create")));

/**
 * Every non-deprecated, in-scope `ts.factory.create*` function is implemented
 * by `@ttsc/factory`.
 */
export const test_factory_completeness = (): void => {
  const real: string[] = realFactoryNames();
  const deprecated: ReadonlySet<string> = deprecatedFactoryNames();
  const ttsc: ReadonlySet<string> = ttscFactoryNames();

  const missing: string[] = real
    .filter((name) => !deprecated.has(name))
    .filter((name) => !isOutOfScope(name))
    .filter((name) => !ttsc.has(name))
    .sort();
  if (missing.length !== 0)
    throw new Error(
      `@ttsc/factory is missing ${missing.length} non-deprecated ts.factory ` +
        `function(s):\n${missing.map((name) => `  - ${name}`).join("\n")}`,
    );
};

/**
 * The {@link OUT_OF_SCOPE} allowlist stays honest: every entry must still be a
 * real `ts.factory` member that `@ttsc/factory` does not implement.
 */
export const test_factory_completeness_allowlist_is_not_stale = (): void => {
  const real = new Set<string>(realFactoryNames());
  const ttsc: ReadonlySet<string> = ttscFactoryNames();

  const notReal: string[] = [...OUT_OF_SCOPE]
    .filter((n) => !real.has(n))
    .sort();
  const implemented: string[] = [...OUT_OF_SCOPE]
    .filter((n) => ttsc.has(n))
    .sort();
  const stale: string[] = [...notReal, ...implemented];
  if (stale.length !== 0)
    throw new Error(
      `OUT_OF_SCOPE allowlist has ${stale.length} stale entr(ies) — remove ` +
        `names no longer missing from @ttsc/factory:\n` +
        `${notReal.map((n) => `  - ${n} (not a ts.factory member)`).join("\n")}\n` +
        `${implemented.map((n) => `  - ${n} (now implemented)`).join("\n")}`,
    );
};
