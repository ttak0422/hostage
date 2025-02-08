BIN := hostage
CURRENT_REVISION = $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS = "-s -w -X github.com/ttak0422/$(BIN)/cli.revision=$(CURRENT_REVISION)"

.PHONY: all
all: build

.PHONY: deps
deps:
	go mod tidy

.PHONY: build
build: deps
	go build -ldflags $(BUILD_LDFLAGS) -o $(BIN)

.PHONY: test
test: build
	go test -v ./...
