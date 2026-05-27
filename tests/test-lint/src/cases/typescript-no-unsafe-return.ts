declare const anyValue: any;

// Positive: `any`-typed expression returned from a function whose
// declared return type is the concrete `number` — the `any` leaks past
// the type boundary.
// expect: typescript/no-unsafe-return error
function asNumber(): number {
  return anyValue;
}

// Positive: same shape with an arrow function whose return type is a
// concrete object shape.
// expect: typescript/no-unsafe-return error
const asObject = (): { id: number } => anyValue;

// Negative: declared return type is `any` itself — the leak is
// explicitly part of the contract.
function asAny(): any {
  return anyValue;
}

JSON.stringify({ a: asNumber(), b: asObject(), c: asAny() });
