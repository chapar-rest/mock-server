package api

import (
	"fmt"
	"net/http"

	"github.com/chapar-rest/mock-server/internal/gen/restapi"
	"github.com/chapar-rest/mock-server/internal/pkg/errx"
	"github.com/chapar-rest/mock-server/internal/pkg/utils/restutils"
)

// GetStatus implements restapi.ServerInterface.
func (c *Controller) GetStatus(w http.ResponseWriter, r *http.Request, code int) {
	if code < 200 || code > 599 {
		c.writeError(w, r, fmt.Errorf("%w: code must be between 200 and 599", errx.ErrInvalidArgument))
		return
	}

	// These statuses must not carry a body.
	if code == http.StatusNoContent || code == http.StatusNotModified {
		w.WriteHeader(code)
		return
	}

	restutils.JSON(w, code, &restapi.StatusResponse{
		Code:        code,
		Description: http.StatusText(code),
	})
}
