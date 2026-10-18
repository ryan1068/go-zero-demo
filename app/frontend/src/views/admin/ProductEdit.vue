<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import {
  deleteSku,
  getProduct,
  listCategories,
  listSkus,
  saveProduct,
  saveSku,
} from '@/api/product'
import { deleteMedia, uploadImage, uploadImages } from '@/api/media'
import type { CategoryInfo, SaveSkuReq, SkuInfo } from '@/types/api'

const route = useRoute()
const router = useRouter()
const productId = computed(() => Number(route.params.id || 0))
const isEdit = computed(() => productId.value > 0)
const categories = ref<CategoryInfo[]>([])
const topCategories = computed(() => categories.value.filter((cat) => cat.parent_id === 0))
const childrenOf = (parentId: number) => categories.value.filter((cat) => cat.parent_id === parentId)
const formRef = ref<FormInstance>()
const loading = ref(false)
const submitting = ref(false)
const uploading = ref(false)
const batchUploading = ref(false)
const imageInput = ref('')

const form = reactive({
  category_id: 0,
  name: '',
  subtitle: '',
  cover_image: '',
  images: [] as string[],
  detail: '',
  priceYuan: 0,
  originalPriceYuan: 0,
  stock: 0,
  unit: '件',
  status: 1,
  is_recommend: 0,
  sort: 0,
})

const rules: FormRules = {
  category_id: [{ required: true, message: '请选择分类', trigger: 'change' }],
  name: [{ required: true, message: '请输入商品名称', trigger: 'blur' }],
  priceYuan: [{ required: true, message: '请输入售价', trigger: 'blur' }],
}

const skus = ref<SkuInfo[]>([])
const skuLoading = ref(false)
const skuDialogVisible = ref(false)
const skuFormRef = ref<FormInstance>()
const skuSubmitting = ref(false)
const editingSkuId = ref(0)
const skuForm = reactive({
  sku_code: '',
  specs: '{}',
  priceYuan: 0,
  originalPriceYuan: 0,
  stock: 0,
  status: 1,
})
const skuRules: FormRules = {
  sku_code: [{ required: true, message: '请输入 SKU 编码', trigger: 'blur' }],
  priceYuan: [{ required: true, message: '请输入售价', trigger: 'blur' }],
}

async function load() {
  loading.value = true
  try {
    categories.value = await listCategories(true)
    if (isEdit.value) {
      const p = await getProduct(productId.value, true)
      Object.assign(form, {
        category_id: p.category_id,
        name: p.name,
        subtitle: p.subtitle,
        cover_image: p.cover_image,
        images: p.images ?? [],
        detail: p.detail,
        priceYuan: p.price / 100,
        originalPriceYuan: p.original_price / 100,
        stock: p.stock,
        unit: p.unit || '件',
        status: p.status,
        is_recommend: p.is_recommend,
        sort: p.sort,
      })
      await loadSkus()
    }
  } finally {
    loading.value = false
  }
}

function appendImages(urls: string[]) {
  const next = urls.filter((url) => url && !form.images.includes(url))
  if (!next.length) return
  form.images.push(...next)
  if (!form.cover_image) form.cover_image = form.images[0]
}

async function handleUpload(options: any) {
  uploading.value = true
  try {
    const res = await uploadImage(options.file)
    if (res.code === 200 && res.url) {
      appendImages([res.url])
      ElMessage.success('图片上传成功')
    } else {
      ElMessage.error(res.msg || '上传失败')
    }
  } catch {
    ElMessage.error('上传失败')
  } finally {
    uploading.value = false
  }
}

async function handleBatchFiles(event: Event) {
  const input = event.target as HTMLInputElement
  const files = Array.from(input.files || [])
  if (!files.length) return
  batchUploading.value = true
  try {
    const res = await uploadImages(files)
    if (res.code === 200 && res.urls?.length) {
      appendImages(res.urls)
      ElMessage.success(`已上传 ${res.urls.length} 张图片`)
    } else {
      ElMessage.error(res.msg || '批量上传失败')
    }
  } catch {
    ElMessage.error('批量上传失败')
  } finally {
    batchUploading.value = false
    input.value = ''
  }
}

function setCover(url: string) {
  form.cover_image = url
  ElMessage.success('已设为封面')
}

async function removeImage(url: string) {
  form.images = form.images.filter((item) => item !== url)
  if (form.cover_image === url) {
    form.cover_image = form.images[0] || ''
  }
  if (url.includes('/uploads/')) {
    const name = url.split('/').pop()
    if (name) {
      try {
        await deleteMedia(name)
      } catch {
        // 文件删除失败不影响商品图片列表。
      }
    }
  }
}

function addImageUrl() {
  const url = imageInput.value.trim()
  if (!url) return
  appendImages([url])
  imageInput.value = ''
}

async function loadSkus() {
  if (!isEdit.value) return
  skuLoading.value = true
  try {
    skus.value = await listSkus(productId.value)
  } finally {
    skuLoading.value = false
  }
}

function openSkuDialog(row?: SkuInfo) {
  editingSkuId.value = row?.id || 0
  Object.assign(skuForm, {
    sku_code: row?.sku_code || '',
    specs: row?.specs || '{}',
    priceYuan: row ? row.price / 100 : form.priceYuan,
    originalPriceYuan: row ? row.original_price / 100 : form.originalPriceYuan,
    stock: row?.stock ?? form.stock,
    status: row?.status ?? 1,
  })
  skuDialogVisible.value = true
}

async function submitSku() {
  if (!skuFormRef.value) return
  const valid = await skuFormRef.value.validate().catch(() => false)
  if (!valid) return
  skuSubmitting.value = true
  try {
    const payload: SaveSkuReq = {
      id: editingSkuId.value || undefined,
      product_id: productId.value,
      sku_code: skuForm.sku_code.trim(),
      specs: skuForm.specs.trim() || '{}',
      price: Math.round(skuForm.priceYuan * 100),
      original_price: Math.round(skuForm.originalPriceYuan * 100),
      stock: skuForm.stock,
      status: skuForm.status,
    }
    const res = await saveSku(payload)
    if (res.code === 200) {
      ElMessage.success('SKU 保存成功')
      skuDialogVisible.value = false
      await loadSkus()
    } else {
      ElMessage.error(res.msg || '保存失败')
    }
  } finally {
    skuSubmitting.value = false
  }
}

async function removeSku(row: SkuInfo) {
  try {
    await ElMessageBox.confirm(`确定删除 SKU ${row.sku_code} 吗？`, '提示', { type: 'warning' })
  } catch {
    return
  }
  const res = await deleteSku(row.id)
  if (res.code === 200) {
    ElMessage.success('已删除')
    await loadSkus()
  } else {
    ElMessage.error(res.msg || '删除失败')
  }
}

async function submit() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    const res = await saveProduct({
      id: isEdit.value ? productId.value : undefined,
      category_id: form.category_id,
      name: form.name.trim(),
      subtitle: form.subtitle.trim(),
      cover_image: form.cover_image,
      images: form.images,
      detail: form.detail,
      price: Math.round(form.priceYuan * 100),
      original_price: Math.round(form.originalPriceYuan * 100),
      stock: form.stock,
      unit: form.unit.trim() || '件',
      status: form.status,
      is_recommend: form.is_recommend,
      sort: form.sort,
    })
    if (res.code === 200) {
      ElMessage.success('保存成功')
      router.push('/product')
    } else {
      ElMessage.error(res.msg || '保存失败')
    }
  } finally {
    submitting.value = false
  }
}

onMounted(load)
</script>

<template>
  <div v-loading="loading">
    <div class="page-header">
      <h3>{{ isEdit ? '编辑商品' : '新增商品' }}</h3>
      <el-button @click="router.push('/product')">返回列表</el-button>
    </div>

    <el-card shadow="never">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="商品名称" prop="name">
              <el-input v-model="form.name" maxlength="128" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="分类" prop="category_id">
              <el-select v-model="form.category_id" placeholder="请选择分类" style="width: 100%">
                <el-option-group v-for="parent in topCategories" :key="parent.id" :label="parent.name">
                  <el-option :label="parent.name" :value="parent.id" />
                  <el-option
                    v-for="child in childrenOf(parent.id)"
                    :key="child.id"
                    :label="`${parent.name} / ${child.name}`"
                    :value="child.id"
                  />
                </el-option-group>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="副标题">
              <el-input v-model="form.subtitle" maxlength="255" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="单位">
              <el-input v-model="form.unit" maxlength="16" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="售价（元）" prop="priceYuan">
              <el-input-number v-model="form.priceYuan" :min="0.01" :precision="2" :step="1" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="原价（元）">
              <el-input-number v-model="form.originalPriceYuan" :min="0" :precision="2" :step="1" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="库存">
              <el-input-number v-model="form.stock" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="排序">
              <el-input-number v-model="form.sort" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="状态">
              <el-switch v-model="form.status" :active-value="1" :inactive-value="0" active-text="上架" inactive-text="下架" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="首页推荐">
              <el-switch v-model="form.is_recommend" :active-value="1" :inactive-value="0" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="商品图片">
              <div class="image-editor">
                <div class="upload-actions">
                  <el-upload :show-file-list="false" :http-request="handleUpload" accept="image/*">
                    <el-button :loading="uploading">上传单张</el-button>
                  </el-upload>
                  <label class="batch-label">
                    <input type="file" accept="image/*" multiple hidden @change="handleBatchFiles" />
                    <el-button :loading="batchUploading">批量上传</el-button>
                  </label>
                  <el-input v-model="imageInput" placeholder="或粘贴图片 URL" style="width: 260px" @keyup.enter="addImageUrl" />
                  <el-button @click="addImageUrl">添加</el-button>
                </div>
                <div class="image-list">
                  <div v-for="image in form.images" :key="image" class="image-item" :class="{ cover: image === form.cover_image }">
                    <el-image :src="image" fit="cover" class="thumb" />
                    <div class="image-actions">
                      <el-button text size="small" type="primary" @click="setCover(image)">封面</el-button>
                      <el-button text size="small" type="danger" @click="removeImage(image)">删除</el-button>
                    </div>
                    <span v-if="image === form.cover_image" class="cover-badge">封面</span>
                  </div>
                  <span v-if="!form.images.length" class="empty-hint">还没有图片，支持一次上传多张</span>
                </div>
                <el-input v-model="form.cover_image" placeholder="封面图 URL" style="width: 520px; margin-top: 10px" />
              </div>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="商品详情">
              <el-input v-model="form.detail" type="textarea" :rows="6" maxlength="20000" show-word-limit />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="submit">保存商品</el-button>
          <el-button @click="router.push('/product')">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card v-if="isEdit" shadow="never" class="sku-card">
      <template #header>
        <div class="sku-header">
          <div>
            <b>SKU / 多规格</b>
            <span class="sku-tip">同一商品可配置颜色、尺寸等规格；下单接入 SKU 仍在下一步</span>
          </div>
          <el-button type="primary" @click="openSkuDialog()">新增 SKU</el-button>
        </div>
      </template>
      <el-table v-loading="skuLoading" :data="skus" border stripe>
        <el-table-column prop="sku_code" label="SKU 编码" min-width="160" />
        <el-table-column prop="specs" label="规格" min-width="220" show-overflow-tooltip />
        <el-table-column label="售价" width="110">
          <template #default="{ row }">¥{{ (row.price / 100).toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="stock" label="库存" width="90" />
        <el-table-column prop="sales" label="销量" width="90" />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '停用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button text type="primary" @click="openSkuDialog(row)">编辑</el-button>
            <el-button text type="danger" @click="removeSku(row)">删除</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty description="还没有 SKU，点击右上角新增" :image-size="70" />
        </template>
      </el-table>
    </el-card>
    <el-card v-else shadow="never" class="sku-card">
      <el-empty description="保存商品后可以添加 SKU 规格" :image-size="70" />
    </el-card>

    <el-dialog v-model="skuDialogVisible" :title="editingSkuId ? '编辑 SKU' : '新增 SKU'" width="560px" destroy-on-close>
      <el-form ref="skuFormRef" :model="skuForm" :rules="skuRules" label-width="100px">
        <el-form-item label="SKU 编码" prop="sku_code">
          <el-input v-model="skuForm.sku_code" placeholder="例如 GO-BOOK-RED-M" />
        </el-form-item>
        <el-form-item label="规格 JSON">
          <el-input v-model="skuForm.specs" placeholder='例如 {"颜色":"胭脂红","尺寸":"M"}' />
        </el-form-item>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="售价（元）" prop="priceYuan">
              <el-input-number v-model="skuForm.priceYuan" :min="0.01" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="原价（元）">
              <el-input-number v-model="skuForm.originalPriceYuan" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="库存">
              <el-input-number v-model="skuForm.stock" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态">
              <el-switch v-model="skuForm.status" :active-value="1" :inactive-value="0" active-text="启用" inactive-text="停用" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="skuDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="skuSubmitting" @click="submitSku">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.page-header h3 { margin: 0; }
.image-editor { width: 100%; }
.upload-actions { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
.batch-label { display: inline-flex; }
.image-list { display: flex; gap: 12px; flex-wrap: wrap; margin-top: 14px; }
.image-item { position: relative; width: 96px; padding: 4px; border: 1px solid #eee; border-radius: 8px; background: #fafafa; }
.image-item.cover { border-color: #a33b4a; box-shadow: 0 0 0 1px #a33b4a inset; }
.thumb { width: 86px; height: 86px; border-radius: 5px; }
.image-actions { display: flex; justify-content: center; }
.image-actions :deep(.el-button) { padding: 2px 4px; font-size: 11px; }
.cover-badge { position: absolute; top: -1px; left: -1px; padding: 1px 5px; border-radius: 6px 0 6px 0; background: #a33b4a; color: #fff; font-size: 10px; }
.empty-hint { color: #999; font-size: 12px; line-height: 96px; }
.sku-card { margin-top: 16px; }
.sku-header { display: flex; align-items: center; justify-content: space-between; gap: 16px; }
.sku-tip { margin-left: 10px; color: #999; font-size: 12px; }
</style>
