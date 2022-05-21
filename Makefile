OUTPUT_DIR      = build
TMP_DIR        := .tmp
RELEASE_VER    := $(shell git rev-parse --short HEAD)
NAME            = default
COVERMODE       = atomic

TEST_PACKAGES      := $(shell go list ./... | grep -v vendor | grep -v fakes)

.PHONY: help
.DEFAULT_GOAL := help

installdeps: ## Install needed dependencies for various middlewares
	go get -t -v ./...

staticcheck:
	staticcheck ./...

generate: ## Run generate for non-vendor packages only
	go list ./... | xargs go generate

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
