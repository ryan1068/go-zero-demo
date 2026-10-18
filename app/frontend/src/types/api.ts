/** 通用 API 响应结构（对应 go-zero 返回格式） */
export interface ApiResponse<T = unknown> {
  code: number
  msg: string
  data: T
}

/** 用户信息 */
export interface UserInfo {
  id: number
  username: string
  nickname: string
  email: string
  phone: string
  avatar: string
  status: number
  remark: string
  created_at: number
  updated_at: number
  role: number
}

export interface CreateUserReq {
  username: string
  password: string
  nickname: string
  email: string
  phone: string
  avatar: string
  status: number
  remark: string
  role: number
}

export interface RegisterReq {
  username: string
  password: string
  nickname: string
  email: string
  phone: string
}

export interface ListUsersParams {
  page?: number
  page_size?: number
  username?: string
  nickname?: string
  email?: string
  phone?: string
  status?: number
}

export interface UserListData {
  list: UserInfo[]
  total: number
  page: number
  page_size: number
}

export interface LoginReq {
  username: string
  password: string
}

export interface LoginData {
  token: string
  user_id: number
  username: string
  role: number
}

export interface LoginResp {
  code: number
  msg: string
  data: LoginData
}

export interface AddressInfo {
  id: number
  user_id: number
  receiver_name: string
  receiver_phone: string
  province: string
  city: string
  district: string
  detail_address: string
  postal_code: string
  tag: string
  is_default: number
  created_at: number
  updated_at: number
}

export interface SaveAddressReq {
  id?: number
  receiver_name: string
  receiver_phone: string
  province: string
  city: string
  district: string
  detail_address: string
  postal_code: string
  tag: string
  is_default: number
}

export interface AddressListData {
  list: AddressInfo[]
}

export interface CategoryInfo {
  id: number
  parent_id: number
  name: string
  icon: string
  sort: number
  status: number
  created_at: number
  updated_at: number
}

export interface SaveCategoryReq {
  id?: number
  parent_id: number
  name: string
  icon: string
  sort: number
  status: number
}

export interface ProductInfo {
  id: number
  category_id: number
  category_name: string
  name: string
  subtitle: string
  cover_image: string
  images: string[]
  detail: string
  price: number
  original_price: number
  stock: number
  sales: number
  unit: string
  status: number
  is_recommend: number
  sort: number
  created_at: number
  updated_at: number
}

export interface ListProductsParams {
  page?: number
  page_size?: number
  category_id?: number
  keyword?: string
  min_price?: number
  max_price?: number
  status?: number
  sort?: string
  low_stock?: number
}

export interface ProductListData {
  list: ProductInfo[]
  total: number
  page: number
  page_size: number
}

export interface SaveProductReq {
  id?: number
  category_id: number
  name: string
  subtitle: string
  cover_image: string
  images: string[]
  detail: string
  price: number
  original_price: number
  stock: number
  unit: string
  status: number
  is_recommend: number
  sort: number
}

export interface ProductStats {
  total: number
  on_sale: number
  off_sale: number
  low_stock: number
}

export interface CartItemInfo {
  id: number
  product_id: number
  product_name: string
  product_cover: string
  price: number
  stock: number
  quantity: number
  selected: number
  product_status: number
  created_at: number
  updated_at: number
}

export interface CartListData {
  list: CartItemInfo[]
  total: number
}

export interface OrderItemInfo {
  id: number
  product_id: number
  product_name: string
  product_cover: string
  price: number
  quantity: number
  total_amount: number
}

export interface OrderInfo {
  id: number
  order_no: string
  user_id: number
  total_amount: number
  freight_amount: number
  discount_amount: number
  pay_amount: number
  status: number
  pay_type: number
  transaction_id: string
  receiver_name: string
  receiver_phone: string
  receiver_province: string
  receiver_city: string
  receiver_district: string
  receiver_address: string
  receiver_postal_code: string
  remark: string
  admin_remark: string
  logistics_company: string
  tracking_no: string
  paid_at: number
  shipped_at: number
  completed_at: number
  canceled_at: number
  cancel_reason: string
  created_at: number
  updated_at: number
  items: OrderItemInfo[]
}

export interface CreateOrderReq {
  address_id: number
  remark: string
  buy_now: number
  product_id: number
  quantity: number
}

export interface ListOrdersParams {
  page?: number
  page_size?: number
  order_no?: string
  status?: number
  start_time?: number
  end_time?: number
}

export interface OrderListData {
  list: OrderInfo[]
  total: number
  page: number
  page_size: number
}

export interface DashboardData {
  today_order_count: number
  today_paid_amount: number
  pending_shipment_count: number
  product_count: number
  low_stock_count: number
}

export interface UploadResp {
  code: number
  msg: string
  url: string
}

export interface AfterSaleInfo {
  id: number
  after_sale_no: string
  order_id: number
  order_item_id: number
  user_id: number
  type: number
  reason: string
  description: string
  images: string[]
  amount: number
  status: number
  admin_remark: string
  tracking_no: string
  apply_at: number
  review_at: number
  buyer_shipped_at: number
  seller_received_at: number
  refunded_at: number
  created_at: number
  updated_at: number
  order_no: string
  product_name: string
  product_cover: string
}

export interface AfterSaleListData {
  list: AfterSaleInfo[]
  total: number
  page: number
  page_size: number
}

export interface HandleAfterSaleParams {
  action: string
  admin_remark?: string
  tracking_no?: string
}

export interface SkuInfo {
  id: number
  product_id: number
  sku_code: string
  specs: string
  price: number
  original_price: number
  stock: number
  sales: number
  status: number
  created_at: number
  updated_at: number
}

export interface SaveSkuReq {
  id?: number
  product_id: number
  sku_code: string
  specs: string
  price: number
  original_price: number
  stock: number
  status: number
}

export interface AuditLogInfo {
  id: number
  admin_id: number
  admin_name: string
  action: string
  target_type: string
  target_id: number
  detail: string
  created_at: number
}

export interface AuditLogListData {
  list: AuditLogInfo[]
  total: number
  page: number
  page_size: number
}

export interface AuditLogListParams {
  page: number
  page_size: number
  action?: string
  target_type?: string
}

export interface OrderTrack {
  time: number
  status: string
  description: string
}

export interface OrderTracksData {
  order_id: number
  order_no: string
  logistics_company: string
  tracking_no: string
  tracks: OrderTrack[]
}
