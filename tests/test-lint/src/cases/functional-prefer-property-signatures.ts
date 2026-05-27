// expect: functional/prefer-property-signatures error
interface Api {
  run(): void;
}

declare const api: Api;
JSON.stringify(api);
