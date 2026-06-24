package gui

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/rushikeshg25/devmem/internal/store"
)

// seedDB writes a small fixture into a temp database and returns its path.
func seedDB(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "t.db")
	s, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	wsID, err := s.UpsertWorkspace("/ws/20-june", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	cleanID, err := s.UpsertRepo(store.Repo{WorkspaceID: wsID, Name: "erpai-report", Path: "/ws/20-june/erpai-report", CurrentBranch: "main"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.UpsertRepo(store.Repo{WorkspaceID: wsID, Name: "dirty-svc", Path: "/ws/20-june/dirty-svc", CurrentBranch: "feat/wip", IsDirty: true, AheadCount: 2}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.InsertCommits([]store.Commit{
		{RepoID: cleanID, Hash: "h1", Message: "fix: timezone handling", CommitTime: time.Unix(1700000000, 0)},
		{RepoID: cleanID, Hash: "h2", Message: "add kafka retries", CommitTime: time.Unix(1700000100, 0)},
	}); err != nil {
		t.Fatal(err)
	}
	return path
}

func newTestService(t *testing.T) *Service {
	t.Helper()
	svc, err := NewService(seedDB(t))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { svc.Close() })
	return svc
}

func TestServiceSearch(t *testing.T) {
	svc := newTestService(t)
	hits, err := svc.Search("timezone")
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 1 {
		t.Fatalf("expected 1 hit, got %d", len(hits))
	}
	if hits[0].Repo != "erpai-report" {
		t.Errorf("unexpected hit: %+v", hits[0])
	}
}

func TestServiceWIP(t *testing.T) {
	svc := newTestService(t)
	repos, err := svc.WIP("")
	if err != nil {
		t.Fatal(err)
	}
	if len(repos) != 1 {
		t.Fatalf("expected 1 at-risk repo, got %d", len(repos))
	}
	if repos[0].Name != "dirty-svc" || repos[0].AheadCount != 2 {
		t.Errorf("unexpected wip repo: %+v", repos[0])
	}
}

func TestServiceTimeline(t *testing.T) {
	svc := newTestService(t)
	hits, err := svc.Timeline()
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 2 {
		t.Fatalf("expected 2 commits, got %d", len(hits))
	}
	// Newest first.
	if hits[0].Hash != "h2" {
		t.Errorf("expected h2 first, got %s", hits[0].Hash)
	}
}
