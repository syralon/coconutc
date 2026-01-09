.PHONY: grpc
# generate grpc code
grpc:
	protoc -I . --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ./{{.ProtoPath}}/*.proto

.PHONY: gateway
# generate grpc-gateway code
gateway:
	protoc -I . --go_out=paths=source_relative:. --grpc-gateway_out=paths=source_relative:. ./{{.ProtoPath}}/*.proto

.PHONY: openapi
# generate openapi documents
openapi:
	protoc -I . --openapi_out=. ./{{.ProtoPath}}/*.proto

.PHONY: proto
# generate all proto code
proto:
	make grpc gateway openapi

.PHONY: generate
generate:
	make proto
	go generate ./...

.PHONY: quickstart
# quick start
quickstart:
	go mod tidy && go generate ./...
	go run ./cmd/{{.Module|basepath}}

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help