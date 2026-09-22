package api

import (
	"net/http"

	"github.com/chapar-rest/mock-server/internal/gen/restapi"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

// DeleteTodo implements restapi.ServerInterface.
func (c *Controller) DeleteTodo(w http.ResponseWriter, r *http.Request, rawTodoId restapi.TodoId, params restapi.DeleteTodoParams) {
	sessionId, err := parseSessionId(params.XSessionId)
	if err != nil {
		c.writeError(w, r, err)
		return
	}

	todoId, err := model.ParseTodoId(rawTodoId)
	if err != nil {
		c.writeError(w, r, err)
		return
	}

	if err := c.service.DeleteTodo(r.Context(), sessionId, todoId); err != nil {
		c.writeError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
