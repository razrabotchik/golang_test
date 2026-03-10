.DEFAULT_GOAL := help
RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(RUN_ARGS):;@:)

run: ## Run
	go run main.go

test: ## Run tests
	go test ./internal/...

gen: ## Generate code
	go generate ./...

bench: ## Run benchmark tests
	go test -bench=. ./internal/...

.PHONY: lint
lint: ## Run golangci-lint
	@echo "Running linter..."
	@golangci-lint run ./...