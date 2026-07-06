FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build
COPY . .
ADD order/rpc/go.mod .
ADD order/rpc/go.sum .
RUN cd order/rpc && go mod download

COPY order/rpc/etc /app/etc
RUN cd order/rpc && go build -ldflags="-s -w" -o /app/order order.go


FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/order /app/order
COPY --from=builder /app/etc /app/etc

CMD ["./order", "-f", "etc/order.yaml"]
