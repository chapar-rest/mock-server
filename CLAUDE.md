# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Public mock API (`mocks.chapar.rest`) for exercising API clients like Chapar. A single Go process serves REST and gRPC on one port; GraphQL, WebSocket and MQTT are planned and must reuse the same service layer.

## Commands

Task automation uses [Task](https://taskfile.dev) (`Taskfile.yaml`), not Make.

```bash
task dev          # run server on :8080 (MOCK_ENV=local)
task              # test + lint
task test         # go test ./...
task lint         # golangci-lint run (v2 config: .golangci.yaml)
task generate     # rm -rf internal/gen && go generate ./...
task tidy         # go mod tidy + go mod vendor (deps are vendored)
task publish      # ko build → ghcr.io/chapar-rest/mock-server, then kubectl set image
```
Single test: `go test ./internal/app -run TestGrpcOnSamePort -v`.

`MOCK_ENV` must be set or the process panics at startup (`internal/pkg/utils/app.go`); tests are exempt.

## Architecture

### Contracts and codegen
- REST: `api/rest/openapi.yaml` → oapi-codegen (`api/config.yaml`, chi-server + models) → `internal/gen/restapi`.
- gRPC: `api/proto/mock/v1/*.proto` (package `mock.v1`) → protoc → `internal/gen/mockv1`.
- Directives are in `generate.go` (repo root). Never hand-edit `internal/gen/`. The generated code is committed.
- The two contracts are maintained by hand and must stay equivalent: when adding an operation, add it to both specs.

### One port, two protocols
`internal/app/http.go` `Controller.ServeHTTP` routes HTTP/2 requests with `Content-Type: application/grpc` to `grpc.Server.ServeHTTP`. Everything else goes to the chi router: `/api/v1/*` → REST, plus `/healthz` and `/readyz`. `cmd/server/main.go` enables cleartext HTTP/2 (h2c) through `http.Server.Protocols`. In k8s, Traefik reaches the pod over h2c because of the `serversscheme: h2c` annotation in `k8s/service.yaml`. gRPC requests skip the chi middleware (CORS, body limit, logging); `internal/app/grpc/api/api.go` has its own interceptors.

### Layers
```
cmd/server/main.go        → env config (envconfig), zap, store, service, HTTP server, worker goroutine, graceful shutdown on SIGTERM
internal/app/rest/api/    → api.<Operation>.go per endpoint; implements restapi.ServerInterface
internal/app/grpc/api/    → api.<Method>.go per RPC; TodoServer / UtilityServer (two types because both Unimplemented*Server embeds define colliding methods)
internal/pkg/service/     → service.<Operation>.go; shared business logic and validation for every transport
internal/pkg/memstore/    → memstore.<Operation>.go; session-scoped in-memory store
internal/pkg/model/       → typed ids/enums (TodoId, SessionId, TodoPriority) with Parse* constructors
internal/pkg/convert/     → convert.<Entity>To{RestApi,GrpcApi}.go
internal/pkg/errx/        → sentinel errors; each transport maps them (rest: api.errors.go writeError, grpc: api.errors.go toStatus)
internal/pkg/worker/      → Agent loop run as a goroutine (not a separate deployment); SessionCleanupAgent expires idle sessions
```
Handler pattern: parse the session (`X-Session-Id` header or `x-session-id` metadata, defaulting to `public`) → parse ids with `model.Parse*` → call `service.*` → convert. Errors are returned as wrapped `errx` sentinels (`fmt.Errorf("%w: ...", errx.ErrInvalidArgument)`), never as transport status codes, except for gRPC `Status` RPC, which fails on purpose.

### Store semantics
- Each session is created and seeded (`service.SampleTodos`) on first touch, and every access refreshes `lastSeen`.
- `MOCK_MAX_SESSIONS` caps private sessions only; `public` is exempt. `MOCK_MAX_TODOS_PER_SESSION` caps each list.
- Every value crossing the store boundary is `Clone()`d.
- `UpdateTodo` takes a mutate func so read, modify and write happen under one lock. The service does validation inside that func.
- Partial updates use `utils.Optional[T]`. REST PATCH distinguishes `"due_at": null` (clear) from an absent field in `patchDueAt`. gRPC uses a `FieldMask`, and an empty mask means a full replace.

### Lint conventions
- `exhaustruct` is on: struct literals set every field. Generated `mockv1` and protobuf well-known types are excluded. Use `//nolint:exhaustruct // reason` when a partial literal is correct.
- Import groups (gci): stdlib → third-party → `github.com/chapar-rest/mock-server`.

## Deploy
Single replica on the k3s cluster, namespace `mock-server`, host `mocks.chapar.rest` (`k8s/`, setup in `k8s/README.md`). `.github/workflows/deploy.yml` runs `task publish` on every push to `main`. It skips docs-only changes, `[skip deploy]` commits and PRs labeled `skip-deploy`. Needs repo secrets `KUBECONFIG` and `PAT`. Image tag is `v1.0.0-<short sha>`. Traefik middlewares in `k8s/middleware.yaml` rate-limit and cap concurrent requests per client, keyed on `CF-Connecting-IP` because traffic arrives through Cloudflare.
