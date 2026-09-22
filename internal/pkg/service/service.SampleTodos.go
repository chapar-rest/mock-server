package service

import (
	"time"

	"github.com/chapar-rest/mock-server/internal/pkg/model"
	"github.com/chapar-rest/mock-server/internal/pkg/utils"
)

// SampleTodos returns the todos every new or reset session starts with, so a
// first GET shows realistic data. Creation times are staggered to give the
// list a stable order.
func SampleTodos(now time.Time) []*model.Todo {
	at := func(offset time.Duration) time.Time { return now.Add(offset) }

	return []*model.Todo{
		{
			Id:          model.NewTodoId(),
			Title:       "Try Chapar",
			Description: "Send a request to the mock server from Chapar.",
			Completed:   true,
			Priority:    model.TodoPriorityHigh,
			Tags:        []string{"chapar", "getting-started"},
			DueAt:       nil,
			CreatedAt:   at(-3 * time.Minute),
			UpdatedAt:   at(-3 * time.Minute),
		},
		{
			Id:          model.NewTodoId(),
			Title:       "Call the gRPC API",
			Description: "Use server reflection to discover mock.v1.TodoService.",
			Completed:   false,
			Priority:    model.TodoPriorityMedium,
			Tags:        []string{"grpc"},
			DueAt:       utils.Ptr(at(24 * time.Hour)),
			CreatedAt:   at(-2 * time.Minute),
			UpdatedAt:   at(-2 * time.Minute),
		},
		{
			Id:          model.NewTodoId(),
			Title:       "Write a test script",
			Description: "",
			Completed:   false,
			Priority:    model.TodoPriorityLow,
			Tags:        []string{},
			DueAt:       nil,
			CreatedAt:   at(-1 * time.Minute),
			UpdatedAt:   at(-1 * time.Minute),
		},
	}
}
