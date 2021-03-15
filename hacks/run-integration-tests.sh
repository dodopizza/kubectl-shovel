#!/bin/bash
set -o errexit

script_dir="$( dirname "${BASH_SOURCE[0]}" )"
project_dir="${script_dir}/.."

kind_context="kind-kind"
current_context=$(kubectl config current-context)

if [ "${current_context}" != "${kind_context}" ]
then
  echo "Your context is wrong. Use ${kind_context}"
  exit 1
fi

echo "Building dumper..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
  go build -v \
  -o ${project_dir}/dumper/dumper \
  ${project_dir}/dumper

image_tag="latest"
image_repository="kubectl-shovel/dumper-integration-tests"

echo "Building dumper's image..."
docker build \
  -t ${image_repository}:${image_tag} \
  -f "${project_dir}/dumper/Dockerfile" \
  "${project_dir}/dumper"
rm "${project_dir}/dumper/dumper"

echo "Loading dumper's image to kind cluster..."
kind load docker-image ${image_repository}:${image_tag}

cli_binary_path="${project_dir}/test/integration/bin"
mkdir -p ${cli_binary_path}

echo "Building cli..."
CGO_ENABLED=0 \
  go build -v \
  -o ${cli_binary_path}/kubectl-shovel \
  ${project_dir}/cli

echo "Running tests..."
go test -v \
  --tags=integration \
  ./test/integration/... | \
  sed "/PASS/s//$(printf "\033[32mPASS\033[0m")/" | \
  sed "/FAIL/s//$(printf "\033[31mFAIL\033[0m")/"

exit ${PIPESTATUS[0]}

