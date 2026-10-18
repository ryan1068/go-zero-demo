#!/usr/bin/env bash
set -uo pipefail

USER_API=${USER_API:-http://localhost:8000}
ORDER_API=${ORDER_API:-http://localhost:8001}
PRODUCT_API=${PRODUCT_API:-http://localhost:8002}
ADMIN_USER=${ADMIN_USER:-admin}
ADMIN_PASSWORD=${ADMIN_PASSWORD:-123456}

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
check_api "product detail" GET "$PRODUCT_API/product/10" 200
check_api "search 龙井" GET "$PRODUCT_API/products?page=1&page_size=5&keyword=%E9%BE%99%E4%BA%95" 200
check_api "suggest 龙井" GET "$PRODUCT_API/products/suggest?keyword=%E9%BE%99%E4%BA%95" 200

if [ "$users" -le 0 ] 2>/dev/null; then
  fail_msg "user count" "expected > 0, got $users"
fi
if [ "$products" -le 0 ] 2>/dev/null; then
  fail_msg "product count" "expected > 0, got $products"
fi

if command -v docker >/dev/null 2>&1 && docker ps --format '{{.Names}}' 2>/dev/null | grep -qx frontend; then
  for spec in user-rpc:9100 user-api:9101 order-rpc:9102 order-api:9103 product-api:9104 product-rpc:9105; do
    name=${spec%%:*}
    port=${spec##*:}
    if [ -n "$(docker exec frontend wget -q -O - "http://$name:$port/metrics" 2>/dev/null | grep -m1 '^# HELP' || true)" ]; then
      pass "metrics $name:$port"
    else
      fail_msg "metrics $name:$port" "unreachable"
    fi
  done
else
  echo "SKIP  metrics: frontend container not running"
fi

if [ "$fail" -ne 0 ]; then
  echo "smoke failed"
  exit 1
fi

echo "smoke passed (users=$users, products=$products, orders=$orders)"
