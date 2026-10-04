package service

import (
	"sync"

	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

// todoEventBuffer is how many events a subscriber may fall behind before it
// is dropped.
const todoEventBuffer = 64

// SubscribeTodos returns the live changes to a session's todos. The channel
// is closed when cancel is called, or when the subscriber falls too far
// behind; a closed channel means the subscription is over.
func (s *Service) SubscribeTodos(sessionId model.SessionId) (events <-chan model.TodoEvent, cancel func()) {
	return s.todoEvents.subscribe(sessionId)
}

// publishTodoEvent tells every subscriber of the session about a change.
func (s *Service) publishTodoEvent(sessionId model.SessionId, typ model.TodoEventType, id model.TodoId, todo *model.Todo) {
	s.todoEvents.publish(sessionId, model.TodoEvent{
		Type:   typ,
		Todo:   todo,
		TodoId: id,
		At:     s.now(),
	})
}

// todoHub fans todo events out to subscribers, per session.
type todoHub struct {
	mu   sync.Mutex
	subs map[model.SessionId]map[*todoSub]struct{}
}

type todoSub struct {
	ch   chan model.TodoEvent
	once sync.Once
}

func (s *todoSub) close() { s.once.Do(func() { close(s.ch) }) }

func newTodoHub() *todoHub {
	return &todoHub{
		mu:   sync.Mutex{},
		subs: map[model.SessionId]map[*todoSub]struct{}{},
	}
}

func (h *todoHub) subscribe(sessionId model.SessionId) (<-chan model.TodoEvent, func()) {
	sub := &todoSub{
		ch:   make(chan model.TodoEvent, todoEventBuffer),
		once: sync.Once{},
	}

	h.mu.Lock()
	if h.subs[sessionId] == nil {
		h.subs[sessionId] = map[*todoSub]struct{}{}
	}
	h.subs[sessionId][sub] = struct{}{}
	h.mu.Unlock()

	return sub.ch, func() {
		h.mu.Lock()
		h.remove(sessionId, sub)
		h.mu.Unlock()
		sub.close()
	}
}

// remove forgets a subscriber; h.mu must be held.
func (h *todoHub) remove(sessionId model.SessionId, sub *todoSub) {
	delete(h.subs[sessionId], sub)
	if len(h.subs[sessionId]) == 0 {
		delete(h.subs, sessionId)
	}
}

func (h *todoHub) publish(sessionId model.SessionId, event model.TodoEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for sub := range h.subs[sessionId] {
		ev := event
		if ev.Todo != nil {
			ev.Todo = ev.Todo.Clone()
		}
		select {
		case sub.ch <- ev:
		default:
			// A subscriber that cannot keep up is cut off rather than
			// allowed to block the writers.
			h.remove(sessionId, sub)
			sub.close()
		}
	}
}
