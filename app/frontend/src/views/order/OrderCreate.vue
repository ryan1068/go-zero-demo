<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useOrderStore } from '@/stores/order'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const store = useOrderStore()
const formRef = ref()
const loading = ref(false)

const form = reactive({
  goods_id: 1001,
  goods_name: '',
  price: 99,
})

const rules = {
  goods_name: [{ required: true, message: '请输入商品名称', trigger: 'blur' }],
  price: [{ required: true, message: '请输入价格', trigger: 'blur' }],
}

async function handleCreate() {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    const res = await store.submitOrder({
      goods_id: form.goods_id,
      goods_name: form.goods_name,
      price: form.price,
      created_at: Math.floor(Date.now() / 1000),
    })
    if (res.code === 200) {
      ElMessage.success('订单创建成功！')
    } else {
      ElMessage.error(res.msg || '创建失败')
    }
  } catch {
    ElMessage.error('服务器连接失败')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="page">
    <el-card shadow="never">
      <template #header>
        <span class="card-title">创建订单</span>
      </template>

      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px" style="max-width: 500px">
        <el-form-item label="商品 ID">
          <el-input-number v-model="form.goods_id" :min="1" />
        </el-form-item>

        <el-form-item label="商品名称" prop="goods_name">
          <el-input v-model="form.goods_name" placeholder="请输入商品名称" />
        </el-form-item>

        <el-form-item label="价格" prop="price">
          <el-input-number v-model="form.price" :min="0" :precision="2" />
        </el-form-item>

        <el-form-item>
          <el-button type="success" :icon="Plus" :loading="loading" @click="handleCreate">
            创建订单
          </el-button>
        </el-form-item>
      </el-form>

      <el-divider />
      <el-text type="info" size="small">
        调用链：order-api → order-rpc → user-rpc → MySQL
      </el-text>
    </el-card>
  </div>
</template>

<style scoped>
.page {
  max-width: 100%;
}
.card-title {
  font-weight: 600;
}
</style>
