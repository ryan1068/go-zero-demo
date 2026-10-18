#!/usr/bin/env bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$repo_root"

set -a
# shellcheck disable=SC1091
. ./.env
set +a

: "${JWT_SECRET:?JWT_SECRET is required in .env}"
: "${MYSQL_ROOT_PASSWORD:?MYSQL_ROOT_PASSWORD is required in .env}"
: "${REDIS_PASSWORD:?REDIS_PASSWORD is required in .env}"

kubectl create namespace mall --dry-run=client -o yaml | kubectl apply -f -

tmp=$(mktemp)
trap 'rm -f "$tmp"' EXIT
chmod 600 "$tmp"
{
  printf 'JWT_SECRET=%s\n' "$JWT_SECRET"
  printf 'MYSQL_ROOT_PASSWORD=%s\n' "$MYSQL_ROOT_PASSWORD"
  printf 'REDIS_PASSWORD=%s\n' "$REDIS_PASSWORD"
} > "$tmp"

kubectl -n mall create secret generic mall-secrets \
  --from-env-file="$tmp" \
  --dry-run=client -o yaml | kubectl apply -f -

echo "secret mall-secrets applied"
