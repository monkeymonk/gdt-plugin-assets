package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/monkeymonk/gdt-assets/internal/policy"
)

func cmdInit(args []string) int {
	root := projectRoot()
	policyPath := filepath.Join(root, policy.FileName)

	if _, err := os.Stat(policyPath); err == nil {
		fmt.Fprintf(os.Stderr, "Policy file already exists: %s\n", policyPath)
		return 1
	}

	content, err := policy.MarshalDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error generating policy: %v\n", err)
		return 1
	}

	if err := os.WriteFile(policyPath, content, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing policy: %v\n", err)
		return 1
	}

	for _, flag := range args {
		if flag == "--with-sample-folders" {
			pol := policy.Default()
			dirs := []string{
				pol.Folders.Images, pol.Folders.Audio,
				pol.Folders.Models, pol.Folders.Vectors,
				pol.Folders.Fonts, pol.Folders.Source,
			}
			for _, d := range dirs {
				os.MkdirAll(filepath.Join(root, d), 0755)
			}
			fmt.Println("Created sample asset folders")
		}
	}

	fmt.Printf("Created %s\n", policyPath)
	return 0
}

