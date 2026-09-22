package app

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	grpcapi "github.com/chapar-rest/mock-server/internal/app/grpc/api"
	"github.com/chapar-rest/mock-server/internal/gen/mockv1"
	"github.com/chapar-rest/mock-server/internal/gen/restapi"
	"github.com/chapar-rest/mock-server/internal/pkg/memstore"
	"github.com/chapar-rest/mock-server/internal/pkg/service"
	"github.com/chapar-rest/mock-server/internal/pkg/utils"
)

// newTestServer starts the full handler on one cleartext port speaking both
// HTTP/1.1 and h2c, as in production.
func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	store := memstore.New(memstore.Config{
		MaxSessions:        10,
		MaxTodosPerSession: 10,
		Seed:               service.SampleTodos,
		Now:                utils.Now,
	})
	svc := service.NewService(zap.NewNop(), store)

	srv := httptest.NewUnstartedServer(NewController(zap.NewNop(), svc, grpcapi.NewServer(zap.NewNop(), svc)))
	var protocols http.Protocols
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)
	srv.Config.Protocols = &protocols
	srv.Start()
	t.Cleanup(srv.Close)
	return srv
}

func doJSON(t *testing.T, method, url, session, body string, out any) *http.Response {
	t.Helper()

	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if session != "" {
		req.Header.Set("X-Session-Id", session)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, url, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			t.Fatalf("decoding %s %s: %v", method, url, err)
		}
	}
	return resp
}

func TestRestTodoLifecycle(t *testing.T) {
	srv := newTestServer(t)
	base := srv.URL + "/api/v1/todos"

	var created restapi.TodoResponse
	resp := doJSON(t, http.MethodPost, base, "rest-test", `{"title":"Buy milk","tags":["home"],"due_at":"2030-01-02T03:04:05Z"}`, &created)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create status = %d, want 201", resp.StatusCode)
	}
	if created.Data.Priority != restapi.TodoPriorityMedium || created.Data.DueAt == nil {
		t.Fatalf("created = %+v, want default priority and a due date", created.Data)
	}
	todoURL := base + "/" + created.Data.Id.String()

	// PATCH with an explicit null clears the due date and leaves the rest.
	var patched restapi.TodoResponse
	doJSON(t, http.MethodPatch, todoURL, "rest-test", `{"completed":true,"due_at":null}`, &patched)
	if !patched.Data.Completed || patched.Data.DueAt != nil || patched.Data.Title != "Buy milk" {
		t.Fatalf("patched = %+v", patched.Data)
	}

	// Other sessions do not see it.
	if resp := doJSON(t, http.MethodGet, todoURL, "someone-else", "", nil); resp.StatusCode != http.StatusNotFound {
		t.Fatalf("cross-session get status = %d, want 404", resp.StatusCode)
	}

	var list restapi.ListTodosResponse
	doJSON(t, http.MethodGet, base+"?completed=true&q=milk", "rest-test", "", &list)
	if list.Total != 1 {
		t.Fatalf("filtered total = %d, want 1", list.Total)
	}

	if resp := doJSON(t, http.MethodDelete, todoURL, "rest-test", "", nil); resp.StatusCode != http.StatusNoContent {
		t.Fatalf("delete status = %d, want 204", resp.StatusCode)
	}
	if resp := doJSON(t, http.MethodGet, todoURL, "rest-test", "", nil); resp.StatusCode != http.StatusNotFound {
		t.Fatalf("get after delete status = %d, want 404", resp.StatusCode)
	}
}

func TestRestValidation(t *testing.T) {
	srv := newTestServer(t)

	tests := []struct {
		name, method, path, body string
	}{
		{"empty title", http.MethodPost, "/api/v1/todos", `{"title":" "}`},
		{"bad priority", http.MethodPost, "/api/v1/todos", `{"title":"x","priority":"urgent"}`},
		{"malformed json", http.MethodPost, "/api/v1/todos", `{`},
		{"non-uuid id", http.MethodGet, "/api/v1/todos/42", ""},
		{"limit too large", http.MethodGet, "/api/v1/todos?limit=1000", ""},
		{"non-integer status", http.MethodGet, "/api/v1/status/abc", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body restapi.Error
			resp := doJSON(t, tt.method, srv.URL+tt.path, "", tt.body, &body)
			if resp.StatusCode != http.StatusBadRequest || body.Error == "" {
				t.Fatalf("status = %d, error = %q; want 400 with a message", resp.StatusCode, body.Error)
			}
		})
	}
}

func TestRestUtilities(t *testing.T) {
	srv := newTestServer(t)

	t.Run("echo", func(t *testing.T) {
		var echo restapi.Echo
		doJSON(t, http.MethodPost, srv.URL+"/api/v1/echo?a=1", "", `{"hello":"world"}`, &echo)
		if echo.Method != http.MethodPost || echo.Query["a"][0] != "1" || echo.Json == nil {
			t.Fatalf("echo = %+v", echo)
		}
	})

	t.Run("status", func(t *testing.T) {
		if resp := doJSON(t, http.MethodGet, srv.URL+"/api/v1/status/418", "", "", nil); resp.StatusCode != http.StatusTeapot {
			t.Fatalf("status = %d, want 418", resp.StatusCode)
		}
	})

	t.Run("redirect chain ends at echo", func(t *testing.T) {
		var echo restapi.Echo
		doJSON(t, http.MethodGet, srv.URL+"/api/v1/redirect/3", "", "", &echo)
		if echo.Path != "/api/v1/echo" {
			t.Fatalf("final path = %q, want /api/v1/echo", echo.Path)
		}
	})

	t.Run("basic auth", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, srv.URL+"/api/v1/auth/basic/alice/pw", nil)
		req.SetBasicAuth("alice", "wrong")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusUnauthorized || resp.Header.Get("WWW-Authenticate") == "" {
			t.Fatalf("status = %d, challenge = %q; want 401 with a challenge", resp.StatusCode, resp.Header.Get("WWW-Authenticate"))
		}
	})

	t.Run("stream", func(t *testing.T) {
		resp, err := http.Get(srv.URL + "/api/v1/stream/3")
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = resp.Body.Close() }()
		body, _ := io.ReadAll(resp.Body)
		if lines := strings.Count(string(body), "\n"); lines != 3 {
			t.Fatalf("stream lines = %d, want 3", lines)
		}
	})
}

func TestGrpcOnSamePort(t *testing.T) {
	srv := newTestServer(t)

	conn, err := grpc.NewClient(strings.TrimPrefix(srv.URL, "http://"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc.NewClient: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "x-session-id", "grpc-test")

	todos := mockv1.NewTodoServiceClient(conn)

	created, err := todos.CreateTodo(ctx, &mockv1.CreateTodoRequest{Title: "From gRPC", Priority: mockv1.TodoPriority_TODO_PRIORITY_HIGH})
	if err != nil {
		t.Fatalf("CreateTodo: %v", err)
	}

	updated, err := todos.UpdateTodo(ctx, &mockv1.UpdateTodoRequest{
		Todo:       &mockv1.Todo{Id: created.GetId(), Completed: true, Title: "ignored: not in mask"},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"completed"}},
	})
	if err != nil {
		t.Fatalf("UpdateTodo: %v", err)
	}
	if !updated.GetCompleted() || updated.GetTitle() != "From gRPC" {
		t.Fatalf("updated = %v, want only completed changed", updated)
	}

	// The REST API shares the same service and store.
	var viaRest restapi.TodoResponse
	doJSON(t, http.MethodGet, srv.URL+"/api/v1/todos/"+created.GetId(), "grpc-test", "", &viaRest)
	if !viaRest.Data.Completed {
		t.Fatalf("REST view of gRPC todo = %+v, want completed", viaRest.Data)
	}

	_, err = todos.GetTodo(ctx, &mockv1.GetTodoRequest{Id: "not-a-uuid"})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("GetTodo(bad id) code = %v, want InvalidArgument", status.Code(err))
	}

	utility := mockv1.NewUtilityServiceClient(conn)

	_, err = utility.Status(ctx, &mockv1.StatusRequest{Code: int32(codes.PermissionDenied)})
	if status.Code(err) != codes.PermissionDenied {
		t.Fatalf("Status code = %v, want PermissionDenied", status.Code(err))
	}

	stream, err := utility.ServerStream(ctx, &mockv1.StreamRequest{Message: "hi", Count: 3})
	if err != nil {
		t.Fatalf("ServerStream: %v", err)
	}
	received := 0
	for {
		_, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatalf("ServerStream.Recv: %v", err)
		}
		received++
	}
	if received != 3 {
		t.Fatalf("received %d stream messages, want 3", received)
	}

	bidi, err := utility.BidiStream(ctx)
	if err != nil {
		t.Fatalf("BidiStream: %v", err)
	}
	if err := bidi.Send(&mockv1.EchoRequest{Message: "ping"}); err != nil {
		t.Fatalf("BidiStream.Send: %v", err)
	}
	reply, err := bidi.Recv()
	if err != nil || reply.GetMessage() != "ping" {
		t.Fatalf("BidiStream.Recv = %v, %v; want ping", reply, err)
	}
	_ = bidi.CloseSend()
}
