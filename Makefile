.DEFAULT_GOAL := help

.PHONY: help githooks
help: ## Show this help
	@echo "\033[36mUsage:\033[0m"
	@echo "make TASK"
	@echo
	@echo "\033[36mTasks:\033[0m"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[33m%-30s\033[0m \r\033[32C%s\n", $$1, $$2}'

githooks: ## Install git hooks
	chmod +x githooks/* && cp -rf githooks/* .git/hooks/
