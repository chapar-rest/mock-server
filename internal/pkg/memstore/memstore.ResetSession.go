package memstore

import (
	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

// ResetSession replaces the session's todos with the seed data.
func (s *Store) ResetSession(sessionId model.SessionId) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sess, err := s.session(sessionId)
	if err != nil {
		return err
	}

	sess.todos = seedTodos(s.cfg.Seed, s.cfg.Now())
	return nil
}
