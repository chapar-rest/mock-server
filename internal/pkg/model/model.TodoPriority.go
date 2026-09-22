package model

import (
	"fmt"

	"github.com/chapar-rest/mock-server/internal/pkg/errx"
)

// TodoPriority is the priority of a todo.
type TodoPriority string

const (
	TodoPriorityLow    TodoPriority = "low"
	TodoPriorityMedium TodoPriority = "medium"
	TodoPriorityHigh   TodoPriority = "high"
)

// ParseTodoPriority parses a string and returns a todo priority.
func ParseTodoPriority(priority string) (TodoPriority, error) {
	switch priority {
	case string(TodoPriorityLow):
		return TodoPriorityLow, nil
	case string(TodoPriorityMedium):
		return TodoPriorityMedium, nil
	case string(TodoPriorityHigh):
		return TodoPriorityHigh, nil
	}
	return "", fmt.Errorf("%w: priority must be one of low, medium, high", errx.ErrInvalidArgument)
}
