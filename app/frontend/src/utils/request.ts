import axios from 'axios'
import type { AxiosInstance, AxiosResponse } from 'axios'

/** 创建 Axios 实例，可通过 baseURL 区分不同服务 */
export function createRequest(baseURL: string): AxiosInstance {
  const instance = axios.create({
    baseURL,
    timeout: 10000,
    headers: {
      'Content-Type': 'application/json',
    },
  })

  // 请求拦截器 — 添加认证 token
  instance.interceptors.request.use((config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  })

  // 响应拦截器 — 统一错误处理
  instance.interceptors.response.use(
    (response: AxiosResponse) => {
      return response
    },
    (error) => {
      if (error.response?.status === 401) {
        localStorage.removeItem('token')
        window.location.href = '/login'
      }
      return Promise.reject(error)
    },
  )

  return instance
}

/** 用户服务请求实例 */
export const userRequest = createRequest('/api/user')
/** 订单服务请求实例 */
export const orderRequest = createRequest('/api/order')
/** 商品服务请求实例 */
export const productRequest = createRequest('/api/product')
/** 媒体服务请求实例 */
export const mediaRequest = createRequest('/api/media')
