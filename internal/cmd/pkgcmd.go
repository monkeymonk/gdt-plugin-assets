package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/monkeymonk/gdt-assets/internal/asset"
	"github.com/monkeymonk/gdt-assets/internal/diagnostic"
	"github.com/monkeymonk/gdt-assets/internal/policy"
	"github.com/monkeymonk/gdt-assets/internal/report"
	"github.com/monkeymonk/gdt-assets/internal/scanner"
)

func cmdPackage(args []string) int {
	pol := policy.LoadOrDefault(filepath.Join(projectRoot(), policy.FileName))
	assets, code := scanAssets(scanner.Options{})
	if code != 0 {
		return code
	}

	diags := &diagnostic.Set{}

	if pol.Package.Release.ForbidSourceFiles {
		for _, a := range assets {
			if a.Type == asset.TypeDocument {
				diags.Add(diagnostic.Diagnostic{
					Path:        a.Path,
					Severity:    diagnostic.Blocker,
					Rule:        "package.source_file",
					Message:     "source/authoring file present in project tree",
					Explanation: "Source files (.blend, .psd, etc.) should not ship in release builds",
				})
			}
		}
	}

	report.FormatDiagnostics(os.Stdout, diags, "table")
	fmt.Printf("\n%s\n", diags.Summary())

	if diags.HasBlockers() {
		fmt.Println("\nPackage validation FAILED")
		return 1
	}
	fmt.Println("\nPackage validation passed")
	return 0
}
