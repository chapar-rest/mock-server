// Package ws serves the WebSocket API. Every endpoint upgrades a plain GET
// under /ws; failures before the upgrade are answered as JSON errors, like
// the REST API.
package ws

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/chapar-rest/mock-server/internal/pkg/errx"
	"github.com/chapar-rest/mock-server/internal/pkg/service"
	"github.com/chapar-rest/mock-server/internal/pkg/utils/restutils"
)

const (
	// Subprotocol is offered to clients that ask for it, so they can check
	// they negotiate one. Clients that ask for none get none.
	Subprotocol = "chapar.v1"

	// readLimit caps one incoming message.
	readLimit = 64 << 10
	// maxLifetime closes every connection after a while: this is a public
	// mock, not a place to keep sockets open for days.
	maxLifetime = 10 * time.Minute
	// pingInterval keeps idle connections alive through proxies that close
	// them after about 100 seconds without traffic.
	pingInterval = 30 * time.Second
	// writeTimeout bounds one write to a slow client.
	writeTimeout = 10 * time.Second

	// maxConnsPerClient and maxConns bound open sockets, which outlive the
	// per-request limits the ingress applies.
	maxConnsPerClient = 5
	maxConns          = 500
)

var errTooManyConnections = errors.New("too many open WebSocket connections")

// Controller serves the WebSocket API.
type Controller struct {
	logger  *zap.Logger
	service *service.Service
	mux     *chi.Mux

	mu       sync.Mutex
	conns    map[*websocket.Conn]string // conn → client key
	perKey   map[string]int
	shutdown bool
}

// NewController creates the WebSocket API.
//
//	/echo    → echoes every message back, text or binary
//	/stream  → sends count events, interval_ms apart, then closes
//	/todos   → live changes to a session's todos
//	/close   → closes at once with the given code and reason
//	/auth    → like /echo, after checking Bearer, Basic or API key credentials
func NewController(logger *zap.Logger, service *service.Service) *Controller {
	c := &Controller{
		logger:   logger,
		service:  service,
		mux:      chi.NewRouter(),
		mu:       sync.Mutex{},
		conns:    map[*websocket.Conn]string{},
		perKey:   map[string]int{},
		shutdown: false,
	}

	c.mux.Get("/echo", c.Echo)
	c.mux.Get("/stream", c.Stream)
	c.mux.Get("/todos", c.Todos)
	c.mux.Get("/close", c.Close)
	c.mux.Get("/auth", c.Auth)

	return c
}

// ServeHTTP serves the API.
func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.mux.ServeHTTP(w, r)
}

// Shutdown tells every open connection the server is going away.
// http.Server.Shutdown does not track upgraded connections, so it is
// registered with RegisterOnShutdown.
func (c *Controller) Shutdown() {
	c.mu.Lock()
	c.shutdown = true
	conns := make([]*websocket.Conn, 0, len(c.conns))
	for conn := range c.conns {
		conns = append(conns, conn)
	}
	c.mu.Unlock()

	var wg sync.WaitGroup
	for _, conn := range conns {
		wg.Go(func() { _ = conn.Close(websocket.StatusGoingAway, "server shutting down") })
	}
	wg.Wait()
}

// session is one accepted connection. ctx ends when the connection's
// lifetime is up or the request is cancelled.
type session struct {
	conn *websocket.Conn
	ctx  context.Context
	done func()
}

// accept upgrades the request. On failure the client has been answered and
// accept returns nil. Callers must call s.done when they finish.
func (c *Controller) accept(w http.ResponseWriter, r *http.Request) *session {
	key := clientKey(r)

	c.mu.Lock()
	if c.shutdown || len(c.conns) >= maxConns || c.perKey[key] >= maxConnsPerClient {
		c.mu.Unlock()
		restutils.Error(w, http.StatusTooManyRequests, errTooManyConnections)
		return nil
	}
	// Reserve the slot before the upgrade so concurrent dials cannot
	// overshoot the limit.
	c.perKey[key]++
	c.mu.Unlock()

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{ //nolint:exhaustruct // ping and pong callbacks keep their defaults.
		Subprotocols: []string{Subprotocol},
		// Any origin may connect, as with the REST API's CORS policy.
		InsecureSkipVerify:   true,
		CompressionMode:      websocket.CompressionNoContextTakeover,
		CompressionThreshold: 0,
	})
	if err != nil {
		// Accept has already written the HTTP error.
		c.release(nil, key)
		return nil
	}
	conn.SetReadLimit(readLimit)

	c.mu.Lock()
	c.conns[conn] = key
	c.mu.Unlock()

	ctx, cancel := context.WithTimeout(r.Context(), maxLifetime)
	go c.keepAlive(ctx, conn)

	return &session{
		conn: conn,
		ctx:  ctx,
		done: func() {
			cancel()
			c.release(conn, key)
			_ = conn.CloseNow()
		},
	}
}

func (c *Controller) release(conn *websocket.Conn, key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if conn != nil {
		delete(c.conns, conn)
	}
	c.perKey[key]--
	if c.perKey[key] <= 0 {
		delete(c.perKey, key)
	}
}

// keepAlive pings the client until ctx ends. When the lifetime is up it
// closes the connection with a reason the client can show.
func (c *Controller) keepAlive(ctx context.Context, conn *websocket.Conn) {
	t := time.NewTicker(pingInterval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			pingCtx, cancel := context.WithTimeout(ctx, writeTimeout)
			_ = conn.Ping(pingCtx)
			cancel()
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				_ = conn.Close(websocket.StatusPolicyViolation, "connection lifetime of 10 minutes reached")
			}
			return
		}
	}
}

// writeJSON sends v as one text message.
func (s *session) writeJSON(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return s.write(websocket.MessageText, data)
}

func (s *session) write(typ websocket.MessageType, data []byte) error {
	ctx, cancel := context.WithTimeout(s.ctx, writeTimeout)
	defer cancel()
	return s.conn.Write(ctx, typ, data)
}

// writeError answers a request that failed before the upgrade.
func (c *Controller) writeError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, errx.ErrInvalidArgument):
		restutils.Error(w, http.StatusBadRequest, err)
	case errors.Is(err, errx.ErrUnauthenticated):
		restutils.Error(w, http.StatusUnauthorized, err)
	case errors.Is(err, errx.ErrSessionLimitReached):
		restutils.Error(w, http.StatusServiceUnavailable, err)
	default:
		c.logger.Error("WebSocket request failed", zap.String("path", r.URL.Path), zap.Error(err))
		restutils.Error(w, http.StatusInternalServerError, errors.New("internal error"))
	}
}

// clientKey identifies the client for the connection limit. Traffic arrives
// through Cloudflare, whose CF-Connecting-IP is the real client; locally the
// remote address is.
func clientKey(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	return cmp.Or(r.Header.Get("CF-Connecting-IP"), host)
}
