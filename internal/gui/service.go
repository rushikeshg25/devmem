package gui

import (
	"github.com/rushikeshg25/devmem/internal/scan"
	"github.com/rushikeshg25/devmem/internal/store"
)

// defaultLimit caps how many rows the views request from the store.
const defaultLimit = 200

// Service is a thin, Fyne-free wrapper over the store and scan packages. It
// holds the open database so the views can query it, and keeps all data access
// in one testable place.
type Service struct {
	store *store.Store
}

// NewService opens the database at dbPath and returns a Service backed by it.
func NewService(dbPath string) (*Service, error) {
	s, err := store.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return &Service{store: s}, nil
}

// Close releases the underlying database connection.
func (svc *Service) Close() error {
	return svc.store.Close()
}

// Search returns commits matching term, newest first.
func (svc *Service) Search(term string) ([]store.SearchHit, error) {
	return svc.store.Search(term, defaultLimit)
}

// WIP returns checkouts with uncommitted, unpushed or stashed work matching
// term — the work most at risk of being lost.
func (svc *Service) WIP(term string) ([]store.RepoStatus, error) {
	return svc.store.SearchWIP(term)
}

// Timeline returns the most recent commits across every indexed repo.
func (svc *Service) Timeline() ([]store.SearchHit, error) {
	return svc.store.Timeline(defaultLimit)
}

// Scan discovers and indexes git checkouts under root, returning a summary.
func (svc *Service) Scan(root string) (scan.Result, error) {
	return scan.Run(svc.store, root)
}
