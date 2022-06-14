#!make

PWD   	?= $(shell pwd)
BUILD 	?= $(PWD)/build
VERSION ?= $(shell git describe --tags --abbrev=0)

all: test build

.PHONY: lint
lint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(call print-target)
	@golangci-lint run --config .golangci.yaml

.PHONY: fmt
fmt:
	$(call print-target)
	gofmt -s -w .

.PHONY: tidy
tidy:
	$(call print-target)
	go mod tidy

.PHONY: clean
clean:
	$(call print-target)
	@rm -rf $(BUILD)
	@rm -rf $(PWD)/examples/**/*.pb.go

.PHONY: build
build: build-plugin build-examples

.PHONY: build-plugin
build-plugin: clean
	$(call print-target)
	@go build $(GOFLAGS) -ldflags="-X 'main.version=$(VERSION)'" -o $(BUILD)/protoc-gen-go-grpcmock ./cmd/protoc-gen-go-grpcmock

.PHONY: build-examples
build-examples: build-examples-testify build-examples-pegomock

.PHONY: build-examples-testify
build-examples-testify:
	$(call print-target)
	@cd examples/helloworld; protoc --go_out=testify --go_opt=paths=source_relative --go-grpc_out=testify --go-grpc_opt=paths=source_relative --plugin=$(BUILD)/protoc-gen-go-grpcmock --go-grpcmock_out=framework=testify,import_package=false:testify --go-grpcmock_opt=paths=source_relative helloworld.proto
	@cd examples/routeguide; protoc --go_out=testify --go_opt=paths=source_relative --go-grpc_out=testify --go-grpc_opt=paths=source_relative --plugin=$(BUILD)/protoc-gen-go-grpcmock --go-grpcmock_out=framework=testify,import_package=false:testify --go-grpcmock_opt=paths=source_relative route_guide.proto

.PHONY: build-examples-pegomock
build-examples-pegomock:
	$(call print-target)
	@cd examples/helloworld; protoc --go_out=pegomock --go_opt=paths=source_relative --go-grpc_out=pegomock --go-grpc_opt=paths=source_relative --plugin=$(BUILD)/protoc-gen-go-grpcmock --go-grpcmock_out=framework=pegomock,import_package=false:pegomock --go-grpcmock_opt=paths=source_relative helloworld.proto
	@cd examples/routeguide; protoc --go_out=pegomock --go_opt=paths=source_relative --go-grpc_out=pegomock --go-grpc_opt=paths=source_relative --plugin=$(BUILD)/protoc-gen-go-grpcmock --go-grpcmock_out=framework=pegomock,import_package=false:pegomock --go-grpcmock_opt=paths=source_relative route_guide.proto

.PHONY: test
test:
	$(call print-target)
	go test -race -cover ./... -timeout 60s

define print-target
    @printf "Executing target: \033[36m$@\033[0m\n"
endef