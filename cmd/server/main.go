package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/chapar-rest/mock-server/internal/app"
	grpcapi "github.com/chapar-rest/mock-server/internal/app/grpc/api"
	"github.com/chapar-rest/mock-server/internal/pkg/memstore"
	"github.com/chapar-rest/mock-server/internal/pkg/service"
	"github.com/chapar-rest/mock-server/internal/pkg/utils"
	"github.com/chapar-rest/mock-server/internal/pkg/worker"
)

// serviceVersion is stamped by ko's ldflags (see .ko.yaml).
var serviceVersion = "snapshot"

const (
	// readHeaderTimeout guards against slow-header clients. There is no
	// overall read or write timeout: uploads, /delay, streams and
	// long-lived gRPC streams legitimately hold a connection open.
	readHeaderTimeout = 10 * time.Second
	idleTimeout       = 120 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type config struct {
	Port     string `envconfig:"PORT" default:"8080"`
	LogLevel string `envconfig:"MOCK_LOG_LEVEL" default:"info"`

	// Sessions idle for longer than SessionTTL are removed by the cleanup
	// agent, which runs every CleanupInterval.
	SessionTTL      time.Duration `envconfig:"MOCK_SESSION_TTL" default:"1h"`
	CleanupInterval time.Duration `envconfig:"MOCK_CLEANUP_INTERVAL" default:"1m"`

	// Together these bound memory: every todo is at most a few KiB.
	MaxSessions        int `envconfig:"MOCK_MAX_SESSIONS" default:"1000"`
	MaxTodosPerSession int `envconfig:"MOCK_MAX_TODOS_PER_SESSION" default:"100"`
}

func main() {
	var cfg config
	utils.MustProcess("", &cfg)
	utils.SetVersion(serviceVersion)

	zapCfg := zap.NewProductionConfig()
	zapCfg.Level = utils.Must(zap.ParseAtomicLevel(cfg.LogLevel))

	logger := utils.Must(zapCfg.Build())

	defer func() {
		if err := logger.Sync(); err != nil {
			log.Println("Failed to sync logger", err)
		}
	}()

	store := memstore.New(memstore.Config{
		MaxSessions:        cfg.MaxSessions,
		MaxTodosPerSession: cfg.MaxTodosPerSession,
		Seed:               service.SampleTodos,
		Now:                utils.Now,
	})

	service := service.NewService(logger.Named("service"), store)

	var protocols http.Protocols
	protocols.SetHTTP1(true)
	// gRPC arrives as cleartext HTTP/2 (h2c) from the ingress, or directly
	// from grpcurl -plaintext in local development.
	protocols.SetUnencryptedHTTP2(true)

	//nolint:exhaustruct // every other field keeps the net/http default.
	server := &http.Server{
		Addr:              net.JoinHostPort("0.0.0.0", cfg.Port),
		Handler:           app.NewController(logger, service, grpcapi.NewServer(logger.Named("grpc"), service)),
		Protocols:         &protocols,
		ReadHeaderTimeout: readHeaderTimeout,
		IdleTimeout:       idleTimeout,
	}

	workerService := worker.NewWorker(logger.Named("worker"), cfg.CleanupInterval)
	if err := workerService.AddAgent("session-cleanup", worker.NewSessionCleanupAgent(logger.Named("session-cleanup"), store, cfg.SessionTTL)); err != nil {
		logger.Fatal("Failed to register session-cleanup agent", zap.Error(err))
	}

	// context with signal; Kubernetes stops pods with SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go workerService.Start(ctx)

	go func() {
		logger.Info("Starting server", zap.String("port", cfg.Port), zap.String("version", utils.Version()), zap.String("env", utils.Env()))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Failed to start server", zap.Error(err))
			stop()
		}
	}()

	<-ctx.Done()
	logger.Info("Shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		// Long-lived streams outlived the grace period; cut them off.
		logger.Warn("Graceful shutdown timed out", zap.Error(err))
		_ = server.Close()
	}
}
