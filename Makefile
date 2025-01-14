GO=go

BUILD_DIR=bin

.PHONY: prepare
prepare:
	mkdir -p $(BUILD_DIR)

.PHONY: build
build: clean prepare tidy gofmt
	go build -o $(BUILD_DIR)/ledger-builder ./cmd/ledgerbuilder && \
	go build -o $(BUILD_DIR)/cligen ./cmd/cligen && \
	go build -o $(BUILD_DIR)/simplesqlorm-gen ./cmd/simplesqlorm

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: gofmt
gofmt:
	go fmt ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: generate
generate:
	go generate ./...

.PHONY: clean
clean:
	rm -rf bin 2>/dev/null || true
	rm go.sum 2>/dev/null || true
