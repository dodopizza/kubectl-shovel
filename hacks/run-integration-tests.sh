#!/usr/bin/env bash
set -eu

[ $# -lt 1 ] && echo "Usage: $(basename $0) <framework|[netcoreapp3.1, net5.0, net6.0]>" && exit 1

framework=$1

if [ "$framework" != "netcoreapp3.1" ] && [ "$framework" == "net5.0" ] && [ "$framework" == "net6.0" ]; then
  echo "Unsupported .net target framework $framework specified, choose from: netcoreapp3.1, net5.0 or net6.0"
  exit 1
fi

flags=""

if [ -z "${CI:-""}" ]; then
  # disable parallel execution on locally for better debugging experience
  flags="-parallel 1"
else
  # restrict parallel degree for ci mode because kind can hang with large amount of pods
  flags="-parallel 4"
fi

go test $flags \
  -v \
  -ldflags="-X github.com/dodopizza/kubectl-shovel/test/integration_test.TargetContainerImage=kubectl-shovel/sample-integration-tests:$framework" \
  -timeout 300s \
  --tags=integration \
  ./test/integration/... |
  sed "/PASS/s//$(printf "\033[32mPASS\033[0m")/" |
  sed "/FAIL/s//$(printf "\033[31mFAIL\033[0m")/"

exit ${PIPESTATUS[0]}
