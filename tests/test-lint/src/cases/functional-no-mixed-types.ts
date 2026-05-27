// expect: functional/no-mixed-types error
type Mixed = {
  value: number;
  compute(): number;
};

declare const m: Mixed;
JSON.stringify(m);
