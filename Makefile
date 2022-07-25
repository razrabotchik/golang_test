.DEFAULT_GOAL := help
RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(RUN_ARGS):;@:)

run: ## Run
	go run main.go

tests: ## Run tests
	go test ./infra/...

bench: ## Run benchmark tests
	go test -bench=. bench_test.go