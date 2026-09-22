package memstore

import (
	"time"
)

// DeleteExpiredSessions removes every session not seen since before, and
// returns how many were removed and how many remain.
func (s *Store) DeleteExpiredSessions(before time.Time) (removed, remaining int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, sess := range s.sessions {
		if sess.lastSeen.Before(before) {
			delete(s.sessions, id)
			removed++
		}
	}
	return removed, len(s.sessions)
}
