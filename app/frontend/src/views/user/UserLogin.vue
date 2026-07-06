<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const store = useUserStore()

const username = ref('')
const password = ref('')
const errorMsg = ref('')
const loading = ref(false)

async function handleLogin() {
  loading.value = true
  errorMsg.value = ''
  try {
    const ok = await store.login(username.value, password.value)
    if (ok) {
      router.push('/')
    } else {
      errorMsg.value = '用户名或密码错误'
    }
  } catch (err: unknown) {
    errorMsg.value = `请求失败: ${err instanceof Error ? err.message : String(err)}`
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="page">
    <h1>用户登录</h1>
    <p class="desc">POST /user/login → user-api → user-rpc → MySQL</p>

    <div class="form">
      <input v-model="username" placeholder="用户名" @keyup.enter="handleLogin" />
      <input v-model="password" type="password" placeholder="密码" @keyup.enter="handleLogin" />
      <button :disabled="!username || !password || loading" @click="handleLogin">
        {{ loading ? '登录中...' : '登录' }}
      </button>
    </div>

    <div v-if="errorMsg" class="error">{{ errorMsg }}</div>

    <p class="hint">测试账号：admin / 123456</p>

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
.form {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}
input {
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
.error {
  margin-top: 1rem;
  padding: 0.75rem;
  background: #fef0f0;
  border: 1px solid #fde2e2;
  border-radius: 6px;
  color: #f56c6c;
  font-size: 0.9rem;
}
.hint {
  margin-top: 1.5rem;
  color: #999;
  font-size: 0.85rem;
}
.back {
  margin-top: 2rem;
}
</style>
