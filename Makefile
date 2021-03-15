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
	./hacks/run-integration-tests.sh

.PHONY: tidy
tidy:
	go mod tidy -v

.PHONY: lint
lint:
	golangci-lint run --tests=false
	golangci-lint run --disable-all -E golint,goimports,misspell
