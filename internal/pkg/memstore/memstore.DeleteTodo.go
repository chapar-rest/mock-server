package memstore

import (
	"github.com/chapar-rest/mock-server/internal/pkg/errx"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

// DeleteTodo removes one todo from the session.
func (s *Store) DeleteTodo(sessionId model.SessionId, id model.TodoId) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sess, err := s.session(sessionId)
	if err != nil {
		return err
	}

	if _, ok := sess.todos[id]; !ok {
		return errx.ErrTodoNotFound
	}
	delete(sess.todos, id)
	return nil
}
