package main

import (
	"os"

	"github.com/monkeymonk/gdt-assets/internal/cmd"
)

func main() {
	if len(os.Args) < 2 {
		cmd.Usage()
		os.Exit(1)
	}
	code := cmd.Run(os.Args[1:])
	if code != 0 {
		os.Exit(code)
	}
}
