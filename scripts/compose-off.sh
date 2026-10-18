#!/usr/bin/env bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$repo_root"

docker compose stop
for c in $(docker compose ps -aq); do
  docker update --restart=no "$c" >/dev/null
done

echo "compose stopped; auto-restart disabled"
