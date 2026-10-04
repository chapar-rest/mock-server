package ws

import (
	"fmt"
	"net/http"

	"github.com/coder/websocket"

	"github.com/chapar-rest/mock-server/internal/pkg/errx"
)

// maxCloseReason is the most a close frame's reason can hold.
const maxCloseReason = 123

// Close closes the connection at once with code (default 1000) and reason,
// so clients can check how they report a close.
func (c *Controller) Close(w http.ResponseWriter, r *http.Request) {
	code, err := intParam(r, "code", int(websocket.StatusNormalClosure))
	if err != nil {
		c.writeError(w, r, err)
		return
	}
	reason := r.URL.Query().Get("reason")
	if err := checkClose(code, reason); err != nil {
		c.writeError(w, r, err)
		return
	}

	s := c.accept(w, r)
	if s == nil {
		return
	}
	defer s.done()

	_ = s.conn.Close(websocket.StatusCode(code), reason)
}

// checkClose accepts the codes an endpoint may send: the defined ones that
// can go on the wire, and the 3000-4999 range for applications.
func checkClose(code int, reason string) error {
	switch {
	case code == 1000, code == 1001, code == 1002, code == 1003,
		code >= 1007 && code <= 1011,
		code >= 3000 && code <= 4999:
	default:
		return fmt.Errorf("%w: code must be 1000-1003, 1007-1011 or 3000-4999", errx.ErrInvalidArgument)
	}
	if len(reason) > maxCloseReason {
		return fmt.Errorf("%w: reason must be at most %d bytes", errx.ErrInvalidArgument, maxCloseReason)
	}
	return nil
}
