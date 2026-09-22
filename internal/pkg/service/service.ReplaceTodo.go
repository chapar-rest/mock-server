package service

import (
	"context"
	"time"

	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

type ReplaceTodoRequest struct {
	SessionId   model.SessionId
	Id          model.TodoId
	Title       string
	Description string
	Completed   bool
	// Empty defaults to medium.
	Priority model.TodoPriority
	Tags     []string
	DueAt    *time.Time
}

// ReplaceTodo overwrites every user-editable field of a todo.
func (s *Service) ReplaceTodo(_ context.Context, req *ReplaceTodoRequest) (*model.Todo, error) {
	return s.store.UpdateTodo(req.SessionId, req.Id, func(t *model.Todo) error {
		t.Title = req.Title
		t.Description = req.Description
		t.Completed = req.Completed
		t.Priority = defaultPriority(req.Priority)
		t.Tags = nonNilTags(req.Tags)
		t.DueAt = req.DueAt
		t.UpdatedAt = s.now()
		return validateTodo(t)
	})
}
