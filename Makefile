.PHONY: build test lint fmt fmt-check check protogen clean \
	run-gateway run-user-api run-user-rpc run-order-api run-order-rpc run-product-api run-product-rpc \
	gen-model gen-check smoke db-backup k8s-up k8s-down k8s-status \
	k8s-images k8s-secrets k8s-config k8s-deploy k8s-restore-db k8s-backup k8s-delete \
	k8s-start k8s-stop k8s-restart compose-off compose-on \
	gen-user-api gen-order-api gen-product-api gen-user-rpc gen-order-rpc \
	frontend-dev frontend-build

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

lint:
	@for d in $(MODULES); do \
		echo ">>> linting $$d..."; \
		(cd $$d && go vet ./...) || exit 1; \
	done

fmt:
	@for d in $(MODULES); do \
		(cd $$d && gofmt -w .); \
	done

fmt-check:
	@for d in $(MODULES); do \
		out=$$(cd $$d && gofmt -l .); \
		if [ -n "$$out" ]; then \
			echo "gofmt required in $$d:"; \
			echo "$$out"; \
			exit 1; \
		fi; \
	done
	@echo "format check passed"

check: fmt-check build test lint
	@echo "all checks passed"

smoke:
	@./scripts/smoke.sh

db-backup:
	@./scripts/db-backup.sh

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

gen-check:
	@./scripts/gen-check.sh

gen-user-rpc:
	cd app/user/rpc && goctl rpc protoc user.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=.

gen-order-rpc:
	cd app/order/rpc && goctl rpc protoc order.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=.

frontend-dev:
	cd app/frontend && npm run dev

frontend-build:
	cd app/frontend && npm run build

k8s-up:
	k3d cluster create --config deploy/k8s/k3d-cluster.yaml --registry-config deploy/k8s/registries.yaml
	kubectl create namespace mall --dry-run=client -o yaml | kubectl apply -f -

k8s-down:
	k3d cluster delete mall

k8s-status:
	@k3d cluster list
	@kubectl get nodes -o wide

k8s-start:
	@k3d cluster start mall 2>/dev/null || true
	@kubectl wait --for=condition=Ready nodes --all --timeout=180s

k8s-stop:
	k3d cluster stop mall

k8s-restart:
	k3d cluster stop mall
	k3d cluster start mall
	kubectl wait --for=condition=Ready nodes --all --timeout=180s

compose-off:
	@./scripts/compose-off.sh

compose-on:
	@./scripts/compose-on.sh

k8s-images:
	@./scripts/k8s-images.sh

k8s-secrets:
	@./scripts/k8s-secrets.sh

k8s-config:
	@./scripts/k8s-config.sh

DB_BACKUP ?= $(shell ls -t backups/mysql-*.sql 2>/dev/null | head -1)
RESTORE ?= 0
FORCE ?= 0

k8s-deploy: k8s-secrets k8s-config
	kubectl apply -f deploy/k8s/base/namespace.yaml
	kubectl apply -f deploy/k8s/base/infra.yaml
	kubectl -n mall rollout status statefulset/mysql --timeout=300s
	kubectl -n mall rollout status deployment/redis deployment/etcd --timeout=300s
ifneq ($(RESTORE),0)
	FORCE=$(FORCE) ./scripts/k8s-restore-db.sh "$(DB_BACKUP)"
endif
	kubectl apply -f deploy/k8s/base/apps.yaml
	kubectl apply -f deploy/k8s/base/frontend.yaml
	@for d in $$(kubectl -n mall get deployment -o name); do kubectl -n mall rollout restart $$d; done
	kubectl -n mall rollout restart statefulset/mysql
	kubectl -n mall wait --for=condition=Available deployment --all --timeout=300s

k8s-restore-db:
	@FORCE=$(FORCE) ./scripts/k8s-restore-db.sh "$(DB_BACKUP)"

k8s-backup:
	@./scripts/k8s-backup.sh

k8s-delete:
	kubectl delete namespace mall --ignore-not-found

clean:
	rm -rf app/*/api/internal/{handler,logic,types} app/*/rpc/pb app/*/rpc/*client
