SRC=$(shell find . -name "*.go")

.PHONY: fmt lint test deps build

default: all

all: fmt lint test build

build: deps
	 go build -o bin/tmplrun ./cmd/tmplrun/

fmt:
	$(info * [checking formatting] **************************************)
	@test -z $(shell gofmt -l $(SRC)) || (gofmt -d $(SRC); exit 1)

lint:
	$(info * [running lint tools] ***************************************)
	golangci-lint run -v

test: deps
	$(info * [running tests] ********************************************)
	go test -v $(shell go list ./... | grep -v /examples$)

deps:
	$(info * {downloading dependencies} *********************************)
	go get -v ./...

