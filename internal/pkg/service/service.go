package service

import (
	"time"

	"go.uber.org/zap"

	"github.com/chapar-rest/mock-server/internal/pkg/memstore"
	"github.com/chapar-rest/mock-server/internal/pkg/utils"
)

// Service holds the business logic shared by every transport (REST, gRPC,
// and later GraphQL, WebSocket and MQTT). Transports translate their wire
// types into service requests and map errx errors onto their status codes.
type Service struct {
	logger *zap.Logger
	store  *memstore.Store
	now    func() time.Time
}

// NewService creates a new service.
func NewService(logger *zap.Logger, store *memstore.Store) *Service {
	return &Service{
		logger: logger,
		store:  store,
		now:    utils.Now,
	}
}
