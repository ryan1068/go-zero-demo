import { orderRequest } from '@/utils/request'
import type {
  AfterSaleInfo,
  AfterSaleListData,
  ApiResponse,
  CreateOrderReq,
  DashboardData,
  HandleAfterSaleParams,
  ListOrdersParams,
  OrderInfo,
  OrderListData,
  OrderTracksData,
} from '@/types/api'

export async function createOrder(data: CreateOrderReq): Promise<ApiResponse<OrderInfo>> {
  const res = await orderRequest.post<ApiResponse<OrderInfo>>('/orders', data)
  return res.data
}

export async function getOrder(id: number): Promise<OrderInfo> {
  const res = await orderRequest.get<ApiResponse<OrderInfo>>(`/order/${id}`)
  return res.data.data
}

export async function listOrders(params: ListOrdersParams): Promise<OrderListData> {
  const res = await orderRequest.get<ApiResponse<OrderListData>>('/orders', { params })
  return res.data.data
}

export async function payOrder(id: number): Promise<ApiResponse<OrderInfo>> {
  const res = await orderRequest.post<ApiResponse<OrderInfo>>(`/order/${id}/pay`)
  return res.data
}

export async function cancelOrder(id: number, reason = ''): Promise<ApiResponse<OrderInfo>> {
  const res = await orderRequest.post<ApiResponse<OrderInfo>>(`/order/${id}/cancel`, { reason })
  return res.data
}

export async function confirmOrder(id: number): Promise<ApiResponse<OrderInfo>> {
  const res = await orderRequest.post<ApiResponse<OrderInfo>>(`/order/${id}/confirm`)
  return res.data
}

export async function adminListOrders(params: ListOrdersParams): Promise<OrderListData> {
  const res = await orderRequest.get<ApiResponse<OrderListData>>('/admin/orders', { params })
  return res.data.data
}

export async function adminGetOrder(id: number): Promise<OrderInfo> {
  const res = await orderRequest.get<ApiResponse<OrderInfo>>(`/admin/order/${id}`)
  return res.data.data
}

export async function shipOrder(
  id: number,
  logistics_company: string,
  tracking_no: string,
): Promise<ApiResponse<OrderInfo>> {
  const res = await orderRequest.post<ApiResponse<OrderInfo>>(`/admin/order/${id}/ship`, {
    logistics_company,
    tracking_no,
  })
  return res.data
}

export async function adminCancelOrder(id: number, reason = ''): Promise<ApiResponse<OrderInfo>> {
  const res = await orderRequest.post<ApiResponse<OrderInfo>>(`/admin/order/${id}/cancel`, { reason })
  return res.data
}

export async function dashboard(): Promise<DashboardData> {
  const res = await orderRequest.get<ApiResponse<DashboardData>>('/admin/dashboard')
  return res.data.data
}

export async function adminListAfterSales(params: {
  page: number
  page_size: number
  status?: number
}): Promise<AfterSaleListData> {
  const res = await orderRequest.get<ApiResponse<AfterSaleListData>>('/admin/after-sales', { params })
  return res.data.data
}

export async function adminGetAfterSale(id: number): Promise<AfterSaleInfo> {
  const res = await orderRequest.get<ApiResponse<AfterSaleInfo>>(`/admin/after-sales/${id}`)
  return res.data.data
}

export async function adminHandleAfterSale(
  id: number,
  data: HandleAfterSaleParams,
): Promise<ApiResponse<AfterSaleInfo>> {
  const res = await orderRequest.post<ApiResponse<AfterSaleInfo>>(`/admin/after-sales/${id}/handle`, data)
  return res.data
}

export async function adminOrderTracks(id: number): Promise<OrderTracksData> {
  const res = await orderRequest.get<ApiResponse<OrderTracksData>>(`/admin/order/${id}/tracks`)
  return res.data.data
}
