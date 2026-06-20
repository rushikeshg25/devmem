package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo <name>",
	Short: "Show every indexed checkout of a repository and its recent commits",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := openStore()
		if err != nil {
			return err
		}
		defer s.Close()

		repos, err := s.FindRepos(args[0])
		if err != nil {
			return err
		}
		if len(repos) == 0 {
			fmt.Printf("No indexed repo named %q\n", args[0])
			return nil
		}

		for _, rs := range repos {
			printRepoStatus(rs)
			commits, err := s.RecentCommits(rs.Path, 5)
			if err != nil {
				return err
			}
			for _, c := range commits {
				fmt.Printf("    %s  %s  %s\n", c.CommitTime.Format("2006-01-02"), c.Hash[:min(7, len(c.Hash))], c.Message)
			}
			fmt.Println()
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(repoCmd)
}
