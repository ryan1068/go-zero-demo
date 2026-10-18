<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getProduct } from '@/api/product'
import { addCart } from '@/api/cart'
import { formatPrice } from '@/utils/format'
import { useUserStore } from '@/stores/user'
import type { ProductInfo } from '@/types/api'

const route = useRoute()
const router = useRouter()
const store = useUserStore()
const product = ref<ProductInfo | null>(null)
const quantity = ref(1)
const loading = ref(false)
const activeImage = ref('')

async function load() {
  loading.value = true
  try {
    product.value = await getProduct(Number(route.params.id))
    activeImage.value = product.value.cover_image || product.value.images?.[0] || ''
  } catch {
    ElMessage.error('商品不存在或已下架')
    router.push('/mall/products')
  } finally {
    loading.value = false
  }
}

function requireLogin(): boolean {
  if (!store.isLoggedIn) {
    router.push({ path: '/login', query: { redirect: route.fullPath } })
    return false
  }
  return true
}

async function handleAddCart() {
  if (!requireLogin() || !product.value) return
  try {
    const res = await addCart(product.value.id, quantity.value)
    if (res.code === 200) ElMessage.success('已加入购物车')
    else ElMessage.error(res.msg || '加入购物车失败')
  } catch {
    ElMessage.error('加入购物车失败')
  }
}

function handleBuyNow() {
  if (!requireLogin() || !product.value) return
  router.push({
    path: '/mall/checkout',
    query: { buy_now: 1, product_id: product.value.id, quantity: quantity.value },
  })
}

onMounted(load)
</script>

<template>
  <div v-loading="loading" class="product-detail-page">
    <div class="breadcrumb">
      <a @click="router.push('/mall')">首页</a>
      <span>/</span>
      <a @click="router.push('/mall/products')">全部好物</a>
      <span v-if="product">/ {{ product.name }}</span>
    </div>

    <section v-if="product" class="detail-card paper-card paper-card--strong">
      <div class="gallery">
        <div class="main-image">
          <img :src="activeImage || '/products/book.svg'" :alt="product.name" />
          <span v-if="product.is_recommend === 1" class="seal-tag main-tag">拾光精选</span>
        </div>
        <div v-if="product.images?.length" class="thumbs">
          <button
            v-for="img in product.images"
            :key="img"
            type="button"
            class="thumb"
            :class="{ active: img === activeImage }"
            @click="activeImage = img"
          >
            <img :src="img" :alt="product.name" />
          </button>
        </div>
      </div>

      <div class="info">
        <p class="eyebrow">SHIGUANG SELECT · 拾光精选</p>
        <h1>{{ product.name }}</h1>
        <p class="subtitle">{{ product.subtitle || '拾光严选 · 品质好物' }}</p>

        <div class="price-panel">
          <div class="price-main">
            <span class="label">拾光价</span>
            <span class="price-text"><span class="price-symbol">¥</span>{{ formatPrice(product.price) }}</span>
            <span v-if="product.original_price > product.price" class="price-original">
              ¥{{ formatPrice(product.original_price) }}
            </span>
          </div>
          <div class="price-meta">
            <span>已售 {{ product.sales }} 件</span>
            <span>库存 {{ product.stock }} {{ product.unit }}</span>
          </div>
        </div>

        <div class="tags">
          <span class="seal-tag seal-tag--gold">正品保障</span>
          <span class="seal-tag seal-tag--gold">满额包邮</span>
          <span class="seal-tag">七天无理由</span>
        </div>

        <div class="quantity-row">
          <span class="field-label">数量</span>
          <el-input-number v-model="quantity" :min="1" :max="Math.min(product.stock, 99)" />
          <span class="stock-hint">{{ product.stock > 0 ? `库存充足` : '暂时缺货' }}</span>
        </div>

        <div class="actions">
          <el-button size="large" class="cart-btn" :disabled="product.stock <= 0" @click="handleAddCart">
            加入购物车
          </el-button>
          <el-button type="primary" size="large" class="buy-btn" :disabled="product.stock <= 0" @click="handleBuyNow">
            立即购买
          </el-button>
        </div>

        <div class="service-list">
          <div><span class="seal-tag seal-tag--gold">发货</span> 下单后 48 小时内发出</div>
          <div><span class="seal-tag seal-tag--gold">售后</span> 七天无理由退换，售后无忧</div>
          <div><span class="seal-tag seal-tag--gold">品质</span> 拾光精选供应链，品质可溯</div>
        </div>
      </div>
    </section>

    <section v-if="product" class="detail-body paper-card">
      <div class="section-heading">
        <div>
          <h2>商品详情</h2>
          <p>DETAILS · 拾光严选</p>
        </div>
      </div>
      <p class="detail-text">{{ product.detail || '暂无更多详情，拾光严选正在补充中。' }}</p>
    </section>
  </div>
</template>

<style scoped>
.product-detail-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.breadcrumb {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--muted);
  font-size: 12px;
}

.breadcrumb a {
  color: var(--muted);
  cursor: pointer;
}

.breadcrumb a:hover {
  color: var(--brand);
}

.detail-card {
  display: grid;
  grid-template-columns: 460px minmax(0, 1fr);
  gap: 40px;
  padding: 30px;
}

.gallery {
  min-width: 0;
}

.main-image {
  position: relative;
  height: 420px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  border: 1px solid var(--line);
  background:
    radial-gradient(circle at 50% 42%, rgba(201, 168, 106, 0.15), transparent 62%),
    #fbf6ee;
}

.main-image img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  padding: 26px;
}

.main-tag {
  position: absolute;
  top: 14px;
  left: 14px;
}

.thumbs {
  display: flex;
  gap: 10px;
  margin-top: 12px;
}

.thumb {
  width: 64px;
  height: 64px;
  padding: 4px;
  border: 1px solid var(--line);
  border-radius: 6px;
  background: #fbf6ee;
  cursor: pointer;
}

.thumb.active {
  border-color: var(--brand);
  box-shadow: 0 0 0 2px rgba(163, 59, 74, 0.12);
}

.thumb img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.info {
  min-width: 0;
}

.eyebrow {
  margin: 0 0 10px;
  color: var(--gold);
  font-size: 11px;
  letter-spacing: 3px;
}

.info h1 {
  margin: 0;
  font-family: var(--font-serif);
  font-size: 32px;
  line-height: 1.35;
  letter-spacing: 2px;
  color: var(--ink);
}

.subtitle {
  margin: 10px 0 20px;
  color: var(--muted);
  font-size: 14px;
}

.price-panel {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 16px;
  padding: 18px 22px;
  border-radius: 8px;
  color: #f7eff0;
  background:
    radial-gradient(circle at 88% 20%, rgba(201, 168, 106, 0.32), transparent 34%),
    linear-gradient(135deg, #7a2733, #a33b4a 55%, #5c1c26);
}

.price-main {
  display: flex;
  align-items: baseline;
  gap: 10px;
}

.label {
  font-size: 12px;
  opacity: 0.78;
  letter-spacing: 2px;
}

.price-main .price-text {
  color: #fffefa;
  font-size: 38px;
}

.price-main .price-original {
  color: rgba(247, 239, 240, 0.55);
}

.price-meta {
  display: flex;
  flex-direction: column;
  gap: 6px;
  text-align: right;
  font-size: 12px;
  opacity: 0.82;
}

.tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin: 18px 0;
}

.quantity-row {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 0;
  border-top: 1px dashed rgba(201, 168, 106, 0.3);
  border-bottom: 1px dashed rgba(201, 168, 106, 0.3);
}

.field-label {
  color: var(--muted);
  font-size: 13px;
}

.stock-hint {
  color: var(--muted);
  font-size: 12px;
}

.actions {
  display: flex;
  gap: 12px;
  margin: 22px 0;
}

.actions .el-button {
  min-width: 160px;
  letter-spacing: 2px;
}

.cart-btn {
  border-color: var(--brand);
  color: var(--brand);
  background: rgba(163, 59, 74, 0.06);
}

.cart-btn:hover {
  background: rgba(163, 59, 74, 0.12);
  border-color: var(--brand);
  color: var(--brand);
}

.service-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  color: var(--ink-soft);
  font-size: 12px;
}

.service-list div {
  display: flex;
  align-items: center;
  gap: 8px;
}

.detail-body {
  padding: 24px 30px 34px;
}

.detail-text {
  margin: 0;
  color: var(--ink-soft);
  line-height: 2;
  font-size: 14px;
  white-space: pre-wrap;
}

@media (max-width: 900px) {
  .detail-card {
    grid-template-columns: 1fr;
    gap: 24px;
    padding: 20px;
  }

  .main-image {
    height: 300px;
  }

  .info h1 {
    font-size: 26px;
  }

  .price-panel {
    flex-direction: column;
    align-items: flex-start;
  }

  .price-meta {
    text-align: left;
  }
}
</style>
