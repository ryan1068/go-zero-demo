<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getUser } from '@/api/user'

const userCount = ref(0)
const orderCount = ref(0)
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    // 尝试查询测试用户确认后端连通性
    await getUser('admin')
    userCount.value = 1
  } catch {
    // ignore
  }
  loading.value = false
})

const stats = [
  { label: '用户总数', value: userCount, icon: 'User', color: '#409eff', bg: '#ecf5ff' },
  { label: '订单总数', value: orderCount, icon: 'Document', color: '#67c23a', bg: '#f0f9eb' },
  { label: 'API 调用', value: '2', icon: 'Connection', color: '#e6a23c', bg: '#fdf6ec' },
  { label: '服务状态', value: '正常', icon: 'CircleCheck', color: '#f56c6c', bg: '#fef0f0' },
]
</script>

<template>
  <div class="dashboard">
    <h3 class="page-title">仪表盘</h3>

    <el-row :gutter="20">
      <el-col v-for="s in stats" :key="s.label" :span="6">
        <el-card shadow="hover" class="stat-card" v-loading="loading">
          <div class="stat-row">
            <div class="stat-icon" :style="{ background: s.bg, color: s.color }">
              <el-icon :size="28"><component :is="s.icon" /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ s.value }}</div>
              <div class="stat-label">{{ s.label }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>微服务架构</template>
          <div class="arch-diagram">
            <div class="arch-node gateway">Gateway<br/>:8888</div>
            <div class="arch-arrow">→</div>
            <div class="arch-node api">user-api<br/>:8000</div>
            <div class="arch-arrow">→</div>
            <div class="arch-node rpc">user-rpc<br/>:9000</div>
            <div class="arch-arrow">→</div>
            <div class="arch-node db">MySQL<br/>:3306</div>
          </div>
          <div class="arch-diagram" style="margin-top: 12px">
            <div class="arch-node api">order-api<br/>:8001</div>
            <div class="arch-arrow">→</div>
            <div class="arch-node rpc">order-rpc<br/>:9001</div>
            <div class="arch-arrow">→</div>
            <div class="arch-node rpc" style="background: #fdf6ec; color: #e6a23c">user-rpc<br/>:9000</div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>技术栈</template>
          <el-descriptions :column="1" border size="small">
            <el-descriptions-item label="前端框架">Vue 3 + TypeScript + Vite</el-descriptions-item>
            <el-descriptions-item label="UI 组件库">Element Plus</el-descriptions-item>
            <el-descriptions-item label="微服务框架">go-zero v1.8.5</el-descriptions-item>
            <el-descriptions-item label="服务发现">etcd</el-descriptions-item>
            <el-descriptions-item label="数据库">MySQL 8.0</el-descriptions-item>
            <el-descriptions-item label="缓存">Redis 6.2</el-descriptions-item>
            <el-descriptions-item label="链路追踪">Jaeger</el-descriptions-item>
            <el-descriptions-item label="监控">Prometheus + Grafana</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.page-title {
  font-size: 18px;
  margin-bottom: 20px;
  color: #303133;
}

.stat-card {
  border-radius: 8px;
}

.stat-row {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

.stat-label {
  font-size: 13px;
  color: #909399;
  margin-top: 2px;
}

.arch-diagram {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.arch-node {
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 12px;
  text-align: center;
  font-weight: 500;
  line-height: 1.5;
}

.arch-node.gateway {
  background: #ecf5ff;
  color: #409eff;
}

.arch-node.api {
  background: #f0f9eb;
  color: #67c23a;
}

.arch-node.rpc {
  background: #fdf6ec;
  color: #e6a23c;
}

.arch-node.db {
  background: #fef0f0;
  color: #f56c6c;
}

.arch-arrow {
  color: #c0c4cc;
  font-size: 18px;
}
</style>
