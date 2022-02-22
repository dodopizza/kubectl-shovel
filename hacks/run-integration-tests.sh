#!/bin/bash
set -o errexit

echo "Running tests..."
go test -v \
  --tags=integration \
  -timeout 300s \
  ./test/integration/... |
  sed "/PASS/s//$(printf "\033[32mPASS\033[0m")/" |
  sed "/FAIL/s//$(printf "\033[31mFAIL\033[0m")/"

exit ${PIPESTATUS[0]}
