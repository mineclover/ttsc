declare function syncWork(): number;
declare function getPromise(): Promise<number>;

// Positive: async function with no `await` in its body.
// expect: typescript/require-await error
async function noAwait(): Promise<number> {
  return syncWork();
}

// Positive: async arrow with synchronous body.
// expect: typescript/require-await error
const noAwaitArrow = async (): Promise<number> => syncWork();

// Positive: async method with no `await`.
class WithMethods {
  // expect: typescript/require-await error
  async run(): Promise<number> {
    return syncWork();
  }
}

// Positive: `await` lives inside a nested non-async function — does not
// count for the outer async function.
// expect: typescript/require-await error
async function nestedAwaitInsideClosure(): Promise<number> {
  const inner = (): Promise<number> => {
    return getPromise();
  };
  return inner();
}

// Negative: async function with `await` in its body.
async function withAwait(): Promise<number> {
  return await getPromise();
}

// Negative: async arrow with `await`.
const withAwaitArrow = async (): Promise<number> => await getPromise();

// Negative: synchronous function — rule does not apply.
function syncFunction(): number {
  return syncWork();
}

// Negative: async generator — exempt by design (uses `yield`).
async function* asyncGenerator(): AsyncGenerator<number> {
  yield syncWork();
}

JSON.stringify({
  noAwait,
  noAwaitArrow,
  WithMethods,
  nestedAwaitInsideClosure,
  withAwait,
  withAwaitArrow,
  syncFunction,
  asyncGenerator,
});
