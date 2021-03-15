#!/bin/bash

if [[ -z ${TEST_DIR} ]]
then
  TEST_DIR="./..."
fi

go test ${TEST_RUN_ARGS}  -v -race -cover \
  -coverprofile cover.out ${TEST_DIR} -timeout 30s | \
  sed "/PASS/s//$(printf "\033[32mPASS\033[0m")/" | \
  sed "/FAIL/s//$(printf "\033[31mFAIL\033[0m")/"

exit ${PIPESTATUS[0]}
