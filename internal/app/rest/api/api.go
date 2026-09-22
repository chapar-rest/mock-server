package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/chapar-rest/mock-server/internal/gen/restapi"
	"github.com/chapar-rest/mock-server/internal/pkg/errx"
	"github.com/chapar-rest/mock-server/internal/pkg/service"
)

var _ restapi.ServerInterface = (*Controller)(nil)

// Controller is a controller for the REST API.
type Controller struct {
	logger  *zap.Logger
	mux     *chi.Mux
	service *service.Service
}

// NewController creates a new controller for the REST API.
func NewController(logger *zap.Logger, service *service.Service) *Controller {
	c := &Controller{
		logger:  logger,
		service: service,
		mux:     chi.NewRouter(),
	}

	restapi.HandlerWithOptions(c,
		restapi.ChiServerOptions{ //nolint:exhaustruct // BaseURL and Middlewares keep their defaults.
			BaseRouter: c.mux,
			// Parameter binding errors are client errors; answer them in the
			// same JSON shape as every other error.
			ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
				c.writeError(w, r, fmt.Errorf("%w: %w", errx.ErrInvalidArgument, err))
			},
		},
	)

	return c
}

// ServeHTTP serves the API.
func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.mux.ServeHTTP(w, r)
}
