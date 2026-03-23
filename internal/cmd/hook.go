package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/monkeymonk/gdt-assets/internal/analyzer"
	"github.com/monkeymonk/gdt-assets/internal/diagnostic"
	"github.com/monkeymonk/gdt-assets/internal/exitcode"
	"github.com/monkeymonk/gdt-assets/internal/policy"
	"github.com/monkeymonk/gdt-assets/internal/refs"
	"github.com/monkeymonk/gdt-assets/internal/scanner"
)

func cmdHook(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: gdt-assets hook <event>")
		return 1
	}

	event := args[0]
	switch event {
	case "after_new":
		return hookAfterNew()
	case "before_run":
		return hookBeforeRun()
	case "before_export":
		return hookBeforeExport()
	default:
		fmt.Printf("OK %s no action\n", event)
		return 0
	}
}

func hookAfterNew() int {
	root := projectRoot()
	policyPath := filepath.Join(root, policy.FileName)
	if _, err := os.Stat(policyPath); err == nil {
		fmt.Println("OK asset policy already exists")
		return 0
	}
	fmt.Println("OK run 'gdt assets init' to set up asset policy")
	return 0
}

func hookBeforeRun() int {
	root := projectRoot()

	broken, err := refs.FindBroken(root)
	if err != nil {
		fmt.Println("WARN could not check asset references")
		return 0
	}
	if len(broken) > 0 {
		fmt.Printf("WARN %d broken asset reference(s) found\n", len(broken))
	} else {
		fmt.Println("OK asset references intact")
	}
	return 0
}

func hookBeforeExport() int {
	root := projectRoot()
	pol := policy.LoadOrDefault(filepath.Join(root, policy.FileName))

	if profile := os.Getenv("GDT_ASSETS_PROFILE"); profile != "" {
		resolved, err := policy.ResolveProfile(pol, profile)
		if err != nil {
			fmt.Printf("WARN invalid profile %q: %v\n", profile, err)
		} else {
			pol = resolved
		}
	}

	assets, err := scanner.Scan(root, scanner.Options{})
	if err != nil {
		fmt.Println("FAIL asset scan failed")
		return 1
	}

	diags := analyzer.RunAll(analyzer.DefaultAnalyzers(), assets, pol)

	if diags.HasBlockers() {
		fmt.Printf("FAIL %d asset blocker(s) prevent export\n", diags.Count(diagnostic.Blocker))
		return exitcode.ErrBlockers
	}
	if diags.HasErrors() {
		fmt.Printf("WARN %d asset error(s) detected\n", diags.Count(diagnostic.Error))
		return exitcode.ErrDiagnostics
	}
	fmt.Printf("OK %d assets validated\n", len(assets))
	return exitcode.OK
}
