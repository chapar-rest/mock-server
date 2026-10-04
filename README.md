Chapar Mock Server
==================

A public mock API at `mocks.chapar.rest` for exercising API clients such as
[Chapar](https://github.com/chapar-rest/chapar). One process serves every
protocol on one port:

| Protocol | Where | Contract |
|----------|-------|----------|
| REST | `https://mocks.chapar.rest/api/v1` | [`api/rest/openapi.yaml`](api/rest/openapi.yaml) |
| gRPC | `mocks.chapar.rest:443` (server reflection enabled) | [`api/proto/mock/v1`](api/proto/mock/v1) |
| WebSocket | `wss://mocks.chapar.rest/ws` | [below](#websocket) |

GraphQL and MQTT are planned.

## Todos

An in-memory todo list with full CRUD on every protocol. Send an
`X-Session-Id` header (gRPC: `x-session-id` metadata) to get your own list;
without it you share the `public` session. Sessions start with sample todos,
expire after an hour of inactivity, and hold at most 100 todos.

```bash
curl -H 'X-Session-Id: me' https://mocks.chapar.rest/api/v1/todos
curl -H 'X-Session-Id: me' -H 'Content-Type: application/json' \
  -d '{"title":"Buy milk","priority":"high"}' https://mocks.chapar.rest/api/v1/todos

grpcurl -H 'x-session-id: me' mocks.chapar.rest:443 mock.v1.TodoService/ListTodos
```

## Utilities

REST: `/echo`, `/status/{code}`, `/delay/{ms}`, `/auth/basic/{user}/{pass}`,
`/auth/bearer`, `/auth/api-key`, `/cookies`, `/cookies/set`, `/redirect/{n}`,
`/upload`, `/form`, `/stream/{n}`, `/bytes/{n}`, `/formats/{json|xml|html|text|csv}`.

gRPC `mock.v1.UtilityService`: `Echo`, `Status`, `Delay`, `ServerStream`,
`ClientStream`, `BidiStream`, `Auth`.

## WebSocket

| Endpoint | Behaviour |
|----------|-----------|
| `/ws/echo` | Sends every message back unchanged, text or binary. |
| `/ws/stream?count=10&interval_ms=1000` | Sends `{"index","count","time"}` `count` times (1-100), `interval_ms` apart (0-5000), then closes with 1000. |
| `/ws/todos` | Live changes to a session's todos, as REST and gRPC calls make them: `{"type":"subscribed"}` first, then `todo.created`, `todo.updated`, `todo.deleted` and `todos.reset` with `id` and `todo`. The session is the `X-Session-Id` header or the `session` query parameter. |
| `/ws/close?code=4000&reason=bye` | Closes at once with that code (1000-1003, 1007-1011, 3000-4999) and reason. |
| `/ws/auth` | Needs a Bearer or Basic `Authorization` header, an `X-API-Key` header, or `access_token` / `api_key` in the query (401 otherwise). Says which credential it accepted, then echoes. |

Every endpoint offers the `chapar.v1` subprotocol and permessage-deflate,
pings every 30 seconds, takes messages up to 64 KiB, and closes connections
after 10 minutes. Each client may hold 5 connections at once. Bad query
parameters fail the handshake with a JSON error, like the REST API.

```bash
websocat wss://mocks.chapar.rest/ws/echo
websocat -H 'X-Session-Id: me' wss://mocks.chapar.rest/ws/todos
```

## Development

Requires Go, [Task](https://taskfile.dev), golangci-lint, and for code
generation oapi-codegen, protoc, protoc-gen-go and protoc-gen-go-grpc.

```bash
task dev       # run on :8080
task           # test + lint
task generate  # regenerate internal/gen from api/
```
