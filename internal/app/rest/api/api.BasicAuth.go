package api

import (
	"net/http"

	"github.com/chapar-rest/mock-server/internal/pkg/model"
	"github.com/chapar-rest/mock-server/internal/pkg/service"
)

// BasicAuth implements restapi.ServerInterface.
func (c *Controller) BasicAuth(w http.ResponseWriter, r *http.Request, username string, password string) {
	c.authenticate(w, r, `Basic realm="mock-server"`, &service.AuthenticateRequest{
		Authorization: r.Header.Get("Authorization"),
		ApiKey:        "",
		Schemes:       []model.AuthScheme{model.AuthSchemeBasic},
		Username:      username,
		Password:      password,
	})
}
