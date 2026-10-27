import * as Effect from "effect/Effect";

// Deliberately invalid Effect usage: TypeScript may accept it, Effect diagnostics must identify it.
export function ignoresEffect(): void {
  Effect.sync(() => console.log("This body will not execute"));
}
