package service

import (
	"context"
	"time"

	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

type CreateTodoRequest struct {
	SessionId   model.SessionId
	Title       string
	Description string
	Completed   bool
	// Empty defaults to medium.
	Priority model.TodoPriority
	Tags     []string
	DueAt    *time.Time
}

func (s *Service) CreateTodo(_ context.Context, req *CreateTodoRequest) (*model.Todo, error) {
	now := s.now()
	todo := &model.Todo{
		Id:          model.NewTodoId(),
		Title:       req.Title,
		Description: req.Description,
		Completed:   req.Completed,
		Priority:    defaultPriority(req.Priority),
		Tags:        nonNilTags(req.Tags),
		DueAt:       req.DueAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := validateTodo(todo); err != nil {
		return nil, err
	}

	return s.store.CreateTodo(req.SessionId, todo)
}

func defaultPriority(p model.TodoPriority) model.TodoPriority {
	if p == "" {
		return model.TodoPriorityMedium
	}
	return p
}

// nonNilTags keeps "no tags" as an empty list, so it serializes as [] rather
// than null.
func nonNilTags(tags []string) []string {
	if tags == nil {
		return []string{}
	}
	return tags
}
