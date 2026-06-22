package main

import (
	"fmt"
	"os"

	"github.com/rushikeshg25/devmem/internal/cmd"
)

// version is overridden at release time via -ldflags.
var version = "dev"

func main() {
	cmd.SetVersion(version)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "devmem:", err)
		os.Exit(1)
	}
}
