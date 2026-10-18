import { orderRequest } from '@/utils/request'
import type { ApiResponse, CartItemInfo, CartListData } from '@/types/api'

export async function listCart(): Promise<CartListData> {
  const res = await orderRequest.get<ApiResponse<CartListData>>('/cart')
  return res.data.data
}

export async function addCart(product_id: number, quantity: number): Promise<ApiResponse<CartItemInfo>> {
  const res = await orderRequest.post<ApiResponse<CartItemInfo>>('/cart/items', { product_id, quantity })
  return res.data
}

export async function updateCart(id: number, quantity: number, selected: number): Promise<ApiResponse<CartItemInfo>> {
  const res = await orderRequest.put<ApiResponse<CartItemInfo>>(`/cart/items/${id}`, { quantity, selected })
  return res.data
}

export async function deleteCart(id: number): Promise<ApiResponse> {
  const res = await orderRequest.delete<ApiResponse>(`/cart/items/${id}`)
  return res.data
}

export async function selectCart(selected: number): Promise<ApiResponse> {
  const res = await orderRequest.put<ApiResponse>('/cart/select', { selected })
  return res.data
}
