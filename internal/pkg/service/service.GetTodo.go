package service

import (
	"context"

	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

func (s *Service) GetTodo(_ context.Context, sessionId model.SessionId, id model.TodoId) (*model.Todo, error) {
	return s.store.GetTodo(sessionId, id)
}
