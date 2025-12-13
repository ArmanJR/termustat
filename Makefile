COMPOSE          ?= docker compose
COMPOSE_FILES    := -f docker-compose.yml -f docker-compose.dev.yml

compose          = $(COMPOSE) $(COMPOSE_FILES)

.PHONY: help up build stop down restart logs lint lint-api lint-frontend

help:           ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?##' $(MAKEFILE_LIST) | \
	  awk 'BEGIN {FS = ":.*?##"}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'

up:             ## Start the stack (attached)
	$(call compose) up

build:          ## Build / rebuild images only
	$(call compose) build

stop:           ## Stop running containers without removing them
	$(call compose) stop

down:           ## Stop containers and remove containers, networks & volumes
	$(call compose) down

restart:        ## down + up
	$(MAKE) down
	$(MAKE) up

logs:           ## Follow logs for all services
	$(call compose) logs -f

lint:           ## Run all linters (API + Frontend)
	$(MAKE) lint-api
	$(MAKE) lint-frontend

lint-api:       ## Lint Go code (API + Engine)
	golangci-lint run

lint-frontend:  ## Lint Frontend code
	cd frontend && npm run lint
