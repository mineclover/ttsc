package linthost

import (
  "path/filepath"
  "strings"
  "testing"

  shimast "github.com/microsoft/typescript-go/shim/ast"
  shimchecker "github.com/microsoft/typescript-go/shim/checker"
)

// TestNoFloatingPromisesDetectsIndexKeyTypeParameters verifies the conservative
// generic proof traverses both halves of an index contract.
//
// A mapped type can retain its method type parameter only in an index key such
// as keyof T while exposing a concrete value type. Looking only at index values
// would misclassify that unresolved parameter shape as concrete.
//
//  1. Declare a generic catch parameter mapped over keyof T.
//  2. Recover the method signature and its original parameter type.
//  3. Assert the latent-type-parameter scan detects T through the index key.
func TestNoFloatingPromisesDetectsIndexKeyTypeParameters(t *testing.T) {
  root := t.TempDir()
  writeFile(t, filepath.Join(root, "tsconfig.json"), `{
  "compilerOptions": {
    "target": "ES2022",
    "module": "commonjs",
    "strict": true
  },
  "files": ["main.ts"]
}
`)
  writeFile(t, filepath.Join(root, "main.ts"), `interface KeyedCatch {
  catch<T>(value: { [K in keyof T]: string }): undefined;
}
declare const keyed: KeyedCatch;
keyed.catch({});
`)

  prog, diags, err := loadProgram(root, "tsconfig.json", loadProgramOptions{
    needsRuleChecker: true,
  })
  if err != nil {
    t.Fatal(err)
  }
  if len(diags) != 0 {
    t.Fatalf("unexpected configuration diagnostics: %#v", diags)
  }
  defer prog.close()
  files := prog.userSourceFiles()
  if len(files) != 1 || prog.checker == nil {
    t.Fatalf("program setup mismatch: files=%d checker=%v", len(files), prog.checker != nil)
  }
  file := files[0]
  offset := strings.Index(file.Text(), "keyed.catch")
  if offset < 0 {
    t.Fatal("mapped-key fixture marker not found")
  }
  node := shimast.GetNodeAtPosition(file, offset, false)
  for node != nil && node.Kind != shimast.KindCallExpression {
    node = node.Parent
  }
  if node == nil {
    t.Fatal("mapped-key fixture call not found")
  }
  call := node.AsCallExpression()
  if call == nil {
    t.Fatal("mapped-key fixture node is not a call")
  }
  access := call.Expression.AsPropertyAccessExpression()
  if access == nil {
    t.Fatal("mapped-key fixture call is not a property access")
  }
  receiver := access.Expression
  receiverType := prog.checker.GetTypeAtLocation(receiver)
  if receiverType == nil {
    t.Fatal("mapped-key receiver has no type")
  }
  property := prog.checker.GetPropertyOfType(receiverType, "catch")
  if property == nil {
    t.Fatal("mapped-key catch property not found")
  }
  propertyType := prog.checker.GetTypeOfSymbolAtLocation(property, call.Expression)
  if propertyType == nil {
    t.Fatal("mapped-key catch property has no type")
  }
  signatures := prog.checker.GetSignaturesOfType(propertyType, shimchecker.SignatureKindCall)
  if len(signatures) != 1 || len(signatures[0].TypeParameters()) != 1 {
    t.Fatalf("mapped-key signature mismatch: signatures=%d", len(signatures))
  }
  parameterType := floatingPromiseParameterType(prog.checker, call, signatures[0], 0)
  if parameterType == nil {
    t.Fatal("mapped-key parameter has no type")
  }
  indexInfos := prog.checker.GetIndexInfosOfType(parameterType)
  if len(indexInfos) == 0 {
    t.Fatal("mapped-key parameter exposes no index info")
  }
  sawLatentKey := false
  for _, indexInfo := range indexInfos {
    if indexInfo == nil {
      t.Fatal("mapped-key parameter exposes a nil index info")
    }
    keyType := indexInfo.KeyType()
    valueType := indexInfo.ValueType()
    if keyType == nil || valueType == nil {
      t.Fatal("mapped-key index info has a nil key or value type")
    }
    if floatingPromiseTypeContainsAnyTypeParameter(
      prog.checker,
      valueType,
      signatures[0].TypeParameters(),
      call.Expression,
      nil,
    ) {
      t.Fatal("mapped index value unexpectedly retains a type parameter")
    }
    if floatingPromiseTypeContainsAnyTypeParameter(
      prog.checker,
      keyType,
      signatures[0].TypeParameters(),
      call.Expression,
      nil,
    ) {
      sawLatentKey = true
    }
  }
  if !sawLatentKey {
    t.Fatal("mapped index key does not retain the method type parameter")
  }
  if !floatingPromiseTypeContainsAnyTypeParameter(
    prog.checker,
    parameterType,
    signatures[0].TypeParameters(),
    call.Expression,
    nil,
  ) {
    t.Fatal("mapped index key lost its method type parameter")
  }
}
