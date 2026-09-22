package api

import (
	"context"
	"time"

	"github.com/chapar-rest/mock-server/internal/gen/mockv1"
)

// Delay implements mockv1.UtilityServiceServer.
func (s *UtilityServer) Delay(ctx context.Context, req *mockv1.DelayRequest) (*mockv1.EchoResponse, error) {
	if err := s.service.Delay(ctx, time.Duration(req.GetMs())*time.Millisecond); err != nil {
		return nil, err
	}

	return echoResponse(ctx, req.GetMessage(), 0, time.Now()), nil
}
