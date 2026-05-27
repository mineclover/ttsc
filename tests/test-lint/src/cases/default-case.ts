// Positive: a `switch` without a `default` clause silently lets unhandled
// discriminants fall through with no value — the rule wants every switch
// to spell out the catch-all path.
function classify(kind: string): string {
  // expect: default-case error
  switch (kind) {
    case "a":
      return "letter-a";
    case "b":
      return "letter-b";
  }
  return "unknown";
}

// Negative: a `switch` that already carries a `default` clause is fine.
function describe(kind: string): string {
  switch (kind) {
    case "a":
      return "letter-a";
    default:
      return "unknown";
  }
}

JSON.stringify({ classify: classify("a"), describe: describe("b") });
