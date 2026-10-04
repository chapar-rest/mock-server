package model

import "time"

// TodoEventType says what happened to a session's todos.
type TodoEventType string

const (
	TodoEventCreated TodoEventType = "todo.created"
	TodoEventUpdated TodoEventType = "todo.updated"
	TodoEventDeleted TodoEventType = "todo.deleted"
	TodoEventReset   TodoEventType = "todos.reset"
)

// TodoEvent is one change to a session's todos, as live subscribers see it.
type TodoEvent struct {
	Type TodoEventType
	// Todo is the todo after the change; nil for deleted and reset.
	Todo *Todo
	// TodoId is the changed todo; empty for reset.
	TodoId TodoId
	At     time.Time
}
