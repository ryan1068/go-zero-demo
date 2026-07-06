<script setup lang="ts">
import { ref } from 'vue'
import { useOrderStore } from '@/stores/order'

const store = useOrderStore()
const orderId = ref<number>(1)
</script>

<template>
  <div class="page">
    <h1>订单查询</h1>
    <p class="desc">GET /order/:id → order-api → order-rpc</p>

    <div class="form-row">
      <input v-model.number="orderId" type="number" placeholder="输入订单 ID" @keyup.enter="store.fetchOrder(orderId)" />
      <button :disabled="!orderId || store.loading" @click="store.fetchOrder(orderId)">
        {{ store.loading ? '查询中...' : '查询' }}
      </button>
    </div>

    <div v-if="store.order" class="result">
      <h3>查询结果</h3>
      <pre>{{ store.order }}</pre>
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
