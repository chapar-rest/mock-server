package api

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/chapar-rest/mock-server/internal/gen/mockv1"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
	"github.com/chapar-rest/mock-server/internal/pkg/service"
)

// Auth implements mockv1.UtilityServiceServer.
func (s *UtilityServer) Auth(ctx context.Context, _ *emptypb.Empty) (*mockv1.AuthResponse, error) {
	result, err := s.service.Authenticate(ctx, &service.AuthenticateRequest{
		Authorization: firstMetadata(ctx, "authorization"),
		ApiKey:        firstMetadata(ctx, "x-api-key"),
		Schemes:       []model.AuthScheme{model.AuthSchemeBearer, model.AuthSchemeBasic, model.AuthSchemeApiKey},
		Username:      "",
		Password:      "",
	})
	if err != nil {
		return nil, err
	}

	return &mockv1.AuthResponse{
		Authenticated: true,
		Scheme:        string(result.Scheme),
		Credential:    result.Credential,
	}, nil
}
