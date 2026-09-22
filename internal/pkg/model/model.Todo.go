package model

import (
	"slices"
	"time"
)

// Todo represents a todo item.
type Todo struct {
	Id          TodoId
	Title       string
	Description string
	Completed   bool
	Priority    TodoPriority
	Tags        []string
	DueAt       *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Clone returns a deep copy, so callers never share mutable state with the
// store.
func (t *Todo) Clone() *Todo {
	c := *t
	c.Tags = slices.Clone(t.Tags)
	if t.DueAt != nil {
		dueAt := *t.DueAt
		c.DueAt = &dueAt
	}
	return &c
}
