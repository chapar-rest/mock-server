package model

import (
	"fmt"
	"regexp"

	"github.com/chapar-rest/mock-server/internal/pkg/errx"
)

// SessionId scopes in-memory data to one client.
type SessionId string

// PublicSessionId is used by requests that do not name a session.
const PublicSessionId SessionId = "public"

var sessionIdPattern = regexp.MustCompile(`^[A-Za-z0-9._-]{1,64}$`)

// ParseSessionId parses a client-supplied session id. An empty value selects
// the shared public session.
func ParseSessionId(s string) (SessionId, error) {
	if s == "" {
		return PublicSessionId, nil
	}
	if !sessionIdPattern.MatchString(s) {
		return SessionId(""), fmt.Errorf("%w: session id must be 1-64 characters of A-Z, a-z, 0-9, '.', '_' or '-'", errx.ErrInvalidArgument)
	}
	return SessionId(s), nil
}

// String returns the string representation of the session id.
func (s SessionId) String() string {
	return string(s)
}
