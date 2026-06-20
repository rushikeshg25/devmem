package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	searchLimit int
	searchWIP   bool
)

var searchCmd = &cobra.Command{
	Use:   "search <term>",
	Short: "Search commit messages, branches and repositories",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := openStore()
		if err != nil {
			return err
		}
		defer s.Close()

		term := args[0]

		if searchWIP {
			repos, err := s.SearchWIP(term)
			if err != nil {
				return err
			}
			if len(repos) == 0 {
				fmt.Printf("No uncommitted/unpushed work matching %q\n", term)
				return nil
			}
			for _, rs := range repos {
				printRepoStatus(rs)
			}
			return nil
		}

		hits, err := s.Search(term, searchLimit)
		if err != nil {
			return err
		}
		if len(hits) == 0 {
			fmt.Printf("No matches for %q\n", term)
			return nil
		}
		for _, h := range hits {
			printHit(h)
		}
		return nil
	},
}

func init() {
	searchCmd.Flags().IntVar(&searchLimit, "limit", 50, "maximum number of results")
	searchCmd.Flags().BoolVar(&searchWIP, "wip", false, "only show dirty/unpushed/stashed work")
	rootCmd.AddCommand(searchCmd)
}
