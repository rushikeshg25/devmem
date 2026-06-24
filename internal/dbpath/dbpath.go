// Package dbpath resolves the location of the devmem SQLite database so the CLI
// and the desktop GUI agree on a single default.
package dbpath

import (
	"os"
	"path/filepath"
)

// Default returns ~/.devmem.db, falling back to devmem.db in the current
// directory when the home directory cannot be determined.
func Default() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "devmem.db"
	}
	return filepath.Join(home, ".devmem.db")
}
