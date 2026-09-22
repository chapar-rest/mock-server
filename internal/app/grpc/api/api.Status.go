package api

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/chapar-rest/mock-server/internal/gen/mockv1"
	"github.com/chapar-rest/mock-server/internal/pkg/errx"
)

// Status implements mockv1.UtilityServiceServer.
func (s *UtilityServer) Status(_ context.Context, req *mockv1.StatusRequest) (*emptypb.Empty, error) {
	if req.GetCode() < int32(codes.OK) || req.GetCode() > int32(codes.Unauthenticated) {
		return nil, fmt.Errorf("%w: code must be between 0 and 16", errx.ErrInvalidArgument)
	}

	code := codes.Code(req.GetCode())
	if code == codes.OK {
		return &emptypb.Empty{}, nil
	}

	message := req.GetMessage()
	if message == "" {
		message = code.String()
	}
	return nil, status.Error(code, message)
}
