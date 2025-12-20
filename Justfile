#!/usr/bin/env -S just --justfile

default:
  @just --list

build:
    time ./build.sh

build-all: tidy
    ./build.sh all

test:
    go test -short ./...

test-sh:
    bin/sh script/test.sh

tidy:
    go mod tidy
    go fmt ./...
    go vet ./...

install: build test
    time CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o "$(go env GOPATH)/bin/shell" -ldflags="-w -extldflags '-static' ${CLI_FLAGS:-}" ./cmd

update:
    go get -u ./...

clean-cache:
    go clean -modcache
