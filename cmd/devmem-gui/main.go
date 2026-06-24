// Command devmem-gui is the desktop front-end for DevMem. It is a separate
// binary from the devmem CLI because Fyne requires CGo, while the CLI is built
// CGo-free for cross-platform releases.
package main

import (
	"github.com/rushikeshg25/devmem/internal/dbpath"
	"github.com/rushikeshg25/devmem/internal/gui"
)

// version is overridden at release time via -ldflags.
var version = "dev"

func main() {
	gui.SetVersion(version)
	gui.Run(dbpath.Default())
}
