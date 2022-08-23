MAKEFILE_PATH := $(abspath $(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
PATH := $(MAKEFILE_PATH):$(PATH)

export GOBIN := $(MAKEFILE_PATH)/bin

PATH := $(GOBIN):$(PATH)

PRE_COMMIT := $(shell command -v pre-commit 2> /dev/null)

install-git-hooks:
ifndef PRE_COMMIT
	$(error "pre-commit is not available. See https://pre-commit.com/#installation for installation instructions")
endif
	@pre-commit install

.PHONY: test
test:
	@go test -count=1 -race -v ./...

.PHONY: test100
test100:
	@go test -count=100 ./...

.PHONY: test-coverage
test-coverage:
	@go test -race -cover ./...

.PHONY: test-coverage-visualize
test-coverage-visualize:
	@go test -race -coverprofile=c.out ./... && go tool cover -html=c.out

lint:
	@echo lint
	@golangci-lint run -c golangci.yml

# ==============================================================================
# Docker

local:
	echo "Starting local docker compose"
	docker-compose -f docker-compose.local.yaml up