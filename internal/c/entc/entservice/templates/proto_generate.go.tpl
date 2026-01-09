// @file: proto_generate.go

//go:generate protoc -I . --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ./{{.ProtoPath}}/*.proto
//go:generate protoc -I . --go_out=paths=source_relative:. --grpc-gateway_out=paths=source_relative:. ./{{.ProtoPath}}/*.proto
//go:generate protoc -I . --openapi_out=. ./{{.ProtoPath}}/*.proto

package {{ .Module | basepath | toPackage}}