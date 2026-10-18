/** 商城与后台使用不同端口，浏览器地址栏的 port 决定当前模式。 */
export const ADMIN_PORT = '8080'
export const MALL_PORT = '8081'

export const isDev = import.meta.env.DEV
export const isMallPort = window.location.port === MALL_PORT
export const isAdminPort = window.location.port === ADMIN_PORT

export function isMallMode(): boolean {
  return !isDev && isMallPort
}

export function isAdminMode(): boolean {
  return !isDev && isAdminPort
}

function portalUrl(port: string, path: string): string {
  return `${window.location.protocol}//${window.location.hostname}:${port}${path}`
}

export function adminUrl(path = '/dashboard'): string {
  return portalUrl(ADMIN_PORT, path)
}

export function mallUrl(path = '/mall'): string {
  return portalUrl(MALL_PORT, path)
}

export function openInNewTab(url: string): void {
  window.open(url, '_blank', 'noopener,noreferrer')
}
