<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { adminCancelOrder, adminGetOrder, adminListOrders, adminOrderTracks, shipOrder } from '@/api/order'
import { formatPrice, formatTime, orderStatusText, orderStatusType } from '@/utils/format'
import type { OrderInfo, OrderTracksData } from '@/types/api'

const orders = ref<OrderInfo[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const filters = reactive({
  order_no: '',
  status: 0,
})

const detailVisible = ref(false)
const detail = ref<OrderInfo | null>(null)
const tracksVisible = ref(false)
const tracksLoading = ref(false)
const tracks = ref<OrderTracksData | null>(null)
const shipVisible = ref(false)
const shipFormRef = ref<FormInstance>()
const shipping = ref(false)
const shipTarget = ref<OrderInfo | null>(null)
const shipForm = reactive({
  logistics_company: '顺丰',
  tracking_no: '',
})
const shipRules: FormRules = {
  logistics_company: [{ required: true, message: '请输入物流公司', trigger: 'blur' }],
  tracking_no: [{ required: true, message: '请输入物流单号', trigger: 'blur' }],
}

async function load() {
  loading.value = true
  try {
    const data = await adminListOrders({
      page: page.value,
      page_size: pageSize.value,
      order_no: filters.order_no.trim() || undefined,
      status: filters.status || undefined,
    })
    orders.value = data.list ?? []
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
  filters.order_no = ''
  filters.status = 0
  page.value = 1
  load()
}

async function openDetail(row: OrderInfo) {
  detail.value = await adminGetOrder(row.id)
  detailVisible.value = true
}

async function openTracks(row: OrderInfo) {
  tracks.value = null
  tracksVisible.value = true
  tracksLoading.value = true
  try {
    tracks.value = await adminOrderTracks(row.id)
  } finally {
    tracksLoading.value = false
  }
}

function openShip(row: OrderInfo) {
  shipTarget.value = row
  shipForm.logistics_company = row.logistics_company || '顺丰'
  shipForm.tracking_no = row.tracking_no || ''
  shipVisible.value = true
}

async function submitShip() {
  if (!shipFormRef.value || !shipTarget.value) return
  const valid = await shipFormRef.value.validate().catch(() => false)
  if (!valid) return
  shipping.value = true
  try {
    const res = await shipOrder(shipTarget.value.id, shipForm.logistics_company.trim(), shipForm.tracking_no.trim())
    if (res.code === 200) {
      ElMessage.success('发货成功')
      shipVisible.value = false
      await load()
    } else {
      ElMessage.error(res.msg || '发货失败')
    }
  } finally {
    shipping.value = false
  }
}

async function cancel(row: OrderInfo) {
  try {
    await ElMessageBox.confirm('确定取消这个订单吗？库存会回补。', '提示', { type: 'warning' })
  } catch {
    return
  }
  const res = await adminCancelOrder(row.id, '管理员取消')
  if (res.code === 200) {
    ElMessage.success('订单已取消')
    await load()
  } else {
    ElMessage.error(res.msg || '取消失败')
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
    <div class="page-header"><h3>订单管理</h3></div>

    <el-card shadow="never" class="filter-card">
      <div class="filters">
        <el-input v-model="filters.order_no" placeholder="订单号" clearable style="width: 220px" @keyup.enter="search" />
        <el-select v-model="filters.status" style="width: 140px">
          <el-option label="全部状态" :value="0" />
          <el-option label="待支付" :value="10" />
          <el-option label="待发货" :value="20" />
          <el-option label="待收货" :value="30" />
          <el-option label="已完成" :value="40" />
          <el-option label="已取消" :value="50" />
        </el-select>
        <el-button type="primary" @click="search">查询</el-button>
        <el-button @click="reset">重置</el-button>
      </div>
    </el-card>

    <el-card shadow="never">
      <el-table :data="orders" border stripe>
        <el-table-column prop="order_no" label="订单号" min-width="200" />
        <el-table-column label="金额" width="120">
          <template #default="{ row }">¥{{ formatPrice(row.pay_amount) }}</template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="orderStatusType(row.status)">{{ orderStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="收货人" width="110">
          <template #default="{ row }">{{ row.receiver_name }}</template>
        </el-table-column>
        <el-table-column label="创建时间" width="180">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="物流" min-width="160">
          <template #default="{ row }">{{ row.logistics_company || '-' }} {{ row.tracking_no }}</template>
        </el-table-column>
        <el-table-column label="操作" width="300">
          <template #default="{ row }">
            <el-button text type="primary" @click="openDetail(row)">详情</el-button>
            <el-button v-if="row.tracking_no" text @click="openTracks(row)">轨迹</el-button>
            <el-button v-if="row.status === 20" text type="success" @click="openShip(row)">发货</el-button>
            <el-button v-if="row.status === 10 || row.status === 20" text type="danger" @click="cancel(row)">取消</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <el-pagination :current-page="page" :page-size="pageSize" :total="total" layout="prev, pager, next, total" background @current-change="handlePage" />
      </div>
    </el-card>

    <el-dialog v-model="detailVisible" title="订单详情" width="760px">
      <template v-if="detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="订单号">{{ detail.order_no }}</el-descriptions-item>
          <el-descriptions-item label="状态">{{ orderStatusText(detail.status) }}</el-descriptions-item>
          <el-descriptions-item label="实付">¥{{ formatPrice(detail.pay_amount) }}</el-descriptions-item>
          <el-descriptions-item label="收货人">{{ detail.receiver_name }}</el-descriptions-item>
          <el-descriptions-item label="电话">{{ detail.receiver_phone }}</el-descriptions-item>
          <el-descriptions-item label="地址" :span="2">{{ detail.receiver_province }}{{ detail.receiver_city }}{{ detail.receiver_district }}{{ detail.receiver_address }}</el-descriptions-item>
          <el-descriptions-item label="备注" :span="2">{{ detail.remark || '-' }}</el-descriptions-item>
          <el-descriptions-item label="物流公司">{{ detail.logistics_company || '-' }}</el-descriptions-item>
          <el-descriptions-item label="物流单号">{{ detail.tracking_no || '-' }}</el-descriptions-item>
        </el-descriptions>
        <el-table :data="detail.items" border style="margin-top: 16px">
          <el-table-column prop="product_name" label="商品" min-width="180" />
          <el-table-column label="单价" width="110">
            <template #default="{ row }">¥{{ formatPrice(row.price) }}</template>
          </el-table-column>
          <el-table-column prop="quantity" label="数量" width="80" />
          <el-table-column label="小计" width="110">
            <template #default="{ row }">¥{{ formatPrice(row.total_amount) }}</template>
          </el-table-column>
        </el-table>
      </template>
    </el-dialog>

    <el-dialog v-model="tracksVisible" title="物流轨迹" width="560px">
      <div v-loading="tracksLoading">
        <template v-if="tracks">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="订单号">{{ tracks.order_no }}</el-descriptions-item>
            <el-descriptions-item label="物流">{{ tracks.logistics_company }} {{ tracks.tracking_no }}</el-descriptions-item>
          </el-descriptions>
          <el-timeline v-if="tracks.tracks.length" style="margin-top: 18px">
            <el-timeline-item
              v-for="track in tracks.tracks"
              :key="`${track.time}-${track.status}`"
              :timestamp="formatTime(track.time)"
              placement="top"
            >
              {{ track.description }}
            </el-timeline-item>
          </el-timeline>
          <el-empty v-else description="暂无物流轨迹" :image-size="70" />
        </template>
      </div>
    </el-dialog>

    <el-dialog v-model="shipVisible" title="订单发货" width="480px">
      <el-form ref="shipFormRef" :model="shipForm" :rules="shipRules" label-width="90px">
        <el-form-item label="物流公司" prop="logistics_company">
          <el-input v-model="shipForm.logistics_company" />
        </el-form-item>
        <el-form-item label="物流单号" prop="tracking_no">
          <el-input v-model="shipForm.tracking_no" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="shipVisible = false">取消</el-button>
        <el-button type="primary" :loading="shipping" @click="submitShip">确认发货</el-button>
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
</style>
