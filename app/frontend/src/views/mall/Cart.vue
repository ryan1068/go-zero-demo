<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { deleteCart, listCart, selectCart, updateCart } from '@/api/cart'
import { formatPrice } from '@/utils/format'
import type { CartItemInfo } from '@/types/api'

const router = useRouter()
const items = ref<CartItemInfo[]>([])
const loading = ref(false)

const selectedItems = computed(() => items.value.filter((item) => item.selected === 1))
const selectedTotal = computed(() => selectedItems.value.reduce((sum, item) => sum + item.price * item.quantity, 0))
const allSelected = computed(() => items.value.length > 0 && items.value.every((item) => item.selected === 1))

async function load() {
  loading.value = true
  try {
    const data = await listCart()
    items.value = data.list ?? []
  } finally {
    loading.value = false
  }
}

async function changeQuantity(item: CartItemInfo, quantity: number) {
  if (!quantity || quantity < 1) return
  const res = await updateCart(item.id, quantity, item.selected)
  if (res.code !== 200) {
    ElMessage.error(res.msg || '修改数量失败')
    await load()
    return
  }
  item.quantity = quantity
}

async function changeSelected(item: CartItemInfo) {
  const res = await updateCart(item.id, item.quantity, item.selected)
  if (res.code !== 200) {
    ElMessage.error(res.msg || '修改选中状态失败')
    await load()
  }
}

async function toggleAll() {
  const target = allSelected.value ? 0 : 1
  const res = await selectCart(target)
  if (res.code !== 200) {
    ElMessage.error(res.msg || '操作失败')
    return
  }
  await load()
}

async function removeItem(item: CartItemInfo) {
  const res = await deleteCart(item.id)
  if (res.code === 200) {
    ElMessage.success('已删除')
    await load()
  } else {
    ElMessage.error(res.msg || '删除失败')
  }
}

function checkout() {
  if (!selectedItems.value.length) {
    ElMessage.warning('请先选择商品')
    return
  }
  router.push('/mall/checkout')
}

onMounted(load)
</script>

<template>
  <div v-loading="loading" class="cart-page">
    <div class="page-head">
      <h1>购物车</h1>
      <p>SHOPPING CART · 把喜欢的好物收进来</p>
    </div>

    <section class="cart-card paper-card">
      <div v-if="items.length" class="cart-head">
        <el-checkbox :model-value="allSelected" @change="toggleAll">全选</el-checkbox>
        <span class="col-product">商品</span>
        <span class="col-price">单价</span>
        <span class="col-quantity">数量</span>
        <span class="col-subtotal">小计</span>
        <span class="col-action">操作</span>
      </div>

      <el-empty v-if="!items.length" description="购物车还是空的，去挑几件喜欢的好物吧">
        <el-button type="primary" @click="router.push('/mall/products')">去逛逛</el-button>
      </el-empty>

      <div v-for="item in items" :key="item.id" class="cart-row">
        <el-checkbox
          v-model="item.selected"
          :true-value="1"
          :false-value="0"
          @change="changeSelected(item)"
        />
        <div class="product-cell" @click="router.push(`/mall/product/${item.product_id}`)">
          <div class="cover-wrap">
            <img :src="item.product_cover || '/products/book.svg'" :alt="item.product_name" />
          </div>
          <div class="product-copy">
            <h3>{{ item.product_name }}</h3>
            <p>拾光严选 · 品质好物</p>
            <span class="seal-tag seal-tag--gold">包邮</span>
          </div>
        </div>
        <div class="col-price price-text"><span class="price-symbol">¥</span>{{ formatPrice(item.price) }}</div>
        <div class="col-quantity">
          <el-input-number
            v-model="item.quantity"
            :min="1"
            :max="Math.min(item.stock, 99)"
            size="small"
            @change="(value: number | undefined) => changeQuantity(item, value || 1)"
          />
        </div>
        <div class="col-subtotal price-text"><span class="price-symbol">¥</span>{{ formatPrice(item.price * item.quantity) }}</div>
        <div class="col-action">
          <el-button text type="danger" @click="removeItem(item)">删除</el-button>
        </div>
      </div>

      <div v-if="items.length" class="cart-foot">
        <div class="foot-left">
          <el-checkbox :model-value="allSelected" @change="toggleAll">全选</el-checkbox>
          <span class="selected-count">已选 {{ selectedItems.length }} 件</span>
        </div>
        <div class="foot-right">
          <div class="total">
            合计：<span class="price-text"><span class="price-symbol">¥</span>{{ formatPrice(selectedTotal) }}</span>
          </div>
          <el-button type="primary" size="large" class="checkout-btn" @click="checkout">去结算</el-button>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.cart-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-head h1 {
  margin: 0;
  font-family: var(--font-serif);
  font-size: 30px;
  letter-spacing: 4px;
}

.page-head p {
  margin: 6px 0 0;
  color: var(--muted);
  font-size: 12px;
  letter-spacing: 2px;
}

.cart-card {
  padding: 0 22px 22px;
}

.cart-head,
.cart-row {
  display: grid;
  grid-template-columns: 52px minmax(260px, 1fr) 120px 150px 120px 80px;
  align-items: center;
  gap: 12px;
}

.cart-head {
  height: 50px;
  color: var(--muted);
  font-size: 12px;
  border-bottom: 1px solid var(--line);
}

.cart-row {
  padding: 18px 0;
  border-bottom: 1px dashed rgba(201, 168, 106, 0.26);
}

.product-cell {
  display: flex;
  align-items: center;
  gap: 14px;
  cursor: pointer;
  min-width: 0;
}

.cover-wrap {
  width: 76px;
  height: 76px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--line);
  border-radius: 6px;
  background: #fbf6ee;
}

.cover-wrap img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  padding: 6px;
}

.product-copy {
  min-width: 0;
}

.product-copy h3 {
  margin: 0 0 6px;
  font-family: var(--font-serif);
  font-size: 16px;
  color: var(--ink);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.product-copy p {
  margin: 0 0 8px;
  color: var(--muted);
  font-size: 12px;
}

.col-price,
.col-subtotal {
  font-size: 18px;
}

.col-action {
  text-align: right;
}

.cart-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding-top: 20px;
}

.foot-left {
  display: flex;
  align-items: center;
  gap: 18px;
}

.selected-count {
  color: var(--muted);
  font-size: 12px;
}

.foot-right {
  display: flex;
  align-items: center;
  gap: 22px;
}

.total {
  color: var(--ink-soft);
  font-size: 14px;
}

.total .price-text {
  font-size: 28px;
}

.checkout-btn {
  min-width: 150px;
  letter-spacing: 4px;
}

@media (max-width: 900px) {
  .cart-head {
    display: none;
  }

  .cart-row {
    grid-template-columns: 32px minmax(0, 1fr);
    grid-template-areas:
      'check product'
      '. price'
      '. quantity'
      '. subtotal'
      '. action';
    row-gap: 10px;
  }

  .cart-row > :nth-child(1) {
    grid-area: check;
  }

  .product-cell {
    grid-area: product;
  }

  .col-price {
    grid-area: price;
  }

  .col-quantity {
    grid-area: quantity;
  }

  .col-subtotal {
    grid-area: subtotal;
  }

  .col-action {
    grid-area: action;
    text-align: left;
  }

  .cart-foot {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
