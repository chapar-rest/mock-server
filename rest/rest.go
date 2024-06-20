package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/chapar-rest/mock-server/api"
)

type Rest struct {
}

func New() *Rest {
	r := &Rest{}
	return r
}

func (ro *Rest) Run() error {
	router := chi.NewRouter()
	router.Use(chimiddleware.StripSlashes)
	router.Use(chimiddleware.RealIP)

	api.HandlerWithOptions(ro, api.ChiServerOptions{
		BaseRouter: router,
	})

	slog.Info("starting server", "port", 8080)
	if err := http.ListenAndServe(":8080", router); err != nil {
		slog.Error("failed to start server", "error", err)
		return err
	}

	return nil
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
