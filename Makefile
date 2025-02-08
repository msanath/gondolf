GO=go

BUILD_DIR=bin

.PHONY: all
all: clean build test

.PHONY: build
build: prepare tidy gofmt
	go build -o $(BUILD_DIR)/ledger-builder ./cmd/ledgerbuilder && \
	go build -o $(BUILD_DIR)/cligen ./cmd/cligen && \
	go build -o $(BUILD_DIR)/simplesqlormgen ./cmd/simplesqlormgen

.PHONY: clean
clean:
	rm -rf bin 2>/dev/null || true
	rm go.sum 2>/dev/null || true

.PHONY: prepare
prepare:
	mkdir -p $(BUILD_DIR)

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: gofmt
gofmt:
	go fmt ./...

.PHONY: test
test: generate
	go test -v ./...

.PHONY: generate
generate:
	go generate ./...

