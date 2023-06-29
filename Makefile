SHELL=/usr/bin/env bash
GO_BUILD_IMAGE?=golang:1.20
VERSION=$(shell git describe --always --tag --dirty)
COMMIT=$(shell git rev-parse --short HEAD)

.PHONY: all
all: build

.PHONY: build
build:
	go build -ldflags="-X 'main.Commit=$(COMMIT)' -X main.Version=$(VERSION)"  -o edge-vertex

.PHONE: clean
clean:
	rm -f edge-vertex

install:
	install -C -m 0755 edge-vertex /usr/local/bin