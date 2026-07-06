<script setup lang="ts">
import { ref } from 'vue'
import { useOrderStore } from '@/stores/order'
import { Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const store = useOrderStore()
const orderId = ref<number | null>(1)
const loading = ref(false)

async function handleSearch() {
  if (!orderId.value) {
    ElMessage.warning('请输入订单 ID')
    return
  }
  loading.value = true
  try {
    await store.fetchOrder(orderId.value)
  } catch {
    ElMessage.error('查询失败')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="page">
    <el-card shadow="never">
      <template #header>
        <span class="card-title">订单查询</span>
      </template>

      <div class="search-bar">
        <el-input-number v-model="orderId" :min="1" placeholder="订单 ID" />
        <el-button type="primary" :icon="Search" :loading="loading" @click="handleSearch">
          查询
        </el-button>
      </div>

      <el-descriptions v-if="store.order" :column="2" border style="margin-top: 16px">
        <el-descriptions-item label="订单 ID">{{ orderId }}</el-descriptions-item>
        <el-descriptions-item label="商品名称">{{ store.order.goodsName || store.order.data }}</el-descriptions-item>
        <el-descriptions-item label="价格">{{ store.order.price || '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag type="success">正常</el-tag>
        </el-descriptions-item>
      </el-descriptions>

      <el-empty v-else description="输入订单 ID 后点击查询" :image-size="80" />
    </el-card>
  </div>
</template>

<style scoped>
.page {
  max-width: 100%;
}
.card-title {
  font-weight: 600;
}
.search-bar {
  display: flex;
  gap: 12px;
  align-items: center;
}
</style>
