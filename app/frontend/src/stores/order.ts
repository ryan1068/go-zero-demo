import { defineStore } from 'pinia'
import { ref } from 'vue'
import { createOrder, getOrder, listOrders } from '@/api/order'
import type { CreateOrderReq, ListOrdersParams, OrderInfo } from '@/types/api'

export const useOrderStore = defineStore('order', () => {
  const order = ref<OrderInfo | null>(null)
  const loading = ref(false)
  const error = ref('')

  const orders = ref<OrderInfo[]>([])
  const total = ref(0)
  const listLoading = ref(false)
  const listError = ref('')

  async function fetchOrder(id: number) {
    loading.value = true
    error.value = ''
    try {
      order.value = await getOrder(id)
    } catch (err) {
      order.value = null
      error.value = '查询失败：订单不存在或无权访问'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchOrders(params: ListOrdersParams) {
    listLoading.value = true
    listError.value = ''
    try {
      const data = await listOrders(params)
      orders.value = data.list ?? []
      total.value = data.total ?? 0
      return data
    } catch (err) {
      orders.value = []
      total.value = 0
      listError.value = '订单列表加载失败'
      throw err
    } finally {
      listLoading.value = false
    }
  }

  async function submitOrder(data: CreateOrderReq) {
    loading.value = true
    try {
      return await createOrder(data)
    } finally {
      loading.value = false
    }
  }

  return {
    order,
    loading,
    error,
    orders,
    total,
    listLoading,
    listError,
    fetchOrder,
    fetchOrders,
    submitOrder,
  }
})
