import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";

import { TtscGraphMemory } from "../model/TtscGraphMemory";
import { loadGraph } from "../model/loadGraph";
import { createServer } from "./createServer";

/**
 * Serve the graph tools over MCP on stdio. The server answers the MCP handshake
 * immediately and builds the resident graph on the first tool call, so a large
 * project cannot make the client give up before tools are advertised.
 */
export async function startServer(options: {
  cwd?: string;
  tsconfig?: string;
  /** Server version reported in the MCP handshake. */
  version: string;
}): Promise<void> {
  let graph: TtscGraphMemory | undefined;
  const server = createServer(() => {
    graph ??= loadGraph({ cwd: options.cwd, tsconfig: options.tsconfig });
    return graph;
  }, options.version);
  const transport = new StdioServerTransport();
  await server.connect(transport);
}
