declare function useEffect(effect: () => void, deps: ReadonlyArray<unknown>): void;

function Component(props: { value: number }) {
  // Positive: deps array omits `props.value`.
  // expect: react/exhaustive-deps error
  useEffect(() => {
    JSON.stringify(props.value);
  }, []);

  // Negative: deps array lists every reactive identifier.
  useEffect(() => {
    JSON.stringify(props.value);
  }, [props.value]);

  return null;
}

JSON.stringify(Component);
