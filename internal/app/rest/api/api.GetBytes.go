package api

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strconv"

	"github.com/chapar-rest/mock-server/internal/pkg/errx"
)

const maxBytes = 100 << 10

// GetBytes implements restapi.ServerInterface.
func (c *Controller) GetBytes(w http.ResponseWriter, r *http.Request, n int) {
	if n < 0 || n > maxBytes {
		c.writeError(w, r, fmt.Errorf("%w: n must be between 0 and %d", errx.ErrInvalidArgument, maxBytes))
		return
	}

	body := make([]byte, n)
	_, _ = rand.Read(body)

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(n))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}
