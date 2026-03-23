package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/monkeymonk/gdt-assets/internal/exitcode"
	"github.com/monkeymonk/gdt-assets/internal/policy"
	"github.com/monkeymonk/gdt-assets/internal/rename"
	"github.com/monkeymonk/gdt-assets/internal/scanner"
)

func cmdRename(args []string) int {
	fs := flag.NewFlagSet("rename", flag.ContinueOnError)
	dryRun := fs.Bool("dry-run", true, "Preview renames without applying")
	apply := fs.Bool("apply", false, "Apply renames")
	rollbackFile := fs.String("rollback", "", "Rollback from manifest file")
	if err := fs.Parse(args); err != nil {
		return exitcode.ErrUsage
	}

	if *rollbackFile != "" {
		ops, err := rename.LoadManifest(*rollbackFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading manifest: %v\n", err)
			return exitcode.ErrUsage
		}
		fmt.Printf("Rolling back %d rename(s)...\n", len(ops))
		errs := rename.Rollback(ops)
		if len(errs) > 0 {
			for _, e := range errs {
				fmt.Fprintf(os.Stderr, "  error: %v\n", e)
			}
			return exitcode.ErrUsage
		}
		fmt.Println("Rollback complete.")
		return exitcode.OK
	}

	if *apply {
		*dryRun = false
	}

	pol := policy.LoadOrDefault(filepath.Join(projectRoot(), policy.FileName))
	assets, code := scanAssets(scanner.Options{})
	if code != 0 {
		return code
	}

	plan := rename.BuildPlan(assets, pol.Naming.Case)
	if len(plan.Ops) == 0 {
		fmt.Println("No renames needed.")
		return exitcode.OK
	}

	if len(plan.Collisions) > 0 {
		fmt.Printf("\n%d collision(s) detected:\n", len(plan.Collisions))
		for _, c := range plan.Collisions {
			fmt.Printf("  COLLISION  %s <- %v\n", c.Target, c.Sources)
		}
		fmt.Println("\nResolve collisions before applying.")
		return exitcode.ErrDiagnostics
	}

	for _, op := range plan.Ops {
		fmt.Printf("  %s -> %s\n", op.OldPath, op.NewPath)
	}
	fmt.Printf("\n%d file(s) to rename\n", len(plan.Ops))

	if *dryRun {
		fmt.Println("\nDry run. Use --apply to execute.")
		return exitcode.OK
	}

	manifestName := fmt.Sprintf(".assets-rollback-%s.json", time.Now().Format("20060102-150405"))
	manifestPath := filepath.Join(projectRoot(), manifestName)
	if err := rename.WriteManifest(manifestPath, plan.Ops); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not write rollback manifest: %v\n", err)
	} else {
		fmt.Printf("Rollback manifest: %s\n", manifestName)
	}

	errs := rename.Apply(plan.Ops)
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "  error: %v\n", e)
		}
		return exitcode.ErrUsage
	}
	fmt.Println("Renames applied.")
	return exitcode.OK
}
