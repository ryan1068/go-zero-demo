.PHONY: build test run clean gen

SERVICES := gateway user-api user-rpc order-api order-rpc

build:
	@for s in $(SERVICES); do \
		echo ">>> building $$s..."; \
		(cd app/$${s/rpc/rpc} && go build ./...) || exit 1; \
		(cd app/$${s/api/api} 2>/dev/null && go build ./...) || true; \
	done
	@echo "all services built"

test:
	@for d in app/user/api app/user/rpc app/order/api app/order/rpc; do \
		echo ">>> testing $$d..."; \
		(cd $$d && go test ./... -v -count=1) || exit 1; \
	done
	@echo "all tests passed"

test-short:
	@for d in app/user/api app/user/rpc app/order/api app/order/rpc; do \
		(cd $$d && go test ./... -count=1) || exit 1; \
	done

lint:
	@for d in app/user/api app/user/rpc app/order/api app/order/rpc app/gateway; do \
		echo ">>> linting $$d..."; \
		(cd $$d && go vet ./...) || exit 1; \
	done

fmt:
	@for d in app/user/api app/user/rpc app/order/api app/order/rpc app/gateway; do \
		(cd $$d && gofmt -w .); \
	done

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

gen-user-api:
	cd app/user/api && goctl api go -api user.api -dir .

gen-order-api:
	cd app/order/api && goctl api go -api order.api -dir .

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