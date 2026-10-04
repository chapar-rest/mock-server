package ws

import (
	"cmp"
	"net/http"
	"time"

	"github.com/coder/websocket"

	"github.com/chapar-rest/mock-server/internal/gen/restapi"
	"github.com/chapar-rest/mock-server/internal/pkg/convert"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

// todoMessage is one message of /todos.
type todoMessage struct {
	Type    string        `json:"type"`
	Session string        `json:"session"`
	Id      string        `json:"id,omitempty"`
	Todo    *restapi.Todo `json:"todo,omitempty"`
	Time    time.Time     `json:"time"`
}

const todosSubscribed = "subscribed"

// Todos sends the live changes to a session's todos, as REST and gRPC calls
// make them. The session is the X-Session-Id header or the session query
// parameter (browsers cannot set headers on a WebSocket), defaulting to
// public. Messages the client sends are ignored.
func (c *Controller) Todos(w http.ResponseWriter, r *http.Request) {
	sessionId, err := parseSessionId(r)
	if err != nil {
		c.writeError(w, r, err)
		return
	}

	s := c.accept(w, r)
	if s == nil {
		return
	}
	defer s.done()

	events, cancel := c.service.SubscribeTodos(sessionId)
	defer cancel()

	ctx := s.discardReads()
	//nolint:exhaustruct // the greeting carries no todo.
	if err := s.writeJSON(&todoMessage{Type: todosSubscribed, Session: sessionId.String(), Time: time.Now().UTC()}); err != nil {
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-events:
			if !ok {
				_ = s.conn.Close(websocket.StatusPolicyViolation, "too slow to read todo events")
				return
			}
			msg := &todoMessage{
				Type:    string(ev.Type),
				Session: sessionId.String(),
				Id:      string(ev.TodoId),
				Todo:    nil,
				Time:    ev.At,
			}
			if ev.Todo != nil {
				todo := convert.TodoToRestApi(ev.Todo)
				msg.Todo = &todo
			}
			if err := s.writeJSON(msg); err != nil {
				return
			}
		}
	}
}

func parseSessionId(r *http.Request) (model.SessionId, error) {
	raw := cmp.Or(r.Header.Get("X-Session-Id"), r.URL.Query().Get("session"))
	if raw == "" {
		return model.PublicSessionId, nil
	}
	return model.ParseSessionId(raw)
}
