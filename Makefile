export GO111MODULE=on

## Install dependencies
.PHONY: deps
deps:
	@go get -v -d

## Run tests
.PHONY: test
test:
	@go test ./...

## Lint
.PHONY: lint
lint:
	@go vet ./...
	@golint --set_exit_status $(go list ./... | grep -v /vendor/)
