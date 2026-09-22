package api

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/chapar-rest/mock-server/internal/gen/mockv1"
)

// Echo implements mockv1.UtilityServiceServer.
func (s *UtilityServer) Echo(ctx context.Context, req *mockv1.EchoRequest) (*mockv1.EchoResponse, error) {
	// Header and trailer let clients check they surface both.
	_ = grpc.SetHeader(ctx, metadata.Pairs("method", "Echo"))
	defer func() {
		_ = grpc.SetTrailer(ctx, metadata.Pairs("timestamp", time.Now().UTC().Format(time.RFC3339Nano)))
	}()

	return echoResponse(ctx, req.GetMessage(), 0, time.Now()), nil
}

func echoResponse(ctx context.Context, message string, index int, at time.Time) *mockv1.EchoResponse {
	return &mockv1.EchoResponse{
		Message:  message,
		Metadata: incomingMetadata(ctx),
		Index:    int32(index),
		Time:     timestamppb.New(at),
	}
}

func incomingMetadata(ctx context.Context) map[string]*mockv1.MetadataValues {
	md, _ := metadata.FromIncomingContext(ctx)
	out := make(map[string]*mockv1.MetadataValues, len(md))
	for key, values := range md {
		out[key] = &mockv1.MetadataValues{Values: values}
	}
	return out
}
