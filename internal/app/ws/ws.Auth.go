package ws

import (
	"cmp"
	"net/http"

	"github.com/chapar-rest/mock-server/internal/pkg/model"
	"github.com/chapar-rest/mock-server/internal/pkg/service"
)

// authMessage is the first message of /auth.
type authMessage struct {
	Authenticated bool   `json:"authenticated"`
	Scheme        string `json:"scheme"`
	Credential    string `json:"credential"`
}

// Auth checks the handshake's credentials, says which were accepted, then
// echoes like /echo. It takes a Bearer or Basic Authorization header, or an
// API key in X-API-Key or the api_key query parameter. Browsers cannot set
// headers on a WebSocket, so access_token in the query counts as a bearer
// token.
func (c *Controller) Auth(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")
	if token := r.URL.Query().Get("access_token"); authorization == "" && token != "" {
		authorization = "Bearer " + token
	}

	result, err := c.service.Authenticate(r.Context(), &service.AuthenticateRequest{
		Authorization: authorization,
		ApiKey:        cmp.Or(r.Header.Get("X-API-Key"), r.URL.Query().Get("api_key")),
		Schemes:       []model.AuthScheme{model.AuthSchemeBearer, model.AuthSchemeBasic, model.AuthSchemeApiKey},
		Username:      "",
		Password:      "",
	})
	if err != nil {
		w.Header().Set("WWW-Authenticate", `Bearer realm="mock-server"`)
		c.writeError(w, r, err)
		return
	}

	s := c.accept(w, r)
	if s == nil {
		return
	}
	defer s.done()

	if err := s.writeJSON(&authMessage{
		Authenticated: true,
		Scheme:        string(result.Scheme),
		Credential:    result.Credential,
	}); err != nil {
		return
	}
	s.echo()
}
