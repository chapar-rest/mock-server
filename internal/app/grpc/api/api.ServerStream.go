package api

import (
	"time"

	"google.golang.org/grpc"

	"github.com/chapar-rest/mock-server/internal/gen/mockv1"
)

// ServerStream implements mockv1.UtilityServiceServer.
func (s *UtilityServer) ServerStream(req *mockv1.StreamRequest, stream grpc.ServerStreamingServer[mockv1.EchoResponse]) error {
	ctx := stream.Context()
	interval := time.Duration(req.GetIntervalMs()) * time.Millisecond

	return s.service.Stream(ctx, int(req.GetCount()), interval, func(index int, at time.Time) error {
		return stream.Send(echoResponse(ctx, req.GetMessage(), index, at))
	})
}
