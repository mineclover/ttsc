Render one or more source files in full: every declaration inside, each with its checker-resolved edges (calls, callers, types), the diagnostics on it, its blast radius, and verbatim source.

Pass paths in `locations`; each file is answered as its own block, in input order.

One call returns a whole file as a graph of its objects and how they relate, which is what a file query is for, not its import surface.
