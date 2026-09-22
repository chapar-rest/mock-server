package service

import (
	"context"
	"time"

	"github.com/chapar-rest/mock-server/internal/pkg/model"
	"github.com/chapar-rest/mock-server/internal/pkg/utils"
)

// UpdateTodoRequest changes only the fields that are present.
type UpdateTodoRequest struct {
	SessionId   model.SessionId
	Id          model.TodoId
	Title       utils.Optional[string]
	Description utils.Optional[string]
	Completed   utils.Optional[bool]
	Priority    utils.Optional[model.TodoPriority]
	Tags        utils.Optional[[]string]
	// A present nil clears the due date.
	DueAt utils.Optional[*time.Time]
}

func (s *Service) UpdateTodo(_ context.Context, req *UpdateTodoRequest) (*model.Todo, error) {
	return s.store.UpdateTodo(req.SessionId, req.Id, func(t *model.Todo) error {
		if v, ok := req.Title.Get(); ok {
			t.Title = v
		}
		if v, ok := req.Description.Get(); ok {
			t.Description = v
		}
		if v, ok := req.Completed.Get(); ok {
			t.Completed = v
		}
		if v, ok := req.Priority.Get(); ok {
			t.Priority = defaultPriority(v)
		}
		if v, ok := req.Tags.Get(); ok {
			t.Tags = nonNilTags(v)
		}
		if v, ok := req.DueAt.Get(); ok {
			t.DueAt = v
		}
		t.UpdatedAt = s.now()
		return validateTodo(t)
	})
}
