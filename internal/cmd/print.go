package cmd

import (
	"fmt"

	"github.com/rushikeshg25/devmem/internal/store"
)

// printHit renders a single commit search/timeline result.
func printHit(h store.SearchHit) {
	branch := h.Branch
	if branch == "" {
		branch = "-"
	}
	fmt.Printf("%s  %s  [%s]\n", h.CommitTime.Format("2006-01-02"), h.Repo, branch)
	fmt.Printf("    %s  %s\n", h.Hash[:min(7, len(h.Hash))], h.Message)
	fmt.Printf("    %s\n\n", h.Workspace)
}

// printRepoStatus renders a repo's current WIP state.
func printRepoStatus(rs store.RepoStatus) {
	flags := ""
	if rs.IsDirty {
		flags += " dirty"
	}
	if rs.AheadCount > 0 {
		flags += fmt.Sprintf(" ahead+%d", rs.AheadCount)
	}
	if rs.StashCount > 0 {
		flags += fmt.Sprintf(" stash:%d", rs.StashCount)
	}
	if rs.IsWorktree {
		flags += " worktree"
	}
	if flags == "" {
		flags = " clean"
	}
	fmt.Printf("%s  [%s]%s\n", rs.Name, rs.CurrentBranch, flags)
	fmt.Printf("    %s\n\n", rs.Path)
}
