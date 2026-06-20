package store

import "time"

// UpsertWorkspace inserts the workspace path (or updates its scan time) and returns its id.
func (s *Store) UpsertWorkspace(path string, scannedAt time.Time) (int64, error) {
	_, err := s.db.Exec(`
		INSERT INTO workspaces (path, last_scanned_at) VALUES (?, ?)
		ON CONFLICT(path) DO UPDATE SET last_scanned_at = excluded.last_scanned_at`,
		path, scannedAt)
	if err != nil {
		return 0, err
	}
	var id int64
	err = s.db.QueryRow(`SELECT id FROM workspaces WHERE path = ?`, path).Scan(&id)
	return id, err
}

// UpsertRepo inserts or refreshes a repo row (keyed by path) and returns its id.
func (s *Store) UpsertRepo(r Repo) (int64, error) {
	_, err := s.db.Exec(`
		INSERT INTO repos (workspace_id, name, path, current_branch, is_dirty, ahead_count, stash_count, is_worktree)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(path) DO UPDATE SET
			workspace_id   = excluded.workspace_id,
			name           = excluded.name,
			current_branch = excluded.current_branch,
			is_dirty       = excluded.is_dirty,
			ahead_count    = excluded.ahead_count,
			stash_count    = excluded.stash_count,
			is_worktree    = excluded.is_worktree`,
		r.WorkspaceID, r.Name, r.Path, r.CurrentBranch, r.IsDirty, r.AheadCount, r.StashCount, r.IsWorktree)
	if err != nil {
		return 0, err
	}
	var id int64
	err = s.db.QueryRow(`SELECT id FROM repos WHERE path = ?`, r.Path).Scan(&id)
	return id, err
}

// InsertCommits adds commits, ignoring ones already indexed for the repo (incremental rescans).
func (s *Store) InsertCommits(commits []Commit) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	stmt, err := tx.Prepare(`
		INSERT OR IGNORE INTO commits (repo_id, hash, refs, author, message, commit_time)
		VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	inserted := 0
	for _, c := range commits {
		res, err := stmt.Exec(c.RepoID, c.Hash, c.Refs, c.Author, c.Message, c.CommitTime)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		if n, _ := res.RowsAffected(); n > 0 {
			inserted++
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return inserted, nil
}
