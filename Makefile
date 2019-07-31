GO_BIN := $(GOPATH)/bin
GOIMPORTS := $(GO_BIN)/goimports
GOLINT := $(GO_BIN)/golint

all: install

install: bootstrap fmt lint
	go install -v

test: bootstrap fmt
	go test ./... -v

bootstrap:
	go get ./...

lint: $(GOLINT)
	golint *.go

fmt: $(GOIMPORTS)
	goimports -w *.go

$(GOIMPORTS):
	go get -u golang.org/x/tools/cmd/goimports

$(GOLINT):
	go get -u golang.org/x/lint/golint

.PHONY: install test fmt lint bootstrap
