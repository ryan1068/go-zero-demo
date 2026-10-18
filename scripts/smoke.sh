#!/usr/bin/env bash
set -uo pipefail

USER_API=${USER_API:-http://localhost:8000}
ORDER_API=${ORDER_API:-http://localhost:8001}
PRODUCT_API=${PRODUCT_API:-http://localhost:8002}
FRONTEND=${FRONTEND:-http://localhost:8081}
MEDIA_API=${MEDIA_API:-$FRONTEND/api/media}
ADMIN_USER=${ADMIN_USER:-admin}
ADMIN_PASSWORD=${ADMIN_PASSWORD:-123456}
PRODUCT_ID=${PRODUCT_ID:-10}

fail=0
LAST_BODY=""

pass() {
  printf 'PASS  %s\n' "$1"
}

fail_msg() {
  printf 'FAIL  %s: %s\n' "$1" "$2"
  fail=1
}

require() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "missing required command: $1"
    exit 2
  }
}

require curl
require jq

login_body=$(curl -sS -X POST "$USER_API/user/login" \
  -H 'Content-Type: application/json' \
  -d "{\"username\":\"$ADMIN_USER\",\"password\":\"$ADMIN_PASSWORD\"}" 2>/dev/null) || {
  echo "FAIL  login: request error"
  exit 1
}
token=$(jq -r '.data.token // empty' <<<"$login_body")
if [ -n "$token" ]; then
  pass "login"
else
  fail_msg "login" "$login_body"
fi

check_api() {
  local name=$1 method=$2 url=$3 expected=$4
  shift 4
  local response http body code
  response=$(curl -sS -w '\n%{http_code}' -X "$method" "$url" "$@" 2>/dev/null) || {
    fail_msg "$name" "request error"
    LAST_BODY="{}"
    return
  }
  http=$(tail -n 1 <<<"$response")
  body=$(sed '$d' <<<"$response")
  code=$(jq -r '.code // 200' <<<"$body" 2>/dev/null || echo "?")
  if [ "$http" = "$expected" ] && [ "$code" = "200" ]; then
    pass "$name"
    LAST_BODY=$body
  else
    fail_msg "$name" "http=$http code=$code body=$body"
    LAST_BODY="{}"
  fi
}

check_http() {
  local name=$1 url=$2 expected=$3
  local http
  http=$(curl -sS -o /dev/null -w '%{http_code}' "$url" 2>/dev/null) || {
    fail_msg "$name" "request error"
    return
  }
  if [ "$http" = "$expected" ]; then
    pass "$name"
  else
    fail_msg "$name" "http=$http"
  fi
}

check_api "users list" GET "$USER_API/users?page=1&page_size=1" 200 -H "Authorization: Bearer $token"
users=$(jq -r '.data.total // 0' <<<"$LAST_BODY")

check_api "admin orders" GET "$ORDER_API/admin/orders?page=1&page_size=1" 200 -H "Authorization: Bearer $token"
orders=$(jq -r '.data.total // 0' <<<"$LAST_BODY")

check_api "dashboard" GET "$ORDER_API/admin/dashboard" 200 -H "Authorization: Bearer $token"
check_api "audit logs" GET "$USER_API/admin/audit-logs?page=1&page_size=1" 200 -H "Authorization: Bearer $token"
check_api "cart" GET "$ORDER_API/cart" 200 -H "Authorization: Bearer $token"

check_api "product list" GET "$PRODUCT_API/products?page=1&page_size=1" 200
products=$(jq -r '.data.total // 0' <<<"$LAST_BODY")

check_api "category list" GET "$PRODUCT_API/categories" 200
check_api "product detail" GET "$PRODUCT_API/product/$PRODUCT_ID" 200
check_api "search 龙井" GET "$PRODUCT_API/products?page=1&page_size=5&keyword=%E9%BE%99%E4%BA%95" 200
check_api "suggest 龙井" GET "$PRODUCT_API/products/suggest?keyword=%E9%BE%99%E4%BA%95" 200

media_png="/tmp/smoke-media-$$.png"
printf '\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x06\x00\x00\x00\x1f\x15\xc4\x89\x00\x00\x00\nIDATx\x9cc\x00\x01\x00\x00\x05\x00\x01\r\n-\xb4\x00\x00\x00\x00IEND\xaeB`\x82' > "$media_png"
check_api "media upload" POST "$MEDIA_API/upload" 200 \
  -H "Authorization: Bearer $token" -F "file=@$media_png"
media_url=$(jq -r '.url // empty' <<<"$LAST_BODY")
if [ -n "$media_url" ]; then
  check_http "media read" "$FRONTEND$media_url" 200
  media_name=${media_url##*/}
  check_api "media delete" DELETE "$MEDIA_API/admin/files/$media_name" 200 \
    -H "Authorization: Bearer $token"
  check_http "media read after delete" "$FRONTEND$media_url" 404
else
  fail_msg "media url" "upload did not return url"
fi
rm -f "$media_png"

if [ "$users" -le 0 ] 2>/dev/null; then
  fail_msg "user count" "expected > 0, got $users"
fi
if [ "$products" -le 0 ] 2>/dev/null; then
  fail_msg "product count" "expected > 0, got $products"
fi

metrics_runtime=""
if command -v docker >/dev/null 2>&1 && docker ps --format '{{.Names}}' 2>/dev/null | grep -qx frontend; then
  metrics_runtime=docker
elif command -v kubectl >/dev/null 2>&1 && kubectl -n mall get deployment frontend >/dev/null 2>&1; then
  metrics_runtime=k8s
fi

if [ -n "$metrics_runtime" ]; then
  for spec in user-rpc:9100 user-api:9101 order-rpc:9102 order-api:9103 \
    product-api:9104 product-rpc:9105 media-api:9106; do
    name=${spec%%:*}
    port=${spec##*:}
    metrics_body=""
    if [ "$metrics_runtime" = docker ]; then
      metrics_body=$(docker exec frontend wget -q -O - "http://$name:$port/metrics" 2>/dev/null || true)
    else
      pod_ip=$(kubectl -n mall get pods -l app="$name" -o jsonpath='{.items[0].status.podIP}' 2>/dev/null || true)
      if [ -n "$pod_ip" ]; then
        metrics_body=$(kubectl -n mall exec deployment/frontend -- \
          wget -q -O - "http://$pod_ip:$port/metrics" 2>/dev/null || true)
      fi
    fi
    if grep -qm1 '^# HELP' <<<"$metrics_body"; then
      pass "metrics $name:$port"
    else
      fail_msg "metrics $name:$port" "unreachable"
    fi
  done
else
  echo "SKIP  metrics: no running frontend in Docker or k8s"
fi

if [ "$fail" -ne 0 ]; then
  echo "smoke failed"
  exit 1
fi

echo "smoke passed (users=$users, products=$products, orders=$orders)"
