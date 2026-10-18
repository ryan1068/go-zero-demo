<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import BrandLogo from '@/components/BrandLogo.vue'
import { useUserStore } from '@/stores/user'
import { isAdminMode, mallUrl, openInNewTab } from '@/utils/portal'

const router = useRouter()
const store = useUserStore()
const formRef = ref<FormInstance>()
const loading = ref(false)

function goMall() {
  if (isAdminMode()) {
    openInNewTab(mallUrl('/mall'))
  } else {
    router.push('/mall')
  }
}

const form = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  nickname: '',
  email: '',
  phone: '',
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
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    {
      validator: (_rule, value, callback) => {
        if (value !== form.password) callback(new Error('两次输入的密码不一致'))
        else callback()
      },
      trigger: 'blur',
    },
  ],
  email: [{ type: 'email', message: '邮箱格式不正确', trigger: 'blur' }],
}

async function handleRegister() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    const res = await store.register({
      username: form.username.trim(),
      password: form.password,
      nickname: form.nickname.trim(),
      email: form.email.trim(),
      phone: form.phone.trim(),
    })
    if (res.code === 200) {
      ElMessage.success('注册成功')
      const ok = await store.login(form.username.trim(), form.password)
      router.push(ok ? '/mall' : '/login')
    } else if (res.code === 409) {
      ElMessage.error('用户名已存在')
    } else {
      ElMessage.error(res.msg || '注册失败')
    }
  } catch {
    ElMessage.error('服务器连接失败')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="auth-page">
    <section class="brand-panel">
      <div class="brand-panel-inner">
        <BrandLogo size="lg" subtitle="SHIGUANG SELECT" light />
        <div class="brand-copy">
          <p class="eyebrow">JOIN SHIGUANG</p>
          <h1>新客入席，<br />好物相伴</h1>
          <p class="desc">注册拾光严选，把喜欢的东西慢慢收进生活里。新客注册即享精选好物与专属服务。</p>
        </div>
        <div class="brand-points">
          <div><span class="seal-tag seal-tag--gold">新客</span><span>注册即享精选好物</span></div>
          <div><span class="seal-tag seal-tag--gold">包邮</span><span>满额包邮，48 小时发货</span></div>
          <div><span class="seal-tag seal-tag--gold">安心</span><span>七天无理由退换</span></div>
        </div>
        <div class="brand-seal">拾光<br />严选</div>
      </div>
    </section>

    <section class="form-panel">
      <div class="form-card">
        <p class="form-eyebrow">CREATE ACCOUNT</p>
        <h2>注册拾光严选</h2>
        <p class="form-sub">填写信息，开启你的拾光之旅</p>

        <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
          <el-row :gutter="14">
            <el-col :span="12">
              <el-form-item label="用户名" prop="username">
                <el-input v-model="form.username" placeholder="3-64 个字符" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="昵称">
                <el-input v-model="form.nickname" placeholder="选填" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="密码" prop="password">
                <el-input v-model="form.password" type="password" show-password placeholder="6-72 个字符" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="确认密码" prop="confirmPassword">
                <el-input v-model="form.confirmPassword" type="password" show-password placeholder="再次输入密码" />
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
          </el-row>
          <el-form-item>
            <el-button type="primary" :loading="loading" class="register-btn" @click="handleRegister">注 册</el-button>
          </el-form-item>
        </el-form>

        <div class="form-links">
          <el-link type="primary" @click="router.push('/login')">已有账号？去登录</el-link>
          <el-link @click="goMall">先去逛逛</el-link>
        </div>
        <div class="form-hint">管理后台请访问 :8080</div>
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
  max-width: 460px;
  padding: 36px 40px;
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
  margin: 8px 0 22px;
  color: var(--muted);
  font-size: 13px;
}

.register-btn {
  width: 100%;
  letter-spacing: 6px;
  height: 44px;
}

.form-links {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 4px;
}

.form-hint {
  margin-top: 20px;
  padding-top: 14px;
  border-top: 1px dashed rgba(201, 168, 106, 0.3);
  color: var(--muted);
  font-size: 11px;
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
    padding: 26px 20px;
  }
}
</style>
