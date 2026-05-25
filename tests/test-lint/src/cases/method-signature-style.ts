interface Service {
  // expect: methodSignatureStyle error
  run(input: string): number;
  keep: (input: string) => number;
}

type Handler = {
  // expect: methodSignatureStyle error
  handle(): void;
  keep: () => void;
};

class Impl {
  run(input: string): number {
    return input.length;
  }
}

JSON.stringify({} as Service);
JSON.stringify({} as Handler);
JSON.stringify(Impl);
