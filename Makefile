BIN_DIR     := $(CURDIR)/bin
BIN_NAME   ?= passenger-tracker

GOPATH = $(shell go env GOPATH)
GOOS=linux
GOARCH=amd64
GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn

# go option
TAGS       :=
TESTS      := .
TESTFLAGS  :=
LDFLAGS    := -w -s
GOFLAGS    :=

all: clean lint test build

# ------------------------------------------------------------------------------
#  make code go ready

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)/$(BIN_NAME)

.PHONY: lint fmt vet mod test test-coverage

lint: fmt vet
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(GOPATH)/bin v1.44.0
	golangci-lint run -v ./internal/... ./pkg/...

test: lint
	go test -count=1 ./... -cover

test-coverage: lint
	go test -cover ./internal/... ./pkg/...

fmt:
	go fmt ./...
	goimports -w ./

vet:
	go vet ./...

mod:
	go mod tidy
	go mod verify

clean-gen:
	grep -rl "THIS FILE IS GENERATED"  * | egrep -v "vendor|Makefile|cmd/generators" | xargs rm
	find . | grep -v vendor | egrep "\.pb\.go|\.pb\.validate\.go|\.gen\.go|generated\.deepcopy\.go|\.pb\.gw\.go" | grep -v mocks | xargs rm -f

generate-proto:
		protoc -I ./proto-vendor -I ./api/proto \
          --go_out=./api/pb --go_opt=paths=source_relative \
          --go-grpc_out=./api/pb --go-grpc_opt=paths=source_relative --go-grpc_opt=require_unimplemented_servers=false \
          --grpc-gateway_out=./api/pb --grpc-gateway_opt=paths=source_relative \
        ./api/proto/v1alpha1/*/*.proto

# ------------------------------------------------------------------------------
#  cli build and install

.PHONY: build install

build: clean fmt vet mod
	GO111MODULE=on CGO_ENABLED=0 go build $(GOFLAGS) -trimpath \
	-tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o '$(BIN_DIR)'/$(BIN_NAME) main.go

install: build
	@install "$(BIN_DIR)/$(BIN_NAME)" "$(GOBIN)/$(BIN_NAME)"