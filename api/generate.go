package api

//go:generate protoc -I . --go-grpc_out=paths=source_relative:./petstorev1 --go_out=paths=source_relative:./petstorev1 ./petstore.proto
//go:generate oapi-codegen --config oapi_config.yaml petstore.yaml
