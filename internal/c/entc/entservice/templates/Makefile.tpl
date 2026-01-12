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

.PHONY: atlas
# generate database migrations by atlas
atlas:
	@LATEST_FILE=$$(ls -1 migrations/$(MODULE)/*.sql 2>/dev/null | sort | tail -n 1); \
	if [ -z "$$LATEST_FILE" ]; then \
	    VERSION="v0.0.1"; \
	else \
	    LATEST_VERSION=$$(basename $$LATEST_FILE | sed -E 's/.*_v([0-9]+\.[0-9]+\.[0-9]+)\.sql/\1/'); \
	    VERSION=v$$(echo $$LATEST_VERSION | awk -F. '{printf "%d.%d.%d", $$1, $$2, $$3+1}'); \
	fi; \
	atlas migrate diff $$VERSION \
		--dir "file://migrations" \
		--to "ent://ent/schema" \
		--dev-url "docker://mysql/8/ent" \
		--format  "{{"{{"}} sql . \"    \"{{"}}"}}"

.PHONY: generate
generate:
	make proto
	go generate ./...

.PHONY: quickstart
# quick start
quickstart:
	@rm -f example.db
	@atlas migrate diff --dir "file://quickstart" --to "ent://ent/schema" --dev-url "sqlite://example.db"
	@atlas migrate apply --url "sqlite://example.db" --dir "file://quickstart"
	@rm -rf ./quickstart
	@echo "=========================================================="
	@echo " The quickstart server is just for example, DON'T USE IT! "
	@echo "=========================================================="
	go mod tidy && go generate ./...
	go run ./cmd/{{.Module|basepath}}

.PHONY: build
# build
build:
	go build -ldflags "-X {{.Module}}/version.BuildTime=$(shell date '+%Y%m%d%H%M%S') -X {{.Module}}/version.Version=$(shell git rev-parse HEAD)" ./cmd/{{.Module|basepath}}

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