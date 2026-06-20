package git

import "strconv"

// CurrentBranch returns the checked-out branch for the working tree at dir.
// For a detached HEAD git returns "HEAD".
func CurrentBranch(dir string) string {
	branch, err := run(dir, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return ""
	}
	return branch
}

// IsDirty reports whether the working tree has uncommitted changes.
func IsDirty(dir string) bool {
	out, err := run(dir, "status", "--porcelain")
	if err != nil {
		return false
	}
	return out != ""
}

// AheadCount returns the number of commits HEAD is ahead of its upstream.
// Returns 0 when there is no upstream configured.
func AheadCount(dir string) int {
	out, err := run(dir, "rev-list", "--count", "@{u}..HEAD")
	if err != nil {
		return 0
	}
	n, _ := strconv.Atoi(out)
	return n
}

// StashCount returns the number of stash entries.
func StashCount(dir string) int {
	out, err := run(dir, "stash", "list")
	if err != nil || out == "" {
		return 0
	}
	return len(splitLines(out))
}
