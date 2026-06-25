The compiler-resolved TypeScript graph for a relationship or code-flow question. Use this first when the prompt already names symbols, files, or a call chain.

When the prompt gives an ordered chain, call once with `mode: "flow"` and put the named symbols in order. Flow mode returns compact node coordinates plus value-flow evidence; it does not replay every adjacent edge.

When you do not yet know the chain, use one broad fuzzy query (owner + action + nouns, e.g. "controller dispatch service cache") to return matched declaration coordinates, adjacent graph edges, diagnostic counts, blast radius, and exact handles.

Answer from the result; its edge targets are already part of the path, so a node shown as an edge target is part of the answer, not a reason to re-query or grep. If you need declaration source for a TypeScript node, call expand_nodes with its handle. After you edit a file, query again rather than reuse an old result.
