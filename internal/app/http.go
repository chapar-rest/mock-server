package app

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	restapi "github.com/chapar-rest/mock-server/internal/app/rest/api"
	"github.com/chapar-rest/mock-server/internal/pkg/service"
	"github.com/chapar-rest/mock-server/internal/pkg/utils/restutils"
)

// maxRequestBody caps every REST request body, uploads included.
const maxRequestBody = 10 << 20

// Controller serves every protocol on one port.
type Controller struct {
	logger *zap.Logger

	restApi    http.Handler
	grpcServer *grpc.Server

	router *chi.Mux
}

// NewController creates the root handler.
//
// REST routes:
//
//	/api/v1/*  → REST API (internal/app/rest/api)
//	/healthz   → liveness probe
//	/readyz    → readiness probe
//
// gRPC requests (HTTP/2 with an application/grpc content type) bypass the
// router and go straight to grpcServer, whatever their path.
func NewController(logger *zap.Logger, service *service.Service, grpcServer *grpc.Server) *Controller {
	router := chi.NewRouter()

	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.RequestSize(maxRequestBody))

	// Browser-based clients must be able to call every endpoint, with any
	// header, from any origin.
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowOriginFunc:    nil,
		AllowedMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders:     []string{"*"},
		ExposedHeaders:     []string{"*"},
		AllowCredentials:   false,
		MaxAge:             300, // Maximum value not ignored by any of major browsers
		OptionsPassthrough: false,
		Debug:              false,
	}))

	c := &Controller{
		logger:     logger,
		restApi:    restapi.NewController(logger.Named("rest"), service),
		grpcServer: grpcServer,
		router:     router,
	}

	c.router.Get("/healthz", probe)
	c.router.Get("/readyz", probe)

	c.router.Mount("/api/v1", c.restApi)

	c.router.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		restutils.JSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
	})
	c.router.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		restutils.JSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	})

	return c
}

// ServeHTTP serves the REST API, or the gRPC API for gRPC requests.
func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
		c.grpcServer.ServeHTTP(w, r)
		return
	}
	c.router.ServeHTTP(w, r)
}

// probe answers both liveness and readiness. All state is in memory, so a
// process that can answer is ready.
func probe(w http.ResponseWriter, _ *http.Request) {
	restutils.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
