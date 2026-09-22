package api

import (
	"errors"
	"io"
	"time"

	"google.golang.org/grpc"

	"github.com/chapar-rest/mock-server/internal/gen/mockv1"
)

// BidiStream implements mockv1.UtilityServiceServer.
func (s *UtilityServer) BidiStream(stream grpc.BidiStreamingServer[mockv1.EchoRequest, mockv1.EchoResponse]) error {
	ctx := stream.Context()
	for index := 0; ; index++ {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}

		if err := stream.Send(echoResponse(ctx, req.GetMessage(), index, time.Now())); err != nil {
			return err
		}
	}
}
