##@ Help

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Protobuf

.PHONY: gen
gen: ## Compile protocol buffers and generate go code.
	protoc -I=. grpc/pb/server.proto --go_out=. --go-grpc_out=.

.PHONY: list
list: ## List gRPC services (make sure gRPC server is listening on port 8546).
	grpcurl -plaintext 127.0.0.1:8546 list v1.System

##@ Lint

.PHONY: tidy
tidy: ## Run go mod tidy against code.
	go mod tidy

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

# shadow reports shadowed variables
# https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/shadow
# go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest
.PHONY: vet
vet: ## Run go vet and shadow against code.
	go vet ./...
	shadow ./...

# golangci-lint runs gofmt, govet, staticcheck and other linters
# https://golangci-lint.run/usage/install/#local-installation
.PHONY: golangci-lint
golangci-lint: ## Run golangci-lint against code.
	golangci-lint run

.PHONY: lint
lint: tidy vet golangci-lint ## Run all of these tools against code.

##@ Build

.PHONY: build
build: gen ## Build binary.
	go build -o bin/mock-server main.go

##@ Test

.PHONY: test
test: ## Send a few gRPC/HTTP requests to the mock server.
	./scripts/test.sh

.PHONY: clean
clean: ## Clean the proof directory.
	rm out/*
