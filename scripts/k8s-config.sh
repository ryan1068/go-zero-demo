#!/usr/bin/env bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$repo_root"

# service name -> config file
pairs=(
  "user-api:app/user/api/etc/user-api.yaml"
  "user-rpc:app/user/rpc/etc/user.yaml"
  "order-api:app/order/api/etc/order.yaml"
  "order-rpc:app/order/rpc/etc/order.yaml"
  "product-api:app/product/api/etc/product.yaml"
  "product-rpc:app/product/rpc/etc/product.yaml"
  "gateway:app/gateway/etc/gateway.yaml"
)

for pair in "${pairs[@]}"; do
  name=${pair%%:*}
  file=${pair#*:}
  base=$(basename "$file")
  kubectl -n mall create configmap "$name-config" \
    --from-file="$base=$file" \
    --dry-run=client -o yaml | kubectl apply -f -
done

kubectl -n mall create configmap frontend-nginx \
  --from-file=default.conf=deploy/k8s/frontend-nginx.conf \
  --dry-run=client -o yaml | kubectl apply -f -

# 空数据目录时 MySQL 会自动执行这些初始化脚本
kubectl -n mall create configmap mysql-init \
  --from-file=01-init-user.sql=deploy/sql/init_user.sql \
  --from-file=02-init-order.sql=deploy/sql/init_order.sql \
  --from-file=03-init-product.sql=deploy/sql/init_product.sql \
  --dry-run=client -o yaml | kubectl apply -f -

echo "configmaps applied"
