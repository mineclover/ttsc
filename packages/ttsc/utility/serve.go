package utility

import (
  "bufio"
  "encoding/json"
  "fmt"
  "io"
  "os"
  "path/filepath"
  "strings"

  shimprinter "github.com/microsoft/typescript-go/shim/printer"
)

// serveRequest is one newline-delimited request the resident host reads from its
// input stream: the absolute or project-relative path of a file whose
// transformed TypeScript the caller wants.
type serveRequest struct {
  File string `json:"file"`
}

// serveResponse is one newline-delimited reply: the transformed TypeScript for
// the requested file and whether the resident program had it.
type serveResponse struct {
  TypeScript string `json:"typescript"`
  Found      bool   `json:"found"`
}

// RunServe is the resident transform host. It transforms the whole project once
// (the expensive compile plus plugin pass), caches every file's transformed
// text, then answers per-file requests read from in by writing one JSON reply
// per line to out, until in reaches EOF. One resident process can serve every
// Metro worker, which removes the per-worker recompile that the per-call
// transform subcommand incurs (samchon/ttsc#255).
//
// in and out are explicit so the request loop is testable; the utility-host
// command wires them to os.Stdin and os.Stdout.
func RunServe(in io.Reader, out io.Writer, args []string) int {
  opts, ok := parseHostOptions("serve", args)
  if !ok {
    return 2
  }
  cache, ok := buildServeCache(opts)
  if !ok {
    return 2
  }
  encoder := json.NewEncoder(out)
  scanner := bufio.NewScanner(in)
  scanner.Buffer(make([]byte, 0, 64*1024), 64*1024*1024)
  for scanner.Scan() {
    line := strings.TrimSpace(scanner.Text())
    if line == "" {
      continue
    }
    var req serveRequest
    if err := json.Unmarshal([]byte(line), &req); err != nil {
      _ = encoder.Encode(serveResponse{})
      continue
    }
    key := apiOutputKey(opts.cwd, resolveServePath(opts.cwd, req.File))
    text, found := cache[key]
    _ = encoder.Encode(serveResponse{TypeScript: text, Found: found})
  }
  if err := scanner.Err(); err != nil {
    fmt.Fprintf(os.Stderr, "ttsc utility serve: read error: %v\n", err)
    return 2
  }
  return 0
}

// buildServeCache runs the whole-project transform once and returns the
// transformed text keyed exactly like the transform subcommand's JSON envelope.
func buildServeCache(opts hostOptions) (map[string]string, bool) {
  prog, _, ok := loadUtilityProgram(opts)
  if !ok {
    return nil, false
  }
  defer prog.Close()
  if err := prog.ApplyLinkedPlugins(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    return nil, false
  }
  printer := shimprinter.NewPrinter(shimprinter.PrinterOptions{}, shimprinter.PrintHandlers{}, nil)
  cache := map[string]string{}
  for _, file := range prog.SourceFiles() {
    cache[apiOutputKey(opts.cwd, file.FileName())] = shimprinter.EmitSourceFile(printer, file)
  }
  return cache, true
}

// resolveServePath turns a request's file into an absolute path so apiOutputKey
// computes the same key buildServeCache stored.
func resolveServePath(cwd, file string) string {
  if filepath.IsAbs(file) {
    return file
  }
  return filepath.Join(cwd, file)
}
