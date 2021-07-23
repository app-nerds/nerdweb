.DEFAULT_GOAL := help
.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test: ## Run all unit tests
	go test ./...

coverage: ## Run tests and display a code-coverage report
	go test ./... -coverprofile=coverageprofile.out
	go tool cover -html=coverageprofile.out

