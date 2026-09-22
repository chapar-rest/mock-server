package api

import (
	"net/http"

	"github.com/chapar-rest/mock-server/internal/pkg/model"
	"github.com/chapar-rest/mock-server/internal/pkg/service"
)

// BearerAuth implements restapi.ServerInterface.
func (c *Controller) BearerAuth(w http.ResponseWriter, r *http.Request) {
	c.authenticate(w, r, `Bearer realm="mock-server"`, &service.AuthenticateRequest{
		Authorization: r.Header.Get("Authorization"),
		ApiKey:        "",
		Schemes:       []model.AuthScheme{model.AuthSchemeBearer},
		Username:      "",
		Password:      "",
	})
}
