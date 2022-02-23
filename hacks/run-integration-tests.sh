#!/bin/bash
set -ox errexit

ci=${CI:-""}
flags=""

if [ -z "$ci" ]; then
  flags="-parallel 1"
else
  flags="-parallel 4"
fi

echo "Running tests..."
go test -v $flags \
  --tags=integration \
  -timeout 300s \
  ./test/integration/... |
  sed "/PASS/s//$(printf "\033[32mPASS\033[0m")/" |
  sed "/FAIL/s//$(printf "\033[31mFAIL\033[0m")/"

exit ${PIPESTATUS[0]}
