// Package store handles SQLite persistence for devmem.
package store

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// Store wraps the SQLite database connection.
type Store struct {
	db *sql.DB
}

const schema = `
CREATE TABLE IF NOT EXISTS workspaces (
    id              INTEGER PRIMARY KEY,
    path            TEXT UNIQUE NOT NULL,
    last_scanned_at DATETIME
);

CREATE TABLE IF NOT EXISTS repos (
    id             INTEGER PRIMARY KEY,
    workspace_id   INTEGER NOT NULL REFERENCES workspaces(id),
    name           TEXT NOT NULL,
    path           TEXT UNIQUE NOT NULL,
    current_branch TEXT,
    is_dirty       BOOLEAN,
    ahead_count    INTEGER,
    stash_count    INTEGER,
    is_worktree    BOOLEAN
);

CREATE TABLE IF NOT EXISTS commits (
    id          INTEGER PRIMARY KEY,
    repo_id     INTEGER NOT NULL REFERENCES repos(id),
    hash        TEXT NOT NULL,
    refs        TEXT,
    author      TEXT,
    message     TEXT,
    commit_time DATETIME,
    UNIQUE(repo_id, hash)
);

CREATE INDEX IF NOT EXISTS idx_commits_repo ON commits(repo_id);
`

// Open opens (creating if needed) the SQLite database at path and applies the schema.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

// Close closes the underlying database.
func (s *Store) Close() error {
	return s.db.Close()
}
