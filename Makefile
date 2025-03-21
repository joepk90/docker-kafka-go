.DEFAULT_GOAL := help

DOCKER_REGISTRY?=jparkkennaby
EXECUTABLE_NAME=docker-kafka-go
GIT_HASH := $(shell git rev-parse HEAD)

help: ## Display help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# dep: ## Fetch dependencies (including private repos)
# 	GOPRIVATE=github.com/jparkkennaby/* GO111MODULE=on go mod download

generate-protos:
	buf generate

dep: ## Fetch dependencies
	go mod download

build: dep ## Build executable
	CGO_ENABLED=0 \
	GOOS=linux \
    GOARCH=amd64 \
	go build -o $(EXECUTABLE_NAME) .

test: ## Runs the tests
	go test -v -count=1 -cover ./...

format: ## Formats the code
	goimports -w .

lint: local/lint-install local/lint-run
local/lint-install: ## Install the linter locally
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.60.3

local/lint-run: ## Run the linter locally
	bin/golangci-lint run

local/docker-build: ## Builds the docker image
	docker build -t $(DOCKER_REGISTRY)/$(EXECUTABLE_NAME):local --build-arg GITHUB_TOKEN=$(GITHUB_TOKEN) -f DockerfileLocalArm64 .

.PHONY: mocks
mocks: ## Generates mock files for testing
	go generate ./...

generate-sql: ## Generates all queries defined in stores/queries
	go run github.com/sqlc-dev/sqlc/cmd/sqlc generate


# DOCKER COMPOSE COMMANDS (LOCAL KAFKA INSTANCE)
# TODO how can we remove the need for the .env file?


dev-start:
	docker-compose --env-file .env.kafka up -d 

dev-stop:
	docker-compose --env-file .env.kafka down -v

dev-watch:
	docker-compose --env-file .env.kafka logs consumer -f

dev-watch-all:
	docker-compose --env-file .env.kafka logs -f

dev-restart:
	docker-compose --env-file .env.kafka restart consumer 

dev-produce-event:
	go run cmd/event-publisher/*.go

