package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/rushikeshg25/devmem/internal/store"
)

// newTimelineView builds the timeline tab: the most recent commits across every
// indexed repo, newest first, with a button to refresh after a scan.
func newTimelineView(svc *Service, win fyne.Window) fyne.CanvasObject {
	var hits []store.SearchHit

	list := widget.NewList(
		func() int { return len(hits) },
		newCommitCell,
		func(i widget.ListItemID, o fyne.CanvasObject) { updateCommitCell(o, hits[i]) },
	)

	status := widget.NewLabel("")

	load := func() {
		res, err := svc.Timeline()
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		hits = res
		list.Refresh()
		status.SetText(fmt.Sprintf("%d recent commit(s)", len(hits)))
	}

	refresh := widget.NewButton("Refresh", load)
	load()

	header := container.NewBorder(nil, nil, nil, refresh, widget.NewLabel("Recent activity across all workspaces"))
	return container.NewBorder(header, status, nil, nil, list)
}
