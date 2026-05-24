package banner_test

import (
  "strings"
  "testing"
)

// TestTypeScriptConfigLoaderSourcePrefersDefaultThenText verifies loader order.
//
// The TypeScript config loader emits a small module that resolves default
// interop wrappers at runtime. It must start from the module default export
// when present, then stop unwrapping once the current value is already a banner
// object.
//
// 1. Generate the TypeScript config loader source.
// 2. Locate the initial default-export selection and the banner-object guard.
// 3. Assert default selection happens before the guard.
func TestTypeScriptConfigLoaderSourcePrefersDefaultThenText(t *testing.T) {
  source := bannerTypeScriptConfigLoaderSource(`"./banner.config.ts"`)
  initial := strings.Index(source, `let current = isObject(value) && hasOwn(value, "default") ? value.default : value;`)
  guard := strings.Index(source, "isBannerObject(current)")
  if initial < 0 || guard < 0 || initial > guard {
    t.Fatalf("loader should select default before testing banner object:\n%s", source)
  }
}
