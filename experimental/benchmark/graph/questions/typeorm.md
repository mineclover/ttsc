Use this checkout's graph evidence first: symbols, signatures, dependency edges, edge ranges, trace steps, and sourceSpan anchors. Do not use web search, external documentation, package docs, or general framework memory.

Trace how `Repository.find()` turns `relations` find options into query-builder joins: `Repository.find` -> `EntityManager.find` -> `SelectQueryBuilder.setFindOptions` -> `applyFindOptions` -> `buildRelations`. Explain how the relation paths are expanded into join aliases and join attributes.
