package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/chapar-rest/mock-server/api/petstorev1"
)

func (g *Grpc) GetPetByID(ctx context.Context, req *petstorev1.PetID) (*petstorev1.Pet, error) {
	// Create trailer in defer to record function return time.
	defer func() {
		trailer := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
		_ = grpc.SetTrailer(ctx, trailer)
	}()

	// Read metadata from client.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "GetPetByID: failed to get metadata")
	}

	for k, v := range md {
		fmt.Printf("Meta data, key: %s, value: %s\n", k, v)
	}

	// Create and send header.
	header := metadata.New(map[string]string{"method": "GetPetByID", "timestamp": time.Now().Format(time.StampNano)})
	_ = grpc.SendHeader(ctx, header)

	return &petstorev1.Pet{
		Id:   req.GetPetId(),
		Name: "dog",
	}, nil
}
