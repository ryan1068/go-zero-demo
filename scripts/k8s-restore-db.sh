#!/usr/bin/env bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$repo_root"

file=${1:-}
if [ -z "$file" ] || [ ! -f "$file" ]; then
  echo "usage: $0 <backup.sql>  (FORCE=1 to overwrite existing data)" >&2
  exit 2
fi
if [ ! -s "$file" ]; then
  echo "backup file is empty: $file" >&2
  exit 2
fi

force=${FORCE:-0}

# 读取当前计数；库还不存在时按空库处理
counts=$(kubectl -n mall exec mysql-0 -- sh -c \
  'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" mysql -uroot -N -e "SELECT (SELECT COUNT(*) FROM looklook_usercenter.user), (SELECT COUNT(*) FROM looklook_order.orders), (SELECT COUNT(*) FROM looklook_product.product);"' \
  2>/dev/null || echo "0 0 0")
read -r users orders products <<<"$counts"

# 初始化脚本会写入 1 个管理员和少量商品；超过这个基线视为已有业务数据
if [ "${users:-0}" -gt 1 ] || [ "${orders:-0}" -gt 0 ]; then
  if [ "$force" != "1" ]; then
    echo "target database already has data (users=$users orders=$orders products=$products)." >&2
    echo "refusing to overwrite; re-run with FORCE=1 if this is intended." >&2
    exit 1
  fi
  pre="backups/mysql-k8s-pre-restore-$(date +%Y%m%d-%H%M%S).sql"
  kubectl -n mall exec mysql-0 -- sh -c \
    'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" mysqldump -uroot --single-transaction --routines --triggers --databases looklook_usercenter looklook_order looklook_product' \
    > "$pre"
  echo "pre-restore backup: $pre"
fi

kubectl -n mall exec -i mysql-0 -- sh -c 'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" mysql -uroot' < "$file"

counts=$(kubectl -n mall exec mysql-0 -- sh -c \
  'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" mysql -uroot -N -e "SELECT (SELECT COUNT(*) FROM looklook_usercenter.user), (SELECT COUNT(*) FROM looklook_order.orders), (SELECT COUNT(*) FROM looklook_product.product);"')
echo "restored $file: users/orders/products = $counts"
