package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/chapar-rest/mock-server/internal/gen/restapi"
	"github.com/chapar-rest/mock-server/internal/pkg/convert"
	"github.com/chapar-rest/mock-server/internal/pkg/errx"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
	"github.com/chapar-rest/mock-server/internal/pkg/service"
	"github.com/chapar-rest/mock-server/internal/pkg/utils"
	"github.com/chapar-rest/mock-server/internal/pkg/utils/restutils"
)

// UpdateTodo implements restapi.ServerInterface.
func (c *Controller) UpdateTodo(w http.ResponseWriter, r *http.Request, rawTodoId restapi.TodoId, params restapi.UpdateTodoParams) {
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.writeError(w, r, fmt.Errorf("%w: reading body: %w", errx.ErrInvalidArgument, err))
		return
	}

	var req restapi.UpdateTodoRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.writeError(w, r, fmt.Errorf("%w: invalid JSON body: %w", errx.ErrInvalidArgument, err))
		return
	}

	dueAt, err := patchDueAt(body)
	if err != nil {
		c.writeError(w, r, err)
		return
	}

	var priority utils.Optional[model.TodoPriority]
	if req.Priority != nil {
		priority = utils.Some(model.TodoPriority(*req.Priority))
	}

	todo, err := c.service.UpdateTodo(r.Context(), &service.UpdateTodoRequest{
		SessionId:   sessionId,
		Id:          todoId,
		Title:       utils.PtrToOptional(req.Title),
		Description: utils.PtrToOptional(req.Description),
		Completed:   utils.PtrToOptional(req.Completed),
		Priority:    priority,
		Tags:        utils.PtrToOptional(req.Tags),
		DueAt:       dueAt,
	})
	if err != nil {
		c.writeError(w, r, err)
		return
	}

	restutils.Success(w, http.StatusOK, &restapi.TodoResponse{Data: convert.TodoToRestApi(todo)})
}

// patchDueAt tells an absent due_at (leave unchanged) apart from an explicit
// null (clear it), which the generated TodoPatch type cannot.
func patchDueAt(body []byte) (utils.Optional[*time.Time], error) {
	var raw struct {
		DueAt json.RawMessage `json:"due_at"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return utils.None[*time.Time](), fmt.Errorf("%w: invalid JSON body: %w", errx.ErrInvalidArgument, err)
	}

	switch string(raw.DueAt) {
	case "":
		return utils.None[*time.Time](), nil
	case "null":
		return utils.Some[*time.Time](nil), nil
	}

	var dueAt time.Time
	if err := json.Unmarshal(raw.DueAt, &dueAt); err != nil {
		return utils.None[*time.Time](), fmt.Errorf("%w: due_at must be an RFC 3339 date-time", errx.ErrInvalidArgument)
	}
	return utils.Some(&dueAt), nil
}
