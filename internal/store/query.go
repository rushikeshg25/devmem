package store

import "time"

// SearchHit is a commit match enriched with its repo/workspace context.
type SearchHit struct {
	Workspace  string
	Repo       string
	Branch     string
	Hash       string
	Message    string
	CommitTime time.Time
}

// Search returns commits whose message or refs match term, or that belong to a
// repo whose name or current branch matches term. Plain substring (LIKE) match.
func (s *Store) Search(term string, limit int) ([]SearchHit, error) {
	like := "%" + term + "%"
	rows, err := s.db.Query(`
		SELECT w.path, r.name, c.refs, c.hash, c.message, c.commit_time
		FROM commits c
		JOIN repos r      ON r.id = c.repo_id
		JOIN workspaces w ON w.id = r.workspace_id
		WHERE c.message LIKE ?
		   OR c.refs LIKE ?
		   OR r.name LIKE ?
		   OR r.current_branch LIKE ?
		ORDER BY c.commit_time DESC
		LIMIT ?`,
		like, like, like, like, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanHits(rows)
}

// RepoStatus is the current state of an indexed repo/worktree.
type RepoStatus struct {
	Workspace     string
	Name          string
	Path          string
	CurrentBranch string
	IsDirty       bool
	AheadCount    int
	StashCount    int
	IsWorktree    bool
}

// FindRepos returns every indexed checkout whose name matches name exactly.
func (s *Store) FindRepos(name string) ([]RepoStatus, error) {
	rows, err := s.db.Query(`
		SELECT w.path, r.name, r.path, r.current_branch, r.is_dirty, r.ahead_count, r.stash_count, r.is_worktree
		FROM repos r
		JOIN workspaces w ON w.id = r.workspace_id
		WHERE r.name = ?
		ORDER BY w.path`, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []RepoStatus
	for rows.Next() {
		var rs RepoStatus
		if err := rows.Scan(&rs.Workspace, &rs.Name, &rs.Path, &rs.CurrentBranch,
			&rs.IsDirty, &rs.AheadCount, &rs.StashCount, &rs.IsWorktree); err != nil {
			return nil, err
		}
		out = append(out, rs)
	}
	return out, rows.Err()
}

// RecentCommits returns the latest commits for a repo path.
func (s *Store) RecentCommits(repoPath string, limit int) ([]SearchHit, error) {
	rows, err := s.db.Query(`
		SELECT w.path, r.name, c.refs, c.hash, c.message, c.commit_time
		FROM commits c
		JOIN repos r      ON r.id = c.repo_id
		JOIN workspaces w ON w.id = r.workspace_id
		WHERE r.path = ?
		ORDER BY c.commit_time DESC
		LIMIT ?`, repoPath, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanHits(rows)
}

// Timeline returns the most recent commits across all repos.
func (s *Store) Timeline(limit int) ([]SearchHit, error) {
	rows, err := s.db.Query(`
		SELECT w.path, r.name, c.refs, c.hash, c.message, c.commit_time
		FROM commits c
		JOIN repos r      ON r.id = c.repo_id
		JOIN workspaces w ON w.id = r.workspace_id
		ORDER BY c.commit_time DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanHits(rows)
}

func scanHits(rows interface {
	Next() bool
	Scan(...any) error
	Err() error
}) ([]SearchHit, error) {
	var out []SearchHit
	for rows.Next() {
		var h SearchHit
		if err := rows.Scan(&h.Workspace, &h.Repo, &h.Branch, &h.Hash, &h.Message, &h.CommitTime); err != nil {
			return nil, err
		}
		out = append(out, h)
	}
	return out, rows.Err()
}
