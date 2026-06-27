I just joined this TypeScript project and have been asked to change behavior on the busiest user-facing path, but I do not yet know which exported API or entry point is the right place to start.

Orient me from the code before choosing the path: identify the public entry points that look user-facing, the internal modules or types they converge on, and the symbols that appear to be reused by several parts of the project. From that evidence, choose one central runtime operation whose implementation crosses the important layers.

Trace that operation in order from the public entry point through the orchestration, state/model, adapter, or engine layers down to the lowest layer where real work happens. For each hop, name the file and symbol, and say how it reaches the next hop: direct call, callback, dependency injection, inheritance, dispatch, or data flow.

Also name two plausible entry points or modules you considered but ruled out, and explain why they are adjacent but not the main path. End with the smallest set of source ranges I should read before editing and one risk that would be easy to miss if I only followed direct calls.

Do not guess; report gaps.
