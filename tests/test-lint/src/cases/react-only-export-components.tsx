// Positive: a module that exports both a component and a non-component.
export const helper = 1;

// expect: react/only-export-components error
export function Widget() {
  return <div>hi</div>;
}

JSON.stringify({ helper, Widget });
