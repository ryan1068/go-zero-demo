import { orderRequest } from '@/utils/request'
import type { ApiResponse, CreateOrderReq, OrderInfo } from '@/types/api'

/**
 * 查询订单详情
 * GET /order/:id
 */
export async function getOrder(id: number): Promise<OrderInfo> {
  const res = await orderRequest.get<ApiResponse<OrderInfo>>(`/order/${id}`)
  return res.data.data
}

/**
 * 创建订单
 * POST /order/create
 */
export async function createOrder(data: CreateOrderReq): Promise<ApiResponse> {
  const res = await orderRequest.post<ApiResponse>('/order/create', data)
  return res.data
}
