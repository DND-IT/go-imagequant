.PHONY: all build bench test help

all: lint lib build ## test, lint and build application

lib: ## build the submodule
	cd libimagequant/imagequant-sys && make

lint: ## Lint the project
	golangci-lint --timeout 300s run ./...

build: ## Build cmd
	cd cmd && CGO_ENABLED=1 go build -o go-imagequant .

bench: ## Run bench
	cd imagequant && go test -mod vendor -bench . -benchmem -run=^$

test: ## Run tests
	go test -v -mod vendor ./...

help: ## Print all possible targets
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {gsub("\\\\n",sprintf("\n%22c",""), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

clean: ## Clean the project
	rm -f cmd/go-imagequant
	go clean
	cd libimagequant/imagequant-sys && make clean