package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/monkeymonk/gdt-assets/internal/analyzer"
	"github.com/monkeymonk/gdt-assets/internal/exitcode"
	"github.com/monkeymonk/gdt-assets/internal/policy"
	"github.com/monkeymonk/gdt-assets/internal/report"
	"github.com/monkeymonk/gdt-assets/internal/scanner"
)

func cmdLint(args []string) int {
	fs := flag.NewFlagSet("lint", flag.ContinueOnError)
	format := fs.String("format", "table", "Output format: table, json, csv")
	profile := fs.String("profile", "", "Policy profile to apply (e.g., mobile, release)")
	if err := fs.Parse(args); err != nil {
		return 1
	}

	subcmd := "all"
	if fs.NArg() > 0 {
		subcmd = fs.Arg(0)
	}

	pol := policy.LoadOrDefault(filepath.Join(projectRoot(), policy.FileName))

	if *profile != "" {
		resolved, err := policy.ResolveProfile(pol, *profile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return exitcode.ErrPolicy
		}
		pol = resolved
	}

	assets, code := scanAssets(scanner.Options{})
	if code != 0 {
		return code
	}

	var analyzers []analyzer.Analyzer
	switch subcmd {
	case "all":
		analyzers = analyzer.DefaultAnalyzers()
	case "names":
		analyzers = []analyzer.Analyzer{&analyzer.NameAnalyzer{}}
	case "structure":
		analyzers = []analyzer.Analyzer{&analyzer.StructureAnalyzer{}}
	case "images":
		analyzers = []analyzer.Analyzer{&analyzer.ImageAnalyzer{}}
	case "audio":
		analyzers = []analyzer.Analyzer{&analyzer.AudioAnalyzer{}}
	case "models":
		analyzers = []analyzer.Analyzer{&analyzer.ModelAnalyzer{}}
	default:
		fmt.Fprintf(os.Stderr, "unknown lint target: %s\n", subcmd)
		return 1
	}

	diags := analyzer.RunAll(analyzers, assets, pol)
	report.FormatDiagnostics(os.Stdout, diags, *format)

	fmt.Printf("\n%s\n", diags.Summary())

	if diags.HasBlockers() {
		return exitcode.ErrBlockers
	}
	if diags.HasErrors() {
		return exitcode.ErrDiagnostics
	}
	return exitcode.OK
}
