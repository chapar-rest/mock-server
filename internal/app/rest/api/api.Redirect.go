package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/chapar-rest/mock-server/internal/pkg/errx"
)

const maxRedirects = 10

// Redirect implements restapi.ServerInterface.
func (c *Controller) Redirect(w http.ResponseWriter, r *http.Request, n int) {
	if n < 1 || n > maxRedirects {
		c.writeError(w, r, fmt.Errorf("%w: n must be between 1 and %d", errx.ErrInvalidArgument, maxRedirects))
		return
	}

	// Relative references resolve against /…/redirect/{n}, so the chain works
	// under any mount prefix.
	location := "../echo"
	if n > 1 {
		location = strconv.Itoa(n - 1)
	}

	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusFound)
}
