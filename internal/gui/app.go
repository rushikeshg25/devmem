// Package gui implements the DevMem desktop application: a Fyne front-end over
// the same store and scan packages the CLI uses.
package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// windowSize is the default size of the main window.
var windowSize = fyne.NewSize(960, 640)

// appVersion is reported in the window; stamped from main at release time.
var appVersion = "dev"

// SetVersion sets the version string shown by the application.
func SetVersion(v string) {
	if v != "" {
		appVersion = v
	}
}

// Run launches the DevMem desktop application and blocks until the window is
// closed. dbPath is the SQLite database to read from and index into.
func Run(dbPath string) {
	a := app.NewWithID("dev.devmem.gui")
	w := a.NewWindow("DevMem " + appVersion)
	w.Resize(windowSize)

	svc, err := NewService(dbPath)
	if err != nil {
		w.SetContent(widget.NewLabel("Failed to open database: " + err.Error()))
		w.ShowAndRun()
		return
	}
	defer svc.Close()

	wipView, reloadWIP := newWIPView(svc, w)
	timelineView, reloadTimeline := newTimelineView(svc, w)
	scanView := newScanView(svc, w, func() {
		reloadWIP()
		reloadTimeline()
	})

	tabs := container.NewAppTabs(
		container.NewTabItem("Search", newSearchView(svc, w)),
		container.NewTabItem("At risk", wipView),
		container.NewTabItem("Timeline", timelineView),
		container.NewTabItem("Scan", scanView),
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	w.SetContent(tabs)
	w.ShowAndRun()
}
