BIN ?= bin/bgo
PKG ?= ./...

.PHONY: all build run test test-watch fmt vet tidy clean

all: build

build:
	@mkdir -p $(dir $(BIN))
	go build -o $(BIN) ./cmd/bgo

run:
	go run ./cmd/bgo $(ARGS)

test:
	go test $(PKG)

test-watch:
	@command -v entr >/dev/null 2>&1 || { echo "entr is required for test-watch"; exit 1; }
	@find . -name '*.go' -not -path './vendor/*' | entr -c go test $(PKG)

fmt:
	gofmt -w .

vet:
	go vet $(PKG)

tidy:
	go mod tidy

clean:
	rm -rf bin
