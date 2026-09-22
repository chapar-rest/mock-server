package api

import (
	"context"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/chapar-rest/mock-server/internal/pkg/errx"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
	"github.com/chapar-rest/mock-server/internal/pkg/utils/restutils"
)

// writeError maps a domain error onto an HTTP status and writes it.
func (c *Controller) writeError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, errx.ErrInvalidArgument):
		restutils.Error(w, http.StatusBadRequest, err)
	case errors.Is(err, errx.ErrUnauthenticated):
		restutils.Error(w, http.StatusUnauthorized, err)
	case errors.Is(err, errx.ErrTodoNotFound):
		restutils.Error(w, http.StatusNotFound, err)
	case errors.Is(err, errx.ErrTodoLimitReached):
		restutils.Error(w, http.StatusConflict, err)
	case errors.Is(err, errx.ErrSessionLimitReached):
		restutils.Error(w, http.StatusServiceUnavailable, err)
	case errors.Is(err, context.Canceled):
		// The client went away; there is nobody to answer.
	default:
		c.logger.Error("Request failed", zap.String("path", r.URL.Path), zap.Error(err))
		restutils.Error(w, http.StatusInternalServerError, errors.New("internal error"))
	}
}

// parseSessionId reads the optional X-Session-Id header.
func parseSessionId(header *string) (model.SessionId, error) {
	if header == nil {
		return model.PublicSessionId, nil
	}
	return model.ParseSessionId(*header)
}
