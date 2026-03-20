// 核心目的：格式化工具函数
// 模块功能：日期格式化、时间范围格式化、数字格式化
// 核心目的：格式化工具函数
// 模块功能：日期格式化、时间范围格式化、数字格式化

/** 格式化日期为 YYYY-MM-DD */
export function formatDate(date: string | Date | null | undefined): string {
  if (!date) return '-'
  const d = typeof date === 'string' ? new Date(date) : date
  if (isNaN(d.getTime())) return '-'
  return d.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' }).replace(/\//g, '-')
}

/** 格式化日期时间为 YYYY-MM-DD HH:mm */
export function formatDateTime(date: string | Date | null | undefined): string {
  if (!date) return '-'
  const d = typeof date === 'string' ? new Date(date) : date
  if (isNaN(d.getTime())) return '-'
  return `${formatDate(d)} ${d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit', hour12: false })}`
}

/** 格式化时间范围：09:00 ~ 09:30 */
export function formatTimeRange(start: string, end: string): string {
  return `${start} ~ ${end}`
}

/** 格式化数字，添加千分位 */
export function formatNumber(n: number | null | undefined, precision = 0): string {
  if (n === null || n === undefined) return '-'
  return n.toLocaleString('zh-CN', { minimumFractionDigits: precision, maximumFractionDigits: precision })
}

/** 分钟转人类可读时长：90 → 1小时30分 */
export function formatDuration(minutes: number): string {
  if (minutes < 60) return `${minutes}分钟`
  const h = Math.floor(minutes / 60)
  const m = minutes % 60
  return m === 0 ? `${h}小时` : `${h}小时${m}分钟`
}
