package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/monkeymonk/gdt-assets/internal/analyzer"
	"github.com/monkeymonk/gdt-assets/internal/policy"
	"github.com/monkeymonk/gdt-assets/internal/report"
	"github.com/monkeymonk/gdt-assets/internal/scanner"
)

func cmdReport(args []string) int {
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	format := fs.String("format", "table", "Output format: table, json, md, csv")
	hash := fs.Bool("hash", false, "Include hashes")
	if err := fs.Parse(args); err != nil {
		return 1
	}

	pol := policy.LoadOrDefault(filepath.Join(projectRoot(), policy.FileName))
	assets, code := scanAssets(scanner.Options{Hash: *hash})
	if code != 0 {
		return code
	}

	fmt.Println("=== Asset Inventory ===")
	report.FormatInventory(os.Stdout, assets, *format)
	fmt.Println()
	report.FormatSummary(os.Stdout, assets)

	fmt.Println("\n=== Lint Results ===")
	diags := analyzer.RunAll(analyzer.DefaultAnalyzers(), assets, pol)
	report.FormatDiagnostics(os.Stdout, diags, *format)
	fmt.Printf("\n%s\n", diags.Summary())

	return 0
}
