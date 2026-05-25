export const marker = 1;

// expect: noUselessEmptyExport error
export {};

JSON.stringify(marker);
