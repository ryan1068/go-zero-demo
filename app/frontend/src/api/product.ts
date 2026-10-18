import { productRequest } from '@/utils/request'
import type {
  ApiResponse,
  CategoryInfo,
  ListProductsParams,
  ProductInfo,
  ProductListData,
  ProductStats,
  SaveCategoryReq,
  SaveProductReq,
  SaveSkuReq,
  SkuInfo,
} from '@/types/api'

export async function listCategories(admin = false): Promise<CategoryInfo[]> {
  const path = admin ? '/admin/categories' : '/categories'
  const res = await productRequest.get<ApiResponse<CategoryInfo[]>>(path)
  return res.data.data ?? []
}

export async function saveCategory(data: SaveCategoryReq): Promise<ApiResponse<CategoryInfo>> {
  const res = data.id
    ? await productRequest.put<ApiResponse<CategoryInfo>>(`/admin/categories/${data.id}`, data)
    : await productRequest.post<ApiResponse<CategoryInfo>>('/admin/categories', data)
  return res.data
}

export async function deleteCategory(id: number): Promise<ApiResponse> {
  const res = await productRequest.delete<ApiResponse>(`/admin/categories/${id}`)
  return res.data
}

export async function listProducts(params: ListProductsParams, admin = false): Promise<ProductListData> {
  const path = admin ? '/admin/products' : '/products'
  const res = await productRequest.get<ApiResponse<ProductListData>>(path, { params })
  return res.data.data
}

export async function getProduct(id: number, admin = false): Promise<ProductInfo> {
  const path = admin ? `/admin/products/${id}` : `/product/${id}`
  const res = await productRequest.get<ApiResponse<ProductInfo>>(path)
  return res.data.data
}

export async function saveProduct(data: SaveProductReq): Promise<ApiResponse<ProductInfo>> {
  const res = data.id
    ? await productRequest.put<ApiResponse<ProductInfo>>(`/admin/products/${data.id}`, data)
    : await productRequest.post<ApiResponse<ProductInfo>>('/admin/products', data)
  return res.data
}

export async function deleteProduct(id: number): Promise<ApiResponse> {
  const res = await productRequest.delete<ApiResponse>(`/admin/products/${id}`)
  return res.data
}

export async function updateProductStatus(id: number, status: number): Promise<ApiResponse> {
  const res = await productRequest.put<ApiResponse>(`/admin/products/${id}/status`, { status })
  return res.data
}

export async function productStats(): Promise<ProductStats> {
  const res = await productRequest.get<ApiResponse<ProductStats>>('/admin/products/stats')
  return res.data.data
}

export async function uploadImage(file: File): Promise<{ code: number; msg: string; url?: string }> {
  const form = new FormData()
  form.append('file', file)
  const res = await productRequest.post<{ code: number; msg: string; url?: string }>('/admin/upload', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
  return res.data
}

export async function uploadImages(files: File[]): Promise<{ code: number; msg: string; urls: string[] }> {
  const form = new FormData()
  for (const file of files) {
    form.append('files', file)
  }
  const res = await productRequest.post<{ code: number; msg: string; urls: string[] }>(
    '/admin/uploads/batch',
    form,
    { headers: { 'Content-Type': 'multipart/form-data' } },
  )
  return res.data
}

export async function deleteUpload(name: string): Promise<ApiResponse> {
  const res = await productRequest.delete<ApiResponse>(`/admin/uploads/${encodeURIComponent(name)}`)
  return res.data
}

export async function listSkus(productId: number): Promise<SkuInfo[]> {
  const res = await productRequest.get<ApiResponse<SkuInfo[]>>(`/admin/products/${productId}/skus`)
  return res.data.data ?? []
}

export async function saveSku(data: SaveSkuReq): Promise<ApiResponse<SkuInfo>> {
  const res = await productRequest.post<ApiResponse<SkuInfo>>('/admin/skus', data)
  return res.data
}

export async function deleteSku(id: number): Promise<ApiResponse> {
  const res = await productRequest.delete<ApiResponse>(`/admin/skus/${id}`)
  return res.data
}
