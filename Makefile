VERSION ?= $(shell git describe --tags --always --dirty || echo "unknown")
LDFLAGS := -s -w -X main.appVersion=$(VERSION)
OUTPUT := build/mcp-liner

.PHONY: all build clean

all: build

build:
	@mkdir -p build
	go build -trimpath -ldflags "$(LDFLAGS)" -o $(OUTPUT) ./cmd/mcp-liner

clean:
	rm -rf build dist mcp_liner/*.so mcp_liner/*.dll mcp_liner/*.h
