package memstore

import (
	"errors"
	"testing"
	"time"

	"github.com/chapar-rest/mock-server/internal/pkg/errx"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

type fakeClock struct{ now time.Time }

func (c *fakeClock) Now() time.Time { return c.now }

func newTestStore(t *testing.T, maxSessions, maxTodos int) (*Store, *fakeClock) {
	t.Helper()
	clock := &fakeClock{now: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
	store := New(Config{
		MaxSessions:        maxSessions,
		MaxTodosPerSession: maxTodos,
		Seed: func(now time.Time) []*model.Todo {
			return []*model.Todo{newTodo("seeded", now)}
		},
		Now: clock.Now,
	})
	return store, clock
}

func newTodo(title string, at time.Time) *model.Todo {
	return &model.Todo{
		Id:          model.NewTodoId(),
		Title:       title,
		Description: "",
		Completed:   false,
		Priority:    model.TodoPriorityMedium,
		Tags:        []string{},
		DueAt:       nil,
		CreatedAt:   at,
		UpdatedAt:   at,
	}
}

func allTodos(t *testing.T, s *Store, id model.SessionId) []*model.Todo {
	t.Helper()
	//nolint:exhaustruct // an unfiltered page.
	todos, _, err := s.ListTodos(id, TodoFilter{Limit: 100})
	if err != nil {
		t.Fatalf("ListTodos(%s): %v", id, err)
	}
	return todos
}

func TestNewSessionIsSeeded(t *testing.T) {
	s, _ := newTestStore(t, 10, 10)

	todos := allTodos(t, s, "a")
	if len(todos) != 1 || todos[0].Title != "seeded" {
		t.Fatalf("new session todos = %+v, want the one seeded todo", todos)
	}
}

func TestSessionsAreIsolated(t *testing.T) {
	s, clock := newTestStore(t, 10, 10)

	if _, err := s.CreateTodo("a", newTodo("only in a", clock.now)); err != nil {
		t.Fatalf("CreateTodo: %v", err)
	}

	if got := len(allTodos(t, s, "a")); got != 2 {
		t.Fatalf("session a has %d todos, want 2", got)
	}
	if got := len(allTodos(t, s, "b")); got != 1 {
		t.Fatalf("session b has %d todos, want 1", got)
	}
}

func TestTodoLimit(t *testing.T) {
	s, clock := newTestStore(t, 10, 2)

	// The seeded todo counts toward the limit.
	if _, err := s.CreateTodo("a", newTodo("second", clock.now)); err != nil {
		t.Fatalf("CreateTodo: %v", err)
	}
	_, err := s.CreateTodo("a", newTodo("third", clock.now))
	if !errors.Is(err, errx.ErrTodoLimitReached) {
		t.Fatalf("CreateTodo over the limit: err = %v, want ErrTodoLimitReached", err)
	}
}

func TestSessionLimitExemptsPublic(t *testing.T) {
	s, _ := newTestStore(t, 1, 10)

	allTodos(t, s, "a")
	//nolint:exhaustruct // an unfiltered page.
	if _, _, err := s.ListTodos("b", TodoFilter{Limit: 1}); !errors.Is(err, errx.ErrSessionLimitReached) {
		t.Fatalf("second private session: err = %v, want ErrSessionLimitReached", err)
	}
	allTodos(t, s, model.PublicSessionId)
}

func TestUpdateTodoErrorStoresNothing(t *testing.T) {
	s, _ := newTestStore(t, 10, 10)
	id := allTodos(t, s, "a")[0].Id

	_, err := s.UpdateTodo("a", id, func(todo *model.Todo) error {
		todo.Title = "changed"
		return errx.ErrInvalidArgument
	})
	if !errors.Is(err, errx.ErrInvalidArgument) {
		t.Fatalf("UpdateTodo: err = %v, want ErrInvalidArgument", err)
	}

	got, err := s.GetTodo("a", id)
	if err != nil {
		t.Fatalf("GetTodo: %v", err)
	}
	if got.Title != "seeded" {
		t.Fatalf("title = %q after a failed update, want unchanged", got.Title)
	}
}

func TestReturnedTodosAreCopies(t *testing.T) {
	s, _ := newTestStore(t, 10, 10)
	todo := allTodos(t, s, "a")[0]

	todo.Title = "mutated by caller"

	got, err := s.GetTodo("a", todo.Id)
	if err != nil {
		t.Fatalf("GetTodo: %v", err)
	}
	if got.Title != "seeded" {
		t.Fatalf("store title = %q, want it unaffected by caller mutation", got.Title)
	}
}

func TestDeleteExpiredSessions(t *testing.T) {
	s, clock := newTestStore(t, 10, 10)

	allTodos(t, s, "old")
	clock.now = clock.now.Add(2 * time.Hour)
	allTodos(t, s, "fresh")

	removed, remaining := s.DeleteExpiredSessions(clock.now.Add(-time.Hour))
	if removed != 1 || remaining != 1 {
		t.Fatalf("removed=%d remaining=%d, want 1 and 1", removed, remaining)
	}

	// An expired session comes back freshly seeded on its next use.
	if got := len(allTodos(t, s, "old")); got != 1 {
		t.Fatalf("recreated session has %d todos, want 1", got)
	}
}
