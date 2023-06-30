.PHONY: all build bench test docker-cmd artifacts docker-lib-ubuntu-amd64 docker-lib-alpine-arm64 help

all: lint build ## test, lint and build application
artifacts: docker-lib-ubuntu-amd64 docker-lib-alpine-arm64

lint: ## Lint the project
	golangci-lint --timeout 300s run ./...

build: ## Build cmd
	cd cmd && CGO_ENABLED=1 go build -o go-imagequant .

bench: ## Run bench
	cd imagequant && go test -mod vendor -bench . -benchmem -run=^$

test: ## Run tests
	go test -v -mod vendor ./...

docker-cmd: ## create cmd as docker alpine based image
	docker buildx build -f docker/amazonlinux/Dockerfile --tag go-imagequant:latest --load .

docker-lib-ubuntu-amd64: ## create ubuntu lib artifacts
	echo "creating ubuntu  lib artifacts ..."
	rm -rf ./lib/ubuntu/22.04/* # cleanup old stuff
	docker buildx build --platform linux/amd64 -f docker/create-ubuntu-artifacts/Dockerfile --output type=local,dest=. .


# alpine build is broken because missing symbol getauxval in alpine libc
docker-lib-alpine-arm64: ## create alpine artifacts
	echo "creating alpine arm64 lib artifacts ..."
	rm -rf lib/alpine/3.18/* # cleanup old stuff
	docker buildx build --platform linux/arm64 -f docker/create-alpine-artifacts/Dockerfile --output type=local,dest=. .




help: ## Print all possible targets
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {gsub("\\\\n",sprintf("\n%22c",""), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)