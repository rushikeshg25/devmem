package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Remove indexed repos whose working directory no longer exists",
	Long: `Deletes index entries for checkouts that have been removed from disk
(e.g. a deleted workspace folder). Commit history that was already indexed for
surviving checkouts is kept, so prior work stays searchable elsewhere.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := openStore()
		if err != nil {
			return err
		}
		defer s.Close()

		refs, err := s.ListRepoRefs()
		if err != nil {
			return err
		}

		removed := 0
		for _, r := range refs {
			if _, err := os.Stat(r.Path); os.IsNotExist(err) {
				if err := s.DeleteRepo(r.ID); err != nil {
					return err
				}
				fmt.Printf("removed %s\n", r.Path)
				removed++
			}
		}

		emptyWs, err := s.DeleteEmptyWorkspaces()
		if err != nil {
			return err
		}
		fmt.Printf("Pruned %d repos and %d empty workspaces\n", removed, emptyWs)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pruneCmd)
}
