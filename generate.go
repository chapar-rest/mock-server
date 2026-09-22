package internal

//go:generate oapi-codegen -package restapi --config ./api/config.yaml -o ./internal/gen/restapi/openapi.gen.go ./api/rest/openapi.yaml
//go:generate protoc -I ./api/proto --go_out=. --go_opt=module=github.com/chapar-rest/mock-server --go-grpc_out=. --go-grpc_opt=module=github.com/chapar-rest/mock-server mock/v1/todo.proto mock/v1/utility.proto
