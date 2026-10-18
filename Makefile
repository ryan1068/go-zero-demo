.PHONY: build test test-short lint fmt protogen run clean gen gen-model gen-user-api gen-order-api gen-product-api

MODULES := app/gateway app/user/api app/user/rpc app/order/api app/order/rpc app/product/api app/product/rpc pkg/auth pkg/errcode

MYSQL_HOST ?= 127.0.0.1
MYSQL_PORT ?= 33069

build:
	@for d in $(MODULES); do \
		echo ">>> building $$d..."; \
		(cd $$d && go build ./...) || exit 1; \
	done
	@echo "all services built"

test:
	@for d in $(MODULES); do \
		echo ">>> testing $$d..."; \
		(cd $$d && go test ./... -count=1) || exit 1; \
	done
	@echo "all tests passed"

test-short:
	@for d in $(MODULES); do \
		(cd $$d && go test ./... -count=1) || exit 1; \
	done

lint:
	@for d in $(MODULES); do \
		echo ">>> linting $$d..."; \
		(cd $$d && go vet ./...) || exit 1; \
	done

fmt:
	@for d in $(MODULES); do \
		(cd $$d && gofmt -w .); \
	done

protogen:
	./scripts/protogen.sh

run-gateway:
	cd app/gateway && go run gateway.go

run-user-api:
	cd app/user/api && go run user.go

run-user-rpc:
	cd app/user/rpc && go run user.go

run-order-api:
	cd app/order/api && go run order.go

run-order-rpc:
	cd app/order/rpc && go run order.go

run-product-api:
	cd app/product/api && go run product.go

run-product-rpc:
	cd app/product/rpc && go run product.go

gen-model:
	@test -n "$$MYSQL_ROOT_PASSWORD" || (echo "MYSQL_ROOT_PASSWORD is required"; exit 1)
	cd app/user/rpc && goctl model mysql datasource -url="root:$$MYSQL_ROOT_PASSWORD@tcp($(MYSQL_HOST):$(MYSQL_PORT))/looklook_usercenter" -table="user,user_address,admin_audit_log" -dir internal/model --style=goZero -i "create_at,created_at,create_time,update_at,updated_at,update_time,deleted_at"
	cd app/order/rpc && goctl model mysql datasource -url="root:$$MYSQL_ROOT_PASSWORD@tcp($(MYSQL_HOST):$(MYSQL_PORT))/looklook_order" -table="cart_item,orders,order_item,payment_record,order_log,after_sale,refund_record" -dir internal/model --style=goZero -i "create_at,created_at,create_time,update_at,updated_at,update_time,deleted_at"
	cd app/product/rpc && goctl model mysql datasource -url="root:$$MYSQL_ROOT_PASSWORD@tcp($(MYSQL_HOST):$(MYSQL_PORT))/looklook_product" -table="category,product,product_sku" -dir internal/model --style=goZero -i "create_at,created_at,create_time,update_at,updated_at,update_time,deleted_at"

gen-user-api:
	cd app/user/api && goctl api go -api user.api -dir .

gen-order-api:
	cd app/order/api && goctl api go -api order.api -dir .

gen-product-api:
	cd app/product/api && goctl api go -api product.api -dir .

gen-user-rpc:
	cd app/user/rpc && goctl rpc protoc user.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=.

gen-order-rpc:
	cd app/order/rpc && goctl rpc protoc order.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=.

frontend-dev:
	cd app/frontend && npm run dev

frontend-build:
	cd app/frontend && npm run build

clean:
	rm -rf app/*/api/internal/{handler,logic,types} app/*/rpc/pb app/*/rpc/*client
