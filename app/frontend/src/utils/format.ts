export function formatPrice(cents: number): string {
  return (Number(cents || 0) / 100).toFixed(2)
}

export function formatTime(seconds?: number): string {
  if (!seconds) return '-'
  return new Date(seconds * 1000).toLocaleString()
}

export function orderStatusText(status: number): string {
  switch (status) {
    case 10: return '待支付'
    case 20: return '待发货'
    case 30: return '待收货'
    case 40: return '已完成'
    case 50: return '已取消'
    default: return '未知'
  }
}

export function orderStatusType(status: number): 'warning' | 'primary' | 'success' | 'info' | 'danger' {
  switch (status) {
    case 10: return 'warning'
    case 20: return 'primary'
    case 30: return 'success'
    case 40: return 'info'
    case 50: return 'danger'
    default: return 'info'
  }
}
