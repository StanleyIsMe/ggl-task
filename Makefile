.PHONY: help test test-race test-leak bench bench-compare lint sec-scan build

help: ## show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

PROJECT_NAME?=ggl-task
SHELL = /bin/bash

########
# test #
########

test: test-race test-leak ## launch all tests

test-race: ## launch all tests with race detection
	go test ./... -cover -race

test-leak: ## launch all tests with leak detection (if possible)
	go test ./... -leak

test-coverage-report:
	go test -v  ./... -cover -race -covermode=atomic -coverprofile=./coverage.out
	go tool cover -html=coverage.out

########
# lint #
########

lint: ## lints the entire codebase
	@golangci-lint run ./... --config=./.golangci.yaml

#######
# sec #
#######

sec-scan: trivy-scan vuln-scan ## scan for security and vulnerability issues

trivy-scan: ## scan for sec issues with trivy (trivy binary needed)
	trivy fs --exit-code 1 --no-progress --severity CRITICAL ./

vuln-scan: ## scan for vulnerability issues with govulncheck (govulncheck binary needed)
	govulncheck ./...

###########
# swagger #
###########

swagger-gen: ## generate swagger docs
	swag init -d ./cmd/api,./internal --parseDependency


#########
# build #
#########

build: ## build docker image
	docker buildx build \
	-f Dockerfile \
	-t $(PROJECT_NAME) \
	--platform linux/arm64 \
	--build-arg GO_VERSION=1.23.3 \
	--build-arg GO_GOOS=linux \
	--build-arg GO_GOARCH=arm64 \
	--build-arg GLOBAL_VAR_PKG=server \
	--build-arg LAST_MAIN_COMMIT_HASH=$(shell git rev-parse HEAD) \
	--build-arg LAST_MAIN_COMMIT_TIME=$(shell git log main -n1 --format='%cd' --date='iso-strict') \
	--progress=plain \
	--load \
	./

#########
# deploy #
#########

up: ## run docker image
	docker run -d --name $(PROJECT_NAME) 

down: ## stop docker image
	docker stop $(PROJECT_NAME)

