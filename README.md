# 🥥Coconutc

A code generate framework based on [ent](https://entgo.io)

## Quickstart:
- `go install github.com/syralon/coconutc/cmd/coconut@latest`
- `coconut new github.com/example/example`
- `cd example`
- `make quickstart`

An example server will start at '0.0.0.0:8000' with http protocol and '0.0.0.0:9000' with grpc protocol.
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

#### **entproto.Field()**
- `entproto.WithFieldImmutable(immutable bool)` 
    Is this field can be modified.
- `entproto.WithFieldSettable(settable bool)`
    Is this field can be set by separate api.
- `entproto.WithFieldFilterable(filterable bool)`
    Is this field can be used in filter arguments.
- `entproto.WithFieldType(fieldType field.Type, repeated ...bool)`
    Reset this filed's type.
#### **entproto.API()**

- `entproto.WithAPIPattern(pattern string) func(a *apiAnnotation)`
    The generated api server pattern prefix, such as '/v1/api'. The api service will not generate if APIPattern is not specified.
- `entproto.WithAPIMethods(methods ...APIMethod) func(a *apiAnnotation)`
    Specify the api methods(GET|CREATE|LIST|UPDATE|DELETE). The all methods will be generated in default.
- `entproto.WithAPIDisableEdge(disable bool) func(a *apiAnnotation)`
    Disable related edges. The edge will be treated as nested data in default.
- `entproto.WithPaginatorStyle(style PaginatorStyle) func(a *apiAnnotation)` 
    Specify the paginator style(ClassicalPaginator|InfinitePaginator) when list schemas. The ClassicalPaginator will be used in default.
The ClassicalPaginator use 'page' and 'pageSize' to separate data list 
#### Example

```go
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/syralon/coconutc/pkg/annotation/entproto"
)

// Application holds the schema definition for the Application entity.
type Application struct {
	ent.Schema
}

// Annotations of the Application.
func (Application) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.API(
			entproto.WithAPIPattern("/v1"),
			entproto.WithPaginatorStyle(entproto.ClassicalPaginator),
		),
	}
}

// Fields of the Application.
func (Application) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("organization_id"),
		field.String("name").Annotations(entproto.Field(entproto.WithFieldSettable(true))),
		field.String("description").Annotations(entproto.Field(entproto.WithFieldSettable(true))),
		field.String("icon").Annotations(entproto.Field(entproto.WithFieldSettable(true))),
		field.String("public_key").Annotations(entproto.Field(entproto.WithFieldSettable(true))),
		field.Strings("white_list").Annotations(
			entproto.Field(entproto.WithFieldSettable(true), entproto.WithFieldType(field.TypeString, true))),
		field.Strings("block_list").Annotations(
			entproto.Field(entproto.WithFieldSettable(true), entproto.WithFieldType(field.TypeString, true))),
		field.Int32("status").Annotations(entproto.Field(entproto.WithFieldSettable(true))),
	}
}
```