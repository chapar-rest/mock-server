package api

import (
	"context"

	"github.com/chapar-rest/mock-server/internal/gen/mockv1"
	"github.com/chapar-rest/mock-server/internal/pkg/convert"
	"github.com/chapar-rest/mock-server/internal/pkg/service"
)

// CreateTodo implements mockv1.TodoServiceServer.
func (s *TodoServer) CreateTodo(ctx context.Context, req *mockv1.CreateTodoRequest) (*mockv1.Todo, error) {
	sessionId, err := sessionIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	todo, err := s.service.CreateTodo(ctx, &service.CreateTodoRequest{
		SessionId:   sessionId,
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Completed:   req.GetCompleted(),
		Priority:    convert.TodoPriorityFromGrpcApi(req.GetPriority()),
		Tags:        req.GetTags(),
		DueAt:       convert.TimeFromGrpcApi(req.GetDueAt()),
	})
	if err != nil {
		return nil, err
	}

	return convert.TodoToGrpcApi(todo), nil
}
