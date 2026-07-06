<script setup lang="ts">
import { ref } from 'vue'
import { useUserStore } from '@/stores/user'

const store = useUserStore()
const username = ref('')
</script>

<template>
  <div class="page">
    <h1>用户查询</h1>
    <p class="desc">调用 user-api → user-rpc 查询用户信息</p>

    <div class="form-row">
      <input v-model="username" placeholder="输入用户名" @keyup.enter="store.fetchUser(username)" />
      <button :disabled="!username || store.loading" @click="store.fetchUser(username)">
        {{ store.loading ? '查询中...' : '查询' }}
      </button>
    </div>

    <div v-if="store.user" class="result">
      <h3>查询结果</h3>
      <pre>{{ store.user }}</pre>
    </div>

    <div class="back">
      <RouterLink to="/">← 返回首页</RouterLink>
    </div>
  </div>
</template>

<style scoped>
.page {
  max-width: 600px;
  margin: 0 auto;
  padding: 2rem;
}
.desc {
  color: #888;
  margin-bottom: 1.5rem;
}
.form-row {
  display: flex;
  gap: 0.75rem;
}
input {
  flex: 1;
  padding: 0.6rem 0.8rem;
  border: 1px solid #ccc;
  border-radius: 6px;
  font-size: 1rem;
}
button {
  padding: 0.6rem 1.5rem;
  background: #409eff;
  color: #fff;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 1rem;
}
button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.result {
  margin-top: 1.5rem;
  padding: 1rem;
  background: #f5f7fa;
  border-radius: 8px;
}
.result pre {
  margin: 0;
  font-size: 0.9rem;
}
.back {
  margin-top: 2rem;
}
</style>
