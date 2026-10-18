import { userRequest } from '@/utils/request'
import type {
  AddressInfo,
  AddressListData,
  ApiResponse,
  AuditLogListData,
  AuditLogListParams,
  CreateUserReq,
  ListUsersParams,
  LoginReq,
  LoginResp,
  RegisterReq,
  SaveAddressReq,
  UserInfo,
  UserListData,
} from '@/types/api'

export async function getUser(name: string): Promise<UserInfo> {
  const res = await userRequest.get<ApiResponse<UserInfo>>(`/user/${encodeURIComponent(name)}`)
  return res.data.data
}

export async function login(data: LoginReq): Promise<LoginResp> {
  const res = await userRequest.post<LoginResp>('/user/login', data)
  return res.data
}

export async function register(data: RegisterReq): Promise<ApiResponse<UserInfo>> {
  const res = await userRequest.post<ApiResponse<UserInfo>>('/user/register', data)
  return res.data
}

export async function listUsers(params: ListUsersParams): Promise<UserListData> {
  const res = await userRequest.get<ApiResponse<UserListData>>('/users', { params })
  return res.data.data
}

export async function createUser(data: CreateUserReq): Promise<ApiResponse<UserInfo>> {
  const res = await userRequest.post<ApiResponse<UserInfo>>('/users', data)
  return res.data
}

export async function updateUserStatus(id: number, status: number): Promise<ApiResponse> {
  const res = await userRequest.put<ApiResponse>(`/admin/users/${id}/status`, { status })
  return res.data
}

export async function updateUserRole(id: number, role: number): Promise<ApiResponse> {
  const res = await userRequest.put<ApiResponse>(`/admin/users/${id}/role`, { role })
  return res.data
}

export async function listAuditLogs(params: AuditLogListParams): Promise<AuditLogListData> {
  const res = await userRequest.get<ApiResponse<AuditLogListData>>('/admin/audit-logs', { params })
  return res.data.data
}

export async function changePassword(data: {
  old_password: string
  new_password: string
}): Promise<ApiResponse> {
  const res = await userRequest.put<ApiResponse>('/user/password', data)
  return res.data
}

export async function logout(): Promise<ApiResponse> {
  const res = await userRequest.post<ApiResponse>('/user/logout')
  return res.data
}

export async function listAddresses(): Promise<AddressListData> {
  const res = await userRequest.get<ApiResponse<AddressListData>>('/addresses')
  return res.data.data
}

export async function createAddress(data: SaveAddressReq): Promise<ApiResponse<AddressInfo>> {
  const res = await userRequest.post<ApiResponse<AddressInfo>>('/addresses', data)
  return res.data
}

export async function updateAddress(id: number, data: SaveAddressReq): Promise<ApiResponse<AddressInfo>> {
  const res = await userRequest.put<ApiResponse<AddressInfo>>(`/addresses/${id}`, data)
  return res.data
}

export async function deleteAddress(id: number): Promise<ApiResponse> {
  const res = await userRequest.delete<ApiResponse>(`/addresses/${id}`)
  return res.data
}

export async function setDefaultAddress(id: number): Promise<ApiResponse> {
  const res = await userRequest.put<ApiResponse>(`/addresses/${id}/default`)
  return res.data
}
