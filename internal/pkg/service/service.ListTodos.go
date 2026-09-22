package service

import (
	"context"
	"fmt"

	"github.com/chapar-rest/mock-server/internal/pkg/errx"
	"github.com/chapar-rest/mock-server/internal/pkg/memstore"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

const (
	defaultListLimit = 20
	maxListLimit     = 100
)

type ListTodosRequest struct {
	SessionId model.SessionId
	Completed *bool
	Priority  *model.TodoPriority
	Query     string
	// Limit of 0 selects the default page size.
	Limit  int
	Offset int
}

type ListTodosResult struct {
	Todos  []*model.Todo
	Total  int
	Limit  int
	Offset int
}

func (s *Service) ListTodos(_ context.Context, req *ListTodosRequest) (*ListTodosResult, error) {
	limit := req.Limit
	if limit == 0 {
		limit = defaultListLimit
	}
	if limit < 1 || limit > maxListLimit {
		return nil, fmt.Errorf("%w: limit must be between 1 and %d", errx.ErrInvalidArgument, maxListLimit)
	}
	if req.Offset < 0 {
		return nil, fmt.Errorf("%w: offset must not be negative", errx.ErrInvalidArgument)
	}

	todos, total, err := s.store.ListTodos(req.SessionId, memstore.TodoFilter{
		Completed: req.Completed,
		Priority:  req.Priority,
		Query:     req.Query,
		Limit:     limit,
		Offset:    req.Offset,
	})
	if err != nil {
		return nil, err
	}

	return &ListTodosResult{
		Todos:  todos,
		Total:  total,
		Limit:  limit,
		Offset: req.Offset,
	}, nil
}
