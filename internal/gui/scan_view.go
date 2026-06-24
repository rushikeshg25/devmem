package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// newScanView builds the scan tab: pick a workspace root, index it in the
// background, and report the result. onScanned is called after a successful
// scan so the other tabs can refresh from the updated database.
func newScanView(svc *Service, win fyne.Window, onScanned func()) fyne.CanvasObject {
	pathLabel := widget.NewLabel("No folder selected.")
	status := widget.NewLabel("")
	progress := widget.NewProgressBarInfinite()
	progress.Hide()

	var scanButton *widget.Button
	var selected string

	runScan := func() {
		if selected == "" {
			return
		}
		scanButton.Disable()
		progress.Show()
		status.SetText("Scanning " + selected + "…")

		go func() {
			res, err := svc.Scan(selected)
			fyne.Do(func() {
				progress.Hide()
				scanButton.Enable()
				if err != nil {
					status.SetText("")
					dialog.ShowError(err, win)
					return
				}
				status.SetText(fmt.Sprintf("Indexed %d repos across %d workspaces (%d new commits)",
					res.Repos, res.Workspaces, res.NewCommits))
				if onScanned != nil {
					onScanned()
				}
			})
		}()
	}

	scanButton = widget.NewButton("Scan", runScan)
	scanButton.Disable()

	pick := widget.NewButton("Choose workspace root…", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if uri == nil {
				return // cancelled
			}
			selected = uri.Path()
			pathLabel.SetText(selected)
			scanButton.Enable()
		}, win)
	})

	controls := container.NewHBox(pick, scanButton)
	intro := widget.NewLabel("Index every git checkout under a workspace root. Re-scanning is safe and incremental.")
	return container.NewVBox(intro, controls, pathLabel, progress, status)
}
