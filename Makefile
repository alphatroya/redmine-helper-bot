GO_BIN := $(GOPATH)/bin
GOIMPORTS := $(GO_BIN)/goimports
GOLANGCI := $(GO_BIN)/golangci-lint

## build: Build an application
.PHONY: build
build: fmt
	go build -v -o main

## install: Install application
.PHONY: install
install:
	go install -v

## run: Run application
.PHONY: run
run: fmt
	go run .

## test: Launch unit tests
.PHONY: test
test:
	go test ./...

## coverage: Launch unit tests
.PHONY: coverage
coverage:
	@go test -v -coverpkg=./... -coverprofile=profile.cov ./... > /dev/null
	@go tool cover -func profile.cov | tail -n 1
	@rm -fr profile.cov

## clean: Cleanup build artefacts
.PHONY:
clean:
	go clean

## generate: Regenerate all required files
generate:
	go generate

## tidy: Cleanup go.sum and go.mod files
.PHONY: tidy
tidy:
	go mod tidy

## lint: Launch project linters
.PHONY: lint
lint: $(GOLANGCI)
	$(GOLANGCI) run

## fmt: Reformat source code
.PHONY: fmt
fmt: $(GOIMPORTS)
	$(GOIMPORTS) -w -l .

$(GOIMPORTS):
	go install golang.org/x/tools/cmd/goimports@master

$(GOLANGCI):
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.40.1

## help: Prints help message
.PHONY: help
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /' | sort
