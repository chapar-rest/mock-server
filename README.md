Chapar Mock Server
==================

A public mock API at `mocks.chapar.rest` for exercising API clients such as
[Chapar](https://github.com/chapar-rest/chapar). One process serves every
protocol on one port:

| Protocol | Where | Contract |
|----------|-------|----------|
| REST | `https://mocks.chapar.rest/api/v1` | [`api/rest/openapi.yaml`](api/rest/openapi.yaml) |
| gRPC | `mocks.chapar.rest:443` (server reflection enabled) | [`api/proto/mock/v1`](api/proto/mock/v1) |

GraphQL, WebSocket and MQTT are planned.

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

## Development

Requires Go, [Task](https://taskfile.dev), golangci-lint, and for code
generation oapi-codegen, protoc, protoc-gen-go and protoc-gen-go-grpc.

```bash
task dev       # run on :8080
task           # test + lint
task generate  # regenerate internal/gen from api/
```
