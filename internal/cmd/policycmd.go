package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/monkeymonk/gdt-assets/internal/policy"
)

func cmdPolicy(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: gdt assets policy <show|validate>")
		return 1
	}

	root := projectRoot()
	policyPath := filepath.Join(root, policy.FileName)

	switch args[0] {
	case "show":
		data, err := os.ReadFile(policyPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "No policy file found at %s\n", policyPath)
			fmt.Fprintln(os.Stderr, "Run: gdt assets init")
			return 1
		}
		fmt.Print(string(data))
		return 0

	case "validate":
		_, err := policy.Load(policyPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Policy validation failed: %v\n", err)
			return 1
		}
		fmt.Println("Policy file is valid.")
		return 0

	default:
		fmt.Fprintf(os.Stderr, "unknown policy subcommand: %s\n", args[0])
		return 1
	}
}
