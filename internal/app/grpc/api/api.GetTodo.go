package api

import (
	"context"

	"github.com/chapar-rest/mock-server/internal/gen/mockv1"
	"github.com/chapar-rest/mock-server/internal/pkg/convert"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

// GetTodo implements mockv1.TodoServiceServer.
func (s *TodoServer) GetTodo(ctx context.Context, req *mockv1.GetTodoRequest) (*mockv1.Todo, error) {
	sessionId, err := sessionIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	todoId, err := model.ParseTodoId(req.GetId())
	if err != nil {
		return nil, err
	}

	todo, err := s.service.GetTodo(ctx, sessionId, todoId)
	if err != nil {
		return nil, err
	}

	return convert.TodoToGrpcApi(todo), nil
}
