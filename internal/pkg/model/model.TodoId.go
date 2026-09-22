package model

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/chapar-rest/mock-server/internal/pkg/errx"
)

// TodoId represents a todo id.
type TodoId string

// NewTodoId returns a new todo id.
func NewTodoId() TodoId {
	return TodoId(uuid.New().String())
}

// ParseTodoId parses a string and returns a todo id.
func ParseTodoId(s string) (TodoId, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return TodoId(""), fmt.Errorf("%w: todo id must be a UUID", errx.ErrInvalidArgument)
	}
	return TodoId(id.String()), nil
}

// String returns the string representation of the todo id.
func (t TodoId) String() string {
	return string(t)
}
