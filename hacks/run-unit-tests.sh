#!/usr/bin/env bash
set -eu

go test \
  -v \
  -race \
  -cover \
  -coverprofile cover.out \
  -timeout 30s \
  ./... |
  sed "/PASS/s//$(printf "\033[32mPASS\033[0m")/" | \
  sed "/FAIL/s//$(printf "\033[31mFAIL\033[0m")/"

exit ${PIPESTATUS[0]}
