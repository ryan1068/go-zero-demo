#!/usr/bin/env bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$repo_root"

mkdir -p backups
out="backups/mysql-k8s-$(date +%Y%m%d-%H%M%S).sql"
tmp="$out.tmp"
trap 'rm -f "$tmp"' EXIT

kubectl -n mall exec mysql-0 -- sh -c \
  'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" mysqldump -uroot --single-transaction --routines --triggers --databases looklook_usercenter looklook_order looklook_product' \
  > "$tmp"

mv "$tmp" "$out"
trap - EXIT
echo "k8s backup written to $out"
