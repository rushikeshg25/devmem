package gui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/rushikeshg25/devmem/internal/store"
)

// dateLayout matches the CLI's date rendering (internal/cmd/print.go).
const dateLayout = "2006-01-02"

// branchOrDash renders an empty branch as a dash, matching the CLI output.
func branchOrDash(branch string) string {
	if branch == "" {
		return "-"
	}
	return branch
}

// shortHash returns the first 7 characters of a commit hash, like the CLI.
func shortHash(hash string) string {
	if len(hash) < 7 {
		return hash
	}
	return hash[:7]
}

// hitTitle is the primary line for a commit row: date, repo and branch.
func hitTitle(h store.SearchHit) string {
	return fmt.Sprintf("%s  %s  [%s]", h.CommitTime.Format(dateLayout), h.Repo, branchOrDash(h.Branch))
}

// hitSubtitle is the secondary line for a commit row: short hash and message.
func hitSubtitle(h store.SearchHit) string {
	return fmt.Sprintf("%s  %s", shortHash(h.Hash), h.Message)
}

// repoFlags summarizes a checkout's at-risk state, mirroring the CLI's
// dirty/ahead/stash/worktree badges (internal/cmd/print.go).
func repoFlags(rs store.RepoStatus) string {
	var flags []string
	if rs.IsDirty {
		flags = append(flags, "dirty")
	}
	if rs.AheadCount > 0 {
		flags = append(flags, fmt.Sprintf("ahead+%d", rs.AheadCount))
	}
	if rs.StashCount > 0 {
		flags = append(flags, fmt.Sprintf("stash:%d", rs.StashCount))
	}
	if rs.IsWorktree {
		flags = append(flags, "worktree")
	}
	if len(flags) == 0 {
		return "clean"
	}
	return strings.Join(flags, " ")
}

// repoTitle is the primary line for a repo row: name, branch and flags.
func repoTitle(rs store.RepoStatus) string {
	return fmt.Sprintf("%s  [%s]  %s", rs.Name, branchOrDash(rs.CurrentBranch), repoFlags(rs))
}

// twoLineCell is a list cell with a bold title over a muted subtitle.
func twoLineCell() fyne.CanvasObject {
	title := widget.NewLabel("")
	title.TextStyle = fyne.TextStyle{Bold: true}
	subtitle := widget.NewLabel("")
	return container.NewVBox(title, subtitle)
}

// setTwoLineCell updates a cell created by twoLineCell.
func setTwoLineCell(o fyne.CanvasObject, title, subtitle string) {
	box := o.(*fyne.Container)
	box.Objects[0].(*widget.Label).SetText(title)
	box.Objects[1].(*widget.Label).SetText(subtitle)
}

// newCommitCell / updateCommitCell render a SearchHit as a two-line row.
func newCommitCell() fyne.CanvasObject { return twoLineCell() }

func updateCommitCell(o fyne.CanvasObject, h store.SearchHit) {
	setTwoLineCell(o, hitTitle(h), hitSubtitle(h))
}
