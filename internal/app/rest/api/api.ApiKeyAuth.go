package api

import (
	"cmp"
	"net/http"

	"github.com/chapar-rest/mock-server/internal/pkg/model"
	"github.com/chapar-rest/mock-server/internal/pkg/service"
)

// ApiKeyAuth implements restapi.ServerInterface.
func (c *Controller) ApiKeyAuth(w http.ResponseWriter, r *http.Request) {
	c.authenticate(w, r, "", &service.AuthenticateRequest{
		Authorization: "",
		ApiKey:        cmp.Or(r.Header.Get("X-API-Key"), r.URL.Query().Get("api_key")),
		Schemes:       []model.AuthScheme{model.AuthSchemeApiKey},
		Username:      "",
		Password:      "",
	})
}
