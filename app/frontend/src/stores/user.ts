import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { createUser, getUser, listUsers, login as loginApi, register as registerApi } from '@/api/user'
import type { CreateUserReq, ListUsersParams, RegisterReq, UserInfo } from '@/types/api'

export const useUserStore = defineStore('user', () => {
  const user = ref<UserInfo | null>(null)
  const token = ref<string>(localStorage.getItem('token') || '')
  const role = ref<number>(Number(localStorage.getItem('role') || 0))
  const username = ref<string>(localStorage.getItem('username') || '')
  const loading = ref(false)

  const users = ref<UserInfo[]>([])
  const total = ref(0)
  const listLoading = ref(false)
  const listError = ref('')

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => role.value >= 1)
  const isSuperAdmin = computed(() => role.value === 1)

  async function fetchUser(name: string) {
    loading.value = true
    try {
      user.value = await getUser(name)
    } finally {
      loading.value = false
    }
  }

  async function fetchUsers(params: ListUsersParams) {
    listLoading.value = true
    listError.value = ''
    try {
      const data = await listUsers(params)
      users.value = data.list ?? []
      total.value = data.total ?? 0
      return data
    } catch (err) {
      users.value = []
      total.value = 0
      listError.value = '用户列表加载失败'
      throw err
    } finally {
      listLoading.value = false
    }
  }

  async function addUser(data: CreateUserReq) {
    return createUser(data)
  }

  async function login(name: string, password: string) {
    loading.value = true
    try {
      const res = await loginApi({ username: name, password })
      if (res.code === 200 && res.data?.token) {
        token.value = res.data.token
        role.value = res.data.role ?? 0
        username.value = res.data.username
        localStorage.setItem('token', res.data.token)
        localStorage.setItem('role', String(role.value))
        localStorage.setItem('username', res.data.username)
        return true
      }
      return false
    } finally {
      loading.value = false
    }
  }

  async function register(data: RegisterReq) {
    return registerApi(data)
  }

  function logout() {
    token.value = ''
    role.value = 0
    username.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('role')
    localStorage.removeItem('username')
  }

  return {
    user,
    token,
    role,
    username,
    loading,
    isLoggedIn,
    isAdmin,
    isSuperAdmin,
    users,
    total,
    listLoading,
    listError,
    fetchUser,
    fetchUsers,
    addUser,
    login,
    register,
    logout,
  }
})
