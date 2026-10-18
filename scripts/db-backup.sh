#!/usr/bin/env bash
set -euo pipefail

: "${MYSQL_ROOT_PASSWORD:?MYSQL_ROOT_PASSWORD is required}"
MYSQL_CONTAINER=${MYSQL_CONTAINER:-mysql}
BACKUP_DIR=${BACKUP_DIR:-backups}

mkdir -p "$BACKUP_DIR"
out="$BACKUP_DIR/mysql-$(date +%Y%m%d-%H%M%S).sql"
tmp="$out.tmp"
trap 'rm -f "$tmp"' EXIT

docker exec -e MYSQL_PWD="$MYSQL_ROOT_PASSWORD" "$MYSQL_CONTAINER" \
  mysqldump -uroot --single-transaction --routines --triggers \
  --databases looklook_usercenter looklook_order looklook_product > "$tmp"

mv "$tmp" "$out"
trap - EXIT

echo "backup written to $out"
