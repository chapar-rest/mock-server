package memstore

import (
	"sync"
	"time"

	"github.com/chapar-rest/mock-server/internal/pkg/errx"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

// Config bounds the store's memory use and supplies each session's initial
// data.
type Config struct {
	// MaxSessions caps the number of live sessions. The public session is
	// exempt so that requests without a session id always work.
	MaxSessions int
	// MaxTodosPerSession caps the todos a single session can hold.
	MaxTodosPerSession int
	// Seed returns the todos a new or reset session starts with.
	Seed func(now time.Time) []*model.Todo
	// Now returns the current time. It drives session expiry.
	Now func() time.Time
}

// Store holds every session's todos in memory. It is safe for concurrent use.
// All values crossing its boundary are copies: callers never share mutable
// state with the store.
type Store struct {
	cfg Config

	mu       sync.Mutex
	sessions map[model.SessionId]*session
}

type session struct {
	todos    map[model.TodoId]*model.Todo
	lastSeen time.Time
}

// New creates an empty store.
func New(cfg Config) *Store {
	//nolint:exhaustruct // a zero sync.Mutex is the usable unlocked mutex.
	return &Store{
		cfg:      cfg,
		sessions: make(map[model.SessionId]*session),
	}
}

// session returns the session with the given id, creating and seeding it when
// it does not exist, and marks it as seen. The caller must hold s.mu.
func (s *Store) session(id model.SessionId) (*session, error) {
	now := s.cfg.Now()

	if sess, ok := s.sessions[id]; ok {
		sess.lastSeen = now
		return sess, nil
	}

	if id != model.PublicSessionId && s.countPrivateSessions() >= s.cfg.MaxSessions {
		return nil, errx.ErrSessionLimitReached
	}

	sess := &session{
		todos:    seedTodos(s.cfg.Seed, now),
		lastSeen: now,
	}
	s.sessions[id] = sess
	return sess, nil
}

// countPrivateSessions counts the sessions subject to MaxSessions. The caller
// must hold s.mu.
func (s *Store) countPrivateSessions() int {
	n := len(s.sessions)
	if _, ok := s.sessions[model.PublicSessionId]; ok {
		n--
	}
	return n
}

func seedTodos(seed func(now time.Time) []*model.Todo, now time.Time) map[model.TodoId]*model.Todo {
	todos := make(map[model.TodoId]*model.Todo)
	if seed == nil {
		return todos
	}
	for _, t := range seed(now) {
		todos[t.Id] = t
	}
	return todos
}
