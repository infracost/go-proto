default: help

.PHONY: help
help: ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST)  | fgrep -v fgrep | sed -e 's/:.*##/:##/' | awk -F':##' '{printf "%-12s %s\n",$$1, $$2}'

.PHONY: lint_install
lint_install: ## Install golangci-lint
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

.PHONY: lint
lint: lint_install ## Run go-lint on the project
	golangci-lint run ./...

.PHONY: test
test: ## Run tests
	go test -v -race ./...

.PHONY: coverage
coverage: ## Run tests with coverage
	go test -coverprofile=coverage.out ./...
