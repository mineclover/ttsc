package graph

import "github.com/samchon/ttsc/packages/ttsc/driver"

// FileDiagnostics returns the resident program's diagnostics for the source file
// at path, in the same code and location tsgo reports. Because the graph rides
// the already-open Program, this is one Program.Diagnostics() call over the warm
// checker, not a second compile and not an external language-server round-trip:
// "what is wrong with this file" is answered from the same handle that built the
// reference graph. Lint findings are merged on top by a later pass.
func FileDiagnostics(prog *driver.Program, path string) []driver.Diagnostic {
  out := make([]driver.Diagnostic, 0)
  for _, diagnostic := range prog.Diagnostics() {
    if diagnostic.File == path {
      out = append(out, diagnostic)
    }
  }
  return out
}
