Compiler-resolved TypeScript relationship discovery. Use this when exact path endpoints are unknown.

It returns matched declaration coordinates, adjacent graph edges, diagnostic counts, blast radius, and handles. Use query_path instead when the task gives exact start/end symbols or an ordered chain.

If source context is required, call expand_nodes with a returned handle. After editing, query again rather than reusing an old result.
