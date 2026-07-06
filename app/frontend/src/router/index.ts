import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/Login.vue'),
    },
    {
      path: '/',
      redirect: '/dashboard',
      component: () => import('@/layouts/AdminLayout.vue'),
      children: [
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('@/views/Dashboard.vue'),
        },
        {
          path: 'user',
          name: 'user',
          component: () => import('@/views/user/UserList.vue'),
        },
        {
          path: 'order',
          name: 'order',
          component: () => import('@/views/order/OrderList.vue'),
        },
        {
          path: 'order/create',
          name: 'order-create',
          component: () => import('@/views/order/OrderCreate.vue'),
        },
      ],
    },
  ],
})

// 路由守卫：未登录跳转到登录页
router.beforeEach((to, _from, next) => {
  const store = useUserStore()
  if (to.path !== '/login' && !store.isLoggedIn) {
    next('/login')
  } else if (to.path === '/login' && store.isLoggedIn) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
