BINARY     := devmem
GUI_BINARY := devmem-gui
PKG        := ./...
GOBIN      ?= $(shell go env GOPATH)/bin

.DEFAULT_GOAL := build

.PHONY: build
build: ## Compile the devmem binary
	go build -o $(BINARY) .

.PHONY: install
install: ## Install devmem into $GOBIN
	go install .

.PHONY: run
run: ## Build and run (use ARGS="scan ~/workspaces")
	go run . $(ARGS)

.PHONY: build-gui
build-gui: ## Compile the devmem-gui desktop binary (requires CGo + GL/X11 headers)
	CGO_ENABLED=1 go build -o $(GUI_BINARY) ./cmd/devmem-gui

.PHONY: run-gui
run-gui: ## Build and run the desktop GUI
	CGO_ENABLED=1 go run ./cmd/devmem-gui

.PHONY: test
test: ## Run the test suite
	go test $(PKG)

.PHONY: cover
cover: ## Run tests with a coverage summary
	go test -cover $(PKG)

.PHONY: fmt
fmt: ## Format all Go source
	gofmt -w .

.PHONY: vet
vet: ## Run go vet
	go vet $(PKG)

.PHONY: tidy
tidy: ## Tidy module dependencies
	go mod tidy

.PHONY: check
check: fmt vet test ## Format, vet and test

.PHONY: clean
clean: ## Remove build artifacts and local databases
	rm -f $(BINARY) $(GUI_BINARY) *.db *.db-journal *.db-wal *.db-shm

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'
