package api

import (
	"net/http"

	"github.com/chapar-rest/mock-server/internal/gen/restapi"
	"github.com/chapar-rest/mock-server/internal/pkg/utils/restutils"
)

// ResetTodos implements restapi.ServerInterface.
func (c *Controller) ResetTodos(w http.ResponseWriter, r *http.Request, params restapi.ResetTodosParams) {
	sessionId, err := parseSessionId(params.XSessionId)
	if err != nil {
		c.writeError(w, r, err)
		return
	}

	result, err := c.service.ResetTodos(r.Context(), sessionId)
	if err != nil {
		c.writeError(w, r, err)
		return
	}

	restutils.Success(w, http.StatusOK, listTodosResponse(result))
}
