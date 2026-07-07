# Vue3 前端开发规范

## 技术栈

- Vue 3 (Composition API + `<script setup>`)
- Element Plus UI 组件库
- Pinia 状态管理
- Vue Router 路由
- Axios 请求封装
- Vite 构建工具

## 组件规范

```vue
<script setup>
import { ref, onMounted } from 'vue'

const data = ref([])

onMounted(async () => {
  data.value = await fetchData()
})
</script>

<template>
  <el-table :data="data" />
</template>
```

- 使用 `<script setup>` 语法糖
- 组件名 PascalCase，文件名 kebab-case
- 页面组件放在 `src/views/`，公共组件放在 `src/components/`

## 状态管理 (Pinia)

```js
// stores/user.js
import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', () => {
  const token = ref('')
  const userInfo = ref(null)

  function login() { /* ... */ }
  function logout() { /* ... */ }

  return { token, userInfo, login, logout }
})
```

- 全局状态用 Pinia store
- 页面级状态用组件内 `ref`/`reactive`
- store 命名：`useXxxStore`

## API 请求

```js
// api/order.js
import request from '@/utils/request'

export function createOrder(data) {
  return request.post('/order/create', data)
}
```

- 统一用 `@/utils/request` 封装的 axios 实例
- API 按模块分文件：`api/user.js`、`api/order.js`
- 错误统一在 request 拦截器中处理

## 路由规范

```js
{
  path: '/order',
  component: () => import('@/views/order/OrderList.vue'),
  meta: { title: '订单管理', auth: true }
}
```

- 路由懒加载
- `meta.auth` 标记需要登录的页面
- 管理后台页面放 `/admin/` 前缀

## 环境变量

```
.env.development → VITE_API_BASE=http://localhost:8888
.env.production  → VITE_API_BASE=https://api.example.com
```

## 禁止事项

- 不要在组件中直接操作 DOM
- 不要硬编码 API 地址
- 不要在模板中写复杂逻辑
- 所有用户输入必须校验