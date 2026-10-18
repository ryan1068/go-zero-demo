<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { listAuditLogs } from '@/api/user'
import { formatTime } from '@/utils/format'
import type { AuditLogInfo } from '@/types/api'

const list = ref<AuditLogInfo[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const loading = ref(false)
const filters = reactive({
  action: '',
  target_type: '',
})

async function load() {
  loading.value = true
  try {
    const data = await listAuditLogs({
      page: page.value,
      page_size: pageSize.value,
      action: filters.action.trim() || undefined,
      target_type: filters.target_type.trim() || undefined,
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
  filters.action = ''
  filters.target_type = ''
  page.value = 1
  load()
}

function handlePage(value: number) {
  page.value = value
  load()
}

onMounted(load)
</script>

<template>
  <div v-loading="loading">
    <div class="page-header"><h3>操作审计</h3></div>

    <el-card shadow="never" class="filter-card">
      <div class="filters">
        <el-input v-model="filters.action" placeholder="动作，例如 create_user" clearable style="width: 220px" @keyup.enter="search" />
        <el-input v-model="filters.target_type" placeholder="目标类型，例如 user" clearable style="width: 200px" @keyup.enter="search" />
        <el-button type="primary" @click="search">查询</el-button>
        <el-button @click="reset">重置</el-button>
      </div>
    </el-card>

    <el-card shadow="never">
      <el-table :data="list" border stripe>
        <el-table-column label="操作人" width="140">
          <template #default="{ row }">{{ row.admin_name || `#${row.admin_id}` }}</template>
        </el-table-column>
        <el-table-column prop="action" label="动作" min-width="180" />
        <el-table-column prop="target_type" label="目标类型" width="120" />
        <el-table-column prop="detail" label="详情" min-width="220" show-overflow-tooltip />
        <el-table-column label="时间" width="180">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <template #empty>
          <el-empty description="暂无审计记录" :image-size="80" />
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
  </div>
</template>

<style scoped>
.page-header { margin-bottom: 16px; }
.page-header h3 { margin: 0; }
.filter-card { margin-bottom: 16px; }
.filters { display: flex; gap: 12px; flex-wrap: wrap; align-items: center; }
.pagination { display: flex; justify-content: flex-end; margin-top: 16px; }
</style>
