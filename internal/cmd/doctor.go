package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/monkeymonk/gdt-assets/internal/policy"
	"github.com/monkeymonk/gdt-assets/internal/scanner"
)

func cmdDoctor(args []string) int {
	if len(args) == 0 || args[0] != "check" {
		fmt.Fprintln(os.Stderr, "usage: gdt-assets doctor check")
		return 1
	}
	ok := true

	root := projectRoot()
	policyPath := filepath.Join(root, policy.FileName)
	if _, err := os.Stat(policyPath); os.IsNotExist(err) {
		fmt.Println("WARN no assets.policy.toml found; run: gdt assets init")
	} else if err != nil {
		fmt.Printf("FAIL cannot read policy file: %v\n", err)
		ok = false
	} else {
		_, err := policy.Load(policyPath)
		if err != nil {
			fmt.Printf("FAIL invalid policy file: %v\n", err)
			ok = false
		} else {
			fmt.Println("OK assets.policy.toml valid")
		}
	}

	assets, err := scanner.Scan(root, scanner.Options{})
	if err != nil {
		fmt.Printf("FAIL asset scan failed: %v\n", err)
		ok = false
	} else {
		fmt.Printf("OK found %d assets\n", len(assets))
	}

	if !ok {
		return 1
	}
	return 0
}
