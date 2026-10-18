#!/usr/bin/env bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
work=$(mktemp -d "$repo_root/.gen-check.XXXXXX")
trap 'rm -rf "$work"' EXIT

command -v goctl >/dev/null 2>&1 || {
  echo "goctl is required (expected v1.8.5)"
  exit 2
}

fail=0

check_app() {
  local src=$1
  local api=$2
  local dest="$work/$(echo "$src" | tr '/' '_')"
  local excludes=()

  case "$src" in
    app/user/api)
      # Hand-written: login injects client IP, logout passes the bearer token,
      # authutil provides the client IP helper.
      excludes=(-x loginhandler.go -x logouthandler.go -x authutil.go)
      ;;
    app/product/api)
      # goctl cannot generate multipart upload or static file responses.
      excludes=(-x serveuploadhandler.go)
      ;;
  esac

  mkdir -p "$dest"
  cp -R "$repo_root/$src/." "$dest/"
  rm -f "$dest"/internal/handler/*.go
  rm -f "$dest"/internal/types/types.go

  (
    cd "$dest"
    goctl api go -api "$api" -dir . >/dev/null
  )

  local path
  for path in internal/handler internal/types internal/logic internal/svc internal/middleware; do
    if [ ! -d "$repo_root/$src/$path" ]; then
      continue
    fi
    if [ ${#excludes[@]} -gt 0 ]; then
      if ! diff -ru "${excludes[@]}" "$repo_root/$src/$path" "$dest/$path"; then
        echo "generated code drift: $src/$path"
        fail=1
      fi
    elif ! diff -ru "$repo_root/$src/$path" "$dest/$path"; then
      echo "generated code drift: $src/$path"
      fail=1
    fi
  done
}

check_app app/user/api user.api
check_app app/order/api order.api
check_app app/product/api product.api

if [ "$fail" -ne 0 ]; then
  echo "generated API code is out of date; run the matching make gen-*-api target"
  exit 1
fi

echo "generated API code is up to date"
