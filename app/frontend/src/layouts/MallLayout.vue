<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import BrandLogo from '@/components/BrandLogo.vue'
import { useUserStore } from '@/stores/user'
import { adminUrl, openInNewTab } from '@/utils/portal'

const router = useRouter()
const store = useUserStore()
const keyword = ref('')

const hotWords = ['Go 语言', '蓝牙耳机', '手冲咖啡']

function search() {
  router.push({ path: '/mall/products', query: { keyword: keyword.value.trim() || undefined } })
}

function searchHot(word: string) {
  keyword.value = word
  search()
}

function handleLogout() {
  store.logout()
  router.push('/login')
}
</script>

<template>
  <div class="mall-shell">
    <div class="topbar">
      <div class="brand-container topbar-inner">
        <span class="slogan">拾光严选 · 把好物收进时光里</span>
        <div class="top-links">
          <template v-if="store.isLoggedIn">
            <span class="welcome">你好，{{ store.username }}</span>
            <RouterLink to="/mall/orders">我的订单</RouterLink>
            <RouterLink to="/mall/addresses">收货地址</RouterLink>
            <a v-if="store.isAdmin" href="javascript:;" @click="openInNewTab(adminUrl('/dashboard'))">管理后台</a>
            <a href="javascript:;" @click="handleLogout">退出</a>
          </template>
          <template v-else>
            <RouterLink to="/login" class="welcome">亲，请登录</RouterLink>
            <RouterLink to="/register">免费注册</RouterLink>
          </template>
        </div>
      </div>
    </div>

    <header class="mall-header">
      <div class="brand-container header-inner">
        <div class="brand" @click="router.push('/mall')">
          <BrandLogo subtitle="SHIGUANG SELECT" />
        </div>

        <div class="search-area">
          <div class="search-box">
            <el-input
              v-model="keyword"
              placeholder="搜索好物，例如：耳机、咖啡、图书"
              clearable
              size="large"
              @keyup.enter="search"
            />
            <el-button type="primary" size="large" @click="search">搜 索</el-button>
          </div>
          <div class="hot-words">
            <span>热搜：</span>
            <a v-for="word in hotWords" :key="word" href="javascript:;" @click="searchHot(word)">{{ word }}</a>
          </div>
        </div>

        <button class="cart-entry" type="button" @click="router.push('/mall/cart')">
          <el-icon :size="22"><ShoppingCart /></el-icon>
          <span>购物车</span>
        </button>
      </div>

      <nav class="mall-nav">
        <div class="brand-container nav-inner">
          <RouterLink to="/mall">首页</RouterLink>
          <RouterLink to="/mall/products">全部商品</RouterLink>
          <RouterLink to="/mall/products?sort=sales">限时好物</RouterLink>
          <RouterLink to="/mall/orders">我的订单</RouterLink>
          <RouterLink to="/mall/addresses">收货地址</RouterLink>
          <span class="nav-note">拾光精选 · 正品好物 · 售后无忧</span>
        </div>
      </nav>
    </header>

    <main class="mall-main">
      <RouterView />
    </main>

    <footer class="mall-footer">
      <div class="brand-container footer-inner">
        <div class="footer-brand">
          <BrandLogo subtitle="把好物收进时光里" light />
        </div>
        <div class="promises">
          <div><b>正品保障</b><span>精选供应链，品质可溯</span></div>
          <div><b>极速发货</b><span>下单后 48 小时内发出</span></div>
          <div><b>售后无忧</b><span>七天无理由退换</span></div>
        </div>
        <div class="copyright">© {{ new Date().getFullYear() }} 拾光严选 · SHIGUANG SELECT</div>
      </div>
    </footer>
  </div>
</template>

<style scoped>
.mall-shell {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.topbar {
  background: var(--ink);
  color: rgba(247, 239, 240, 0.78);
  font-size: 12px;
}

.topbar-inner {
  height: 34px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
}

.slogan {
  letter-spacing: 1px;
  font-family: var(--font-serif);
  color: var(--gold);
}

.top-links {
  display: flex;
  align-items: center;
  gap: 16px;
}

.top-links a,
.top-links .welcome {
  color: rgba(247, 239, 240, 0.78);
  font-size: 12px;
}

.top-links a:hover {
  color: #fff;
}

.mall-header {
  position: sticky;
  top: 0;
  z-index: 20;
  background: rgba(255, 253, 248, 0.96);
  backdrop-filter: blur(8px);
  border-bottom: 1px solid var(--line);
  box-shadow: 0 4px 18px rgba(46, 35, 38, 0.06);
}

.header-inner {
  height: 92px;
  display: flex;
  align-items: center;
  gap: 36px;
}

.brand {
  cursor: pointer;
  flex-shrink: 0;
}

.search-area {
  flex: 1;
  max-width: 620px;
}

.search-box {
  display: flex;
  gap: 0;
}

.search-box :deep(.el-input__wrapper) {
  border-radius: 6px 0 0 6px;
  box-shadow: 0 0 0 2px var(--brand) inset;
  background: #fffefa;
}

.search-box :deep(.el-input__inner) {
  font-size: 14px;
}

.search-box .el-button {
  border-radius: 0 6px 6px 0;
  padding: 0 26px;
  font-family: var(--font-serif);
  letter-spacing: 4px;
}

.hot-words {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 8px;
  color: var(--muted);
  font-size: 12px;
}

.hot-words a {
  color: var(--ink-soft);
}

.hot-words a:hover {
  color: var(--brand);
}

.cart-entry {
  display: flex;
  align-items: center;
  gap: 8px;
  height: 44px;
  padding: 0 20px;
  border: 1px solid var(--line);
  border-radius: 6px;
  background: var(--card);
  color: var(--ink);
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s ease;
}

.cart-entry:hover {
  border-color: var(--brand);
  color: var(--brand);
}

.mall-nav {
  background: var(--brand);
}

.nav-inner {
  height: 42px;
  display: flex;
  align-items: center;
  gap: 30px;
}

.nav-inner a {
  color: rgba(255, 253, 248, 0.9);
  font-size: 14px;
  letter-spacing: 1px;
}

.nav-inner a:hover,
.nav-inner a.router-link-active {
  color: #fff;
  font-weight: 600;
}

.nav-note {
  margin-left: auto;
  color: rgba(255, 253, 248, 0.72);
  font-size: 12px;
  font-family: var(--font-serif);
  letter-spacing: 1px;
}

.mall-main {
  flex: 1;
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px 20px 48px;
}

.mall-footer {
  background: var(--ink);
  color: rgba(247, 239, 240, 0.75);
  padding: 36px 0 26px;
}

.footer-inner {
  display: flex;
  flex-direction: column;
  gap: 26px;
}

.footer-brand {
  display: flex;
  justify-content: center;
}

.promises {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  text-align: center;
  border-top: 1px solid rgba(201, 168, 106, 0.25);
  border-bottom: 1px solid rgba(201, 168, 106, 0.25);
  padding: 22px 0;
}

.promises div {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.promises b {
  font-family: var(--font-serif);
  color: var(--gold);
  font-size: 16px;
  letter-spacing: 2px;
}

.promises span {
  font-size: 12px;
}

.copyright {
  text-align: center;
  font-size: 12px;
  color: rgba(247, 239, 240, 0.45);
  letter-spacing: 1px;
}

@media (max-width: 900px) {
  .header-inner {
    height: auto;
    padding: 16px 0;
    flex-wrap: wrap;
    gap: 16px;
  }

  .search-area {
    order: 3;
    flex-basis: 100%;
    max-width: none;
  }

  .cart-entry {
    margin-left: auto;
  }

  .nav-inner {
    overflow-x: auto;
    gap: 20px;
  }

  .nav-note {
    display: none;
  }

  .promises {
    grid-template-columns: 1fr;
  }

  .mall-main {
    padding: 18px 14px 36px;
  }
}
</style>
