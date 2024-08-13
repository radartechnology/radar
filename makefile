.DEFAULT_GOAL := help

.PHONY: help
help: ## display command overview
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[35m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: check
check: ## run quality tools for the server
	cd server && sqlc generate
	cd server && go fmt ./...
	cd server && perfsprint --fix ./... # for some reason golangci-lint does not fix all issues
	cd server && wsl --fix ./...
	cd server && golangci-lint run
