<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { deleteCategory, listCategories, saveCategory } from '@/api/product'
import type { CategoryInfo, SaveCategoryReq } from '@/types/api'

type CategoryNode = CategoryInfo & { children?: CategoryInfo[] }

const categories = ref<CategoryInfo[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const form = reactive<SaveCategoryReq>({
  id: undefined,
  parent_id: 0,
  name: '',
  icon: '',
  sort: 0,
  status: 1,
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }],
}

const treeData = computed<CategoryNode[]>(() => {
  const roots = categories.value.filter((cat) => cat.parent_id === 0)
  return roots.map((root) => {
    const children = categories.value.filter((cat) => cat.parent_id === root.id)
    return children.length ? { ...root, children } : { ...root }
  })
})

const topCategories = computed(() => categories.value.filter((cat) => cat.parent_id === 0))
const parentOptions = computed(() => topCategories.value.filter((cat) => cat.id !== form.id))
const parentDisabled = computed(
  () => !!form.id && categories.value.some((cat) => cat.parent_id === form.id),
)

function parentName(row: CategoryInfo): string {
  if (!row.parent_id) return '一级分类'
  return categories.value.find((cat) => cat.id === row.parent_id)?.name || '-'
}

function hasChildren(row: CategoryInfo): boolean {
  return categories.value.some((cat) => cat.parent_id === row.id)
}

async function load() {
  loading.value = true
  try {
    categories.value = await listCategories(true)
  } finally {
    loading.value = false
  }
}

function openCreate() {
  Object.assign(form, { id: undefined, parent_id: 0, name: '', icon: '', sort: 0, status: 1 })
  dialogVisible.value = true
}

function openEdit(row: CategoryInfo) {
  Object.assign(form, {
    id: row.id,
    parent_id: hasChildren(row) ? 0 : row.parent_id,
    name: row.name,
    icon: row.icon,
    sort: row.sort,
    status: row.status,
  })
  dialogVisible.value = true
}

async function submit() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    const res = await saveCategory({
      id: form.id,
      parent_id: parentDisabled.value ? 0 : form.parent_id,
      name: form.name.trim(),
      icon: form.icon.trim(),
      sort: form.sort,
      status: form.status,
    })
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

async function remove(row: CategoryInfo) {
  if (hasChildren(row)) {
    ElMessage.warning('该分类下还有子分类，请先删除子分类')
    return
  }
  try {
    await ElMessageBox.confirm('确定删除这个分类吗？', '提示', { type: 'warning' })
  } catch {
    return
  }
  const res = await deleteCategory(row.id)
  if (res.code === 200) {
    ElMessage.success('已删除')
    await load()
  } else {
    ElMessage.error(res.msg || '删除失败')
  }
}

onMounted(load)
</script>

<template>
  <div v-loading="loading">
    <div class="page-header">
      <h3>分类管理</h3>
      <el-button type="primary" @click="openCreate">新增分类</el-button>
    </div>

    <el-card shadow="never">
      <el-table
        :data="treeData"
        row-key="id"
        :tree-props="{ children: 'children' }"
        default-expand-all
        border
        stripe
      >
        <el-table-column prop="name" label="分类名称" min-width="220" />
        <el-table-column label="父分类" min-width="140">
          <template #default="{ row }">{{ parentName(row) }}</template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="90" />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160">
          <template #default="{ row }">
            <el-button text type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button text type="danger" @click="remove(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑分类' : '新增分类'" width="480px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
        <el-form-item label="分类名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="父分类">
          <el-select v-model="form.parent_id" :disabled="parentDisabled" style="width: 100%">
            <el-option label="无（一级分类）" :value="0" />
            <el-option v-for="cat in parentOptions" :key="cat.id" :label="cat.name" :value="cat.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="图标地址">
          <el-input v-model="form.icon" placeholder="选填" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="form.status" :active-value="1" :inactive-value="0" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.page-header h3 { margin: 0; }
</style>
