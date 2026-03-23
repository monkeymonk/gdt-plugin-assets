package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/monkeymonk/gdt-assets/internal/policy"
	"github.com/monkeymonk/gdt-assets/internal/rename"
	"github.com/monkeymonk/gdt-assets/internal/scanner"
)

func cmdRename(args []string) int {
	fs := flag.NewFlagSet("rename", flag.ContinueOnError)
	dryRun := fs.Bool("dry-run", true, "Preview renames without applying")
	apply := fs.Bool("apply", false, "Apply renames")
	if err := fs.Parse(args); err != nil {
		return 1
	}

	if *apply {
		*dryRun = false
	}

	pol := policy.LoadOrDefault(filepath.Join(projectRoot(), policy.FileName))
	assets, code := scanAssets(scanner.Options{})
	if code != 0 {
		return code
	}

	ops := rename.Plan(assets, pol.Naming.Case)
	if len(ops) == 0 {
		fmt.Println("No renames needed.")
		return 0
	}

	for _, op := range ops {
		fmt.Printf("  %s -> %s\n", op.OldPath, op.NewPath)
	}
	fmt.Printf("\n%d file(s) to rename\n", len(ops))

	if *dryRun {
		fmt.Println("\nDry run. Use --apply to execute.")
		return 0
	}

	errs := rename.Apply(ops)
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "  error: %v\n", e)
		}
		return 1
	}
	fmt.Println("Renames applied.")
	return 0
}
