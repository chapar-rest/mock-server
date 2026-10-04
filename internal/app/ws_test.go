package app

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"

	"github.com/chapar-rest/mock-server/internal/app/ws"
	"github.com/chapar-rest/mock-server/internal/gen/restapi"
)

func dialWS(t *testing.T, ctx context.Context, url string, opts *websocket.DialOptions) (*websocket.Conn, *http.Response) {
	t.Helper()

	conn, resp, err := websocket.Dial(ctx, url, opts)
	if err != nil {
		status := 0
		if resp != nil {
			status = resp.StatusCode
		}
		t.Fatalf("dial %s: %v (status %d)", url, err, status)
	}
	t.Cleanup(func() { _ = conn.CloseNow() })
	return conn, resp
}

func wsURL(srv string, path string) string {
	return "ws" + strings.TrimPrefix(srv, "http") + path
}

func readJSON(t *testing.T, ctx context.Context, conn *websocket.Conn, out any) {
	t.Helper()

	typ, data, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if typ != websocket.MessageText {
		t.Fatalf("message type = %v, want text", typ)
	}
	if err := json.Unmarshal(data, out); err != nil {
		t.Fatalf("decoding %q: %v", data, err)
	}
}

func TestWebSocketEcho(t *testing.T) {
	srv := newTestServer(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//nolint:exhaustruct // only the subprotocol matters here.
	conn, _ := dialWS(t, ctx, wsURL(srv.URL, "/ws/echo"), &websocket.DialOptions{Subprotocols: []string{ws.Subprotocol}})
	if conn.Subprotocol() != ws.Subprotocol {
		t.Fatalf("subprotocol = %q, want %q", conn.Subprotocol(), ws.Subprotocol)
	}

	for _, msg := range []struct {
		typ  websocket.MessageType
		data []byte
	}{
		{websocket.MessageText, []byte(`{"hello":"world"}`)},
		{websocket.MessageBinary, []byte{0, 1, 2, 255}},
	} {
		if err := conn.Write(ctx, msg.typ, msg.data); err != nil {
			t.Fatalf("write: %v", err)
		}
		typ, data, err := conn.Read(ctx)
		if err != nil {
			t.Fatalf("read: %v", err)
		}
		if typ != msg.typ || string(data) != string(msg.data) {
			t.Fatalf("echo = %v %q, want %v %q", typ, data, msg.typ, msg.data)
		}
	}
}

func TestWebSocketStream(t *testing.T) {
	srv := newTestServer(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _ := dialWS(t, ctx, wsURL(srv.URL, "/ws/stream?count=3&interval_ms=50"), nil)
	// Messages from the client are ignored, not a reason to close.
	if err := conn.Write(ctx, websocket.MessageText, []byte("hello")); err != nil {
		t.Fatalf("write: %v", err)
	}
	for i := range 3 {
		var ev struct {
			Index int `json:"index"`
			Count int `json:"count"`
		}
		readJSON(t, ctx, conn, &ev)
		if ev.Index != i || ev.Count != 3 {
			t.Fatalf("event %d = %+v", i, ev)
		}
	}
	_, _, err := conn.Read(ctx)
	if websocket.CloseStatus(err) != websocket.StatusNormalClosure {
		t.Fatalf("after the stream: %v, want a normal close", err)
	}

	// Bad parameters are refused before the upgrade.
	resp, err := http.Get(srv.URL + "/ws/stream?count=1000")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("count=1000 status = %d, want 400", resp.StatusCode)
	}
}

func TestWebSocketTodos(t *testing.T) {
	srv := newTestServer(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _ := dialWS(t, ctx, wsURL(srv.URL, "/ws/todos?session=ws-test"), nil)
	var msg struct {
		Type    string        `json:"type"`
		Session string        `json:"session"`
		Id      string        `json:"id"`
		Todo    *restapi.Todo `json:"todo"`
	}
	readJSON(t, ctx, conn, &msg)
	if msg.Type != "subscribed" || msg.Session != "ws-test" {
		t.Fatalf("greeting = %+v", msg)
	}

	// Another session's change is not seen; this session's is.
	doJSON(t, http.MethodPost, srv.URL+"/api/v1/todos", "other", `{"title":"Not mine"}`, nil)
	var created restapi.TodoResponse
	doJSON(t, http.MethodPost, srv.URL+"/api/v1/todos", "ws-test", `{"title":"Live"}`, &created)

	readJSON(t, ctx, conn, &msg)
	if msg.Type != "todo.created" || msg.Todo == nil || msg.Todo.Title != "Live" || msg.Id != created.Data.Id.String() {
		t.Fatalf("event = %+v", msg)
	}

	doJSON(t, http.MethodDelete, srv.URL+"/api/v1/todos/"+created.Data.Id.String(), "ws-test", "", nil)
	msg.Todo = nil
	readJSON(t, ctx, conn, &msg)
	if msg.Type != "todo.deleted" || msg.Todo != nil || msg.Id != created.Data.Id.String() {
		t.Fatalf("event = %+v", msg)
	}
}

func TestWebSocketClose(t *testing.T) {
	srv := newTestServer(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _ := dialWS(t, ctx, wsURL(srv.URL, "/ws/close?code=4001&reason=bye"), nil)
	_, _, err := conn.Read(ctx)
	var ce websocket.CloseError
	if !errors.As(err, &ce) || ce.Code != 4001 || ce.Reason != "bye" {
		t.Fatalf("read = %v, want close 4001 bye", err)
	}

	resp, err := http.Get(srv.URL + "/ws/close?code=1005")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("code=1005 status = %d, want 400", resp.StatusCode)
	}
}

func TestWebSocketAuth(t *testing.T) {
	srv := newTestServer(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, resp, err := websocket.Dial(ctx, wsURL(srv.URL, "/ws/auth"), nil)
	if err == nil || resp == nil || resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("dial without credentials: err %v, resp %v; want 401", err, resp)
	}

	//nolint:exhaustruct // only the header matters here.
	conn, _ := dialWS(t, ctx, wsURL(srv.URL, "/ws/auth"), &websocket.DialOptions{
		HTTPHeader: http.Header{"Authorization": []string{"Bearer secret"}},
	})
	var auth struct {
		Scheme     string `json:"scheme"`
		Credential string `json:"credential"`
	}
	readJSON(t, ctx, conn, &auth)
	if auth.Scheme != "bearer" || auth.Credential != "secret" {
		t.Fatalf("auth = %+v", auth)
	}

	conn, _ = dialWS(t, ctx, wsURL(srv.URL, "/ws/auth?api_key=k1"), nil)
	readJSON(t, ctx, conn, &auth)
	if auth.Scheme != "api_key" || auth.Credential != "k1" {
		t.Fatalf("auth = %+v", auth)
	}
}

func TestWebSocketConnectionLimit(t *testing.T) {
	srv := newTestServer(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for range 5 {
		dialWS(t, ctx, wsURL(srv.URL, "/ws/echo"), nil)
	}
	_, resp, err := websocket.Dial(ctx, wsURL(srv.URL, "/ws/echo"), nil)
	if err == nil || resp == nil || resp.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("sixth dial: err %v, resp %v; want 429", err, resp)
	}
}

func TestWebSocketShutdown(t *testing.T) {
	srv := newTestServer(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _ := dialWS(t, ctx, wsURL(srv.URL, "/ws/echo"), nil)
	// Shutdown waits for the close handshake, which needs this side reading.
	go srv.Config.Handler.(*Controller).Shutdown()

	_, _, err := conn.Read(ctx)
	if websocket.CloseStatus(err) != websocket.StatusGoingAway {
		t.Fatalf("read after shutdown = %v, want going away", err)
	}
}
