package api

import (
	"context"

	"github.com/chapar-rest/mock-server/internal/gen/mockv1"
	"github.com/chapar-rest/mock-server/internal/pkg/convert"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
	"github.com/chapar-rest/mock-server/internal/pkg/service"
)

// ListTodos implements mockv1.TodoServiceServer.
func (s *TodoServer) ListTodos(ctx context.Context, req *mockv1.ListTodosRequest) (*mockv1.ListTodosResponse, error) {
	sessionId, err := sessionIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var priority *model.TodoPriority
	if req.GetPriority() != mockv1.TodoPriority_TODO_PRIORITY_UNSPECIFIED {
		p, err := model.ParseTodoPriority(string(convert.TodoPriorityFromGrpcApi(req.GetPriority())))
		if err != nil {
			return nil, err
		}
		priority = &p
	}

	result, err := s.service.ListTodos(ctx, &service.ListTodosRequest{
		SessionId: sessionId,
		Completed: req.Completed,
		Priority:  priority,
		Query:     req.GetQuery(),
		Limit:     int(req.GetLimit()),
		Offset:    int(req.GetOffset()),
	})
	if err != nil {
		return nil, err
	}

	return listTodosResponse(result), nil
}

func listTodosResponse(result *service.ListTodosResult) *mockv1.ListTodosResponse {
	return &mockv1.ListTodosResponse{
		Todos:  convert.Slice(result.Todos, convert.TodoToGrpcApi),
		Total:  int32(result.Total),
		Limit:  int32(result.Limit),
		Offset: int32(result.Offset),
	}
}
