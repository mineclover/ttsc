// The canonical graph data model: the wire contract `ttscgraph dump` emits and
// the MCP server loads. Pure types (plus the schema-version constant) so typia
// can derive validators and tool schemas from them at build time, and so the Go
// `dump.go` writer has one TypeScript source of truth to mirror.

export * from "./Confidence";
export * from "./EdgeKind";
export * from "./IComponentMetadata";
export * from "./IDecoratorFact";
export * from "./IEvidence";
export * from "./IGraphDiagnostic";
export * from "./IGraphDump";
export * from "./IGraphEdge";
export * from "./IGraphNode";
export * from "./IRouteMetadata";
export * from "./NodeKind";
export * from "./NodeModifier";
export * from "./Provenance";
