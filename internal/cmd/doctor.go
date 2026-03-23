package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/monkeymonk/gdt-assets/internal/analyzer"
	"github.com/monkeymonk/gdt-assets/internal/asset"
	"github.com/monkeymonk/gdt-assets/internal/diagnostic"
	"github.com/monkeymonk/gdt-assets/internal/exitcode"
	"github.com/monkeymonk/gdt-assets/internal/policy"
	"github.com/monkeymonk/gdt-assets/internal/refs"
	"github.com/monkeymonk/gdt-assets/internal/scanner"
)

func cmdDoctor(args []string) int {
	if len(args) == 0 || args[0] != "check" {
		fmt.Fprintln(os.Stderr, "usage: gdt-assets doctor check")
		return exitcode.ErrUsage
	}

	root := projectRoot()
	ok := true

	// 1. Policy check
	policyPath := filepath.Join(root, policy.FileName)
	var pol *policy.Policy
	if _, err := os.Stat(policyPath); os.IsNotExist(err) {
		fmt.Println("WARN no assets.policy.toml found; run: gdt assets init")
	} else if err != nil {
		fmt.Printf("FAIL cannot read policy file: %v\n", err)
		ok = false
	} else {
		loaded, err := policy.Load(policyPath)
		if err != nil {
			fmt.Printf("FAIL invalid policy file: %v\n", err)
			ok = false
		} else {
			fmt.Println("OK   assets.policy.toml valid")
			pol = loaded
		}
	}

	// 2. Asset scan
	assets, err := scanner.Scan(root, scanner.Options{})
	if err != nil {
		fmt.Printf("FAIL asset scan failed: %v\n", err)
		ok = false
	} else {
		byType := make(map[asset.AssetType]int)
		for _, a := range assets {
			byType[a.Type]++
		}
		fmt.Printf("OK   found %d assets", len(assets))
		if len(byType) > 0 {
			fmt.Print(" (")
			first := true
			for t := asset.AssetType(1); t <= asset.TypeEngineResource; t++ {
				if n, exists := byType[t]; exists {
					if !first {
						fmt.Print(", ")
					}
					fmt.Printf("%s: %d", t, n)
					first = false
				}
			}
			fmt.Print(")")
		}
		fmt.Println()
	}

	// 3. Lint summary (only if policy loaded)
	if pol != nil && len(assets) > 0 {
		diags := analyzer.RunAll(analyzer.DefaultAnalyzers(), assets, pol)
		blockers := diags.Count(diagnostic.Blocker)
		errors := diags.Count(diagnostic.Error)
		warnings := diags.Count(diagnostic.Warning)
		if blockers > 0 {
			fmt.Printf("FAIL %d blocker(s), %d error(s), %d warning(s)\n", blockers, errors, warnings)
			ok = false
		} else if errors > 0 {
			fmt.Printf("WARN %d error(s), %d warning(s)\n", errors, warnings)
		} else if warnings > 0 {
			fmt.Printf("OK   %d warning(s), no errors\n", warnings)
		} else {
			fmt.Println("OK   no lint issues")
		}
	}

	// 4. Broken refs
	broken, err := refs.FindBroken(root)
	if err != nil {
		fmt.Println("WARN could not check references")
	} else if len(broken) > 0 {
		fmt.Printf("WARN %d broken reference(s)\n", len(broken))
	} else {
		fmt.Println("OK   no broken references")
	}

	if !ok {
		return exitcode.ErrDiagnostics
	}
	return exitcode.OK
}
