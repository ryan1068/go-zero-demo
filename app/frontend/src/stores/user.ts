import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getUser, login as loginApi } from '@/api/user'
import type { UserInfo } from '@/types/api'

export const useUserStore = defineStore('user', () => {
  const user = ref<UserInfo | null>(null)
  const token = ref<string>(localStorage.getItem('token') || '')
  const loading = ref(false)

  const isLoggedIn = computed(() => !!token.value)

  async function fetchUser(name: string) {
    loading.value = true
    try {
      user.value = await getUser(name)
    } finally {
      loading.value = false
    }
  }

  async function login(username: string, password: string) {
    loading.value = true
    try {
      const res = await loginApi({ username, password })
      if (res.code === 200 && res.token) {
        token.value = res.token
        localStorage.setItem('token', res.token)
        return true
      }
      return false
    } finally {
      loading.value = false
    }
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
  }

  return { user, token, loading, isLoggedIn, fetchUser, login, logout }
})
