#!/usr/bin/env bash
set -eu

[ $# -lt 4 ] && echo "Usage: $(basename $0) <directory> <context> <arch|[amd64, arm64]> <framework|[netcoreapp3.1, net5.0, net6.0]>" && exit 1

directory=$1
context=$2
arch=$3
framework=$4

current_context=$(kubectl config current-context)
if [ "${current_context}" != "$context" ]; then
  echo "Current context must be $context. Set current context with command \"kubectl config set-context $context\""
  exit 1
fi

cluster_os="linux"
if [ "$arch" == "x86_64" ]; then
  cluster_arch="amd64"
elif [ "$arch" == "arm64" ]; then
  cluster_arch="arm64"
else
  echo "Unsupported arch $arch, choose from: x86_64 or arm64"
  exit 1
fi

# dumper options
dumper_image_tag="latest"
dumper_image_repository="kubectl-shovel/dumper-integration-tests"
dumper_context="${directory}/dumper"
dumper_binary="$dumper_context/bin/dumper"

# sample app options
sample_image_tag="$framework"
sample_image_repository="kubectl-shovel/sample-integration-tests"
sample_context="${directory}/sample"

if [ "$framework" == "netcoreapp3.1" ]; then
  sample_sdk_image_tag="3.1-focal"
  sample_runtime_image_tag="3.1-focal"
elif [ "$framework" == "net5.0" ]; then
  sample_sdk_image_tag="5.0-focal"
  sample_runtime_image_tag="5.0-focal"
elif [ "$framework" == "net6.0" ]; then
  sample_sdk_image_tag="6.0-focal"
  sample_runtime_image_tag="6.0-focal"
elif [ "$framework" == "net7.0" ]; then
  sample_sdk_image_tag="7.0-jammy"
  sample_runtime_image_tag="7.0-jammy"
elif [ "$framework" == "net8.0" ]; then
  sample_sdk_image_tag="8.0-jammy"
  sample_runtime_image_tag="8.0-jammy"
elif [ "$framework" == "net9.0" ]; then
  sample_sdk_image_tag="9.0-noble"
  sample_runtime_image_tag="9.0-noble"
else
  echo "Unsupported .net target framework $framework specified, choose from: netcoreapp3.1, net5.0, net6.0, net7.0, net8.0, net9.0"
  exit 1
fi


echo "Building dumper binary ($cluster_os/$cluster_arch):"
GOOS=$cluster_os GOARCH=$cluster_arch CGO_ENABLED=0 \
  go build \
  -v \
  -o "$dumper_binary" \
  "./dumper"

echo "Building dumper docker image ($cluster_os/$cluster_arch):"
docker buildx build \
  --platform "$cluster_os/$cluster_arch" \
  --progress plain \
  --load \
  -t "$dumper_image_repository:$dumper_image_tag" \
  -f "$dumper_context/Dockerfile" \
  "${directory}"
rm "$dumper_binary"

echo "Building sample docker image ($cluster_os/$cluster_arch):"
docker buildx build \
  --platform "$cluster_os/$cluster_arch" \
  --progress plain \
  --load \
  --build-arg SDK_IMAGE_TAG="$sample_sdk_image_tag" \
  --build-arg RUNTIME_IMAGE_TAG="$sample_runtime_image_tag" \
  --build-arg FRAMEWORK="$framework" \
  -t "$sample_image_repository:$sample_image_tag" \
  -t "$sample_image_repository:latest" \
  -f "$sample_context/Dockerfile" \
  "$sample_context"

images=(
  "$dumper_image_repository:$dumper_image_tag"
  "$sample_image_repository:$sample_image_tag"
  "$sample_image_repository:latest"
)

for image in "${images[@]}"; do
  echo "Loading image to cluster ($image):"
  kind load docker-image "$image"
done

