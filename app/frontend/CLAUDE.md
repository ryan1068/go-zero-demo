# Vue3 前端开发规范

## 技术栈

Vue 3（Composition API + `<script setup>`）、Element Plus、Pinia、Vue Router、Axios、Vite。

## 组件与目录

- 一律使用 `<script setup>`；组件名与文件名 PascalCase（如 `ProductCard.vue`）
- 页面放 `src/views/`，公共组件放 `src/components/`，布局放 `src/layouts/`
- 不要在组件中直接操作 DOM，不在模板中写复杂逻辑
- 所有用户输入必须校验

## 状态与请求

- 全局状态用 Pinia，store 命名 `useXxxStore`；页面级状态用组件内 `ref`/`reactive`
- 统一用 `@/utils/request` 的 axios 实例，API 按模块分文件（`src/api/user.ts`）
- 错误统一在 request 拦截器中处理；不要硬编码 API 地址

## 路由与样式

- 路由懒加载，`meta.auth` 标记需登录页面；管理后台页面放 `src/views/admin/`
- 商城页面使用 `src/styles/brand.css` 的品牌变量，不写死颜色
- 搜索/筛选不要用原始 ID（用户 ID、操作人 ID、分类 ID），用名称、单号或状态
