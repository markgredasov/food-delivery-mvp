// Package store holds the restaurant-simulator's in-memory, thread-safe
// bookkeeping of which orders it has already started processing.
//
// A plain mutex-protected map (rather than a real database) is an explicit
// MVP simplification: the simulator is a stub/simulated third-party
// integration, not a system of record, so its own state does not need to
// survive a restart.
package store

import "sync"

// Store tracks which order ids have already been claimed for processing, so
// the same order delivered twice (once via webhook push, once via polling
// fallback) is only ever accepted/progressed once.
type Store struct {
	mu      sync.Mutex
	claimed map[string]struct{}
}

// New builds an empty Store.
func New() *Store {
	return &Store{claimed: make(map[string]struct{})}
}

// TryClaim atomically marks orderID as claimed and reports whether this
// call was the first to claim it (true) or it was already claimed (false).
func (s *Store) TryClaim(orderID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.claimed[orderID]; ok {
		return false
	}
	s.claimed[orderID] = struct{}{}
	return true
}
