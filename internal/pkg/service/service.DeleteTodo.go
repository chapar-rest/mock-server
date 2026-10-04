package service

import (
	"context"

	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

func (s *Service) DeleteTodo(_ context.Context, sessionId model.SessionId, id model.TodoId) error {
	if err := s.store.DeleteTodo(sessionId, id); err != nil {
		return err
	}
	s.publishTodoEvent(sessionId, model.TodoEventDeleted, id, nil)
	return nil
}
