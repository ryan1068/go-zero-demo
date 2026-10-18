#!/usr/bin/env bash
set -euo pipefail

# Regenerate protobuf message + gRPC stubs for user/order/product.
#
# Requirements:
#   protoc
#   protoc-gen-go
#   protoc-gen-go-grpc
#
# Install on macOS:
#   brew install protobuf
#   GOBIN=/tmp/protogen-build GOPROXY=https://goproxy.cn,direct \
#     go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.6
#   GOBIN=/tmp/protogen-build GOPROXY=https://goproxy.cn,direct \
#     go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
#
# Then run:
#   ./scripts/protogen.sh
#
# Set PROTOGEN_BIN to the directory containing the two plugins if they are not
# in PATH and not in /tmp/protogen-build. Set PROTOC if protoc is not in PATH.

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PROTOC="${PROTOC:-protoc}"

if ! command -v "$PROTOC" >/dev/null 2>&1; then
  echo "protoc not found; install it with: brew install protobuf" >&2
  exit 1
fi

if [ -z "${PROTOGEN_BIN:-}" ]; then
  if [ -x /tmp/protogen-build/protoc-gen-go ] && [ -x /tmp/protogen-build/protoc-gen-go-grpc ]; then
    PROTOGEN_BIN=/tmp/protogen-build
  else
    PROTOGEN_BIN="$(dirname "$(command -v protoc-gen-go)")"
  fi
fi

for plugin in protoc-gen-go protoc-gen-go-grpc; do
  if [ ! -x "$PROTOGEN_BIN/$plugin" ]; then
    echo "$plugin not found in $PROTOGEN_BIN" >&2
    exit 1
  fi
done

OUT_DIR="$(mktemp -d)"
trap 'rm -rf "$OUT_DIR"' EXIT

for svc in user order product; do
  echo ">>> generating $svc.proto"
  mkdir -p "$OUT_DIR/$svc"
  (
    cd "$ROOT/app/$svc/rpc"
    "$PROTOC" \
      --plugin=protoc-gen-go="$PROTOGEN_BIN/protoc-gen-go" \
      --plugin=protoc-gen-go-grpc="$PROTOGEN_BIN/protoc-gen-go-grpc" \
      --go_out="$OUT_DIR/$svc" --go_opt=paths=source_relative \
      --go-grpc_out="$OUT_DIR/$svc" --go-grpc_opt=paths=source_relative \
      "$svc.proto"
    cp "$OUT_DIR/$svc/$svc.pb.go" "pb/$svc/"
    cp "$OUT_DIR/$svc/${svc}_grpc.pb.go" "pb/$svc/"
  )
done

cat <<'EOF'
protobuf stubs regenerated.

Note: the hand-written client wrappers and server method bindings are not
regenerated. If you added an RPC, update:
  app/<service>/rpc/<service>client/*.go
  app/<service>/rpc/internal/server/*.go
EOF
