package cmd

import (
	"fmt"

	"github.com/rushikeshg25/devmem/internal/scan"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan <root>",
	Short: "Discover and index git repositories under a workspace root",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := openStore()
		if err != nil {
			return err
		}
		defer s.Close()

		res, err := scan.Run(s, args[0])
		if err != nil {
			return err
		}
		fmt.Printf("Indexed %d repos across %d workspaces (%d new commits)\n",
			res.Repos, res.Workspaces, res.NewCommits)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
