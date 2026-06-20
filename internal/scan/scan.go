// Package scan orchestrates repository discovery, git metadata collection and persistence.
package scan

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/rushikeshg25/devmem/internal/git"
	"github.com/rushikeshg25/devmem/internal/store"
)

// Result summarizes a scan run.
type Result struct {
	Workspaces int
	Repos      int
	NewCommits int
}

// Run discovers git checkouts under root, collects their metadata and stores it.
func Run(s *store.Store, root string) (Result, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return Result{}, err
	}

	locs, err := git.Discover(root)
	if err != nil {
		return Result{}, err
	}

	now := time.Now().UTC()
	seenWorkspaces := map[string]int64{}
	var res Result

	for _, loc := range locs {
		wsPath := workspaceFor(root, loc.Path)
		wsID, ok := seenWorkspaces[wsPath]
		if !ok {
			wsID, err = s.UpsertWorkspace(wsPath, now)
			if err != nil {
				return res, err
			}
			seenWorkspaces[wsPath] = wsID
			res.Workspaces++
		}

		repoID, err := s.UpsertRepo(store.Repo{
			WorkspaceID:   wsID,
			Name:          filepath.Base(loc.Path),
			Path:          loc.Path,
			CurrentBranch: git.CurrentBranch(loc.Path),
			IsDirty:       git.IsDirty(loc.Path),
			AheadCount:    git.AheadCount(loc.Path),
			StashCount:    git.StashCount(loc.Path),
			IsWorktree:    loc.IsWorktree,
		})
		if err != nil {
			return res, err
		}
		res.Repos++

		entries, err := git.Log(loc.Path)
		if err != nil {
			continue // unreadable history shouldn't abort the whole scan
		}
		commits := make([]store.Commit, 0, len(entries))
		for _, e := range entries {
			commits = append(commits, store.Commit{
				RepoID:     repoID,
				Hash:       e.Hash,
				Refs:       e.Refs,
				Author:     e.Author,
				Message:    e.Message,
				CommitTime: e.CommitTime,
			})
		}
		n, err := s.InsertCommits(commits)
		if err != nil {
			return res, err
		}
		res.NewCommits += n
	}

	return res, nil
}

// workspaceFor returns the workspace folder for a repo: the immediate child of
// root that contains it, or root itself when the repo sits directly under root.
func workspaceFor(root, repoPath string) string {
	rel, err := filepath.Rel(root, repoPath)
	if err != nil {
		return root
	}
	segments := strings.Split(rel, string(filepath.Separator))
	if len(segments) <= 1 {
		return root
	}
	return filepath.Join(root, segments[0])
}
