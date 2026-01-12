# 🥥Coconutc

A code generate framework based on [ent](https://entgo.io)

## Quickstart:
- `go install github.com/syralon/coconutc/cmd/coconut@latest`
- `coconut new github.com/example/example`
- `cd example`
- `make quickstart`

There is a grpc server listening on '0.0.0.0:9000' and an http server listening on '0.0.0.0:8000'. 
There are 5 apis for the model `example`: 
```
example:
  - id: int
  - created_at: time
  - updated_at: time
  - name: string
  - status: int
  
GET    /v1/example
POST   /v1/example
GET    /v1/example/{id}
PUT    /v1/example/{id}
DELETE /v1/example/{id}
```

## Requirements

- [Go](https://go.dev/dl/)
- [atlas](https://atlasgo.io/)
- [protoc](https://github.com/protocolbuffers/protobuf/releases)
- protoc-go: `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
- protoc-gen-go-grpc: `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`
- protoc-gen-grpc-gateway: `go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@latest`
- protoc-gen-openapi: `go install github.com/googleapis/gnostic/apps/protoc-gen-openapi@latest`

## Usage:
- Run `coconut add Model1 Model2 'Model3'(name:string,status:int)` add new model. Or modify `ent/schema` directly.
- Run `coconut generate`, the services code will be auto generated. 
- Run `go mod tidy && go generate ./... && go run ./cmd/<YOUR_MODULE_NAME>` to start server

## Commands
- `new`
- `add`
- `proto`
- `service`
- `generate`

## Annotations