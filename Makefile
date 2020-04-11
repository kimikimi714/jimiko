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
	@golint --set_exit_status ./...

