package memstore

import (
	"cmp"
	"slices"
	"strings"

	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

// TodoFilter narrows ListTodos. Nil and empty fields match everything.
type TodoFilter struct {
	Completed *bool
	Priority  *model.TodoPriority
	// Query is matched case-insensitively against title and description.
	Query  string
	Limit  int
	Offset int
}

// ListTodos returns one page of the session's todos matching filter, oldest
// first, along with the total number of matches.
func (s *Store) ListTodos(sessionId model.SessionId, filter TodoFilter) ([]*model.Todo, int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sess, err := s.session(sessionId)
	if err != nil {
		return nil, 0, err
	}

	query := strings.ToLower(filter.Query)
	matches := make([]*model.Todo, 0, len(sess.todos))
	for _, t := range sess.todos {
		if filter.Completed != nil && t.Completed != *filter.Completed {
			continue
		}
		if filter.Priority != nil && t.Priority != *filter.Priority {
			continue
		}
		if query != "" &&
			!strings.Contains(strings.ToLower(t.Title), query) &&
			!strings.Contains(strings.ToLower(t.Description), query) {
			continue
		}
		matches = append(matches, t)
	}

	slices.SortFunc(matches, func(a, b *model.Todo) int {
		return cmp.Or(a.CreatedAt.Compare(b.CreatedAt), cmp.Compare(a.Id, b.Id))
	})

	total := len(matches)
	start := min(filter.Offset, total)
	end := min(start+filter.Limit, total)

	page := make([]*model.Todo, 0, end-start)
	for _, t := range matches[start:end] {
		page = append(page, t.Clone())
	}
	return page, total, nil
}
