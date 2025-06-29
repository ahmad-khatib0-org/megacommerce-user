mailer-mocks: ## generate mocks for the mailer pkg
	mockery --config internal/mailer/.mockery.yaml

store-mocks: ## generate mocks for the store pkg
	mockery --config internal/store/.mockery.yaml

worker-mocks: ## generate mocks for the worker pkg
	mockery --config internal/worker/.mockery.yaml

.PHONY:
	mailer-mocks
	store-mocks
	worker-mocks
