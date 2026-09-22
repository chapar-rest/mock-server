package api

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/chapar-rest/mock-server/internal/gen/mockv1"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

// DeleteTodo implements mockv1.TodoServiceServer.
func (s *TodoServer) DeleteTodo(ctx context.Context, req *mockv1.DeleteTodoRequest) (*emptypb.Empty, error) {
	sessionId, err := sessionIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	todoId, err := model.ParseTodoId(req.GetId())
	if err != nil {
		return nil, err
	}

	if err := s.service.DeleteTodo(ctx, sessionId, todoId); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
