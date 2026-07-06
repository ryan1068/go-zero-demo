# ── 构建阶段 ──
FROM node:20-alpine AS builder

WORKDIR /app

# 利用 Docker 缓存层：先装依赖
COPY package.json package-lock.json ./
RUN npm ci

# 再拷贝源码构建
COPY . .
RUN npm run build


# ── 运行阶段 ──
FROM nginx:alpine

# 复制构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制 Nginx 配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
