// Package git collects metadata from git repositories by shelling out to the git binary.
package git

import (
	"io/fs"
	"os"
	"path/filepath"
)

// Location is a discovered git checkout — a normal repo or a linked worktree.
type Location struct {
	Path       string // absolute path to the working tree
	IsWorktree bool   // true when .git is a file (linked worktree / submodule)
}

// Discover walks root and returns every git checkout found, without descending
// into a checkout once located (so nested files are not re-walked).
func Discover(root string) ([]Location, error) {
	var locs []Location
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // skip unreadable entries rather than aborting the whole scan
		}
		if !d.IsDir() {
			return nil
		}
		gitPath := filepath.Join(path, ".git")
		info, statErr := os.Stat(gitPath)
		if statErr != nil {
			return nil // no .git here, keep walking
		}
		locs = append(locs, Location{Path: path, IsWorktree: !info.IsDir()})
		return filepath.SkipDir // don't descend into a checkout
	})
	if err != nil {
		return nil, err
	}
	return locs, nil
}
