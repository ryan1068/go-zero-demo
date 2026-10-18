<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import BrandLogo from '@/components/BrandLogo.vue'
import ProductCard from '@/components/ProductCard.vue'
import { listCategories, listProducts } from '@/api/product'
import { formatPrice } from '@/utils/format'
import type { CategoryInfo, ProductInfo } from '@/types/api'

const router = useRouter()
const categories = ref<CategoryInfo[]>([])
const products = ref<ProductInfo[]>([])
const loading = ref(false)

const flashProducts = computed(() => products.value.slice(0, 3))
const featuredProducts = computed(() => products.value.slice(3))
const topCategories = computed(() => categories.value.filter((cat) => cat.parent_id === 0))
const childrenOf = (parentId: number) => categories.value.filter((cat) => cat.parent_id === parentId)

async function load() {
  loading.value = true
  try {
    const [cats, list] = await Promise.all([
      listCategories(),
      listProducts({ page: 1, page_size: 12, sort: 'default' }),
    ])
    categories.value = cats
    products.value = list.list ?? []
  } finally {
    loading.value = false
  }
}

function goCategory(id: number) {
  router.push({ path: '/mall/products', query: { category_id: id } })
}

onMounted(load)
</script>

<template>
  <div v-loading="loading" class="home">
    <section class="home-hero">
      <aside class="category-rail paper-card">
        <div class="rail-title">全部分类</div>
        <div v-for="cat in topCategories" :key="cat.id" class="rail-wrap">
          <a href="javascript:;" class="rail-item" @click="goCategory(cat.id)">
            <span>{{ cat.name }}</span>
            <span class="arrow">›</span>
          </a>
          <div v-if="childrenOf(cat.id).length" class="rail-flyout">
            <div class="flyout-head">{{ cat.name }}</div>
            <a
              v-for="child in childrenOf(cat.id)"
              :key="child.id"
              href="javascript:;"
              class="flyout-item"
              @click="goCategory(child.id)"
            >
              {{ child.name }}
            </a>
          </div>
        </div>
        <a href="javascript:;" class="rail-item" @click="router.push('/mall/products')">
          <span>全部好物</span>
          <span class="arrow">›</span>
        </a>
        <div class="rail-foot">拾光精选 · 品质可溯</div>
      </aside>

      <div class="hero-banner">
        <div class="banner-seal">拾光</div>
        <div class="banner-copy">
          <p class="eyebrow">SHIGUANG SELECT · 拾光精选</p>
          <h1>把好物<br />收进时光里</h1>
          <p class="desc">图书 · 数码 · 食品 · 生活好物<br />慢慢挑，认真选，买得安心。</p>
          <el-button size="large" class="banner-btn" @click="router.push('/mall/products')">
            立即逛逛
          </el-button>
        </div>
        <div class="banner-grain" />
      </div>

      <aside class="hero-side">
        <div class="side-card paper-card">
          <div class="side-title">今日推荐</div>
          <template v-if="products.length">
            <div class="today-product" @click="router.push(`/mall/product/${products[0].id}`)">
              <img :src="products[0].cover_image || '/products/book.svg'" :alt="products[0].name" />
              <div>
                <h3>{{ products[0].name }}</h3>
                <div class="price-text"><span class="price-symbol">¥</span>{{ formatPrice(products[0].price) }}</div>
              </div>
            </div>
          </template>
          <p v-else class="side-empty">今日好物正在上架</p>
          <div class="gold-line" />
          <div class="side-note">拾光严选 · 每日精选</div>
        </div>

        <div class="side-card paper-card service-card">
          <div class="side-title">服务承诺</div>
          <div class="service-item"><span class="seal-tag seal-tag--gold">正品</span><span>精选供应链，品质可溯</span></div>
          <div class="service-item"><span class="seal-tag seal-tag--gold">包邮</span><span>满额包邮，48 小时发货</span></div>
          <div class="service-item"><span class="seal-tag seal-tag--gold">售后</span><span>七天无理由退换</span></div>
        </div>
      </aside>
    </section>

    <section class="section">
      <div class="section-heading">
        <div>
          <h2>限时好物</h2>
          <p>LIMITED TIME · 今日限定好价</p>
        </div>
        <el-link type="primary" @click="router.push('/mall/products?sort=sales')">更多好物 ›</el-link>
      </div>
      <div class="product-grid">
        <ProductCard v-for="product in flashProducts" :key="product.id" :product="product" />
      </div>
    </section>

    <section class="brand-story">
      <BrandLogo size="lg" subtitle="把好物收进时光里" light />
      <div class="story-copy">
        <h2>拾光，是一种认真生活的态度</h2>
        <p>我们相信，好的东西值得被慢慢挑选，也值得被长久使用。拾光严选从图书、数码到食品，只把经得起时间的好物带给你。</p>
      </div>
      <div class="story-seal">拾光<br />严选</div>
    </section>

    <section class="section">
      <div class="section-heading">
        <div>
          <h2>拾光精选</h2>
          <p>SHIGUANG PICKS · 认真挑过的每一件</p>
        </div>
        <el-link type="primary" @click="router.push('/mall/products')">查看全部 ›</el-link>
      </div>
      <el-empty v-if="!featuredProducts.length" description="好物正在上架中" />
      <div v-else class="product-grid">
        <ProductCard v-for="product in featuredProducts" :key="product.id" :product="product" />
      </div>
    </section>
  </div>
</template>

<style scoped>
.home {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.home-hero {
  display: grid;
  grid-template-columns: 210px minmax(0, 1fr) 270px;
  gap: 16px;
}

.category-rail {
  position: relative;
  z-index: 5;
  padding: 14px 12px;
  display: flex;
  flex-direction: column;
}

.rail-wrap {
  position: relative;
}

.rail-flyout {
  display: none;
  position: absolute;
  left: 100%;
  top: -1px;
  width: 240px;
  padding: 12px 14px;
  background: var(--card);
  border: 1px solid var(--line);
  border-radius: 10px;
  box-shadow: var(--shadow);
  z-index: 30;
}

.rail-wrap:hover .rail-flyout {
  display: block;
}

.flyout-head {
  padding-bottom: 8px;
  margin-bottom: 6px;
  border-bottom: 1px solid var(--line);
  font-family: var(--font-serif);
  font-size: 15px;
  font-weight: 700;
  letter-spacing: 2px;
  color: var(--ink);
}

.flyout-item {
  display: block;
  padding: 7px 6px;
  border-radius: 4px;
  color: var(--ink-soft);
  font-size: 13px;
}

.flyout-item:hover {
  color: var(--brand);
  background: rgba(163, 59, 74, 0.05);
}

.rail-title {
  font-family: var(--font-serif);
  font-size: 17px;
  font-weight: 700;
  letter-spacing: 2px;
  padding: 4px 10px 12px;
  border-bottom: 1px solid var(--line);
}

.rail-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 11px 10px;
  color: var(--ink-soft);
  font-size: 14px;
  border-bottom: 1px dashed rgba(201, 168, 106, 0.22);
}

.rail-item:hover {
  color: var(--brand);
  background: rgba(163, 59, 74, 0.04);
}

.arrow {
  color: var(--gold);
}

.rail-foot {
  margin-top: auto;
  padding: 16px 10px 4px;
  color: var(--muted);
  font-size: 11px;
  letter-spacing: 1px;
}

.hero-banner {
  position: relative;
  overflow: hidden;
  min-height: 330px;
  border-radius: 12px;
  padding: 42px 48px;
  color: #f7eff0;
  background:
    radial-gradient(circle at 82% 22%, rgba(201, 168, 106, 0.32), transparent 34%),
    linear-gradient(135deg, #7a2733 0%, #a33b4a 48%, #5c1c26 100%);
  box-shadow: var(--shadow);
}

.banner-copy {
  position: relative;
  z-index: 2;
  max-width: 420px;
}

.eyebrow {
  margin: 0 0 12px;
  font-size: 11px;
  letter-spacing: 4px;
  color: rgba(247, 239, 240, 0.72);
}

.hero-banner h1 {
  margin: 0;
  font-family: var(--font-serif);
  font-size: 42px;
  line-height: 1.28;
  letter-spacing: 6px;
  font-weight: 700;
}

.desc {
  margin: 18px 0 26px;
  font-size: 13px;
  line-height: 1.9;
  color: rgba(247, 239, 240, 0.8);
}

.banner-btn {
  background: #f7eff0;
  border-color: #f7eff0;
  color: var(--brand-deep);
  font-family: var(--font-serif);
  letter-spacing: 4px;
}

.banner-btn:hover {
  background: #fff;
  border-color: #fff;
  color: var(--brand-deep);
}

.banner-seal {
  position: absolute;
  right: 32px;
  bottom: 26px;
  width: 104px;
  height: 104px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid rgba(247, 239, 240, 0.4);
  border-radius: 14px;
  color: rgba(247, 239, 240, 0.48);
  font-family: var(--font-serif);
  font-size: 30px;
  letter-spacing: 8px;
  transform: rotate(-6deg);
}

.banner-grain {
  position: absolute;
  inset: 0;
  opacity: 0.18;
  background-image: radial-gradient(rgba(255, 255, 255, 0.55) 0.6px, transparent 0.6px);
  background-size: 7px 7px;
  pointer-events: none;
}

.hero-side {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.side-card {
  padding: 16px;
}

.side-title {
  font-family: var(--font-serif);
  font-size: 16px;
  font-weight: 700;
  letter-spacing: 2px;
  margin-bottom: 12px;
}

.today-product {
  display: flex;
  gap: 10px;
  align-items: center;
  cursor: pointer;
}

.today-product img {
  width: 66px;
  height: 66px;
  object-fit: contain;
  background: #fbf6ee;
  border-radius: 6px;
}

.today-product h3 {
  margin: 0 0 6px;
  font-family: var(--font-serif);
  font-size: 14px;
  color: var(--ink);
  line-height: 1.4;
}

.today-product .price-text {
  font-size: 20px;
}

.side-empty {
  color: var(--muted);
  font-size: 12px;
}

.side-note {
  margin-top: 12px;
  color: var(--muted);
  font-size: 11px;
  letter-spacing: 1px;
}

.service-card {
  flex: 1;
}

.service-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 9px 0;
  color: var(--ink-soft);
  font-size: 12px;
  border-bottom: 1px dashed rgba(201, 168, 106, 0.22);
}

.service-item:last-child {
  border-bottom: none;
}

.section {
  min-height: 40px;
}

.product-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}

.brand-story {
  display: flex;
  align-items: center;
  gap: 34px;
  padding: 30px 38px;
  border-radius: 12px;
  background:
    radial-gradient(circle at 88% 20%, rgba(163, 59, 74, 0.24), transparent 30%),
    var(--ink);
  color: rgba(247, 239, 240, 0.82);
  box-shadow: var(--shadow);
}

.story-copy {
  flex: 1;
}

.story-copy h2 {
  margin: 0 0 10px;
  font-family: var(--font-serif);
  color: var(--gold);
  font-size: 22px;
  letter-spacing: 2px;
}

.story-copy p {
  margin: 0;
  line-height: 1.9;
  font-size: 13px;
}

.story-seal {
  width: 92px;
  height: 92px;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  border: 1px solid rgba(201, 168, 106, 0.55);
  border-radius: 10px;
  color: var(--gold);
  font-family: var(--font-serif);
  font-size: 20px;
  line-height: 1.5;
  letter-spacing: 4px;
  transform: rotate(-4deg);
}

@media (max-width: 1050px) {
  .home-hero {
    grid-template-columns: 190px minmax(0, 1fr);
  }

  .hero-side {
    grid-column: 1 / -1;
    flex-direction: row;
  }

  .hero-side .side-card {
    flex: 1;
  }
}

@media (max-width: 900px) {
  .product-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 720px) {
  .home-hero {
    grid-template-columns: 1fr;
  }

  .product-grid {
    grid-template-columns: 1fr;
  }

  .category-rail {
    flex-direction: row;
    flex-wrap: wrap;
    gap: 4px;
  }

  .rail-flyout {
    display: none !important;
  }

  .rail-title,
  .rail-foot {
    flex-basis: 100%;
  }

  .rail-item {
    border-bottom: none;
    padding: 8px 10px;
  }

  .hero-banner {
    min-height: 280px;
    padding: 30px 24px;
  }

  .hero-banner h1 {
    font-size: 32px;
  }

  .banner-seal {
    width: 72px;
    height: 72px;
    font-size: 20px;
  }

  .hero-side {
    flex-direction: column;
  }

  .brand-story {
    flex-direction: column;
    text-align: center;
    padding: 26px 20px;
  }
}
</style>
