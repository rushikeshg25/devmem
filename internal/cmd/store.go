package cmd

import "github.com/rushikeshg25/devmem/internal/store"

// openStore opens the database at the resolved --db path.
func openStore() (*store.Store, error) {
	return store.Open(dbPath)
}
