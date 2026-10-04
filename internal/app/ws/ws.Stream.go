package ws

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/coder/websocket"

	"github.com/chapar-rest/mock-server/internal/pkg/errx"
)

const (
	defaultStreamCount    = 10
	defaultStreamInterval = time.Second
)

// streamEvent is one message of /stream.
type streamEvent struct {
	Index int       `json:"index"`
	Count int       `json:"count"`
	Time  time.Time `json:"time"`
}

// Stream sends count events (default 10), interval_ms apart (default 1000),
// then closes normally. Messages the client sends are ignored.
func (c *Controller) Stream(w http.ResponseWriter, r *http.Request) {
	count, err := intParam(r, "count", defaultStreamCount)
	if err != nil {
		c.writeError(w, r, err)
		return
	}
	intervalMs, err := intParam(r, "interval_ms", int(defaultStreamInterval.Milliseconds()))
	if err != nil {
		c.writeError(w, r, err)
		return
	}
	interval := time.Duration(intervalMs) * time.Millisecond
	if err := c.service.CheckStream(count, interval); err != nil {
		c.writeError(w, r, err)
		return
	}

	s := c.accept(w, r)
	if s == nil {
		return
	}
	defer s.done()

	// Nothing is read, but control frames (pongs, the client's close) still
	// need a reader; CloseRead runs one and ends ctx when the client leaves.
	ctx := s.conn.CloseRead(s.ctx)
	err = c.service.Stream(ctx, count, interval, func(index int, at time.Time) error {
		return s.writeJSON(&streamEvent{Index: index, Count: count, Time: at})
	})
	if err == nil {
		_ = s.conn.Close(websocket.StatusNormalClosure, "stream complete")
	}
}

// intParam reads an optional integer query parameter.
func intParam(r *http.Request, name string, fallback int) (int, error) {
	raw := r.URL.Query().Get(name)
	if raw == "" {
		return fallback, nil
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%w: %s must be an integer", errx.ErrInvalidArgument, name)
	}
	return v, nil
}
