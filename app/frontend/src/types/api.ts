/** 通用 API 响应结构（对应 go-zero 返回格式） */
export interface ApiResponse<T = unknown> {
  code: number
  msg: string
  data: T
}

/** 用户信息 */
export interface UserInfo {
  name: string
  message?: string
}

/** 用户登录请求 */
export interface LoginReq {
  username: string
  password: string
}

/** 用户登录响应 */
export interface LoginResp {
  code: number
  msg: string
  token?: string
}

/** 订单信息 */
export interface OrderInfo {
  id: number
  goodsId?: number
  goodsName?: string
  price?: number
  data?: string
}

/** 创建订单请求 */
export interface CreateOrderReq {
  goods_id: number
  goods_name: string
  price: number
  created_at: number
}
