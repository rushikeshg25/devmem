package store

// RepoRef is a minimal repo identity used for prune checks.
type RepoRef struct {
	ID   int64
	Path string
}

// ListRepoRefs returns the id and path of every indexed repo.
func (s *Store) ListRepoRefs() ([]RepoRef, error) {
	rows, err := s.db.Query(`SELECT id, path FROM repos`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []RepoRef
	for rows.Next() {
		var r RepoRef
		if err := rows.Scan(&r.ID, &r.Path); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// DeleteRepo removes a repo and its commits.
func (s *Store) DeleteRepo(id int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM commits WHERE repo_id = ?`, id); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(`DELETE FROM repos WHERE id = ?`, id); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// DeleteEmptyWorkspaces removes workspaces that no longer contain any repos and
// returns the number deleted.
func (s *Store) DeleteEmptyWorkspaces() (int, error) {
	res, err := s.db.Exec(`
		DELETE FROM workspaces
		WHERE id NOT IN (SELECT DISTINCT workspace_id FROM repos)`)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return int(n), nil
}
