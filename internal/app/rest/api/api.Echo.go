package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/chapar-rest/mock-server/internal/gen/restapi"
	"github.com/chapar-rest/mock-server/internal/pkg/errx"
	"github.com/chapar-rest/mock-server/internal/pkg/utils/restutils"
)

// maxEchoBody is how much of the request body the echo response repeats.
const maxEchoBody = 64 << 10

// EchoGet implements restapi.ServerInterface.
func (c *Controller) EchoGet(w http.ResponseWriter, r *http.Request) { c.echo(w, r) }

// EchoPost implements restapi.ServerInterface.
func (c *Controller) EchoPost(w http.ResponseWriter, r *http.Request) { c.echo(w, r) }

// EchoPut implements restapi.ServerInterface.
func (c *Controller) EchoPut(w http.ResponseWriter, r *http.Request) { c.echo(w, r) }

// EchoPatch implements restapi.ServerInterface.
func (c *Controller) EchoPatch(w http.ResponseWriter, r *http.Request) { c.echo(w, r) }

// EchoDelete implements restapi.ServerInterface.
func (c *Controller) EchoDelete(w http.ResponseWriter, r *http.Request) { c.echo(w, r) }

func (c *Controller) echo(w http.ResponseWriter, r *http.Request) {
	echo, err := echoRequest(r)
	if err != nil {
		c.writeError(w, r, err)
		return
	}
	restutils.Success(w, http.StatusOK, echo)
}

// echoRequest describes the request as the server received it.
func echoRequest(r *http.Request) (*restapi.Echo, error) {
	body, err := io.ReadAll(io.LimitReader(r.Body, maxEchoBody))
	if err != nil {
		return nil, fmt.Errorf("%w: reading body: %w", errx.ErrInvalidArgument, err)
	}

	var parsed any
	if len(body) > 0 && json.Valid(body) {
		_ = json.Unmarshal(body, &parsed)
	}

	// RemoteAddr is already the client address: the RealIP middleware
	// rewrites it from X-Forwarded-For / X-Real-IP set by the ingress.
	remoteIp := r.RemoteAddr
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		remoteIp = host
	}

	return &restapi.Echo{
		Method:   r.Method,
		Url:      r.URL.RequestURI(),
		Path:     r.URL.Path,
		Protocol: r.Proto,
		Host:     r.Host,
		RemoteIp: remoteIp,
		Query:    r.URL.Query(),
		Headers:  r.Header,
		Body:     string(body),
		Json:     parsed,
	}, nil
}
