<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { cancelOrder, confirmOrder, getOrder, payOrder } from '@/api/order'
import { formatPrice, formatTime, orderStatusText, orderStatusType } from '@/utils/format'
import type { OrderInfo } from '@/types/api'

const route = useRoute()
const router = useRouter()
const order = ref<OrderInfo | null>(null)
const loading = ref(false)

const addressText = computed(() => {
  if (!order.value) return ''
  const o = order.value
  return `${o.receiver_province}${o.receiver_city}${o.receiver_district}${o.receiver_address}`
})

async function load() {
  loading.value = true
  try {
    order.value = await getOrder(Number(route.params.id))
  } catch {
    ElMessage.error('订单不存在或无权访问')
    router.push('/mall/orders')
  } finally {
    loading.value = false
  }
}

async function pay() {
  if (!order.value) return
  const res = await payOrder(order.value.id)
  if (res.code === 200) {
    ElMessage.success('支付成功')
    await load()
  } else {
    ElMessage.error(res.msg || '支付失败')
  }
}

async function cancel() {
  if (!order.value) return
  try {
    await ElMessageBox.confirm('确定取消这个订单吗？', '提示', { type: 'warning' })
  } catch {
    return
  }
  const res = await cancelOrder(order.value.id, '用户取消')
  if (res.code === 200) {
    ElMessage.success('订单已取消')
    await load()
  } else {
    ElMessage.error(res.msg || '取消失败')
  }
}

async function confirm() {
  if (!order.value) return
  const res = await confirmOrder(order.value.id)
  if (res.code === 200) {
    ElMessage.success('已确认收货')
    await load()
  } else {
    ElMessage.error(res.msg || '操作失败')
  }
}

onMounted(load)
</script>

<template>
  <div v-loading="loading" class="order-detail-page">
    <div class="breadcrumb">
      <a @click="router.push('/mall')">首页</a>
      <span>/</span>
      <a @click="router.push('/mall/orders')">我的订单</a>
      <span v-if="order">/ {{ order.order_no }}</span>
    </div>

    <template v-if="order">
      <section class="status-hero paper-card">
        <div>
          <p class="eyebrow">ORDER STATUS · 订单状态</p>
          <h1>{{ orderStatusText(order.status) }}</h1>
          <p class="order-no">订单号：{{ order.order_no }}</p>
        </div>
        <div class="hero-actions">
          <el-tag :type="orderStatusType(order.status)" size="large" effect="plain">{{ orderStatusText(order.status) }}</el-tag>
          <el-button v-if="order.status === 10" type="primary" @click="pay">去支付</el-button>
          <el-button v-if="order.status === 10 || order.status === 20" @click="cancel">取消订单</el-button>
          <el-button v-if="order.status === 30" type="primary" @click="confirm">确认收货</el-button>
          <el-button @click="router.push('/mall/orders')">返回订单列表</el-button>
        </div>
      </section>

      <div class="detail-grid">
        <div class="detail-main">
          <section class="paper-card section-card">
            <div class="card-title">收货信息</div>
            <div class="receiver">
              <div class="receiver-line">
                <b>{{ order.receiver_name }}</b>
                <span>{{ order.receiver_phone }}</span>
                <span class="seal-tag seal-tag--gold">收货</span>
              </div>
              <p>{{ addressText }}</p>
              <p v-if="order.remark" class="remark">备注：{{ order.remark }}</p>
            </div>
          </section>

          <section class="paper-card section-card">
            <div class="card-title"><span>商品明细</span><span class="muted">共 {{ order.items.length }} 件</span></div>
            <div class="item-list">
              <div v-for="item in order.items" :key="item.id" class="item-row">
                <div class="cover-wrap">
                  <img :src="item.product_cover || '/products/book.svg'" :alt="item.product_name" />
                </div>
                <div class="item-copy">
                  <h3>{{ item.product_name }}</h3>
                  <p>¥{{ formatPrice(item.price) }} × {{ item.quantity }}</p>
                </div>
                <div class="item-total price-text"><span class="price-symbol">¥</span>{{ formatPrice(item.total_amount) }}</div>
              </div>
            </div>
          </section>

          <section class="paper-card section-card">
            <div class="card-title">物流信息</div>
            <div class="logistics">
              <div><span class="label">物流公司</span><b>{{ order.logistics_company || '待发货' }}</b></div>
              <div><span class="label">物流单号</span><b>{{ order.tracking_no || '-' }}</b></div>
              <div><span class="label">发货时间</span><b>{{ formatTime(order.shipped_at) }}</b></div>
              <div><span class="label">完成时间</span><b>{{ formatTime(order.completed_at) }}</b></div>
            </div>
          </section>
        </div>

        <aside class="summary-card paper-card paper-card--strong">
          <div class="summary-title">订单金额</div>
          <div class="summary-line"><span>商品总额</span><b>¥{{ formatPrice(order.total_amount) }}</b></div>
          <div class="summary-line"><span>运费</span><b>¥{{ formatPrice(order.freight_amount) }}</b></div>
          <div class="summary-line"><span>优惠</span><b>¥{{ formatPrice(order.discount_amount) }}</b></div>
          <div class="gold-line" />
          <div class="summary-total">
            <span>实付金额</span>
            <span class="price-text"><span class="price-symbol">¥</span>{{ formatPrice(order.pay_amount) }}</span>
          </div>
          <div class="timeline">
            <div><span>创建时间</span><b>{{ formatTime(order.created_at) }}</b></div>
            <div><span>支付时间</span><b>{{ formatTime(order.paid_at) }}</b></div>
            <div><span>发货时间</span><b>{{ formatTime(order.shipped_at) }}</b></div>
            <div><span>完成时间</span><b>{{ formatTime(order.completed_at) }}</b></div>
          </div>
        </aside>
      </div>
    </template>
  </div>
</template>

<style scoped>
.order-detail-page {
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

.status-hero {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  padding: 28px 32px;
  background:
    radial-gradient(circle at 90% 20%, rgba(201, 168, 106, 0.2), transparent 35%),
    var(--card);
}

.eyebrow {
  margin: 0 0 8px;
  color: var(--gold);
  font-size: 11px;
  letter-spacing: 3px;
}

.status-hero h1 {
  margin: 0;
  font-family: var(--font-serif);
  font-size: 34px;
  letter-spacing: 5px;
  color: var(--ink);
}

.order-no {
  margin: 10px 0 0;
  color: var(--muted);
  font-size: 12px;
}

.hero-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  justify-content: flex-end;
}

.detail-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 300px;
  gap: 20px;
  align-items: start;
}

.detail-main {
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
  margin-bottom: 14px;
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

.receiver-line {
  display: flex;
  align-items: center;
  gap: 12px;
  color: var(--ink);
}

.receiver-line b {
  font-family: var(--font-serif);
  font-size: 16px;
  letter-spacing: 1px;
}

.receiver p {
  margin: 8px 0 0;
  color: var(--muted);
  font-size: 13px;
  line-height: 1.7;
}

.receiver .remark {
  color: var(--ink-soft);
}

.item-list {
  display: flex;
  flex-direction: column;
}

.item-row {
  display: grid;
  grid-template-columns: 64px minmax(0, 1fr) 110px;
  align-items: center;
  gap: 14px;
  padding: 12px 0;
  border-bottom: 1px dashed rgba(201, 168, 106, 0.24);
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
  margin: 0 0 6px;
  font-family: var(--font-serif);
  font-size: 15px;
}

.item-copy p {
  margin: 0;
  color: var(--muted);
  font-size: 12px;
}

.item-total {
  text-align: right;
  font-size: 17px;
}

.logistics {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.logistics div {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.logistics .label {
  color: var(--muted);
  font-size: 12px;
}

.logistics b {
  color: var(--ink);
  font-weight: 500;
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

.timeline {
  display: flex;
  flex-direction: column;
  gap: 9px;
  margin-top: 4px;
  padding-top: 12px;
  border-top: 1px dashed rgba(201, 168, 106, 0.28);
}

.timeline div {
  display: flex;
  justify-content: space-between;
  gap: 10px;
  font-size: 12px;
}

.timeline span {
  color: var(--muted);
}

.timeline b {
  color: var(--ink-soft);
  font-weight: 400;
}

@media (max-width: 900px) {
  .status-hero {
    align-items: flex-start;
    flex-direction: column;
  }

  .hero-actions {
    justify-content: flex-start;
  }

  .detail-grid {
    grid-template-columns: 1fr;
  }

  .summary-card {
    position: static;
  }

  .logistics {
    grid-template-columns: 1fr;
  }
}
</style>
