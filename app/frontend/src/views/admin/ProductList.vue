<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { deleteProduct, listCategories, listProducts, updateProductStatus } from '@/api/product'
import { formatPrice } from '@/utils/format'
import type { CategoryInfo, ProductInfo } from '@/types/api'

const router = useRouter()
const products = ref<ProductInfo[]>([])
const categories = ref<CategoryInfo[]>([])
const topCategories = computed(() => categories.value.filter((cat) => cat.parent_id === 0))
const childrenOf = (parentId: number) => categories.value.filter((cat) => cat.parent_id === parentId)
const categoryOptions = computed(() =>
  topCategories.value.map((parent) => ({
    value: parent.id,
    label: parent.name,
    children: childrenOf(parent.id).map((child) => ({ value: child.id, label: child.name })),
  })),
)
const categoryPath = ref<number[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const filters = reactive({
  keyword: '',
  category_id: 0,
  status: 0,
  low_stock: 0,
})

async function load() {
  loading.value = true
  try {
    const data = await listProducts({
      page: page.value,
      page_size: pageSize.value,
      keyword: filters.keyword.trim() || undefined,
      category_id: filters.category_id || undefined,
      status: filters.status || undefined,
      low_stock: filters.low_stock || undefined,
    }, true)
    products.value = data.list ?? []
    total.value = data.total ?? 0
  } finally {
    loading.value = false
  }
}

function search() {
  page.value = 1
  load()
}

function handleCategoryChange(value: number[] | null) {
  const path = Array.isArray(value) ? value : []
  filters.category_id = path.length ? Number(path[path.length - 1]) : 0
  search()
}

function reset() {
  filters.keyword = ''
  filters.category_id = 0
  filters.status = 0
  filters.low_stock = 0
  categoryPath.value = []
  page.value = 1
  load()
}

async function toggleStatus(row: ProductInfo) {
  const next = row.status === 1 ? 0 : 1
  const res = await updateProductStatus(row.id, next)
  if (res.code === 200) {
    ElMessage.success(next === 1 ? '已上架' : '已下架')
    await load()
  } else {
    ElMessage.error(res.msg || '操作失败')
  }
}

async function remove(row: ProductInfo) {
  try {
    await ElMessageBox.confirm(`确定删除商品「${row.name}」吗？`, '提示', { type: 'warning' })
  } catch {
    return
  }
  const res = await deleteProduct(row.id)
  if (res.code === 200) {
    ElMessage.success('已删除')
    await load()
  } else {
    ElMessage.error(res.msg || '删除失败')
  }
}

function handlePage(value: number) {
  page.value = value
  load()
}

onMounted(async () => {
  categories.value = await listCategories(true)
  await load()
})
</script>

<template>
  <div v-loading="loading">
    <div class="page-header">
      <h3>商品管理</h3>
      <el-button type="primary" @click="router.push('/product/create')">新增商品</el-button>
    </div>

    <el-card shadow="never" class="filter-card">
      <div class="filters">
        <el-input v-model="filters.keyword" placeholder="商品名称" clearable style="width: 200px" @keyup.enter="search" />
        <el-cascader
          v-model="categoryPath"
          :options="categoryOptions"
          :props="{ checkStrictly: true, expandTrigger: 'hover' }"
          placeholder="全部分类"
          clearable
          style="width: 220px"
          @change="handleCategoryChange"
        />
        <el-select v-model="filters.status" style="width: 130px">
          <el-option label="全部状态" :value="0" />
          <el-option label="上架" :value="1" />
          <el-option label="下架" :value="2" />
        </el-select>
        <el-checkbox v-model="filters.low_stock" :true-value="1" :false-value="0">只看低库存</el-checkbox>
        <el-button type="primary" @click="search">查询</el-button>
        <el-button @click="reset">重置</el-button>
      </div>
    </el-card>

    <el-card shadow="never">
      <el-table :data="products" border stripe>
        <el-table-column label="封面" width="80">
          <template #default="{ row }">
            <img :src="row.cover_image || '/products/book.svg'" class="cover" />
          </template>
        </el-table-column>
        <el-table-column prop="name" label="商品名称" min-width="180" />
        <el-table-column prop="category_name" label="分类" width="100" />
        <el-table-column label="售价" width="110">
          <template #default="{ row }">¥{{ formatPrice(row.price) }}</template>
        </el-table-column>
        <el-table-column prop="stock" label="库存" width="90" />
        <el-table-column prop="sales" label="销量" width="90" />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '上架' : '下架' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220">
          <template #default="{ row }">
            <el-button text type="primary" @click="router.push(`/product/${row.id}/edit`)">编辑</el-button>
            <el-button text @click="toggleStatus(row)">{{ row.status === 1 ? '下架' : '上架' }}</el-button>
            <el-button text type="danger" @click="remove(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <el-pagination :current-page="page" :page-size="pageSize" :total="total" layout="prev, pager, next, total" background @current-change="handlePage" />
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.page-header h3 { margin: 0; }
.filter-card { margin-bottom: 16px; }
.filters { display: flex; gap: 12px; flex-wrap: wrap; align-items: center; }
.cover { width: 46px; height: 46px; object-fit: contain; background: #fafafa; border-radius: 6px; }
.pagination { display: flex; justify-content: flex-end; margin-top: 16px; }
</style>
