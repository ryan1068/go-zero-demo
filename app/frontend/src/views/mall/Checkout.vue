<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { listAddresses } from '@/api/user'
import { listCart } from '@/api/cart'
import { getProduct } from '@/api/product'
import { createOrder } from '@/api/order'
import { formatPrice } from '@/utils/format'
import type { AddressInfo } from '@/types/api'

interface CheckoutItem {
  product_id: number
  name: string
  cover: string
  price: number
  quantity: number
  total: number
}

const route = useRoute()
const router = useRouter()
const addresses = ref<AddressInfo[]>([])
const items = ref<CheckoutItem[]>([])
const loading = ref(false)
const submitting = ref(false)
const buyNow = Number(route.query.buy_now || 0) === 1
const productId = Number(route.query.product_id || 0)
const quantity = Number(route.query.quantity || 1)
const form = reactive({
  address_id: 0,
  remark: '',
})

const total = computed(() => items.value.reduce((sum, item) => sum + item.total, 0))
const totalQuantity = computed(() => items.value.reduce((sum, item) => sum + item.quantity, 0))

async function load() {
  loading.value = true
  try {
    const addressData = await listAddresses()
    addresses.value = addressData.list ?? []
    const defaultAddress = addresses.value.find((item) => item.is_default === 1) || addresses.value[0]
    if (defaultAddress) form.address_id = defaultAddress.id

    if (buyNow) {
      const product = await getProduct(productId)
      items.value = [{
        product_id: product.id,
        name: product.name,
        cover: product.cover_image,
        price: product.price,
        quantity,
        total: product.price * quantity,
      }]
    } else {
      const cart = await listCart()
      items.value = (cart.list ?? [])
        .filter((item) => item.selected === 1)
        .map((item) => ({
          product_id: item.product_id,
          name: item.product_name,
          cover: item.product_cover,
          price: item.price,
          quantity: item.quantity,
          total: item.price * item.quantity,
        }))
    }
  } finally {
    loading.value = false
  }
}

async function submit() {
  if (!form.address_id) {
    ElMessage.warning('请先选择收货地址')
    return
  }
  if (!items.value.length) {
    ElMessage.warning('没有可结算的商品')
    return
  }
  submitting.value = true
  try {
    const res = await createOrder({
      address_id: form.address_id,
      remark: form.remark.trim(),
      buy_now: buyNow ? 1 : 0,
      product_id: buyNow ? productId : 0,
      quantity: buyNow ? quantity : 0,
    })
    if (res.code === 200 && res.data) {
      ElMessage.success('订单创建成功')
      router.push(`/mall/order/${res.data.id}`)
    } else {
      ElMessage.error(res.msg || '下单失败')
    }
  } catch {
    ElMessage.error('下单失败')
  } finally {
    submitting.value = false
  }
}

onMounted(load)
</script>

<template>
  <div v-loading="loading" class="checkout-page">
    <div class="page-head">
      <h1>确认订单</h1>
      <p>CHECKOUT · 确认好物，收进时光</p>
    </div>

    <el-alert v-if="!addresses.length" type="warning" :closable="false" show-icon class="address-alert">
      还没有收货地址，请先
      <el-link type="primary" @click="router.push('/mall/addresses')">添加收货地址</el-link>
    </el-alert>

    <div class="checkout-grid">
      <div class="checkout-main">
        <section class="paper-card section-card">
          <div class="card-title">
            <span>收货地址</span>
            <el-link type="primary" @click="router.push('/mall/addresses')">管理地址</el-link>
          </div>
          <div class="address-grid">
            <button
              v-for="addr in addresses"
              :key="addr.id"
              type="button"
              class="address-card"
              :class="{ active: form.address_id === addr.id }"
              @click="form.address_id = addr.id"
            >
              <div class="addr-line">
                <b>{{ addr.receiver_name }}</b>
                <span>{{ addr.receiver_phone }}</span>
                <span v-if="addr.is_default === 1" class="seal-tag">默认</span>
              </div>
              <div class="addr-detail">{{ addr.province }}{{ addr.city }}{{ addr.district }}{{ addr.detail_address }}</div>
            </button>
          </div>
        </section>

        <section class="paper-card section-card">
          <div class="card-title"><span>商品清单</span><span class="muted">共 {{ totalQuantity }} 件</span></div>
          <div class="item-list">
            <div v-for="item in items" :key="item.product_id" class="item-row">
              <div class="cover-wrap">
                <img :src="item.cover || '/products/book.svg'" :alt="item.name" />
              </div>
              <div class="item-copy">
                <h3>{{ item.name }}</h3>
                <p>拾光严选 · 品质好物</p>
              </div>
              <div class="item-price price-text"><span class="price-symbol">¥</span>{{ formatPrice(item.price) }}</div>
              <div class="item-qty">× {{ item.quantity }}</div>
              <div class="item-total price-text"><span class="price-symbol">¥</span>{{ formatPrice(item.total) }}</div>
            </div>
          </div>
        </section>

        <section class="paper-card section-card">
          <div class="card-title"><span>订单备注</span></div>
          <el-input v-model="form.remark" type="textarea" :rows="3" maxlength="255" show-word-limit placeholder="选填，例如：请尽快发货" />
        </section>
      </div>

      <aside class="summary-card paper-card paper-card--strong">
        <div class="summary-title">订单汇总</div>
        <div class="summary-line"><span>商品件数</span><b>{{ totalQuantity }} 件</b></div>
        <div class="summary-line"><span>商品总额</span><b>¥{{ formatPrice(total) }}</b></div>
        <div class="summary-line"><span>运费</span><b>¥0.00</b></div>
        <div class="gold-line" />
        <div class="summary-total">
          <span>应付金额</span>
          <span class="price-text"><span class="price-symbol">¥</span>{{ formatPrice(total) }}</span>
        </div>
        <el-button type="primary" size="large" class="submit-btn" :loading="submitting" @click="submit">
          提交订单
        </el-button>
        <div class="summary-note">提交即表示同意《拾光严选交易服务协议》</div>
      </aside>
    </div>
  </div>
</template>

<style scoped>
.checkout-page {
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

.address-alert {
  border-radius: 8px;
}

.checkout-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 300px;
  gap: 20px;
  align-items: start;
}

.checkout-main {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.section-card {
  padding: 20px 22px;
}

.card-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
  font-family: var(--font-serif);
  font-size: 17px;
  font-weight: 700;
  letter-spacing: 2px;
}

.card-title .muted {
  color: var(--muted);
  font-size: 12px;
  font-family: var(--font-sans);
  font-weight: 400;
}

.address-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 12px;
}

.address-card {
  text-align: left;
  padding: 14px 16px;
  border: 1px solid var(--line);
  border-radius: 8px;
  background: #fbf6ee;
  cursor: pointer;
  font-family: var(--font-sans);
  transition: all 0.18s ease;
}

.address-card:hover {
  border-color: var(--brand);
}

.address-card.active {
  border-color: var(--brand);
  background: rgba(163, 59, 74, 0.06);
  box-shadow: 0 0 0 2px rgba(163, 59, 74, 0.1);
}

.addr-line {
  display: flex;
  align-items: center;
  gap: 10px;
  color: var(--ink);
  font-size: 14px;
}

.addr-line b {
  font-family: var(--font-serif);
  letter-spacing: 1px;
}

.addr-detail {
  margin-top: 8px;
  color: var(--muted);
  font-size: 12px;
  line-height: 1.6;
}

.item-list {
  display: flex;
  flex-direction: column;
}

.item-row {
  display: grid;
  grid-template-columns: 64px minmax(0, 1fr) 110px 70px 110px;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px dashed rgba(201, 168, 106, 0.26);
}

.item-row:last-child {
  border-bottom: none;
}

.cover-wrap {
  width: 64px;
  height: 64px;
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
  padding: 5px;
}

.item-copy h3 {
  margin: 0 0 5px;
  font-family: var(--font-serif);
  font-size: 15px;
  color: var(--ink);
}

.item-copy p {
  margin: 0;
  color: var(--muted);
  font-size: 12px;
}

.item-price,
.item-total {
  font-size: 16px;
}

.item-qty {
  color: var(--muted);
  font-size: 13px;
}

.summary-card {
  position: sticky;
  top: 160px;
  padding: 22px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.summary-title {
  font-family: var(--font-serif);
  font-size: 18px;
  font-weight: 700;
  letter-spacing: 3px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--line);
}

.summary-line {
  display: flex;
  justify-content: space-between;
  color: var(--ink-soft);
  font-size: 13px;
}

.summary-line b {
  color: var(--ink);
  font-weight: 500;
}

.summary-total {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  color: var(--ink-soft);
  font-size: 13px;
}

.summary-total .price-text {
  font-size: 28px;
}

.submit-btn {
  width: 100%;
  letter-spacing: 4px;
  margin-top: 4px;
}

.summary-note {
  color: var(--muted);
  font-size: 11px;
  text-align: center;
  line-height: 1.6;
}

@media (max-width: 900px) {
  .checkout-grid {
    grid-template-columns: 1fr;
  }

  .summary-card {
    position: static;
  }

  .item-row {
    grid-template-columns: 56px minmax(0, 1fr);
    grid-template-areas:
      'cover copy'
      'cover price'
      'cover meta';
    row-gap: 5px;
  }

  .cover-wrap {
    grid-area: cover;
    width: 56px;
    height: 56px;
  }

  .item-copy {
    grid-area: copy;
  }

  .item-price {
    grid-area: price;
  }

  .item-qty {
    grid-area: meta;
  }

  .item-total {
    display: none;
  }
}
</style>
