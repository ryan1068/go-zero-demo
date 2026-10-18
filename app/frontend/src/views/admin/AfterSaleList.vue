<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { adminGetAfterSale, adminHandleAfterSale, adminListAfterSales } from '@/api/order'
import { formatPrice, formatTime } from '@/utils/format'
import type { AfterSaleInfo } from '@/types/api'

const statusMap: Record<number, { text: string; type: 'success' | 'warning' | 'info' | 'danger' }> = {
  10: { text: '待审核', type: 'warning' },
  20: { text: '已通过', type: 'success' },
  30: { text: '已拒绝', type: 'danger' },
  40: { text: '买家已寄回', type: 'info' },
  50: { text: '商家已收货', type: 'info' },
  60: { text: '已退款', type: 'success' },
  70: { text: '已关闭', type: 'info' },
}

const typeText = (type: number) => (type === 2 ? '退货退款' : '仅退款')

const list = ref<AfterSaleInfo[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const filters = reactive({ status: 0 })

const detailVisible = ref(false)
const detail = ref<AfterSaleInfo | null>(null)

async function load() {
  loading.value = true
  try {
    const data = await adminListAfterSales({
      page: page.value,
      page_size: pageSize.value,
      status: filters.status || undefined,
    })
    list.value = data.list ?? []
    total.value = data.total ?? 0
  } finally {
    loading.value = false
  }
}

function search() {
  page.value = 1
  load()
}

function reset() {
  filters.status = 0
  page.value = 1
  load()
}

async function openDetail(row: AfterSaleInfo) {
  detail.value = await adminGetAfterSale(row.id)
  detailVisible.value = true
}

async function handle(row: AfterSaleInfo, action: string) {
  let remark = ''
  if (action === 'reject') {
    try {
      const result = await ElMessageBox.prompt('请输入拒绝原因', '拒绝售后', {
        inputPlaceholder: '例如：商品已使用，不支持退款',
        inputValidator: (value) => (value && value.trim() ? true : '请填写拒绝原因'),
      })
      remark = result.value.trim()
    } catch {
      return
    }
  } else {
    const labels: Record<string, string> = {
      approve: '确定同意这个售后申请吗？',
      receive: '确定已收到买家退回的商品吗？',
      refund: '确定退款吗？退款会写入退款记录。',
    }
    try {
      await ElMessageBox.confirm(labels[action] || '确定执行这个操作吗？', '提示', { type: 'warning' })
    } catch {
      return
    }
  }
  const res = await adminHandleAfterSale(row.id, { action, admin_remark: remark })
  if (res.code === 200) {
    ElMessage.success('操作成功')
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
  <div v-loading="loading">
    <div class="page-header"><h3>售后管理</h3></div>

    <el-card shadow="never" class="filter-card">
      <div class="filters">
        <el-select v-model="filters.status" style="width: 160px">
          <el-option label="全部状态" :value="0" />
          <el-option v-for="(item, key) in statusMap" :key="key" :label="item.text" :value="Number(key)" />
        </el-select>
        <el-button type="primary" @click="search">查询</el-button>
        <el-button @click="reset">重置</el-button>
      </div>
    </el-card>

    <el-card shadow="never">
      <el-table :data="list" border stripe>
        <el-table-column prop="after_sale_no" label="售后单号" min-width="210" />
        <el-table-column prop="order_no" label="订单号" min-width="200" show-overflow-tooltip />
        <el-table-column prop="product_name" label="商品" min-width="150" show-overflow-tooltip />
        <el-table-column label="类型" width="100">
          <template #default="{ row }">{{ typeText(row.type) }}</template>
        </el-table-column>
        <el-table-column label="退款金额" width="120">
          <template #default="{ row }">¥{{ formatPrice(row.amount) }}</template>
        </el-table-column>
        <el-table-column label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="statusMap[row.status]?.type || 'info'">{{ statusMap[row.status]?.text || row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="申请时间" width="170">
          <template #default="{ row }">{{ formatTime(row.apply_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button text type="primary" @click="openDetail(row)">详情</el-button>
            <template v-if="row.status === 10">
              <el-button text type="success" @click="handle(row, 'approve')">同意</el-button>
              <el-button text type="danger" @click="handle(row, 'reject')">拒绝</el-button>
            </template>
            <el-button v-if="row.type === 2 && row.status === 40" text type="success" @click="handle(row, 'receive')">
              确认收货
            </el-button>
            <el-button
              v-if="(row.type === 1 && row.status === 20) || (row.type === 2 && row.status === 50)"
              text
              type="danger"
              @click="handle(row, 'refund')"
            >
              退款
            </el-button>
            <el-tag v-if="row.type === 2 && row.status === 20" size="small" type="info">等待买家寄回</el-tag>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty description="暂无售后申请" :image-size="80" />
        </template>
      </el-table>
      <div class="pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          background
          @current-change="handlePage"
        />
      </div>
    </el-card>

    <el-dialog v-model="detailVisible" title="售后详情" width="760px">
      <template v-if="detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="售后单号">{{ detail.after_sale_no }}</el-descriptions-item>
          <el-descriptions-item label="订单号">{{ detail.order_no }}</el-descriptions-item>
          <el-descriptions-item label="商品">{{ detail.product_name }}</el-descriptions-item>
          <el-descriptions-item label="类型">{{ typeText(detail.type) }}</el-descriptions-item>
          <el-descriptions-item label="退款金额">¥{{ formatPrice(detail.amount) }}</el-descriptions-item>
          <el-descriptions-item label="状态">{{ statusMap[detail.status]?.text }}</el-descriptions-item>
          <el-descriptions-item label="原因" :span="2">{{ detail.reason }}</el-descriptions-item>
          <el-descriptions-item label="问题描述" :span="2">{{ detail.description || '-' }}</el-descriptions-item>
          <el-descriptions-item label="商家备注" :span="2">{{ detail.admin_remark || '-' }}</el-descriptions-item>
          <el-descriptions-item label="退回物流">{{ detail.tracking_no || '-' }}</el-descriptions-item>
          <el-descriptions-item label="申请时间">{{ formatTime(detail.apply_at) }}</el-descriptions-item>
        </el-descriptions>
        <div v-if="detail.images?.length" class="image-row">
          <el-image
            v-for="image in detail.images"
            :key="image"
            :src="image"
            :preview-src-list="detail.images"
            fit="cover"
            class="evidence"
          />
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.page-header { margin-bottom: 16px; }
.page-header h3 { margin: 0; }
.filter-card { margin-bottom: 16px; }
.filters { display: flex; gap: 12px; flex-wrap: wrap; align-items: center; }
.pagination { display: flex; justify-content: flex-end; margin-top: 16px; }
.image-row { display: flex; gap: 10px; margin-top: 16px; flex-wrap: wrap; }
.evidence { width: 88px; height: 88px; border-radius: 6px; border: 1px solid #eee; }
</style>
