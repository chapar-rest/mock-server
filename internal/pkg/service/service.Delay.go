package service

import (
	"context"
	"fmt"
	"time"

	"github.com/chapar-rest/mock-server/internal/pkg/errx"
)

// MaxDelay bounds how long a delay request may hold a connection.
const MaxDelay = 10 * time.Second

// Delay waits for d, or until ctx is done.
func (s *Service) Delay(ctx context.Context, d time.Duration) error {
	if d < 0 || d > MaxDelay {
		return fmt.Errorf("%w: delay must be between 0 and %d ms", errx.ErrInvalidArgument, MaxDelay.Milliseconds())
	}

	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
