package service

import (
	"context"

	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

func (s *Service) DeleteTodo(_ context.Context, sessionId model.SessionId, id model.TodoId) error {
	return s.store.DeleteTodo(sessionId, id)
}
