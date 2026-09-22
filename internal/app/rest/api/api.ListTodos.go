package api

import (
	"net/http"

	"github.com/chapar-rest/mock-server/internal/gen/restapi"
	"github.com/chapar-rest/mock-server/internal/pkg/convert"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
	"github.com/chapar-rest/mock-server/internal/pkg/service"
	"github.com/chapar-rest/mock-server/internal/pkg/utils"
	"github.com/chapar-rest/mock-server/internal/pkg/utils/restutils"
)

// ListTodos implements restapi.ServerInterface.
func (c *Controller) ListTodos(w http.ResponseWriter, r *http.Request, params restapi.ListTodosParams) {
	sessionId, err := parseSessionId(params.XSessionId)
	if err != nil {
		c.writeError(w, r, err)
		return
	}

	var priority *model.TodoPriority
	if params.Priority != nil {
		p, err := model.ParseTodoPriority(string(*params.Priority))
		if err != nil {
			c.writeError(w, r, err)
			return
		}
		priority = &p
	}

	result, err := c.service.ListTodos(r.Context(), &service.ListTodosRequest{
		SessionId: sessionId,
		Completed: params.Completed,
		Priority:  priority,
		Query:     utils.Value(params.Q),
		Limit:     utils.Value(params.Limit),
		Offset:    utils.Value(params.Offset),
	})
	if err != nil {
		c.writeError(w, r, err)
		return
	}

	restutils.Success(w, http.StatusOK, listTodosResponse(result))
}

func listTodosResponse(result *service.ListTodosResult) *restapi.ListTodosResponse {
	return &restapi.ListTodosResponse{
		Data:   convert.Slice(result.Todos, convert.TodoToRestApi),
		Total:  result.Total,
		Limit:  result.Limit,
		Offset: result.Offset,
	}
}
