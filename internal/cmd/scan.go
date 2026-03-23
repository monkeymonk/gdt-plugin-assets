package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/monkeymonk/gdt-assets/internal/report"
	"github.com/monkeymonk/gdt-assets/internal/scanner"
)

func cmdScan(args []string) int {
	fs := flag.NewFlagSet("scan", flag.ContinueOnError)
	format := fs.String("format", "table", "Output format: table, json, md, csv")
	types := fs.String("type", "", "Filter by types (comma-separated: image,audio,model)")
	hash := fs.Bool("hash", false, "Compute SHA256 hashes")
	profile := fs.String("profile", "", "Policy profile (reserved for future use)")
	_ = profile
	if err := fs.Parse(args); err != nil {
		return 1
	}

	opts := scanner.Options{Hash: *hash}
	if *types != "" {
		opts.Types = strings.Split(*types, ",")
	}

	assets, code := scanAssets(opts)
	if code != 0 {
		return code
	}

	report.FormatInventory(os.Stdout, assets, *format)
	fmt.Println()
	report.FormatSummary(os.Stdout, assets)
	return 0
}
