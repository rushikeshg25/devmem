package gui

import (
	"os/exec"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// openInFileManager opens path in the platform's file manager. On Linux this is
// xdg-open, which Ubuntu/GNOME route to Files (Nautilus).
func openInFileManager(path string) error {
	var name string
	switch runtime.GOOS {
	case "darwin":
		name = "open"
	case "windows":
		name = "explorer"
	default:
		name = "xdg-open"
	}
	return exec.Command(name, path).Start()
}

// showPathActions pops a dialog letting the user copy a checkout's path or open
// it in the file manager. title describes the selected row.
func showPathActions(win fyne.Window, title, path string) {
	pathLabel := widget.NewLabel(path)
	pathLabel.Wrapping = fyne.TextWrapBreak

	copyBtn := widget.NewButton("Copy path", func() {
		fyne.CurrentApp().Clipboard().SetContent(path)
	})
	openBtn := widget.NewButton("Open folder", func() {
		if err := openInFileManager(path); err != nil {
			dialog.ShowError(err, win)
		}
	})

	content := container.NewVBox(pathLabel, container.NewHBox(copyBtn, openBtn))
	dialog.ShowCustom(title, "Close", content, win)
}
