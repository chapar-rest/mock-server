package memstore

import (
	"github.com/chapar-rest/mock-server/internal/pkg/errx"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

// GetTodo returns one todo from the session.
func (s *Store) GetTodo(sessionId model.SessionId, id model.TodoId) (*model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sess, err := s.session(sessionId)
	if err != nil {
		return nil, err
	}

	t, ok := sess.todos[id]
	if !ok {
		return nil, errx.ErrTodoNotFound
	}
	return t.Clone(), nil
}
