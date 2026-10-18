<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ProductCard from '@/components/ProductCard.vue'
import { listCategories, listProducts } from '@/api/product'
import type { CategoryInfo, ProductInfo } from '@/types/api'

const route = useRoute()
const router = useRouter()
const categories = ref<CategoryInfo[]>([])
const topCategories = computed(() => categories.value.filter((cat) => cat.parent_id === 0))
const selectedCategory = computed(() => categories.value.find((cat) => cat.id === filters.category_id))
const activeParentId = computed(() =>
  selectedCategory.value ? selectedCategory.value.parent_id || selectedCategory.value.id : 0,
)
const products = ref<ProductInfo[]>([])
const total = ref(0)
const loading = ref(false)
const page = ref(1)
const pageSize = ref(12)
const filters = reactive({
  keyword: typeof route.query.keyword === 'string' ? route.query.keyword : '',
  category_id: Number(route.query.category_id || 0),
  sort: typeof route.query.sort === 'string' ? route.query.sort : 'default',
})

const sortOptions = [
  { label: '综合', value: 'default' },
  { label: '销量', value: 'sales' },
  { label: '价格从低到高', value: 'price_asc' },
  { label: '价格从高到低', value: 'price_desc' },
  { label: '最新上架', value: 'newest' },
]

async function load() {
  loading.value = true
  try {
    const data = await listProducts({
      page: page.value,
      page_size: pageSize.value,
      category_id: filters.category_id || undefined,
      keyword: filters.keyword.trim() || undefined,
      sort: filters.sort,
    })
    products.value = data.list ?? []
    total.value = data.total ?? 0
  } finally {
    loading.value = false
  }
}

function search() {
  page.value = 1
  router.replace({
    path: '/mall/products',
    query: {
      keyword: filters.keyword.trim() || undefined,
      category_id: filters.category_id || undefined,
      sort: filters.sort === 'default' ? undefined : filters.sort,
    },
  })
  load()
}

function reset() {
  filters.keyword = ''
  filters.category_id = 0
  filters.sort = 'default'
  page.value = 1
  router.replace({ path: '/mall/products' })
  load()
}

function selectCategory(id: number) {
  filters.category_id = id
  search()
}

function selectSort(value: string) {
  filters.sort = value
  search()
}

function handlePage(value: number) {
  page.value = value
  load()
}

onMounted(async () => {
  categories.value = await listCategories()
  await load()
})
</script>

<template>
  <div v-loading="loading" class="product-list-page">
    <div class="page-head">
      <div>
        <h1>全部好物</h1>
        <p>SHIGUANG PICKS · 认真挑过的每一件</p>
      </div>
      <div class="result-count">共 {{ total }} 件好物</div>
    </div>

    <section class="filter-panel paper-card">
      <div class="filter-row">
        <span class="filter-label">分类</span>
        <div class="chips">
          <button type="button" class="chip" :class="{ active: !filters.category_id }" @click="selectCategory(0)">全部</button>
          <button
            v-for="cat in topCategories"
            :key="cat.id"
            type="button"
            class="chip"
            :class="{ active: activeParentId === cat.id }"
            @click="selectCategory(cat.id)"
          >
            {{ cat.name }}
          </button>
        </div>
      </div>

      <div class="filter-row">
        <span class="filter-label">搜索</span>
        <div class="search-row">
          <el-input
            v-model="filters.keyword"
            placeholder="输入商品名称"
            clearable
            @keyup.enter="search"
          />
          <el-button type="primary" @click="search">搜 索</el-button>
          <el-button text @click="reset">重置</el-button>
        </div>
      </div>

      <div class="filter-row sort-row">
        <span class="filter-label">排序</span>
        <div class="sort-tabs">
          <button
            v-for="option in sortOptions"
            :key="option.value"
            type="button"
            :class="{ active: filters.sort === option.value }"
            @click="selectSort(option.value)"
          >
            {{ option.label }}
          </button>
        </div>
      </div>
    </section>

    <el-empty v-if="!products.length" description="没有找到好物，换个关键词试试">
      <el-button type="primary" @click="reset">看看全部好物</el-button>
    </el-empty>

    <div v-else class="product-grid">
      <ProductCard v-for="product in products" :key="product.id" :product="product" />
    </div>

    <div class="pagination">
      <el-pagination
        :current-page="page"
        :page-size="pageSize"
        :total="total"
        layout="prev, pager, next, total"
        background
        @current-change="handlePage"
      />
    </div>
  </div>
</template>

<style scoped>
.product-list-page {
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
  color: var(--ink);
}

.page-head p {
  margin: 6px 0 0;
  color: var(--muted);
  font-size: 12px;
  letter-spacing: 2px;
}

.result-count {
  color: var(--muted);
  font-size: 12px;
}

.filter-panel {
  padding: 18px 22px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.filter-row {
  display: flex;
  align-items: center;
  gap: 16px;
  min-height: 32px;
}

.filter-label {
  width: 44px;
  flex-shrink: 0;
  color: var(--muted);
  font-size: 13px;
}

.chips {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.chip {
  border: 1px solid var(--line);
  background: transparent;
  color: var(--ink-soft);
  padding: 6px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
  font-family: var(--font-sans);
  transition: all 0.18s ease;
}

.chip:hover {
  border-color: var(--brand);
  color: var(--brand);
}

.chip.active {
  background: var(--brand);
  border-color: var(--brand);
  color: #fffefa;
}

.search-row {
  display: flex;
  gap: 10px;
  flex: 1;
  max-width: 560px;
}

.sort-row {
  border-top: 1px dashed rgba(201, 168, 106, 0.28);
  padding-top: 14px;
}

.sort-tabs {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.sort-tabs button {
  border: none;
  background: transparent;
  color: var(--ink-soft);
  padding: 5px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
}

.sort-tabs button:hover {
  color: var(--brand);
}

.sort-tabs button.active {
  background: rgba(163, 59, 74, 0.08);
  color: var(--brand);
  font-weight: 600;
}

.product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 16px;
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 4px;
}

@media (max-width: 720px) {
  .page-head {
    align-items: flex-start;
    flex-direction: column;
  }

  .filter-row {
    align-items: flex-start;
    flex-direction: column;
    gap: 8px;
  }

  .filter-label {
    width: auto;
  }

  .search-row {
    width: 100%;
  }
}
</style>
