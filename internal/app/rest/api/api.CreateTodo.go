package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chapar-rest/mock-server/internal/gen/restapi"
	"github.com/chapar-rest/mock-server/internal/pkg/convert"
	"github.com/chapar-rest/mock-server/internal/pkg/errx"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
	"github.com/chapar-rest/mock-server/internal/pkg/service"
	"github.com/chapar-rest/mock-server/internal/pkg/utils"
	"github.com/chapar-rest/mock-server/internal/pkg/utils/restutils"
)

// CreateTodo implements restapi.ServerInterface.
func (c *Controller) CreateTodo(w http.ResponseWriter, r *http.Request, params restapi.CreateTodoParams) {
	sessionId, err := parseSessionId(params.XSessionId)
	if err != nil {
		c.writeError(w, r, err)
		return
	}

	var req restapi.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.writeError(w, r, fmt.Errorf("%w: invalid JSON body: %w", errx.ErrInvalidArgument, err))
		return
	}

	todo, err := c.service.CreateTodo(r.Context(), &service.CreateTodoRequest{
		SessionId:   sessionId,
		Title:       req.Title,
		Description: utils.Value(req.Description),
		Completed:   utils.Value(req.Completed),
		Priority:    model.TodoPriority(utils.Value(req.Priority)),
		Tags:        utils.Value(req.Tags),
		DueAt:       req.DueAt,
	})
	if err != nil {
		c.writeError(w, r, err)
		return
	}

	restutils.Success(w, http.StatusCreated, &restapi.TodoResponse{Data: convert.TodoToRestApi(todo)})
}
