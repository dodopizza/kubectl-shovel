GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all
all: help

.PHONY: cover
cover:
	go tool cover -html ./cover.out

.PHONY: doc
doc:
	./hacks/run-doc-generation.sh

.PHONY: lint
lint:
	golangci-lint run

.PHONY: build-cli
build-cli:
	go build -v -o ./cli/bin/kubectl-shovel ./cli

.PHONY: build-dumper
build-dumper:
	go build -v -o ./dumper/bin/dumper ./dumper

.PHONY: prepare
prepare: tidy lint doc test

.PHONY: setup
setup:
	kind create cluster

.PHONY: test
test: test-unit test-integration

.PHONY: test-unit
test-unit:
	TEST_RUN_ARGS="$(TEST_RUN_ARGS)" TEST_DIR="$(TEST_DIR)" ./hacks/run-unit-tests.sh

.PHONY: test-integration
test-integration:
	./hacks/run-integration-tests.sh amd64

.PHONY: test-integration-arm64
test-integration-arm64:
	./hacks/run-integration-tests.sh arm64

.PHONY: tidy
tidy:
	go mod tidy -v

.PHONY: help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@echo "  ${YELLOW}cover                  ${RESET} Open html coverage report in browser"
	@echo "  ${YELLOW}doc                    ${RESET} Run doc generation"
	@echo "  ${YELLOW}lint                   ${RESET} Run linters via golangci-lint"
	@echo "  ${YELLOW}prepare                ${RESET} Run all available checks and generators"
	@echo "  ${YELLOW}setup                  ${RESET} Setup local environment. Create kind cluster"
	@echo "  ${YELLOW}test                   ${RESET} Run all available tests"
	@echo "  ${YELLOW}test-integration       ${RESET} Run all integration tests (for amd64 arch)"
	@echo "  ${YELLOW}test-integration-arm64 ${RESET} Run all integration tests (for arm64 arch)"
	@echo "  ${YELLOW}test-unit              ${RESET} Run all unit tests"
	@echo "  ${YELLOW}tidy                   ${RESET} Run tidy for go module to remove unused dependencies"
