package api

import (
	"errors"
	"fmt"
	"io"

	"google.golang.org/grpc"

	"github.com/chapar-rest/mock-server/internal/gen/mockv1"
	"github.com/chapar-rest/mock-server/internal/pkg/errx"
)

// maxClientStreamMessages bounds the memory one client stream can hold.
const maxClientStreamMessages = 1000

// ClientStream implements mockv1.UtilityServiceServer.
func (s *UtilityServer) ClientStream(stream grpc.ClientStreamingServer[mockv1.EchoRequest, mockv1.ClientStreamResponse]) error {
	messages := make([]string, 0)
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return stream.SendAndClose(&mockv1.ClientStreamResponse{
				Count:    int32(len(messages)),
				Messages: messages,
			})
		}
		if err != nil {
			return err
		}

		if len(messages) >= maxClientStreamMessages {
			return fmt.Errorf("%w: at most %d messages per stream", errx.ErrInvalidArgument, maxClientStreamMessages)
		}
		messages = append(messages, req.GetMessage())
	}
}
