package memstore

import (
	"github.com/chapar-rest/mock-server/internal/pkg/errx"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

// UpdateTodo applies update to a copy of the todo and stores the result. The
// read and the write happen under one lock, so concurrent updates never lose
// each other's changes. When update returns an error nothing is stored.
func (s *Store) UpdateTodo(sessionId model.SessionId, id model.TodoId, update func(*model.Todo) error) (*model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sess, err := s.session(sessionId)
	if err != nil {
		return nil, err
	}

	current, ok := sess.todos[id]
	if !ok {
		return nil, errx.ErrTodoNotFound
	}

	next := current.Clone()
	if err := update(next); err != nil {
		return nil, err
	}
	// The id is the map key; an update must not move the todo.
	next.Id = id

	sess.todos[id] = next
	return next.Clone(), nil
}
