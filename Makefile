GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)
MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_DIR := $(dir $(MAKEFILE_PATH))
ARCH := $(shell uname -m)

.PHONY: all
all: help

.PHONY: prepare
prepare: tidy lint doc build test

.PHONY: cover
cover:
	go tool cover -html ./cover.out

.PHONY: doc
doc: build-cli
	cd cli && HOME="/home/user" ./bin/kubectl-shovel doc && cd -

.PHONY: lint
lint:
	golangci-lint run

.PHONY: tidy
tidy:
	go mod tidy -v

.PHONY: build
build: build-cli build-dumper

.PHONY: build-cli
build-cli:
	go build -v -o ./cli/bin/kubectl-shovel ./cli

.PHONY: build-dumper
build-dumper:
	go build -v -o ./dumper/bin/dumper ./dumper

.PHONY: test
test: test-unit test-integration

.PHONY: test-unit
test-unit:
	go test \
      -v \
      -race \
      -cover \
      -coverprofile cover.out \
      -timeout 30s \
      ./...

.PHONY: test-integration-setup
test-integration-setup:
	kind create cluster --name "kind"

.PHONY: test-integration-prepare
test-integration-prepare:
	FRAMEWORK=$(FRAMEWORK) ./hacks/prepare-integration-tests.sh "$(CURRENT_DIR)" "kind-kind" "$(ARCH)" "$(FRAMEWORK)"

.PHONY: test-integration
test-integration: test-integration-prepare
	FRAMEWORK=$(FRAMEWORK) go test \
      -v \
      -ldflags="-X github.com/dodopizza/kubectl-shovel/test/integration_test.TargetContainerImage=kubectl-shovel/sample-integration-tests:$(FRAMEWORK)" \
      -parallel 1 \
      -timeout 600s \
      --tags=integration \
      ./test/integration/...

.PHONY: help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@echo "  ${YELLOW}prepare                   ${RESET} Run all available checks and generators"
	@echo "  ${YELLOW}cover                     ${RESET} Open html coverage report in browser"
	@echo "  ${YELLOW}doc                       ${RESET} Run doc generation"
	@echo "  ${YELLOW}lint                      ${RESET} Run linters via golangci-lint"
	@echo "  ${YELLOW}tidy                      ${RESET} Run tidy for go module to remove unused dependencies"
	@echo "  ${YELLOW}build                     ${RESET} Build all components"
	@echo "  ${YELLOW}build-cli                 ${RESET} Build cli component of shovel"
	@echo "  ${YELLOW}build-dumper              ${RESET} Build dumper component of shovel"
	@echo "  ${YELLOW}test                      ${RESET} Run all available tests"
	@echo "  ${YELLOW}test-unit                 ${RESET} Run all unit tests"
	@echo "  ${YELLOW}test-integration          ${RESET} Run all integration tests"
	@echo "  ${YELLOW}test-integration-setup    ${RESET} Setup integration tests environment. Create kind cluster"
	@echo "  ${YELLOW}test-integration-prepare  ${RESET} Prepare integration tests (build required images and load to kind cluster)"
