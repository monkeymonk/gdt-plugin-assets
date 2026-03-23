package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/monkeymonk/gdt-assets/internal/exitcode"
	"github.com/monkeymonk/gdt-assets/internal/refs"
	"github.com/monkeymonk/gdt-assets/internal/rename"
)

func cmdRefs(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: gdt assets refs <check|repair>")
		return exitcode.ErrUsage
	}

	root := projectRoot()

	switch args[0] {
	case "check":
		broken, err := refs.FindBroken(root)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return exitcode.ErrUsage
		}
		if len(broken) == 0 {
			fmt.Println("No broken references found.")
			return exitcode.OK
		}
		for _, b := range broken {
			fmt.Printf("  BROKEN  %s:%d -> %s\n", b.Source, b.Line, b.Target)
		}
		fmt.Printf("\n%d broken reference(s)\n", len(broken))
		return exitcode.ErrDiagnostics

	case "repair":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: gdt assets refs repair --from-manifest <rollback.json>")
			fmt.Fprintln(os.Stderr, "       Rewrites res:// references based on a rename manifest.")
			return exitcode.ErrUsage
		}

		fs := flag.NewFlagSet("refs-repair", flag.ContinueOnError)
		manifest := fs.String("from-manifest", "", "Rename manifest to derive ref mappings from")
		dryRun := fs.Bool("dry-run", true, "Preview repairs without applying")
		applyFlag := fs.Bool("apply", false, "Apply repairs")
		if err := fs.Parse(args[1:]); err != nil {
			return exitcode.ErrUsage
		}
		if *applyFlag {
			*dryRun = false
		}

		if *manifest == "" {
			fmt.Fprintln(os.Stderr, "error: --from-manifest is required")
			return exitcode.ErrUsage
		}

		renameOps, err := rename.LoadManifest(*manifest)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading manifest: %v\n", err)
			return exitcode.ErrUsage
		}

		pairs := make([]refs.RenamePair, len(renameOps))
		for i, op := range renameOps {
			pairs[i] = refs.RenamePair{OldPath: op.OldPath, NewPath: op.NewPath}
		}

		engineFiles, err := refs.FindEngineFiles(root)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error scanning project: %v\n", err)
			return exitcode.ErrUsage
		}

		repairOps, err := refs.PlanRepairFromRenames(root, pairs, engineFiles)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error planning repair: %v\n", err)
			return exitcode.ErrUsage
		}

		if len(repairOps) == 0 {
			fmt.Println("No references to repair.")
			return exitcode.OK
		}

		for _, op := range repairOps {
			rel, _ := filepath.Rel(root, op.File)
			fmt.Printf("  %s:%d  %s -> %s\n", rel, op.Line, op.OldRef, op.NewRef)
		}
		fmt.Printf("\n%d reference(s) to repair\n", len(repairOps))

		if *dryRun {
			fmt.Println("\nDry run. Use --apply to execute.")
			return exitcode.OK
		}

		if err := refs.ApplyRepair(repairOps); err != nil {
			fmt.Fprintf(os.Stderr, "error applying repairs: %v\n", err)
			return exitcode.ErrUsage
		}
		fmt.Println("References repaired.")
		return exitcode.OK

	default:
		fmt.Fprintf(os.Stderr, "unknown refs subcommand: %s\n", args[0])
		return exitcode.ErrUsage
	}
}
