package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var timelineLimit int

var timelineCmd = &cobra.Command{
	Use:   "timeline",
	Short: "Show recent activity across all indexed repositories",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := openStore()
		if err != nil {
			return err
		}
		defer s.Close()

		hits, err := s.Timeline(timelineLimit)
		if err != nil {
			return err
		}
		if len(hits) == 0 {
			fmt.Println("Nothing indexed yet — run `devmem scan <root>` first")
			return nil
		}
		for _, h := range hits {
			printHit(h)
		}
		return nil
	},
}

func init() {
	timelineCmd.Flags().IntVar(&timelineLimit, "limit", 30, "maximum number of commits")
	rootCmd.AddCommand(timelineCmd)
}
