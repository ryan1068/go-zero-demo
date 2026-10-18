#!/usr/bin/env bash
set -euo pipefail

# 用法：
#   MYSQL_ROOT_PASSWORD=xxx ./genModel.sh <usercenter|order|product> <表名，逗号分隔> [目标目录]
# 示例：
#   MYSQL_ROOT_PASSWORD=xxx ./genModel.sh usercenter user,user_address app/user/rpc/internal/model
#
# 也可以直接使用仓库根目录的 make gen-model。

db=${1:?usage: ./genModel.sh <usercenter|order|product> <tables> [dir]}
tables=${2:?usage: ./genModel.sh <usercenter|order|product> <tables> [dir]}
modeldir=${3:-internal/model}

: "${MYSQL_ROOT_PASSWORD:?MYSQL_ROOT_PASSWORD is required}"
host=${MYSQL_HOST:-127.0.0.1}
port=${MYSQL_PORT:-33069}

echo "generating models for looklook_${db}: ${tables} -> ${modeldir}"
goctl model mysql datasource \
  -url="root:${MYSQL_ROOT_PASSWORD}@tcp(${host}:${port})/looklook_${db}" \
  -table="${tables}" \
  -dir="${modeldir}" \
  --style=goZero \
  -i "create_at,created_at,create_time,update_at,updated_at,update_time,deleted_at"
