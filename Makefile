HAS_BREW := $(shell command -v brew 2>/dev/null)

BIN := $(PWD)/bin

all: build

build:
	GOBIN=$(BIN) go install ./...

install-protobuf:
ifdef HAS_BREW
	brew install protobuf
else
	$(error "No supported package manager found")
endif

check-protobuf:
	@scripts/check-installed.sh protoc

DEV_TOOLS := \
	gopls 		golang.org/x/tools/gopls \
	goimports 	golang.org/x/tools/cmd/goimports \
	gofumpt 	mvdan.cc/gofumpt \
	golines 	github.com/segmentio/golines \
	golangci-lint 	github.com/golangci/golangci-lint/cmd/golangci-lint

check-tools: check-protobuf
	@echo $(DEV_TOOLS) | xargs -n2 scripts/check-installed.sh

install-tools:
	@echo $(DEV_TOOLS) | xargs -n2 scripts/go-install.sh

go.mod: tools/tools.go
	go mod tidy
	touch go.mod

mod-download:
	go mod download $$(go list -m -f '{{ .Path }}' all)

build-tools: go.mod mod-download
	GOBIN=$(BIN) go install $$(go list -f '{{join .Imports " "}}' tools/tools.go)

# TODO: remove if not needed
# proto: export GOPATH := $(BIN):$(GOPATH)
proto: export PATH := $(BIN):$(PATH)
proto:
	protoc -I api/ --go_out=. --go-grpc_out=. api/service.proto

lint:
	golangci-lint run

fmt:
	golines --base-formatter=gofumpt --ignore-generated --list-files ./...

config-schema: build
	$(BIN)/claug config-schema > config.schema.json
