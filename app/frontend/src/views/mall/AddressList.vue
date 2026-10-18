<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { createAddress, deleteAddress, listAddresses, setDefaultAddress, updateAddress } from '@/api/user'
import type { AddressInfo, SaveAddressReq } from '@/types/api'

const addresses = ref<AddressInfo[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const form = reactive<SaveAddressReq>({
  id: undefined,
  receiver_name: '',
  receiver_phone: '',
  province: '',
  city: '',
  district: '',
  detail_address: '',
  postal_code: '',
  tag: '',
  is_default: 0,
})

const rules: FormRules = {
  receiver_name: [{ required: true, message: '请输入收货人', trigger: 'blur' }],
  receiver_phone: [{ required: true, message: '请输入手机号', trigger: 'blur' }],
  detail_address: [{ required: true, message: '请输入详细地址', trigger: 'blur' }],
}

async function load() {
  loading.value = true
  try {
    const data = await listAddresses()
    addresses.value = data.list ?? []
  } finally {
    loading.value = false
  }
}

function openCreate() {
  Object.assign(form, {
    id: undefined,
    receiver_name: '',
    receiver_phone: '',
    province: '',
    city: '',
    district: '',
    detail_address: '',
    postal_code: '',
    tag: '',
    is_default: addresses.value.length ? 0 : 1,
  })
  dialogVisible.value = true
}

function openEdit(addr: AddressInfo) {
  Object.assign(form, {
    id: addr.id,
    receiver_name: addr.receiver_name,
    receiver_phone: addr.receiver_phone,
    province: addr.province,
    city: addr.city,
    district: addr.district,
    detail_address: addr.detail_address,
    postal_code: addr.postal_code,
    tag: addr.tag,
    is_default: addr.is_default,
  })
  dialogVisible.value = true
}

async function submit() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    const payload: SaveAddressReq = {
      receiver_name: form.receiver_name.trim(),
      receiver_phone: form.receiver_phone.trim(),
      province: form.province.trim(),
      city: form.city.trim(),
      district: form.district.trim(),
      detail_address: form.detail_address.trim(),
      postal_code: form.postal_code.trim(),
      tag: form.tag.trim(),
      is_default: form.is_default,
    }
    const res = form.id
      ? await updateAddress(form.id, payload)
      : await createAddress(payload)
    if (res.code === 200) {
      ElMessage.success('保存成功')
      dialogVisible.value = false
      await load()
    } else {
      ElMessage.error(res.msg || '保存失败')
    }
  } finally {
    submitting.value = false
  }
}

async function remove(addr: AddressInfo) {
  try {
    await ElMessageBox.confirm('确定删除这个地址吗？', '提示', { type: 'warning' })
  } catch {
    return
  }
  const res = await deleteAddress(addr.id)
  if (res.code === 200) {
    ElMessage.success('已删除')
    await load()
  } else {
    ElMessage.error(res.msg || '删除失败')
  }
}

async function makeDefault(addr: AddressInfo) {
  const res = await setDefaultAddress(addr.id)
  if (res.code === 200) {
    ElMessage.success('已设为默认地址')
    await load()
  } else {
    ElMessage.error(res.msg || '操作失败')
  }
}

onMounted(load)
</script>

<template>
  <div v-loading="loading" class="address-page">
    <div class="page-head">
      <div>
        <h1>收货地址</h1>
        <p>ADDRESS BOOK · 把好物送到你手中</p>
      </div>
      <el-button type="primary" @click="openCreate">新增地址</el-button>
    </div>

    <el-empty v-if="!addresses.length" description="还没有收货地址，添加一个吧">
      <el-button type="primary" @click="openCreate">新增地址</el-button>
    </el-empty>

    <div v-else class="address-grid">
      <article v-for="addr in addresses" :key="addr.id" class="address-card paper-card hover-lift">
        <div class="card-top">
          <div class="receiver">
            <b>{{ addr.receiver_name }}</b>
            <span>{{ addr.receiver_phone }}</span>
          </div>
          <span v-if="addr.is_default === 1" class="seal-tag">默认地址</span>
        </div>
        <p class="detail">{{ addr.province }}{{ addr.city }}{{ addr.district }}{{ addr.detail_address }}</p>
        <div class="card-meta">
          <span v-if="addr.tag" class="seal-tag seal-tag--gold">{{ addr.tag }}</span>
          <span v-if="addr.postal_code" class="postal">邮编 {{ addr.postal_code }}</span>
        </div>
        <div class="card-actions">
          <el-button v-if="addr.is_default !== 1" text type="primary" @click="makeDefault(addr)">设为默认</el-button>
          <el-button text @click="openEdit(addr)">编辑</el-button>
          <el-button text type="danger" @click="remove(addr)">删除</el-button>
        </div>
      </article>
    </div>

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑地址' : '新增地址'" width="560px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="收货人" prop="receiver_name">
              <el-input v-model="form.receiver_name" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="手机号" prop="receiver_phone">
              <el-input v-model="form.receiver_phone" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="省">
              <el-input v-model="form.province" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="市">
              <el-input v-model="form.city" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="区/县">
              <el-input v-model="form.district" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="详细地址" prop="detail_address">
              <el-input v-model="form.detail_address" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="邮编">
              <el-input v-model="form.postal_code" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="标签">
              <el-input v-model="form.tag" placeholder="家/公司/学校" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="默认地址">
              <el-switch v-model="form.is_default" :active-value="1" :inactive-value="0" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.address-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 16px;
}

.page-head h1 {
  margin: 0;
  font-family: var(--font-serif);
  font-size: 30px;
  letter-spacing: 4px;
}

.page-head p {
  margin: 6px 0 0;
  color: var(--muted);
  font-size: 12px;
  letter-spacing: 2px;
}

.address-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
}

.address-card {
  padding: 18px 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.card-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.receiver {
  display: flex;
  align-items: baseline;
  gap: 10px;
  color: var(--ink);
}

.receiver b {
  font-family: var(--font-serif);
  font-size: 18px;
  letter-spacing: 1px;
}

.receiver span {
  color: var(--muted);
  font-size: 13px;
}

.detail {
  margin: 0;
  color: var(--ink-soft);
  font-size: 13px;
  line-height: 1.7;
  min-height: 44px;
}

.card-meta {
  display: flex;
  align-items: center;
  gap: 10px;
}

.postal {
  color: var(--muted);
  font-size: 11px;
}

.card-actions {
  display: flex;
  justify-content: flex-end;
  gap: 2px;
  border-top: 1px dashed rgba(201, 168, 106, 0.28);
  padding-top: 10px;
}

@media (max-width: 720px) {
  .page-head {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
