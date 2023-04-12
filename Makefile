###
# Params.
###

PROJECT_FULL_NAME := pubsub

HAS_GODOC := $(shell command -v godoc;)
HAS_GOLANGCI := $(shell command -v golangci-lint;)
HAS_DOCKER_COMPOSE := $(shell command -v docker-compose;)

default: ci

###
# Validate variables.
###

ENV_FILE ?= testing.env

ifneq ($(filter development.env integration.env testing.env,$(ENV_FILE)),)
else
$(error ENV_FILE must be either "development.env" or "integration.env" or "testing.env")
endif

###
# Entries.
###

ci: lint test coverage
ci-integration: lint test-integration coverage
ci-integration-local: lint test-integration-local coverage

coverage:
	@go tool cover -func=coverage.out && echo "Coverage OK"

deps:
	@go install golang.org/x/tools/cmd/godoc@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.2

infra-start:
ifndef HAS_DOCKER_COMPOSE
	@echo "Could not find docker-compose, please install it"
endif
	@echo "Starting required infrastructure"
	@docker-compose -f resources/docker-compose.yml up --remove-orphans -d nats

infra-stop:
ifndef HAS_DOCKER_COMPOSE
	@echo "Could not find docker-compose, please install it"
endif
	@echo "Stopping required infrastructure"
	@docker-compose -f resources/docker-compose.yml down --remove-orphans

doc:
ifndef HAS_GODOC
	@echo "Could not find godoc, installing it and any other missing tool(s)"
	@make deps
endif
	@echo "Open localhost:6060/pkg/github.com/thalesfsp/$(PROJECT_FULL_NAME)/ in your browser\n"
	@godoc -http :6060

lint:
ifndef HAS_GOLANGCI
	@echo "Could not find golangci-list, installing it and any other missing tool(s)"
	@make deps
endif
	@golangci-lint run -v -c .golangci.yml && echo "Lint OK"

test:
	@ENVIRONMENT="testing" configurer l d -f testing.env -- go test -timeout 30s -short -v -race -cover \
	-coverprofile=coverage.out ./... && echo "Test OK"

test-integration:
	@ENVIRONMENT="integration" configurer l d -f integration.env -- go test -timeout 120s -v -race \
	-cover -coverprofile=coverage.out ./... && echo "Integration test OK"

test-integration-local:
	@ENVIRONMENT="integration" configurer l d -f integration-local.env -- go test -timeout 120s -v -race \
	-cover -coverprofile=coverage.out ./... && echo "Integration test OK"

.PHONY: ci \
	ci-integration \
	ci-integration-local \
	coverage \
	deps \
	infra-start \
	infra-stop \
	doc \
	lint \
	test \
	test-integration \
	test-integration-local
