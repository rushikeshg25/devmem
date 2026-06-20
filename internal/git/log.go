package git

import (
	"strconv"
	"strings"
	"time"
)

// LogEntry is one parsed commit from git log.
type LogEntry struct {
	Hash       string
	Refs       string
	Author     string
	Message    string
	CommitTime time.Time
}

// unitSep matches the %x1f separator used in the pretty format below.
const unitSep = "\x1f"

// logFormat lays out fields in a single line separated by the unit separator,
// which never appears in commit text — so parsing stays robust.
const logFormat = "%H" + "%x1f" + "%D" + "%x1f" + "%an" + "%x1f" + "%ct" + "%x1f" + "%s"

// Log returns all commits reachable from any ref in the repo at dir.
func Log(dir string) ([]LogEntry, error) {
	out, err := run(dir, "log", "--all", "--pretty=format:"+logFormat)
	if err != nil {
		return nil, err
	}
	return parseLog(out), nil
}

// parseLog turns raw git-log output into entries. Malformed lines are skipped.
func parseLog(out string) []LogEntry {
	var entries []LogEntry
	for _, line := range splitLines(out) {
		fields := strings.Split(line, unitSep)
		if len(fields) != 5 {
			continue
		}
		ts, _ := strconv.ParseInt(fields[3], 10, 64)
		entries = append(entries, LogEntry{
			Hash:       fields[0],
			Refs:       strings.TrimSpace(fields[1]),
			Author:     fields[2],
			Message:    fields[4],
			CommitTime: time.Unix(ts, 0).UTC(),
		})
	}
	return entries
}
