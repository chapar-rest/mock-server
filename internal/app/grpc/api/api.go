package api

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/chapar-rest/mock-server/internal/gen/mockv1"
	"github.com/chapar-rest/mock-server/internal/pkg/service"
)

var (
	_ mockv1.TodoServiceServer    = (*TodoServer)(nil)
	_ mockv1.UtilityServiceServer = (*UtilityServer)(nil)
)

// Controller holds what every gRPC handler needs.
type Controller struct {
	logger  *zap.Logger
	service *service.Service
}

// TodoServer implements mock.v1.TodoService.
type TodoServer struct {
	mockv1.UnimplementedTodoServiceServer
	*Controller
}

// UtilityServer implements mock.v1.UtilityService.
type UtilityServer struct {
	mockv1.UnimplementedUtilityServiceServer
	*Controller
}

// NewServer creates the gRPC server with both services and server
// reflection registered. It is served through http.Handler on the main port;
// see internal/app/http.go.
func NewServer(logger *zap.Logger, service *service.Service) *grpc.Server {
	c := &Controller{
		logger:  logger,
		service: service,
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(c.unaryInterceptor),
		grpc.ChainStreamInterceptor(c.streamInterceptor),
	)

	mockv1.RegisterTodoServiceServer(server, &TodoServer{
		UnimplementedTodoServiceServer: mockv1.UnimplementedTodoServiceServer{},
		Controller:                     c,
	})
	mockv1.RegisterUtilityServiceServer(server, &UtilityServer{
		UnimplementedUtilityServiceServer: mockv1.UnimplementedUtilityServiceServer{},
		Controller:                        c,
	})
	reflection.Register(server)

	return server
}

// unaryInterceptor maps domain errors onto status codes, turns panics into
// INTERNAL and logs every call.
func (c *Controller) unaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	start := time.Now()
	defer func() {
		if p := recover(); p != nil {
			c.logger.Error("Panic in gRPC handler", zap.String("method", info.FullMethod), zap.Any("panic", p))
			err = status.Error(codes.Internal, "internal error")
		}
		c.logCall(info.FullMethod, start, err)
	}()

	resp, err = handler(ctx, req)
	return resp, c.toStatus(info.FullMethod, err)
}

// streamInterceptor is unaryInterceptor for streaming calls.
func (c *Controller) streamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	start := time.Now()
	defer func() {
		if p := recover(); p != nil {
			c.logger.Error("Panic in gRPC handler", zap.String("method", info.FullMethod), zap.Any("panic", p))
			err = status.Error(codes.Internal, "internal error")
		}
		c.logCall(info.FullMethod, start, err)
	}()

	return c.toStatus(info.FullMethod, handler(srv, ss))
}

func (c *Controller) logCall(method string, start time.Time, err error) {
	c.logger.Info("gRPC call",
		zap.String("method", method),
		zap.String("code", status.Code(err).String()),
		zap.Duration("duration", time.Since(start)),
	)
}
