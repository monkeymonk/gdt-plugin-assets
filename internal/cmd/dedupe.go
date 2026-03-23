package cmd

import (
	"flag"
	"fmt"

	"github.com/monkeymonk/gdt-assets/internal/dedupe"
	"github.com/monkeymonk/gdt-assets/internal/scanner"
)

func cmdDedupe(args []string) int {
	fs := flag.NewFlagSet("dedupe", flag.ContinueOnError)
	byName := fs.Bool("name", false, "Find duplicates by filename")
	if err := fs.Parse(args); err != nil {
		return 1
	}

	assets, code := scanAssets(scanner.Options{Hash: !*byName})
	if code != 0 {
		return code
	}

	var groups []dedupe.DupeGroup
	if *byName {
		groups = dedupe.FindByName(assets)
	} else {
		groups = dedupe.FindExact(assets)
	}

	if len(groups) == 0 {
		fmt.Println("No duplicates found.")
		return 0
	}

	for _, g := range groups {
		label := g.Key
		if len(label) > 16 {
			label = label[:16]
		}
		fmt.Printf("  Duplicate group [%s]:\n", label)
		for _, p := range g.Paths {
			fmt.Printf("    %s\n", p)
		}
	}
	fmt.Printf("\n%d duplicate group(s)\n", len(groups))
	return 0
}
