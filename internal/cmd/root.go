package cmd

import (
	"fmt"
	"os"

	"github.com/monkeymonk/gdt-assets/internal/asset"
	"github.com/monkeymonk/gdt-assets/internal/scanner"
)

func Usage() {
	fmt.Fprintln(os.Stderr, `gdt-assets - Asset operations, hygiene, and validation for game teams

Usage: gdt-assets <command> [subcommand] [options]

Commands:
  init       Initialize asset policy for a project
  scan       Inventory project assets
  lint       Run policy checks on assets
  report     Generate asset reports
  rename     Batch rename assets
  refs       Check and repair asset references
  dedupe     Detect duplicate assets
  package    Validate release packaging
  policy     Inspect and manage policy rules
  doctor     Check plugin health
  hook       Run lifecycle hook (internal)
  help       Show this help message

Environment:
  GDT_PROJECT_ROOT  Path to the Godot project root
  GDT_HOME          Path to gdt home directory`)
}

func Run(args []string) int {
	switch args[0] {
	case "init":
		return cmdInit(args[1:])
	case "scan":
		return cmdScan(args[1:])
	case "lint":
		return cmdLint(args[1:])
	case "report":
		return cmdReport(args[1:])
	case "rename":
		return cmdRename(args[1:])
	case "refs":
		return cmdRefs(args[1:])
	case "dedupe":
		return cmdDedupe(args[1:])
	case "package":
		return cmdPackage(args[1:])
	case "policy":
		return cmdPolicy(args[1:])
	case "doctor":
		return cmdDoctor(args[1:])
	case "hook":
		return cmdHook(args[1:])
	case "completions":
		return cmdCompletions(args[1:])
	case "help", "--help", "-h":
		Usage()
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", args[0])
		Usage()
		return 1
	}
}

func scanAssets(opts scanner.Options) ([]asset.Asset, int) {
	root := projectRoot()
	assets, err := scanner.Scan(root, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "scan error: %v\n", err)
		return nil, 1
	}
	return assets, 0
}

func projectRoot() string {
	root := os.Getenv("GDT_PROJECT_ROOT")
	if root != "" {
		return root
	}
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}
