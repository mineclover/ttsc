The compiler-resolved TypeScript graph for a relationship or code-flow question. Use this first when the prompt already names symbols, files, or a call chain.

One broad fuzzy query (owner + action + nouns, e.g. "controller dispatch service cache") returns matched declaration coordinates, adjacent graph edges, diagnostic counts, blast radius, and exact handles: the whole cluster in one call, so you do not query one symbol at a time.

Answer from the result; its edge targets are already part of the path, so a node shown as an edge target is part of the answer, not a reason to re-query or grep. If you need declaration source for a TypeScript node, call expand_nodes with its handle. After you edit a file, query again rather than reuse an old result.
