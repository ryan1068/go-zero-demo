<script setup lang="ts">
import { reactive, ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import BrandLogo from '@/components/BrandLogo.vue'
import { useUserStore } from '@/stores/user'
import { changePassword, logout as logoutApi } from '@/api/user'
import { mallUrl, openInNewTab } from '@/utils/portal'

const router = useRouter()
const route = useRoute()
const store = useUserStore()

const isCollapse = ref(false)
const passwordVisible = ref(false)
const passwordFormRef = ref<FormInstance>()
const changingPassword = ref(false)
const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: '',
})
const passwordRules: FormRules = {
  old_password: [{ required: true, message: '请输入原密码', trigger: 'blur' }],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 72, message: '密码长度为 6-72 个字符', trigger: 'blur' },
  ],
  confirm_password: [
    {
      validator: (_rule, value, callback) => {
        if (value !== passwordForm.new_password) callback(new Error('两次输入的密码不一致'))
        else callback()
      },
      trigger: 'blur',
    },
  ],
}

const menuItems = [
  { path: '/dashboard', title: '数据看板', icon: 'Odometer' },
  { path: '/category', title: '分类管理', icon: 'Menu' },
  { path: '/product', title: '商品管理', icon: 'Goods' },
  { path: '/order', title: '订单管理', icon: 'Document' },
  { path: '/after-sale', title: '售后管理', icon: 'RefreshLeft' },
  { path: '/audit', title: '操作审计', icon: 'DocumentChecked' },
  { path: '/user', title: '用户管理', icon: 'User' },
]

const activeMenu = computed(() => route.path)

function handleSelect(path: string) {
  router.push(path)
}

async function handleLogout() {
  try {
    await logoutApi()
  } catch {
    // 退出接口失败也要清掉本地登录态。
  }
  store.logout()
  router.push('/login')
}

async function submitPassword() {
  if (!passwordFormRef.value) return
  const valid = await passwordFormRef.value.validate().catch(() => false)
  if (!valid) return
  changingPassword.value = true
  try {
    const res = await changePassword({
      old_password: passwordForm.old_password,
      new_password: passwordForm.new_password,
    })
    if (res.code === 200) {
      ElMessage.success('密码修改成功')
      passwordVisible.value = false
      passwordForm.old_password = ''
      passwordForm.new_password = ''
      passwordForm.confirm_password = ''
    } else {
      ElMessage.error(res.msg || '修改失败')
    }
  } catch {
    ElMessage.error('服务器连接失败')
  } finally {
    changingPassword.value = false
  }
}
</script>

<template>
  <el-container class="layout">
    <!-- 侧边栏 -->
    <el-aside :width="isCollapse ? '64px' : '220px'" class="sidebar">
      <div class="logo" @click="router.push('/dashboard')">
        <BrandLogo v-if="!isCollapse" subtitle="管理后台" light />
        <div v-else class="mini-seal">拾光</div>
      </div>

      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapse"
        background-color="#2e2326"
        text-color="#c9b8ba"
        active-text-color="#c97a86"
        @select="handleSelect"
      >
        <el-menu-item v-for="item in menuItems" :key="item.path" :index="item.path">
          <el-icon><component :is="item.icon" /></el-icon>
          <template #title>{{ item.title }}</template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <!-- 右侧主体 -->
    <el-container>
      <!-- 顶栏 -->
      <el-header class="header">
        <div class="header-left">
          <el-icon class="collapse-btn" @click="isCollapse = !isCollapse" :size="20">
            <Fold v-if="!isCollapse" />
            <Expand v-else />
          </el-icon>
          <span class="header-title">拾光严选 · 管理后台</span>
        </div>
        <div class="header-right">
          <el-button text type="primary" size="small" @click="openInNewTab(mallUrl('/mall'))">去商城</el-button>
          <el-button text size="small" @click="passwordVisible = true">修改密码</el-button>
          <el-tag type="success" size="small" effect="plain">管理员</el-tag>
          <span class="username">{{ store.username || '管理员' }}</span>
          <el-button type="danger" text size="small" @click="handleLogout">退出登录</el-button>
        </div>
      </el-header>

      <!-- 内容区 -->
      <el-main>
        <RouterView />
      </el-main>
    </el-container>
  </el-container>

  <el-dialog v-model="passwordVisible" title="修改密码" width="460px" destroy-on-close>
    <el-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" label-width="90px">
      <el-form-item label="原密码" prop="old_password">
        <el-input v-model="passwordForm.old_password" type="password" show-password />
      </el-form-item>
      <el-form-item label="新密码" prop="new_password">
        <el-input v-model="passwordForm.new_password" type="password" show-password placeholder="6-72 个字符" />
      </el-form-item>
      <el-form-item label="确认密码" prop="confirm_password">
        <el-input v-model="passwordForm.confirm_password" type="password" show-password />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="passwordVisible = false">取消</el-button>
      <el-button type="primary" :loading="changingPassword" @click="submitPassword">保存</el-button>
    </template>
  </el-dialog>
</template>

<style scoped>
.layout {
  height: 100vh;
}

.sidebar {
  background: #2e2326;
  overflow: hidden;
  transition: width 0.3s;
}

.logo {
  height: 78px;
  display: flex;
  align-items: center;
  padding: 0 16px;
  cursor: pointer;
  border-bottom: 1px solid rgba(201, 168, 106, 0.22);
}

.mini-seal {
  width: 38px;
  height: 38px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(247, 239, 240, 0.55);
  border-radius: 6px;
  background: var(--brand);
  color: #f7eff0;
  font-family: var(--font-serif);
  font-size: 12px;
  letter-spacing: 1px;
}

.el-menu {
  border-right: none;
}

.header {
  background: var(--card);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  border-bottom: 1px solid var(--line);
  box-shadow: 0 2px 12px rgba(46, 35, 38, 0.05);
  z-index: 1;
}

.header-title {
  margin-left: 12px;
  font-family: var(--font-serif);
  font-size: 15px;
  letter-spacing: 2px;
  color: var(--ink);
}

.header-left {
  display: flex;
  align-items: center;
}

.collapse-btn {
  cursor: pointer;
  color: #666;
}

.collapse-btn:hover {
  color: var(--brand);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.username {
  color: var(--ink);
  font-size: 14px;
}

.el-main {
  background: transparent;
  padding: 22px;
}
</style>
