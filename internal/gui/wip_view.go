package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/rushikeshg25/devmem/internal/store"
)

// newWIPView builds the at-risk tab: checkouts with uncommitted, unpushed or
// stashed work — the work most likely to be lost when a workspace is deleted.
// It returns the content and a reload func callers can invoke after a scan.
func newWIPView(svc *Service, win fyne.Window) (fyne.CanvasObject, func()) {
	var repos []store.RepoStatus

	list := widget.NewList(
		func() int { return len(repos) },
		twoLineCell,
		func(i widget.ListItemID, o fyne.CanvasObject) {
			rs := repos[i]
			setTwoLineCell(o, repoTitle(rs), rs.Path)
		},
	)

	status := widget.NewLabel("Showing all at-risk checkouts. Filter by name or branch above.")

	load := func(term string) {
		res, err := svc.WIP(term)
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		repos = res
		list.Refresh()
		status.SetText(fmt.Sprintf("%d at-risk checkout(s)", len(repos)))
	}

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Filter at-risk work by repo or branch…")
	entry.OnSubmitted = load

	// reload re-reads the unfiltered list, e.g. after a scan.
	reload := func() { load("") }

	// Populate immediately so the tab is useful before any filter is typed.
	reload()

	return container.NewBorder(entry, status, nil, nil, list), reload
}
