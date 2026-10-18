<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import BrandLogo from '@/components/BrandLogo.vue'
import { useUserStore } from '@/stores/user'
import { isAdminMode, isMallMode, mallUrl, openInNewTab } from '@/utils/portal'
import { User, Lock } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const store = useUserStore()

const formRef = ref()
const loading = ref(false)
const errorMsg = ref(route.query.reason === 'not_admin' ? '请使用管理员账号登录' : '')
const adminMode = computed(() => isAdminMode())
const mallMode = computed(() => isMallMode())

const form = reactive({
  username: 'admin',
  password: '123456',
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

async function handleLogin() {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  errorMsg.value = ''
  try {
    const ok = await store.login(form.username, form.password)
    if (!ok) {
      errorMsg.value = '用户名或密码错误'
      return
    }

    if (mallMode.value) {
      const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : ''
      router.push(redirect || '/mall')
      return
    }

    if (adminMode.value && !store.isAdmin) {
      store.logout()
      errorMsg.value = '该账号不是管理员，请使用商城端口登录'
      return
    }

    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : ''
    router.push(redirect || (store.isAdmin ? '/dashboard' : '/mall'))
  } catch {
    errorMsg.value = '服务器连接失败，请检查后端服务是否启动'
  } finally {
    loading.value = false
  }
}

function goMall() {
  if (adminMode.value) {
    openInNewTab(mallUrl('/mall'))
  } else {
    router.push('/mall')
  }
}
</script>

<template>
  <div class="auth-page">
    <section class="brand-panel">
      <div class="brand-panel-inner">
        <BrandLogo size="lg" :subtitle="adminMode ? '管理后台' : 'SHIGUANG SELECT'" light />
        <div class="brand-copy">
          <p class="eyebrow">{{ adminMode ? 'SHIGUANG ADMIN' : 'SHIGUANG SELECT' }}</p>
          <h1 v-if="adminMode">掌柜，<br />欢迎回来</h1>
          <h1 v-else>把好物<br />收进时光里</h1>
          <p class="desc">
            <template v-if="adminMode">拾光严选管理后台，查看订单、商品与用户，把每一笔好生意照顾好。</template>
            <template v-else>图书、数码、食品与生活好物。慢慢挑，认真选，买得安心。</template>
          </p>
        </div>
        <div class="brand-points">
          <div><span class="seal-tag seal-tag--gold">精选</span><span>认真挑过的每一件好物</span></div>
          <div><span class="seal-tag seal-tag--gold">正品</span><span>精选供应链，品质可溯</span></div>
          <div><span class="seal-tag seal-tag--gold">安心</span><span>七天无理由退换</span></div>
        </div>
        <div class="brand-seal">拾光<br />严选</div>
      </div>
    </section>

    <section class="form-panel">
      <div class="form-card">
        <p class="form-eyebrow">{{ adminMode ? 'ADMIN LOGIN' : 'WELCOME BACK' }}</p>
        <h2>{{ adminMode ? '登录管理后台' : '登录拾光严选' }}</h2>
        <p class="form-sub">{{ adminMode ? '使用管理员账号继续' : '登录后继续挑选好物' }}</p>

        <el-form ref="formRef" :model="form" :rules="rules" size="large" @keyup.enter="handleLogin">
          <el-form-item prop="username">
            <el-input v-model="form.username" placeholder="用户名" :prefix-icon="User" />
          </el-form-item>
          <el-form-item prop="password">
            <el-input v-model="form.password" type="password" placeholder="密码" :prefix-icon="Lock" show-password />
          </el-form-item>
          <el-alert v-if="errorMsg" :title="errorMsg" type="error" show-icon :closable="false" class="form-error" />
          <el-form-item>
            <el-button type="primary" :loading="loading" class="login-btn" @click="handleLogin">登 录</el-button>
          </el-form-item>
        </el-form>

        <div class="form-links">
          <el-link v-if="!adminMode" type="primary" @click="router.push('/register')">没有账号？去注册</el-link>
          <el-link @click="goMall">去商城逛逛</el-link>
        </div>

        <div class="form-hint">
          <template v-if="adminMode">管理员测试账号：admin / 123456 · 商城端：:8081</template>
          <template v-else>管理后台请访问 :8080</template>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: grid;
  grid-template-columns: 1.05fr 0.95fr;
  background: var(--paper);
}

.brand-panel {
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 56px 64px;
  color: #f7eff0;
  background:
    radial-gradient(circle at 78% 18%, rgba(201, 168, 106, 0.3), transparent 30%),
    radial-gradient(circle at 16% 86%, rgba(163, 59, 74, 0.32), transparent 34%),
    linear-gradient(145deg, #1d1719 0%, #2e2326 48%, #4a2229 100%);
}

.brand-panel::after {
  content: '';
  position: absolute;
  inset: 0;
  opacity: 0.16;
  background-image: radial-gradient(rgba(247, 239, 240, 0.6) 0.6px, transparent 0.6px);
  background-size: 8px 8px;
  pointer-events: none;
}

.brand-panel-inner {
  position: relative;
  z-index: 2;
  max-width: 480px;
  display: flex;
  flex-direction: column;
  gap: 40px;
}

.brand-copy .eyebrow {
  margin: 0 0 14px;
  color: var(--gold);
  font-size: 11px;
  letter-spacing: 4px;
}

.brand-copy h1 {
  margin: 0;
  font-family: var(--font-serif);
  font-size: 46px;
  line-height: 1.35;
  letter-spacing: 6px;
  font-weight: 700;
}

.brand-copy .desc {
  margin: 20px 0 0;
  color: rgba(247, 239, 240, 0.7);
  font-size: 13px;
  line-height: 1.9;
  max-width: 380px;
}

.brand-points {
  display: flex;
  flex-direction: column;
  gap: 12px;
  color: rgba(247, 239, 240, 0.76);
  font-size: 13px;
}

.brand-points div {
  display: flex;
  align-items: center;
  gap: 10px;
}

.brand-seal {
  position: absolute;
  right: -20px;
  bottom: -8px;
  width: 120px;
  height: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  border: 1px solid rgba(201, 168, 106, 0.45);
  border-radius: 14px;
  color: rgba(201, 168, 106, 0.7);
  font-family: var(--font-serif);
  font-size: 24px;
  line-height: 1.5;
  letter-spacing: 6px;
  transform: rotate(-7deg);
}

.form-panel {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
  background:
    radial-gradient(circle at 88% 8%, rgba(201, 168, 106, 0.16), transparent 30%),
    var(--paper);
}

.form-card {
  width: 100%;
  max-width: 400px;
  padding: 40px;
  background: var(--card);
  border: 1px solid var(--line);
  border-radius: 14px;
  box-shadow: var(--shadow);
}

.form-eyebrow {
  margin: 0 0 10px;
  color: var(--gold);
  font-size: 11px;
  letter-spacing: 3px;
}

.form-card h2 {
  margin: 0;
  font-family: var(--font-serif);
  font-size: 28px;
  letter-spacing: 3px;
  color: var(--ink);
}

.form-sub {
  margin: 8px 0 26px;
  color: var(--muted);
  font-size: 13px;
}

.form-error {
  margin-bottom: 16px;
  border-radius: 6px;
}

.login-btn {
  width: 100%;
  letter-spacing: 6px;
  height: 44px;
}

.form-links {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 6px;
}

.form-hint {
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px dashed rgba(201, 168, 106, 0.3);
  color: var(--muted);
  font-size: 11px;
  line-height: 1.7;
  text-align: center;
}

@media (max-width: 860px) {
  .auth-page {
    grid-template-columns: 1fr;
  }

  .brand-panel {
    padding: 36px 24px;
  }

  .brand-copy h1 {
    font-size: 34px;
  }

  .brand-points,
  .brand-seal {
    display: none;
  }

  .brand-panel-inner {
    gap: 24px;
  }

  .form-panel {
    padding: 24px 16px 40px;
  }

  .form-card {
    padding: 28px 22px;
  }
}
</style>
