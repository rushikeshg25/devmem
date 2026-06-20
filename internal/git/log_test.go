package git

import "testing"

func TestParseLog(t *testing.T) {
	// Two records using the unit separator; second message contains spaces and a colon.
	raw := "abc123\x1fHEAD -> main, origin/main\x1fAlice\x1f1700000000\x1ffix: timezone handling in exports\n" +
		"def456\x1f\x1fBob\x1f1700000100\x1fadd kafka retries"

	got := parseLog(raw)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}

	if got[0].Hash != "abc123" {
		t.Errorf("hash: got %q", got[0].Hash)
	}
	if got[0].Refs != "HEAD -> main, origin/main" {
		t.Errorf("refs: got %q", got[0].Refs)
	}
	if got[0].Author != "Alice" {
		t.Errorf("author: got %q", got[0].Author)
	}
	if got[0].Message != "fix: timezone handling in exports" {
		t.Errorf("message: got %q", got[0].Message)
	}
	if got[0].CommitTime.Unix() != 1700000000 {
		t.Errorf("time: got %d", got[0].CommitTime.Unix())
	}

	if got[1].Refs != "" {
		t.Errorf("expected empty refs for second commit, got %q", got[1].Refs)
	}
}

func TestParseLogSkipsMalformed(t *testing.T) {
	if got := parseLog("not-a-valid-line"); got != nil {
		t.Errorf("expected nil for malformed input, got %v", got)
	}
	if got := parseLog(""); got != nil {
		t.Errorf("expected nil for empty input, got %v", got)
	}
}
