package service

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"

	"go.uber.org/zap"

	"github.com/chapar-rest/mock-server/internal/pkg/errx"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

func basic(username, password string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
}

func TestAuthenticate(t *testing.T) {
	s := NewService(zap.NewNop(), nil)

	tests := []struct {
		name    string
		req     AuthenticateRequest
		want    *model.AuthResult
		wantErr error
	}{
		{
			name:    "basic matching path credentials",
			req:     AuthenticateRequest{Authorization: basic("alice", "s3cret"), ApiKey: "", Schemes: []model.AuthScheme{model.AuthSchemeBasic}, Username: "alice", Password: "s3cret"},
			want:    &model.AuthResult{Scheme: model.AuthSchemeBasic, Credential: "alice"},
			wantErr: nil,
		},
		{
			name:    "basic with wrong password",
			req:     AuthenticateRequest{Authorization: basic("alice", "nope"), ApiKey: "", Schemes: []model.AuthScheme{model.AuthSchemeBasic}, Username: "alice", Password: "s3cret"},
			want:    nil,
			wantErr: errx.ErrUnauthenticated,
		},
		{
			name:    "bearer any token",
			req:     AuthenticateRequest{Authorization: "Bearer abc.def", ApiKey: "", Schemes: []model.AuthScheme{model.AuthSchemeBearer}, Username: "", Password: ""},
			want:    &model.AuthResult{Scheme: model.AuthSchemeBearer, Credential: "abc.def"},
			wantErr: nil,
		},
		{
			name:    "bearer with empty token",
			req:     AuthenticateRequest{Authorization: "Bearer ", ApiKey: "", Schemes: []model.AuthScheme{model.AuthSchemeBearer}, Username: "", Password: ""},
			want:    nil,
			wantErr: errx.ErrUnauthenticated,
		},
		{
			name:    "bearer sent where only api key is accepted",
			req:     AuthenticateRequest{Authorization: "Bearer abc", ApiKey: "", Schemes: []model.AuthScheme{model.AuthSchemeApiKey}, Username: "", Password: ""},
			want:    nil,
			wantErr: errx.ErrUnauthenticated,
		},
		{
			name:    "api key",
			req:     AuthenticateRequest{Authorization: "", ApiKey: "k-123", Schemes: []model.AuthScheme{model.AuthSchemeBearer, model.AuthSchemeApiKey}, Username: "", Password: ""},
			want:    &model.AuthResult{Scheme: model.AuthSchemeApiKey, Credential: "k-123"},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Authenticate(context.Background(), &tt.req)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.want != nil && (got == nil || *got != *tt.want) {
				t.Fatalf("result = %+v, want %+v", got, tt.want)
			}
		})
	}
}
