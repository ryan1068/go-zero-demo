#!/usr/bin/env bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cluster=${CLUSTER:-mall}

cd "$repo_root"

images=(
  gateway:local
  user-rpc:local
  user-api:local
  order-rpc:local
  order-api:local
  product-rpc:local
  product-api:local
  frontend:local
)

echo ">>> building images"
docker compose build "${images[@]%%:*}"

nodes=$(kubectl get nodes -o jsonpath='{.items[*].metadata.name}')
if [ -z "$nodes" ]; then
  echo "no cluster nodes found; run make k8s-up first" >&2
  exit 1
fi

for img in "${images[@]}"; do
  tar="/tmp/k8s-image-$(echo "$img" | tr '/:' '__').tar"
  platform=$(docker image inspect --format '{{.Os}}/{{.Architecture}}' "$img")
  echo ">>> exporting $img ($platform)"
  docker save --platform "$platform" "$img" -o "$tar"
  for node in $nodes; do
    docker cp "$tar" "$node:/tmp/image.tar"
    docker exec "$node" ctr -n k8s.io images import /tmp/image.tar >/dev/null
  done
  rm -f "$tar"
done

echo "images imported into: $nodes"
