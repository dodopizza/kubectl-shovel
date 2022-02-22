#!/bin/bash
set -o errexit

script_dir="$(dirname "${BASH_SOURCE[0]}")"
project_dir="${script_dir}/.."

kind_context="kind-kind"
current_context=$(kubectl config current-context)

if [ "${current_context}" != "${kind_context}" ]; then
  echo "Your context is wrong. Use ${kind_context}"
  exit 1
fi

arch=${1}
if [ "$arch" == "x86_64" ]; then
  arch="amd64"
fi

if [ "$arch" != "amd64" ] && [ "$arch" != "arm64" ]; then
  echo "Unsupported arch, choose from: amd64 or arm64"
  exit 1
fi

echo "Building dumper..."
GOOS=linux GOARCH=$arch CGO_ENABLED=0 \
  go build -v \
  -o ${project_dir}/dumper/bin/dumper \
  ${project_dir}/dumper

image_tag="latest"
image_repository="kubectl-shovel/dumper-integration-tests"

echo "Building dumper's image..."
docker buildx build \
  --platform "linux/$arch" \
  --progress plain \
  --load \
  -t ${image_repository}:${image_tag} \
  -f "${project_dir}/dumper/Dockerfile" \
  "${project_dir}/dumper"
rm "${project_dir}/dumper/bin/dumper"

echo "Loading dumper's image to kind cluster..."
kind load docker-image ${image_repository}:${image_tag}