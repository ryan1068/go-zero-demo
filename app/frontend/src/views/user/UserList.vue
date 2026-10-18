<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { Plus, Refresh, Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { updateUserRole, updateUserStatus } from '@/api/user'
import type { CreateUserReq, ListUsersParams, UserInfo } from '@/types/api'

const store = useUserStore()

const totalText = computed(() => String(store.total))

const page = ref(1)
const pageSize = ref(10)
const filters = reactive({
  username: '',
  nickname: '',
  email: '',
  phone: '',
  status: 0,
})

const dialogVisible = ref(false)
const formRef = ref<FormInstance>()
const submitting = ref(false)
const form = reactive<CreateUserReq>({
  username: '',
  password: '',
  nickname: '',
  email: '',
  phone: '',
  avatar: '',
  status: 1,
  role: 0,
  remark: '',
})

const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 64, message: '用户名长度为 3-64 个字符', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 72, message: '密码长度为 6-72 个字符', trigger: 'blur' },
  ],
  email: [{ type: 'email', message: '邮箱格式不正确', trigger: 'blur' }],
}

function statusFor(status: number) {
  return status === 1
    ? { text: '正常', type: 'success' as const }
    : { text: '禁用', type: 'danger' as const }
}

function roleFor(role: number) {
  if (role === 1) return { text: '超级管理员', type: 'warning' as const }
  if (role === 2) return { text: '运营', type: 'success' as const }
  return { text: '普通用户', type: 'info' as const }
}

const roleDialogVisible = ref(false)
const roleSaving = ref(false)
const roleTarget = ref<UserInfo | null>(null)
const roleValue = ref(0)

function openRoleDialog(row: UserInfo) {
  roleTarget.value = row
  roleValue.value = row.role
  roleDialogVisible.value = true
}

async function submitRole() {
  if (!roleTarget.value) return
  roleSaving.value = true
  try {
    const res = await updateUserRole(roleTarget.value.id, roleValue.value)
    if (res.code === 200) {
      ElMessage.success('角色已更新')
      roleDialogVisible.value = false
      await loadUsers()
    } else {
      ElMessage.error(res.msg || '更新失败')
    }
  } finally {
    roleSaving.value = false
  }
}

function formatTime(unixSeconds?: number): string {
  if (!unixSeconds) return '-'
  return new Date(unixSeconds * 1000).toLocaleString()
}

function buildParams(): ListUsersParams {
  const params: ListUsersParams = { page: page.value, page_size: pageSize.value }
  if (filters.username.trim()) params.username = filters.username.trim()
  if (filters.nickname.trim()) params.nickname = filters.nickname.trim()
  if (filters.email.trim()) params.email = filters.email.trim()
  if (filters.phone.trim()) params.phone = filters.phone.trim()
  if (filters.status) params.status = filters.status
  return params
}

async function loadUsers() {
  try {
    await store.fetchUsers(buildParams())
  } catch {
    ElMessage.error(store.listError || '用户列表加载失败')
  }
}

function handleSearch() {
  page.value = 1
  loadUsers()
}

function handleReset() {
  filters.username = ''
  filters.nickname = ''
  filters.email = ''
  filters.phone = ''
  filters.status = 0
  page.value = 1
  pageSize.value = 10
  loadUsers()
}

function handlePageChange(value: number) {
  page.value = value
  loadUsers()
}

function handleSizeChange(value: number) {
  pageSize.value = value
  page.value = 1
  loadUsers()
}

function openCreate() {
  Object.assign(form, {
    username: '',
    password: '',
    nickname: '',
    email: '',
    phone: '',
    avatar: '',
    status: 1,
    role: 0,
    remark: '',
  })
  dialogVisible.value = true
}

async function submitCreate() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const res = await store.addUser({
      username: form.username.trim(),
      password: form.password,
      nickname: form.nickname.trim(),
      email: form.email.trim(),
      phone: form.phone.trim(),
      avatar: form.avatar.trim(),
      status: form.status,
      role: form.role,
      remark: form.remark.trim(),
    })
    if (res.code === 200) {
      ElMessage.success('用户创建成功')
      dialogVisible.value = false
      page.value = 1
      await loadUsers()
    } else if (res.code === 409) {
      ElMessage.error('用户名已存在')
    } else {
      ElMessage.error(res.msg || '创建失败')
    }
  } catch {
    ElMessage.error('服务器连接失败')
  } finally {
    submitting.value = false
  }
}

async function toggleStatus(row: UserInfo) {
  const next = row.status === 1 ? 0 : 1
  const res = await updateUserStatus(row.id, next)
  if (res.code === 200) {
    ElMessage.success(next === 1 ? '已启用' : '已禁用')
    await loadUsers()
  } else {
    ElMessage.error(res.msg || '操作失败')
  }
}

onMounted(loadUsers)
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h3 class="page-title">用户管理</h3>
        <p class="page-desc">支持按用户名、昵称、邮箱、手机号和状态筛选，并可新增用户</p>
      </div>
      <el-button type="primary" :icon="Plus" @click="openCreate">新增用户</el-button>
    </div>

    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" class="filter-form" @submit.prevent>
        <el-form-item label="用户名">
          <el-input v-model="filters.username" placeholder="前缀匹配" clearable class="filter-input" @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="filters.nickname" placeholder="前缀匹配" clearable class="filter-input" @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="filters.email" placeholder="前缀匹配" clearable class="filter-input" @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="filters.phone" placeholder="前缀匹配" clearable class="filter-input" @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" class="status-select">
            <el-option label="全部" :value="0" />
            <el-option label="正常" :value="1" />
            <el-option label="禁用" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" :loading="store.listLoading" @click="handleSearch">查询</el-button>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-alert
      v-if="store.listError"
      :title="store.listError"
      type="error"
      show-icon
      :closable="false"
      class="error-alert"
    />

    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span class="card-title">用户列表</span>
          <span class="total-text">共 {{ totalText }} 个用户</span>
        </div>
      </template>

      <el-table v-loading="store.listLoading" :data="store.users" border stripe>
        <el-table-column label="头像" width="80">
          <template #default="{ row }">
            <el-avatar :size="34" :src="row.avatar || undefined">
              {{ row.username.slice(0, 1).toUpperCase() }}
            </el-avatar>
          </template>
        </el-table-column>
        <el-table-column prop="username" label="用户名" min-width="140" />
        <el-table-column prop="nickname" label="昵称" min-width="120" />
        <el-table-column prop="email" label="邮箱" min-width="190" show-overflow-tooltip />
        <el-table-column prop="phone" label="手机号" width="140" />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="statusFor(row.status).type">{{ statusFor(row.status).text }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="角色" width="120">
          <template #default="{ row }">
            <el-tag :type="roleFor(row.role).type">{{ roleFor(row.role).text }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="160" show-overflow-tooltip />
        <el-table-column label="创建时间" min-width="170">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="更新时间" min-width="170">
          <template #default="{ row }">{{ formatTime(row.updated_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button text :type="row.status === 1 ? 'danger' : 'success'" @click="toggleStatus(row)">
              {{ row.status === 1 ? '禁用' : '启用' }}
            </el-button>
            <el-button v-if="store.isSuperAdmin" text type="primary" @click="openRoleDialog(row)">改角色</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty description="暂无符合条件的用户" :image-size="80" />
        </template>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="store.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          background
          @current-change="handlePageChange"
          @size-change="handleSizeChange"
        />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" title="新增用户" width="600px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="form.username" placeholder="3-64 个字符" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="密码" prop="password">
              <el-input v-model="form.password" type="password" show-password placeholder="6-72 个字符" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="昵称">
              <el-input v-model="form.nickname" placeholder="选填" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态">
              <el-select v-model="form.status" style="width: 100%">
                <el-option label="正常" :value="1" />
                <el-option label="禁用" :value="0" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="角色">
              <el-select v-model="form.role" style="width: 100%">
                <el-option label="普通用户" :value="0" />
                <el-option label="运营" :value="2" />
                <el-option label="超级管理员" :value="1" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="form.email" placeholder="选填" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="手机号">
              <el-input v-model="form.phone" placeholder="选填" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="头像地址">
              <el-input v-model="form.avatar" placeholder="选填，例如 /avatars/alice.svg" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="备注">
              <el-input v-model="form.remark" type="textarea" :rows="3" maxlength="255" show-word-limit placeholder="选填" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitCreate">创建</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="roleDialogVisible" title="修改角色" width="420px">
      <el-form label-width="80px">
        <el-form-item label="用户">
          <span>{{ roleTarget?.username }}</span>
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="roleValue" style="width: 100%">
            <el-option label="普通用户" :value="0" />
            <el-option label="运营" :value="2" />
            <el-option label="超级管理员" :value="1" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="roleDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="roleSaving" @click="submitRole">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.page {
  max-width: 100%;
}

.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 16px;
}

.page-title {
  margin: 0;
  font-size: 18px;
  color: #303133;
}

.page-desc {
  margin: 6px 0 0;
  color: #909399;
  font-size: 13px;
}

.filter-card {
  margin-bottom: 16px;
}

.filter-form {
  display: flex;
  flex-wrap: wrap;
  gap: 4px 0;
}

.filter-input {
  width: 170px;
}

.status-select {
  width: 110px;
}

.error-alert {
  margin-bottom: 16px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-title {
  font-weight: 600;
  font-size: 15px;
}

.total-text {
  color: #909399;
  font-size: 13px;
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
