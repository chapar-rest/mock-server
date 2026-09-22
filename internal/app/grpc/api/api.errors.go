package api

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/chapar-rest/mock-server/internal/pkg/errx"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

// toStatus maps a domain error onto a gRPC status. Errors that already carry
// a status pass through unchanged.
func (c *Controller) toStatus(method string, err error) error {
	if err == nil {
		return nil
	}
	if _, ok := status.FromError(err); ok {
		return err
	}

	switch {
	case errors.Is(err, errx.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, errx.ErrUnauthenticated):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, errx.ErrTodoNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, errx.ErrTodoLimitReached):
		return status.Error(codes.ResourceExhausted, err.Error())
	case errors.Is(err, errx.ErrSessionLimitReached):
		return status.Error(codes.Unavailable, err.Error())
	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, err.Error())
	}

	c.logger.Error("gRPC call failed", zap.String("method", method), zap.Error(err))
	return status.Error(codes.Internal, "internal error")
}

// sessionIdFromContext reads the optional x-session-id metadata.
func sessionIdFromContext(ctx context.Context) (model.SessionId, error) {
	return model.ParseSessionId(firstMetadata(ctx, "x-session-id"))
}

// firstMetadata returns the first value of an incoming metadata key, or "".
func firstMetadata(ctx context.Context, key string) string {
	if values := metadata.ValueFromIncomingContext(ctx, key); len(values) > 0 {
		return values[0]
	}
	return ""
}
