.PHONY: all build bench docker-alpine docker-libimagequant help

all: lint build ## test, lint and build application

lint: ## Lint the project
	golangci-lint --timeout 300s run ./...

build: ## Build cmd
	cd cmd && CGO_ENABLED=1 go build -o go-imagequant .

bench: ## Run bench
	cd imagequant && go test -mod vendor -bench . -run=^$

docker-alpine: ## create cmd as docker image
	docker buildx build -f docker/alpine/Dockerfile --tag go-imagequant:latest .

docker-libimagequant: ## create feed alpine image with libimagequant already installed
	docker buildx build -f docker/alpine-libimagequant/Dockerfile --tag feed-alpine:1.16.12-alpine3.15 .

help: ## Print all possible targets
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {gsub("\\\\n",sprintf("\n%22c",""), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)