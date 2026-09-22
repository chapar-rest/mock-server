package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/chapar-rest/mock-server/internal/gen/restapi"
)

const streamInterval = 100 * time.Millisecond

// Stream implements restapi.ServerInterface.
func (c *Controller) Stream(w http.ResponseWriter, r *http.Request, n int) {
	rc := http.NewResponseController(w)
	enc := json.NewEncoder(w)

	started := false
	err := c.service.Stream(r.Context(), n, streamInterval, func(index int, at time.Time) error {
		if !started {
			w.Header().Set("Content-Type", "application/x-ndjson")
			w.WriteHeader(http.StatusOK)
			started = true
		}
		if err := enc.Encode(&restapi.StreamEvent{Index: index, Time: at}); err != nil {
			return err
		}
		return rc.Flush()
	})

	// Once the first line is out the status is committed; a failure then can
	// only end the stream early.
	if err != nil && !started {
		c.writeError(w, r, err)
	}
}
