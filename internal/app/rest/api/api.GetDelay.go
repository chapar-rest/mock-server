package api

import (
	"net/http"
	"time"
)

// GetDelay implements restapi.ServerInterface.
func (c *Controller) GetDelay(w http.ResponseWriter, r *http.Request, ms int) {
	if err := c.service.Delay(r.Context(), time.Duration(ms)*time.Millisecond); err != nil {
		c.writeError(w, r, err)
		return
	}
	c.echo(w, r)
}
