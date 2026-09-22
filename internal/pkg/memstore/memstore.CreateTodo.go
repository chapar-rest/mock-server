package memstore

import (
	"github.com/chapar-rest/mock-server/internal/pkg/errx"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

// CreateTodo stores a new todo in the session.
func (s *Store) CreateTodo(sessionId model.SessionId, todo *model.Todo) (*model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sess, err := s.session(sessionId)
	if err != nil {
		return nil, err
	}

	if len(sess.todos) >= s.cfg.MaxTodosPerSession {
		return nil, errx.ErrTodoLimitReached
	}

	sess.todos[todo.Id] = todo.Clone()
	return todo.Clone(), nil
}
