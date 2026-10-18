<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { adminListOrders, dashboard } from '@/api/order'
import { productStats } from '@/api/product'
import { formatPrice, formatTime, orderStatusText, orderStatusType } from '@/utils/format'
import type { DashboardData, OrderInfo, ProductStats } from '@/types/api'

const router = useRouter()
const data = ref<DashboardData>({
  today_order_count: 0,
  today_paid_amount: 0,
  pending_shipment_count: 0,
  product_count: 0,
  low_stock_count: 0,
})
const stats = ref<ProductStats>({ total: 0, on_sale: 0, off_sale: 0, low_stock: 0 })
const recent = ref<OrderInfo[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const [d, s, orders] = await Promise.all([
      dashboard(),
      productStats(),
      adminListOrders({ page: 1, page_size: 5 }),
    ])
    data.value = d
    stats.value = s
    recent.value = orders.list ?? []
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div v-loading="loading">
    <h3 class="page-title">数据看板</h3>

    <el-row :gutter="18">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-value">{{ data.today_order_count }}</div>
          <div class="stat-label">今日订单</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-value">¥{{ formatPrice(data.today_paid_amount) }}</div>
          <div class="stat-label">今日销售额</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-value">{{ data.pending_shipment_count }}</div>
          <div class="stat-label">待发货订单</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-value">{{ data.product_count }}</div>
          <div class="stat-label">商品总数</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="18" class="second-row">
      <el-col :span="8">
        <el-card shadow="hover">
          <template #header>商品状态</template>
          <div class="mini-stat"><span>上架商品</span><b>{{ stats.on_sale }}</b></div>
          <div class="mini-stat"><span>下架商品</span><b>{{ stats.off_sale }}</b></div>
          <div class="mini-stat"><span>低库存商品</span><b class="warn">{{ stats.low_stock }}</b></div>
        </el-card>
      </el-col>
      <el-col :span="16">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>最近订单</span>
              <el-button text type="primary" @click="router.push('/order')">查看全部</el-button>
            </div>
          </template>
          <el-table :data="recent" size="small">
            <el-table-column prop="order_no" label="订单号" min-width="190" />
            <el-table-column label="金额" width="100">
              <template #default="{ row }">¥{{ formatPrice(row.pay_amount) }}</template>
            </el-table-column>
            <el-table-column label="状态" width="90">
              <template #default="{ row }">
                <el-tag :type="orderStatusType(row.status)" size="small">{{ orderStatusText(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="时间" width="170">
              <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.page-title { margin: 0 0 18px; color: #303133; }
.stat-card { text-align: center; }
.stat-value { font-size: 28px; font-weight: 700; color: #303133; }
.stat-label { margin-top: 6px; color: #909399; font-size: 13px; }
.second-row { margin-top: 18px; }
.mini-stat { display: flex; justify-content: space-between; padding: 10px 0; border-bottom: 1px solid #f2f2f2; color: #606266; }
.mini-stat:last-child { border-bottom: none; }
.mini-stat b { color: #303133; }
.mini-stat b.warn { color: #e6a23c; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
</style>
