declare const anyValue: any;

// Positive: assigning `any` to a typed local discards the static guarantee.
// expect: typescript/no-unsafe-assignment error
const num: number = anyValue;

// Positive: assigning `any` to a typed local is unsafe regardless of target.
// expect: typescript/no-unsafe-assignment error
const str: string = anyValue;

// Positive: same hazard via a reassignment to a typed binding.
let label: string = "ok";
// expect: typescript/no-unsafe-assignment error
label = anyValue;

// Negative: a properly typed initializer is fine.
const safe: number = 1;

JSON.stringify({ num, str, label, safe });
