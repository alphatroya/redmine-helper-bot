GO_BIN := $(GOPATH)/bin
GOIMPORTS := $(GO_BIN)/goimports
GOLINT := $(GO_BIN)/golangci-lint

all: install

install: bootstrap fmt
	go install -v

test: bootstrap fmt
	go test ./... -v

coverage: bootstrap fmt
	go test -v -race -coverprofile=coverage.txt -covermode=atomic

bootstrap:
	go get ./...

lint: $(GOLINT)
	golangci-lint run

fmt: $(GOIMPORTS)
	goimports -w *.go

$(GOIMPORTS):
	go get -u golang.org/x/tools/cmd/goimports

$(GOLINT):
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: install test fmt lint bootstrap
