package worker

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/chapar-rest/mock-server/internal/pkg/memstore"
	"github.com/chapar-rest/mock-server/internal/pkg/utils"
)

// SessionCleanupAgent frees sessions that have been idle longer than ttl.
type SessionCleanupAgent struct {
	logger *zap.Logger
	store  *memstore.Store
	ttl    time.Duration
	now    func() time.Time
}

func NewSessionCleanupAgent(logger *zap.Logger, store *memstore.Store, ttl time.Duration) *SessionCleanupAgent {
	return &SessionCleanupAgent{
		logger: logger,
		store:  store,
		ttl:    ttl,
		now:    utils.Now,
	}
}

func (a *SessionCleanupAgent) Run(_ context.Context) error {
	removed, remaining := a.store.DeleteExpiredSessions(a.now().Add(-a.ttl))
	if removed > 0 {
		a.logger.Info("Removed expired sessions", zap.Int("removed", removed), zap.Int("remaining", remaining))
	}
	return nil
}
