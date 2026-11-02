# ADR 0004：API 不直接对外暴露

- 状态：Accepted
- 日期：2026-11-02

## 背景

登录限流按客户端 IP 计数，而客户端 IP 来自 `X-Real-IP` / `X-Forwarded-For`。只要请求能绕过 nginx 直连 API，就可以每次伪造一个新的 IP 头，绕过 5 次/15 分钟的登录失败限制。

Compose 和 k8s 原本都把 user/order/product API 发布到宿主机，k8s 还使用 NodePort，因此这个问题在两种环境都成立。

## 决策

1. **Compose：API 端口只绑本机**
   - 8000–8002 改为 `127.0.0.1:port:port`
   - 前端容器仍通过服务名访问，localhost 开发不受影响

2. **k8s：API 只保留 ClusterIP**
   - user-api / order-api / product-api 从 NodePort 改为 ClusterIP
   - 删除 k3d 的 8000–8002 宿主机映射
   - 对外只保留 frontend（8080/8081）和 gateway（8888）
   - 需要直连时使用 `kubectl -n mall port-forward`

3. **客户端 IP 解析**
   - 优先 `X-Real-IP`（nginx 会用 `$remote_addr` 覆盖）
   - 其次从右向左取 `X-Forwarded-For` 的最后一个合法 IP
   - 最后回退 `RemoteAddr`

## 后果

- 正常前端流量全部经过 nginx，限流按真实客户端 IP 生效
- k8s 下不能再从宿主机直接 curl API，需要用 port-forward
- gateway 的 `/user/login` 仍然直接映射到 user-rpc，不经过 user-api 的限流；这是已知遗留风险，后续在 user-rpc 或 gateway 补限流

## 备选

- 信任代理 CIDR 列表：需要区分 Compose 与 k8s 网段，配置更复杂
- nginx 与 API 共享代理密钥头：需要给 nginx 配置注入密钥，改动面更大
