import { userRequest } from '@/utils/request'
import type { ApiResponse, LoginReq, LoginResp, UserInfo } from '@/types/api'

/**
 * 查询用户信息
 * GET /user/:name
 */
export async function getUser(name: string): Promise<UserInfo> {
  const res = await userRequest.get<ApiResponse<UserInfo>>(`/user/${encodeURIComponent(name)}`)
  return res.data.data
}

/**
 * 用户登录
 * POST /user/login
 */
export async function login(data: LoginReq): Promise<LoginResp> {
  const res = await userRequest.post<LoginResp>('/user/login', data)
  return res.data
}
