package api

import (
	"net/http"

	"github.com/chapar-rest/mock-server/internal/gen/restapi"
	"github.com/chapar-rest/mock-server/internal/pkg/convert"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
	"github.com/chapar-rest/mock-server/internal/pkg/utils/restutils"
)

// GetTodo implements restapi.ServerInterface.
func (c *Controller) GetTodo(w http.ResponseWriter, r *http.Request, rawTodoId restapi.TodoId, params restapi.GetTodoParams) {
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

	todo, err := c.service.GetTodo(r.Context(), sessionId, todoId)
	if err != nil {
		c.writeError(w, r, err)
		return
	}

	restutils.Success(w, http.StatusOK, &restapi.TodoResponse{Data: convert.TodoToRestApi(todo)})
}
