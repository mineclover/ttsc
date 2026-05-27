// Positive: the `get value` and `set value` accessors are split apart by an
// unrelated method, so a reader scanning the class has to chase the pair
// across the body — the rule wants getter/setter pairs to sit together.
class Splayed {
  private state = 0;

  get value(): number {
    return this.state;
  }

  other(): void {
    this.state += 1;
  }

  // expect: grouped-accessor-pairs error
  set value(next: number) {
    this.state = next;
  }
}

// Negative: the matching `get` and `set` declarations are adjacent, which
// is the layout the rule wants every accessor pair to follow.
class Grouped {
  private state = 0;

  get value(): number {
    return this.state;
  }
  set value(next: number) {
    this.state = next;
  }
}

JSON.stringify({ splayed: new Splayed(), grouped: new Grouped() });
