<script setup lang="ts">
import { useRouter } from 'vue-router'
import { formatPrice } from '@/utils/format'
import type { ProductInfo } from '@/types/api'

const props = defineProps<{ product: ProductInfo }>()
const router = useRouter()
</script>

<template>
  <article class="product-card hover-lift" @click="router.push(`/mall/product/${props.product.id}`)">
    <div class="image-wrap">
      <img :src="props.product.cover_image || '/products/book.svg'" :alt="props.product.name" loading="lazy" />
      <span v-if="props.product.is_recommend === 1" class="seal-tag corner">拾光精选</span>
      <span v-if="props.product.stock > 0 && props.product.stock <= 10" class="seal-tag seal-tag--gold corner right">
        仅剩 {{ props.product.stock }}
      </span>
      <span v-if="props.product.stock <= 0" class="sold-out">已售罄</span>
    </div>
    <div class="body">
      <h3>{{ props.product.name }}</h3>
      <p class="subtitle">{{ props.product.subtitle || '拾光严选 · 品质好物' }}</p>
      <div class="meta">
        <div class="price price-text">
          <span class="price-symbol">¥</span>{{ formatPrice(props.product.price) }}
          <span v-if="props.product.original_price > props.product.price" class="price-original">
            ¥{{ formatPrice(props.product.original_price) }}
          </span>
        </div>
        <span class="sales">已售 {{ props.product.sales }}</span>
      </div>
      <div class="tags">
        <span class="seal-tag seal-tag--gold">包邮</span>
        <span class="seal-tag">正品</span>
        <span v-if="props.product.category_name" class="category">{{ props.product.category_name }}</span>
      </div>
    </div>
  </article>
</template>

<style scoped>
.product-card {
  background: var(--card);
  border: 1px solid var(--line);
  border-radius: 10px;
  overflow: hidden;
  cursor: pointer;
}

.image-wrap {
  position: relative;
  height: 190px;
  display: flex;
  align-items: center;
  justify-content: center;
  background:
    radial-gradient(circle at 50% 40%, rgba(201, 168, 106, 0.14), transparent 60%),
    #fbf6ee;
}

.image-wrap img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  padding: 14px;
}

.corner {
  position: absolute;
  top: 10px;
  left: 10px;
}

.corner.right {
  left: auto;
  right: 10px;
}

.sold-out {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(46, 35, 38, 0.55);
  color: #f7eff0;
  font-family: var(--font-serif);
  font-size: 20px;
  letter-spacing: 4px;
}

.body {
  padding: 14px 14px 15px;
}

h3 {
  margin: 0;
  font-family: var(--font-serif);
  font-size: 17px;
  font-weight: 700;
  color: var(--ink);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.subtitle {
  margin: 6px 0 10px;
  color: var(--muted);
  font-size: 12px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.meta {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 8px;
}

.price {
  font-size: 22px;
}

.sales {
  color: var(--muted);
  font-size: 11px;
}

.tags {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 11px;
}

.category {
  margin-left: auto;
  color: var(--muted);
  font-size: 11px;
}
</style>
