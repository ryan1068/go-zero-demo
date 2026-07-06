<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>

<template>
  <div class="home">
    <h1>go-zero-demo 微服务管理平台</h1>
    <p class="subtitle">Vue 3 + TypeScript + Vite 前端项目</p>

    <div class="card-grid">
      <RouterLink to="/user" class="card">
        <h2>👤 用户查询</h2>
        <p>根据用户名查询用户信息，调用 user-api → user-rpc</p>
      </RouterLink>

      <RouterLink to="/login" class="card">
        <h2>🔐 用户登录</h2>
        <p>用户名密码登录，调用 user-api → user-rpc</p>
      </RouterLink>

      <RouterLink to="/order" class="card">
        <h2>📦 订单查询</h2>
        <p>根据订单 ID 查询订单详情，调用 order-api → order-rpc</p>
      </RouterLink>

      <RouterLink to="/order/create" class="card">
        <h2>🛒 创建订单</h2>
        <p>创建新订单，order-rpc 内部调用 user-rpc 验证用户</p>
      </RouterLink>
    </div>
  </div>
</template>

<style scoped>
.home {
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem;
  text-align: center;
}

h1 {
  font-size: 2rem;
  margin-bottom: 0.5rem;
}

.subtitle {
  color: #888;
  margin-bottom: 2.5rem;
}

.card-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.25rem;
}

.card {
  display: block;
  padding: 1.5rem;
  border: 1px solid #e0e0e0;
  border-radius: 10px;
  text-decoration: none;
  color: inherit;
  transition: box-shadow 0.2s, transform 0.2s;
}

.card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.card h2 {
  margin: 0 0 0.5rem;
  font-size: 1.2rem;
}

.card p {
  margin: 0;
  color: #666;
  font-size: 0.9rem;
}
</style>
