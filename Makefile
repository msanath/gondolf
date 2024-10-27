GO=go

BUILD_DIR=bin

.PHONY: prepare
prepare:
	mkdir -p $(BUILD_DIR)

.PHONY: build
build: clean prepare tidy gofmt
	go build -o $(BUILD_DIR)/ledger-builder ./cmd/ledgerbuilder

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: gofmt
gofmt:
	go fmt ./...

.PHONY: clean
clean:
	rm -rf bin 2>/dev/null || true
	rm go.sum 2>/dev/null || true
