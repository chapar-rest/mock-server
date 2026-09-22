package api

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/chapar-rest/mock-server/internal/gen/mockv1"
)

// ResetTodos implements mockv1.TodoServiceServer.
func (s *TodoServer) ResetTodos(ctx context.Context, _ *emptypb.Empty) (*mockv1.ListTodosResponse, error) {
	sessionId, err := sessionIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	result, err := s.service.ResetTodos(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	return listTodosResponse(result), nil
}
