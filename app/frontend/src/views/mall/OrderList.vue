<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { cancelOrder, confirmOrder, listOrders, payOrder } from '@/api/order'
import { formatPrice, formatTime, orderStatusText, orderStatusType } from '@/utils/format'
import type { OrderInfo } from '@/types/api'

const router = useRouter()
const orders = ref<OrderInfo[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const status = ref(0)
const loading = ref(false)

const statusTabs = [
  { label: '全部', value: 0 },
  { label: '待支付', value: 10 },
  { label: '待发货', value: 20 },
  { label: '待收货', value: 30 },
  { label: '已完成', value: 40 },
  { label: '已取消', value: 50 },
]

async function load() {
  loading.value = true
  try {
    const data = await listOrders({ page: page.value, page_size: pageSize.value, status: status.value || undefined })
    orders.value = data.list ?? []
    total.value = data.total ?? 0
  } finally {
    loading.value = false
  }
}

function selectStatus(value: number) {
  status.value = value
  page.value = 1
  load()
}

async function pay(order: OrderInfo) {
  const res = await payOrder(order.id)
  if (res.code === 200) {
    ElMessage.success('支付成功')
    await load()
  } else {
    ElMessage.error(res.msg || '支付失败')
  }
}

async function cancel(order: OrderInfo) {
  try {
    await ElMessageBox.confirm('确定取消这个订单吗？', '提示', { type: 'warning' })
  } catch {
    return
  }
  const res = await cancelOrder(order.id, '用户取消')
  if (res.code === 200) {
    ElMessage.success('订单已取消')
    await load()
  } else {
    ElMessage.error(res.msg || '取消失败')
  }
}

async function confirm(order: OrderInfo) {
  const res = await confirmOrder(order.id)
  if (res.code === 200) {
    ElMessage.success('已确认收货')
    await load()
  } else {
    ElMessage.error(res.msg || '操作失败')
  }
}

function handlePage(value: number) {
  page.value = value
  load()
}

onMounted(load)
</script>

<template>
  <div v-loading="loading" class="order-list-page">
    <div class="page-head">
      <h1>我的订单</h1>
      <p>MY ORDERS · 每一件好物，都有迹可循</p>
    </div>

    <section class="status-tabs paper-card">
      <button
        v-for="tab in statusTabs"
        :key="tab.value"
        type="button"
        class="status-tab"
        :class="{ active: status === tab.value }"
        @click="selectStatus(tab.value)"
      >
        {{ tab.label }}
      </button>
    </section>

    <el-empty v-if="!orders.length" description="还没有订单，去拾光严选逛逛吧">
      <el-button type="primary" @click="router.push('/mall/products')">去逛逛</el-button>
    </el-empty>

    <section v-for="order in orders" :key="order.id" class="order-card paper-card">
      <div class="order-head">
        <div class="order-meta">
          <span class="order-no">订单号：{{ order.order_no }}</span>
          <span class="order-time">{{ formatTime(order.created_at) }}</span>
        </div>
        <el-tag :type="orderStatusType(order.status)" effect="plain">{{ orderStatusText(order.status) }}</el-tag>
      </div>

      <div class="order-items">
        <div v-for="item in order.items" :key="item.id" class="order-item" @click="router.push(`/mall/product/${item.product_id}`)">
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

      <div class="order-foot">
        <div class="order-summary">
          <span>共 {{ order.items.length }} 件好物</span>
          <span class="amount">实付 <b class="price-text"><span class="price-symbol">¥</span>{{ formatPrice(order.pay_amount) }}</b></span>
        </div>
        <div class="order-actions">
          <el-button v-if="order.status === 10" type="primary" @click="pay(order)">去支付</el-button>
          <el-button v-if="order.status === 10 || order.status === 20" @click="cancel(order)">取消订单</el-button>
          <el-button v-if="order.status === 30" type="primary" @click="confirm(order)">确认收货</el-button>
          <el-button @click="router.push(`/mall/order/${order.id}`)">订单详情</el-button>
        </div>
      </div>
    </section>

    <div v-if="orders.length" class="pagination">
      <el-pagination
        :current-page="page"
        :page-size="pageSize"
        :total="total"
        layout="prev, pager, next, total"
        background
        @current-change="handlePage"
      />
    </div>
  </div>
</template>

<style scoped>
.order-list-page {
  display: flex;
  flex-direction: column;
  gap: 18px;
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

.status-tabs {
  display: flex;
  gap: 4px;
  padding: 8px;
  flex-wrap: wrap;
}

.status-tab {
  border: none;
  background: transparent;
  padding: 9px 22px;
  border-radius: 5px;
  color: var(--ink-soft);
  cursor: pointer;
  font-size: 13px;
  letter-spacing: 1px;
  font-family: var(--font-sans);
}

.status-tab:hover {
  color: var(--brand);
}

.status-tab.active {
  background: var(--brand);
  color: #fffefa;
}

.order-card {
  padding: 0 22px 20px;
}

.order-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  height: 54px;
  border-bottom: 1px solid var(--line);
}

.order-meta {
  display: flex;
  align-items: center;
  gap: 18px;
  min-width: 0;
}

.order-no {
  font-size: 13px;
  color: var(--ink);
}

.order-time {
  color: var(--muted);
  font-size: 12px;
}

.order-items {
  padding: 6px 0;
}

.order-item {
  display: grid;
  grid-template-columns: 64px minmax(0, 1fr) 110px;
  align-items: center;
  gap: 14px;
  padding: 12px 0;
  border-bottom: 1px dashed rgba(201, 168, 106, 0.24);
  cursor: pointer;
}

.order-item:last-child {
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
  color: var(--ink);
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

.order-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding-top: 14px;
  border-top: 1px solid var(--line);
}

.order-summary {
  display: flex;
  align-items: baseline;
  gap: 18px;
  color: var(--muted);
  font-size: 12px;
}

.amount b {
  font-size: 22px;
}

.order-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.pagination {
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 720px) {
  .order-head {
    height: auto;
    padding: 12px 0;
    align-items: flex-start;
    flex-direction: column;
  }

  .order-foot {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
