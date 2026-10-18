import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { isAdminMode, isMallMode } from '@/utils/portal'

const adminPrefixes = ['/dashboard', '/category', '/product', '/order', '/after-sale', '/audit', '/user']

function isAdminPath(path: string): boolean {
  return adminPrefixes.some((prefix) => path === prefix || path.startsWith(prefix + '/'))
}

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/Login.vue'),
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/Register.vue'),
    },
    {
      path: '/mall',
      component: () => import('@/layouts/MallLayout.vue'),
      children: [
        { path: '', name: 'mall-home', component: () => import('@/views/mall/Home.vue') },
        { path: 'products', name: 'mall-products', component: () => import('@/views/mall/ProductList.vue') },
        { path: 'product/:id', name: 'mall-product', component: () => import('@/views/mall/ProductDetail.vue') },
        { path: 'cart', name: 'mall-cart', component: () => import('@/views/mall/Cart.vue'), meta: { auth: true } },
        { path: 'checkout', name: 'mall-checkout', component: () => import('@/views/mall/Checkout.vue'), meta: { auth: true } },
        { path: 'orders', name: 'mall-orders', component: () => import('@/views/mall/OrderList.vue'), meta: { auth: true } },
        { path: 'order/:id', name: 'mall-order', component: () => import('@/views/mall/OrderDetail.vue'), meta: { auth: true } },
        { path: 'addresses', name: 'mall-addresses', component: () => import('@/views/mall/AddressList.vue'), meta: { auth: true } },
      ],
    },
    {
      path: '/',
      component: () => import('@/layouts/AdminLayout.vue'),
      meta: { admin: true },
      children: [
        { path: '', redirect: '/dashboard' },
        { path: 'dashboard', name: 'dashboard', component: () => import('@/views/Dashboard.vue') },
        { path: 'category', name: 'category', component: () => import('@/views/admin/CategoryList.vue') },
        { path: 'product', name: 'product', component: () => import('@/views/admin/ProductList.vue') },
        { path: 'product/create', name: 'product-create', component: () => import('@/views/admin/ProductEdit.vue') },
        { path: 'product/:id/edit', name: 'product-edit', component: () => import('@/views/admin/ProductEdit.vue') },
        { path: 'order', name: 'order', component: () => import('@/views/admin/OrderList.vue') },
        { path: 'after-sale', name: 'after-sale', component: () => import('@/views/admin/AfterSaleList.vue') },
        { path: 'audit', name: 'audit', component: () => import('@/views/admin/AuditLogList.vue') },
        { path: 'user', name: 'user', component: () => import('@/views/user/UserList.vue') },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/mall',
    },
  ],
})

router.beforeEach((to, _from, next) => {
  const store = useUserStore()

  // 商城端口只允许商城页面，后台端口只允许后台页面。
  if (isMallMode() && (to.meta.admin || isAdminPath(to.path))) {
    next('/mall')
    return
  }
  if (isAdminMode() && (to.path === '/mall' || to.path.startsWith('/mall/'))) {
    next('/dashboard')
    return
  }
  if (isAdminMode() && to.path === '/register') {
    next('/login')
    return
  }

  if (to.meta.auth && !store.isLoggedIn) {
    next({ path: '/login', query: { redirect: to.fullPath } })
    return
  }
  if (to.meta.admin && (!store.isLoggedIn || !store.isAdmin)) {
    if (isAdminMode() && store.isLoggedIn && !store.isAdmin) {
      // 后台端口遇到普通用户/旧 token：清掉登录态，留在后台登录页。
      store.logout()
      next({ path: '/login', query: { reason: 'not_admin' } })
      return
    }
    next(store.isLoggedIn ? '/mall' : '/login')
    return
  }
  if ((to.path === '/login' || to.path === '/register') && store.isLoggedIn) {
    next(isAdminMode() ? '/dashboard' : '/mall')
    return
  }
  next()
})

export default router
