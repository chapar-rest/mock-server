package api

import (
	"net/http"

	"github.com/chapar-rest/mock-server/internal/gen/restapi"
	"github.com/chapar-rest/mock-server/internal/pkg/service"
	"github.com/chapar-rest/mock-server/internal/pkg/utils/restutils"
)

// authenticate answers an auth utility endpoint. challenge, when set, is sent
// as WWW-Authenticate on failure so clients can prompt for credentials.
func (c *Controller) authenticate(w http.ResponseWriter, r *http.Request, challenge string, req *service.AuthenticateRequest) {
	result, err := c.service.Authenticate(r.Context(), req)
	if err != nil {
		if challenge != "" {
			w.Header().Set("WWW-Authenticate", challenge)
		}
		c.writeError(w, r, err)
		return
	}

	restutils.Success(w, http.StatusOK, &restapi.AuthResponse{
		Authenticated: true,
		Scheme:        string(result.Scheme),
		Credential:    result.Credential,
	})
}
