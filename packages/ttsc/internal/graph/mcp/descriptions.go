package mcp

import (
	_ "embed"
	"strings"
)

// Tool and schema descriptions are embedded from Markdown files so prompt
// tuning stays reviewable text instead of long Go string literals.
//
//go:embed descriptions/graph_explore.md
var graphExploreDescriptionMarkdown string

//go:embed descriptions/graph_explore_query.md
var graphExploreQueryDescriptionMarkdown string

//go:embed descriptions/graph_diagnostics.md
var graphDiagnosticsDescriptionMarkdown string

//go:embed descriptions/graph_diagnostics_file.md
var graphDiagnosticsFileDescriptionMarkdown string

var (
	graphExploreDescription         = strings.TrimSpace(graphExploreDescriptionMarkdown)
	graphExploreQueryDescription    = strings.TrimSpace(graphExploreQueryDescriptionMarkdown)
	graphDiagnosticsDescription     = strings.TrimSpace(graphDiagnosticsDescriptionMarkdown)
	graphDiagnosticsFileDescription = strings.TrimSpace(graphDiagnosticsFileDescriptionMarkdown)
)
