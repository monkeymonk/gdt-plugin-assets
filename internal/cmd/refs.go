package cmd

import (
	"fmt"
	"os"

	"github.com/monkeymonk/gdt-assets/internal/refs"
)

func cmdRefs(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: gdt assets refs <check|repair>")
		return 1
	}

	root := projectRoot()

	switch args[0] {
	case "check":
		broken, err := refs.FindBroken(root)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		if len(broken) == 0 {
			fmt.Println("No broken references found.")
			return 0
		}
		for _, b := range broken {
			fmt.Printf("  BROKEN  %s:%d -> %s\n", b.Source, b.Line, b.Target)
		}
		fmt.Printf("\n%d broken reference(s)\n", len(broken))
		return 1

	case "repair":
		fmt.Println("Interactive repair not yet implemented.")
		fmt.Println("Use: gdt assets refs check to identify broken references.")
		return 0

	default:
		fmt.Fprintf(os.Stderr, "unknown refs subcommand: %s\n", args[0])
		return 1
	}
}
