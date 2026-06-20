package scan

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/rushikeshg25/devmem/internal/store"
)

// git runs a git command in dir and fails the test on error.
func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(cmd.Environ(),
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t",
		"GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

// buildFixture creates an origin repo cloned into two workspace folders, mirroring
// the real "same repo in many date folders" setup. Returns the workspace root.
func buildFixture(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
	base := t.TempDir()

	origin := filepath.Join(base, "origin")
	runGit(t, base, "init", "-q", origin)
	runGit(t, origin, "commit", "-q", "--allow-empty", "-m", "fix: timezone handling")
	runGit(t, origin, "commit", "-q", "--allow-empty", "-m", "add kafka retries")

	root := filepath.Join(base, "ws")
	runGit(t, base, "clone", "-q", origin, filepath.Join(root, "20-june", "repoA"))
	runGit(t, base, "clone", "-q", origin, filepath.Join(root, "20-june2", "repoA"))

	// A linked worktree under the first workspace.
	runGit(t, filepath.Join(root, "20-june", "repoA"),
		"worktree", "add", "-q", filepath.Join(root, "20-june", "repoA-wt"), "-b", "hotfix")
	return root
}

func TestScanIndexesDuplicateHashesPerWorkspace(t *testing.T) {
	root := buildFixture(t)
	s, err := store.Open(filepath.Join(t.TempDir(), "t.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	res, err := Run(s, root)
	if err != nil {
		t.Fatal(err)
	}
	if res.Repos != 3 { // two clones + one worktree
		t.Errorf("expected 3 repos, got %d", res.Repos)
	}

	// The two repoA clones share commit hashes; both must be indexed independently.
	repos, err := s.FindRepos("repoA")
	if err != nil {
		t.Fatal(err)
	}
	if len(repos) != 2 {
		t.Fatalf("expected repoA in 2 workspaces, got %d", len(repos))
	}

	// "timezone" should match in both clones (4 commits: 2 hashes x 2 repos).
	hits, err := s.Search("timezone", 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) == 0 {
		t.Fatal("expected timezone matches, got none")
	}
}

func TestScanIsIdempotent(t *testing.T) {
	root := buildFixture(t)
	s, err := store.Open(filepath.Join(t.TempDir(), "t.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	first, err := Run(s, root)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Run(s, root)
	if err != nil {
		t.Fatal(err)
	}
	if second.NewCommits != 0 {
		t.Errorf("rescan should add 0 commits, added %d", second.NewCommits)
	}
	if first.NewCommits == 0 {
		t.Error("first scan should have indexed commits")
	}
}

func TestScanFlagsWorktree(t *testing.T) {
	root := buildFixture(t)
	s, err := store.Open(filepath.Join(t.TempDir(), "t.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	if _, err := Run(s, root); err != nil {
		t.Fatal(err)
	}
	wt, err := s.FindRepos("repoA-wt")
	if err != nil {
		t.Fatal(err)
	}
	if len(wt) != 1 {
		t.Fatalf("expected 1 worktree checkout, got %d", len(wt))
	}
	if !wt[0].IsWorktree {
		t.Error("expected IsWorktree=true")
	}
	if wt[0].CurrentBranch != "hotfix" {
		t.Errorf("expected branch hotfix, got %q", wt[0].CurrentBranch)
	}
}
