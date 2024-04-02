.PHONY: all build bench test docker-cmd artifacts docker-lib-ubuntu-amd64 help

all: lint build ## test, lint and build application
artifacts: docker-lib-ubuntu-amd64 

lint: ## Lint the project
	golangci-lint --timeout 300s run ./...

build: ## Build cmd
	cd cmd && CGO_ENABLED=1 go build -o go-imagequant .

bench: ## Run bench
	cd imagequant && go test -mod vendor -bench . -benchmem -run=^$

test: ## Run tests
	go test -v -mod vendor ./...

docker-cmd-amazon2: ## create cmd as docker amazon based image
	docker buildx build --platform linux/amd64 -f docker/amazonlinux/Dockerfile --tag go-imagequant:latest --load .

docker-cmd-alpine: ## create cmd as docker amazon based image
	docker buildx build -f docker/alpine/Dockerfile --tag go-imagequant:latest --load .


docker-lib-ubuntu-amd64: ## create ubuntu lib artifacts
	echo "creating ubuntu  lib artifacts ..."
	rm -rf ./lib/ubuntu/22.04/* # cleanup old stuff
	docker buildx build --platform linux/amd64 -f docker/create-ubuntu-artifacts/Dockerfile --output type=local,dest=. .

help: ## Print all possible targets
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {gsub("\\\\n",sprintf("\n%22c",""), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)