class Repeated {
  constructor(
    public value: string,
    private readonly count: number,
    normal: string,
  ) {
    // expect: noUnnecessaryParameterPropertyAssignment error
    this.value = value;
    // expect: noUnnecessaryParameterPropertyAssignment error
    this.count = count;
    this.normal = normal;
  }
}

class ChangedFirst {
  constructor(public value: string) {
    this.value = value.trim();
    this.value = value;
  }
}

JSON.stringify([Repeated, ChangedFirst]);
