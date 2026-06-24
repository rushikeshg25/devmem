// Package cmd defines the devmem CLI commands.
package cmd

import (
	"github.com/rushikeshg25/devmem/internal/dbpath"
	"github.com/spf13/cobra"
)

var dbPath string

// SetVersion sets the version reported by `devmem --version`. It is wired up
// from main and stamped at release time via -ldflags.
func SetVersion(v string) {
	if v != "" {
		rootCmd.Version = v
	}
}

// rootCmd is the base command for devmem.
var rootCmd = &cobra.Command{
	Use:   "devmem",
	Short: "A searchable memory layer for your git workspaces",
	Long: `devmem indexes git repositories across your workspace folders and lets you
search commits, branches and repositories — so you can recall where you did
a piece of work, even after the workspace is gone.`,
	Version: "dev",
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dbPath, "db", dbpath.Default(), "path to the devmem sqlite database")
}
