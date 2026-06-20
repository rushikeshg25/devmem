package git

import (
	"bytes"
	"os/exec"
	"strings"
)

// run executes git in dir with the given args and returns trimmed stdout.
func run(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimRight(out.String(), "\n"), nil
}

// splitLines splits s into non-empty lines.
func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}
