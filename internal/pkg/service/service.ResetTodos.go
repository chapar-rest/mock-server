package service

import (
	"context"

	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

// ResetTodos restores the session's sample todos and returns the first page.
func (s *Service) ResetTodos(ctx context.Context, sessionId model.SessionId) (*ListTodosResult, error) {
	if err := s.store.ResetSession(sessionId); err != nil {
		return nil, err
	}

	//nolint:exhaustruct // an unfiltered first page is the zero-valued filter.
	return s.ListTodos(ctx, &ListTodosRequest{SessionId: sessionId})
}
