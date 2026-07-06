import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getOrder, createOrder } from '@/api/order'
import type { CreateOrderReq, OrderInfo } from '@/types/api'

export const useOrderStore = defineStore('order', () => {
  const order = ref<OrderInfo | null>(null)
  const loading = ref(false)

  async function fetchOrder(id: number) {
    loading.value = true
    try {
      order.value = await getOrder(id)
    } finally {
      loading.value = false
    }
  }

  async function submitOrder(data: CreateOrderReq) {
    loading.value = true
    try {
      const res = await createOrder(data)
      return res
    } finally {
      loading.value = false
    }
  }

  return { order, loading, fetchOrder, submitOrder }
})
