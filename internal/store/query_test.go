package store

import (
	"path/filepath"
	"testing"
	"time"
)

// newTestStore returns a Store backed by a temp-file database.
func newTestStore(t *testing.T) *Store {
	t.Helper()
	s, err := Open(filepath.Join(t.TempDir(), "t.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestSearchMatchesMessageAndIsScopedToRepo(t *testing.T) {
	s := newTestStore(t)
	wsID, err := s.UpsertWorkspace("/ws/20-june", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	repoID, err := s.UpsertRepo(Repo{WorkspaceID: wsID, Name: "erpai-report", Path: "/ws/20-june/erpai-report", CurrentBranch: "main"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.InsertCommits([]Commit{
		{RepoID: repoID, Hash: "h1", Message: "fix: timezone handling", CommitTime: time.Unix(1700000000, 0)},
		{RepoID: repoID, Hash: "h2", Message: "add kafka retries", CommitTime: time.Unix(1700000100, 0)},
	})
	if err != nil {
		t.Fatal(err)
	}

	hits, err := s.Search("timezone", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 1 {
		t.Fatalf("expected 1 hit, got %d", len(hits))
	}
	if hits[0].Repo != "erpai-report" || hits[0].Workspace != "/ws/20-june" {
		t.Errorf("unexpected context: %+v", hits[0])
	}
}

func TestInsertCommitsIgnoresDuplicates(t *testing.T) {
	s := newTestStore(t)
	wsID, _ := s.UpsertWorkspace("/ws", time.Now())
	repoID, _ := s.UpsertRepo(Repo{WorkspaceID: wsID, Name: "r", Path: "/ws/r"})

	c := []Commit{{RepoID: repoID, Hash: "dup", Message: "m", CommitTime: time.Unix(1, 0)}}
	if n, _ := s.InsertCommits(c); n != 1 {
		t.Fatalf("first insert: expected 1, got %d", n)
	}
	if n, _ := s.InsertCommits(c); n != 0 {
		t.Fatalf("re-insert: expected 0, got %d", n)
	}
}
