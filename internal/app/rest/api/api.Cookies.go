package api

import (
	"net/http"

	"github.com/chapar-rest/mock-server/internal/gen/restapi"
	"github.com/chapar-rest/mock-server/internal/pkg/utils/restutils"
)

// GetCookies implements restapi.ServerInterface.
func (c *Controller) GetCookies(w http.ResponseWriter, r *http.Request) {
	restutils.Success(w, http.StatusOK, &restapi.CookiesResponse{Cookies: requestCookies(r)})
}

// SetCookies implements restapi.ServerInterface.
func (c *Controller) SetCookies(w http.ResponseWriter, r *http.Request) {
	cookies := requestCookies(r)
	for name, values := range r.URL.Query() {
		value := values[len(values)-1]
		//nolint:exhaustruct // a plain session cookie; every other attribute keeps its default.
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    value,
			Path:     "/",
			HttpOnly: true,
		})
		cookies[name] = value
	}

	restutils.Success(w, http.StatusOK, &restapi.CookiesResponse{Cookies: cookies})
}

func requestCookies(r *http.Request) map[string]string {
	cookies := make(map[string]string)
	for _, cookie := range r.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}
	return cookies
}
