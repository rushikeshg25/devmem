package store

import "time"

// Repo is a git repository (or worktree) discovered inside a workspace.
type Repo struct {
	ID            int64
	WorkspaceID   int64
	Name          string
	Path          string
	CurrentBranch string
	IsDirty       bool
	AheadCount    int
	StashCount    int
	IsWorktree    bool
}

// Commit is a single commit indexed from a repo.
type Commit struct {
	RepoID     int64
	Hash       string
	Refs       string
	Author     string
	Message    string
	CommitTime time.Time
}
