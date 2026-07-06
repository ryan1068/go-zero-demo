<script setup lang="ts">
import { ref } from 'vue'
import { useUserStore } from '@/stores/user'
import { Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const store = useUserStore()
const username = ref('')
const loading = ref(false)

async function handleSearch() {
  if (!username.value) {
    ElMessage.warning('请输入用户名')
    return
  }
  loading.value = true
  try {
    await store.fetchUser(username.value)
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
        <span class="card-title">用户管理</span>
      </template>

      <div class="search-bar">
        <el-input
          v-model="username"
          placeholder="输入用户名搜索"
          :prefix-icon="Search"
          clearable
          style="width: 300px"
          @keyup.enter="handleSearch"
        />
        <el-button type="primary" :icon="Search" :loading="loading" @click="handleSearch">
          查询
        </el-button>
      </div>

      <el-table v-if="store.user" :data="[store.user]" border stripe style="margin-top: 16px">
        <el-table-column prop="name" label="用户名" min-width="150" />
        <el-table-column prop="message" label="返回信息" min-width="250" />
      </el-table>

      <el-empty v-else description="输入用户名后点击查询" :image-size="80" />
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
}
</style>
