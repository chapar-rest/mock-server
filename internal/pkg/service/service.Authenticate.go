package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"slices"
	"strings"

	"github.com/chapar-rest/mock-server/internal/pkg/errx"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

type AuthenticateRequest struct {
	// Authorization is the raw Authorization header or metadata value.
	Authorization string
	// ApiKey is the API key sent in a header, query parameter or metadata.
	ApiKey string
	// Schemes lists the accepted schemes, in order of preference.
	Schemes []model.AuthScheme
	// Username and Password, when set, are the only Basic credentials
	// accepted. When empty, any non-empty Basic credentials are accepted.
	Username string
	Password string
}

// Authenticate accepts the first credential matching one of req.Schemes. Any
// non-empty bearer token or API key is valid: these endpoints exist to
// exercise how a client sends credentials, not to protect anything.
func (s *Service) Authenticate(_ context.Context, req *AuthenticateRequest) (*model.AuthResult, error) {
	scheme, credential, _ := strings.Cut(req.Authorization, " ")
	credential = strings.TrimSpace(credential)

	for _, accepted := range req.Schemes {
		switch accepted {
		case model.AuthSchemeBasic:
			if !strings.EqualFold(scheme, "basic") {
				continue
			}
			username, ok := checkBasic(credential, req.Username, req.Password)
			if ok {
				return &model.AuthResult{Scheme: model.AuthSchemeBasic, Credential: username}, nil
			}
		case model.AuthSchemeBearer:
			if strings.EqualFold(scheme, "bearer") && credential != "" {
				return &model.AuthResult{Scheme: model.AuthSchemeBearer, Credential: credential}, nil
			}
		case model.AuthSchemeApiKey:
			if req.ApiKey != "" {
				return &model.AuthResult{Scheme: model.AuthSchemeApiKey, Credential: req.ApiKey}, nil
			}
		}
	}

	return nil, fmt.Errorf("%w: expected %s credentials", errx.ErrUnauthenticated, describeSchemes(req.Schemes))
}

// checkBasic decodes Basic credentials and matches them against the expected
// username and password, if any.
func checkBasic(encoded, wantUsername, wantPassword string) (string, bool) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", false
	}
	username, password, ok := strings.Cut(string(decoded), ":")
	if !ok || username == "" {
		return "", false
	}
	if wantUsername != "" && (username != wantUsername || password != wantPassword) {
		return "", false
	}
	return username, true
}

func describeSchemes(schemes []model.AuthScheme) string {
	names := make([]string, 0, len(schemes))
	for _, s := range schemes {
		names = append(names, string(s))
	}
	slices.Sort(names)
	return strings.Join(names, " or ")
}
