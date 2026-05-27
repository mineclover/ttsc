// Positive: `Foo.Bar` referenced from inside `namespace Foo` — the
// `Foo.` qualifier names the enclosing scope and is redundant.
namespace Foo {
  export const Bar = 1;
  // expect: typescript/no-unnecessary-qualifier error
  export const alias = Foo.Bar;
}

// Positive: enum member referenced from inside the same `enum` body
// via `enum E { X, Y = E.X }`.
enum Color {
  Red = 1,
  // expect: typescript/no-unnecessary-qualifier error
  Crimson = Color.Red,
}

// Negative: `Foo.Bar` referenced from outside `namespace Foo` — the
// qualifier is the only way to reach the member.
const outside = Foo.Bar;

// Negative: an unrelated qualified access whose head does not name
// any enclosing scope.
namespace Outer {
  export namespace Inner {
    export const value = 1;
  }
  // The qualifier `Inner` is the path into the inner namespace — the
  // enclosing scope is `Outer`, not `Inner`, so this is fine.
  export const fine = Inner.value;
}

JSON.stringify({ outside, outer: Outer.Inner.value, alias: Foo.alias });
