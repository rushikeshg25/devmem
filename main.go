package main

import (
	"fmt"
	"os"

	"github.com/rushikeshg25/devmem/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "devmem:", err)
		os.Exit(1)
	}
}
