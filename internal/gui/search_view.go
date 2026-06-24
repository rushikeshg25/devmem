package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/rushikeshg25/devmem/internal/store"
)

// newSearchView builds the commit-search tab: a query box over a results list.
func newSearchView(svc *Service, win fyne.Window) fyne.CanvasObject {
	var hits []store.SearchHit

	list := widget.NewList(
		func() int { return len(hits) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			h := hits[i]
			o.(*widget.Label).SetText(fmt.Sprintf("%s  %s  [%s]  %s",
				h.CommitTime.Format("2006-01-02"), h.Repo, branchOrDash(h.Branch), h.Message))
		},
	)

	status := widget.NewLabel("Type a term and press Enter to search.")

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Search commits, branches, repos…")
	entry.OnSubmitted = func(term string) {
		res, err := svc.Search(term)
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		hits = res
		list.Refresh()
		status.SetText(fmt.Sprintf("%d result(s) for %q", len(hits), term))
	}

	header := container.NewBorder(nil, nil, nil, nil, entry)
	return container.NewBorder(header, status, nil, nil, list)
}

// branchOrDash renders an empty branch as a dash, matching the CLI output.
func branchOrDash(branch string) string {
	if branch == "" {
		return "-"
	}
	return branch
}
