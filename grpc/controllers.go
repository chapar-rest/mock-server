package grpc

import (
	"context"

	"github.com/chapar-rest/mock-server/api/petstorev1"
)

func (g *Grpc) GetPetByID(ctx context.Context, req *petstorev1.PetID) (*petstorev1.Pet, error) {
	return &petstorev1.Pet{
		Id:   req.GetPetId(),
		Name: "dog",
	}, nil
}
