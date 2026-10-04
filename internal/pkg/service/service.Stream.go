package service

import (
	"context"
	"fmt"
	"time"

	"github.com/chapar-rest/mock-server/internal/pkg/errx"
)

const (
	maxStreamCount    = 100
	maxStreamInterval = 5 * time.Second
)

// CheckStream reports whether Stream accepts count and interval, for
// transports that must answer before the stream starts.
func (s *Service) CheckStream(count int, interval time.Duration) error {
	if count < 1 || count > maxStreamCount {
		return fmt.Errorf("%w: count must be between 1 and %d", errx.ErrInvalidArgument, maxStreamCount)
	}
	if interval < 0 || interval > maxStreamInterval {
		return fmt.Errorf("%w: interval must be between 0 and %d ms", errx.ErrInvalidArgument, maxStreamInterval.Milliseconds())
	}
	return nil
}

// Stream calls emit count times, interval apart, starting immediately. It
// stops early when emit fails or ctx is done.
func (s *Service) Stream(ctx context.Context, count int, interval time.Duration, emit func(index int, at time.Time) error) error {
	if err := s.CheckStream(count, interval); err != nil {
		return err
	}

	for i := range count {
		if i > 0 {
			if err := s.Delay(ctx, interval); err != nil {
				return err
			}
		}
		if err := emit(i, s.now()); err != nil {
			return err
		}
	}
	return nil
}
