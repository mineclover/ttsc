// expect: functional/no-return-void error
function log(): void {
  return;
}

JSON.stringify(log);
