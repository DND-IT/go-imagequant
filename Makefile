.PHONY: all build lint test bench generate docker-alpine help

all: test lint build ## test, lint and build application


build: ## Build cmd
	cd cmd && CGO_ENABLED=1 go build -o go-imagequant .

bench: ## Run bench
	cd imagequant && go test -mod vendor -bench . -run=^$

docker-alpine: ## create cmd as docker image
	docker buildx build -f docker/alpine/Dockerfile --tag go-imagequant:latest .

docker-libimagequant: ## create current alpine image with libimagequant already installed
	docker buildx build -f docker/alpine-libimagequant/Dockerfile --tag feed-alpine:latest .
