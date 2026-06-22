// Command ttscgraph serves a checker-resolved code graph and diagnostics to
// coding agents over the Model Context Protocol (JSON-RPC 2.0 on stdio). It
// builds one resident tsgo Program for the project and answers every tool call
// from that warm handle, so a query is a method call on an already-built checker
// rather than a fresh compile or an external language-server round-trip.
//
// The JavaScript launcher (@ttsc/graph) resolves the per-platform native binary
// and spawns `ttscgraph --stdio`; an agent's MCP client drives it over stdio.
// Everything here is deliberately small: flag parsing, version metadata, and a
// single delegation to the resident MCP server.
package main

import (
  "flag"
  "fmt"
  "io"
  "os"
  "runtime"
  "strings"

  "github.com/samchon/ttsc/packages/ttsc/driver"
  "github.com/samchon/ttsc/packages/ttsc/internal/graph/mcp"
)

// Build metadata; overwritten via -ldflags in release builds.
var (
  version = "0.0.0-dev"
  commit  = "dev"
  date    = "unknown"
)

// Package-level streams so command tests can capture I/O without patching the
// os globals.
var (
  stdout io.Writer = os.Stdout
  stderr io.Writer = os.Stderr
  stdin  io.Reader = os.Stdin
)

// Seams command tests substitute to avoid building a real program or touching
// the process working directory.
var (
  loadProgram = driver.LoadProgram
  getwd       = os.Getwd
)

func main() {
  os.Exit(run(os.Args[1:]))
}

// run dispatches top-level flags and returns an exit code. Called by main with
// os.Args[1:] and overridden in tests with a synthetic argument slice.
func run(args []string) int {
  if len(args) > 0 {
    switch args[0] {
    case "-h", "--help", "help":
      printHelp(stdout)
      return 0
    case "-v", "--version", "version":
      printVersion(stdout)
      return 0
    }
  }
  return runServe(args)
}

// runServe parses serve flags, loads the resident program, and serves MCP over
// stdio. It returns 0 on a clean EOF shutdown, 1 on a load or runtime error, and
// 2 on invalid invocation.
func runServe(args []string) int {
  fs := flag.NewFlagSet("ttscgraph", flag.ContinueOnError)
  fs.SetOutput(stderr)
  _ = fs.Bool("stdio", true, "serve MCP over stdin/stdout (the only transport)")
  cwdFlag := fs.String("cwd", "", "project root (defaults to process cwd)")
  tsconfigFlag := fs.String("tsconfig", "tsconfig.json", "project tsconfig path")
  if err := fs.Parse(args); err != nil {
    return 2
  }

  cwd := strings.TrimSpace(*cwdFlag)
  if cwd == "" {
    resolved, err := getwd()
    if err != nil {
      fmt.Fprintf(stderr, "ttscgraph: could not resolve working directory: %v\n", err)
      return 2
    }
    cwd = resolved
  }

  mcp.Version = version
  prog, _, err := loadProgram(cwd, strings.TrimSpace(*tsconfigFlag), driver.LoadProgramOptions{})
  if err != nil {
    fmt.Fprintf(stderr, "ttscgraph: could not load project: %v\n", err)
    return 1
  }
  if prog == nil {
    fmt.Fprintf(stderr, "ttscgraph: could not load project %q\n", strings.TrimSpace(*tsconfigFlag))
    return 1
  }
  defer func() { _ = prog.Close() }()

  if err := mcp.NewServer(prog).Serve(stdin, stdout); err != nil {
    fmt.Fprintf(stderr, "ttscgraph: %v\n", err)
    return 1
  }
  return 0
}

func printVersion(w io.Writer) {
  fmt.Fprintf(
    w,
    "ttscgraph %s (commit %s, built %s, %s/%s, go %s)\n",
    version,
    commit,
    date,
    runtime.GOOS,
    runtime.GOARCH,
    runtime.Version(),
  )
}

func printHelp(w io.Writer) {
  fmt.Fprintln(w, strings.TrimSpace(`
ttscgraph — checker-resolved code graph + diagnostics over MCP for ttsc.

Usage:
  ttscgraph --stdio
  ttscgraph --version
  ttscgraph --help

Options:
  --stdio              Serve MCP over stdin/stdout (the only transport).
  --cwd <dir>          Project root (defaults to the process working directory).
  --tsconfig <path>    Project tsconfig path (default: tsconfig.json).

Typical embedding:
  An agent's MCP client spawns ttscgraph through the @ttsc/graph launcher, which
  resolves the per-platform native binary. ttscgraph builds one resident tsgo
  Program for the project and answers graph_explore / graph_diagnostics calls
  from that warm checker. Usage guidance is delivered in the MCP initialize
  response; nothing is written to your agent config files.
`))
}
