// 可用性检测结果的展示文案

const REASON_TEXT = {
  http_status: 'HTTP 状态码异常',
  timeout: '请求超时',
  dns: '域名解析失败',
  connection_refused: '连接被拒绝',
  tls: 'TLS 握手或证书校验失败',
  invalid_url: '地址格式无效',
  network: '网络错误'
}

const STATE_TEXT = {
  up: '可用',
  down: '不可用',
  pending: '等待检测'
}

function formatTime(value) {
  if (!value) return ''
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return ''
  return d.toLocaleTimeString('zh-CN', { hour12: false })
}

export function probeStateText(status) {
  return STATE_TEXT[status?.state] || ''
}

// 一行说明，用于 tooltip：可用 · HTTP 200 · 87 ms · 10:05:12 检测
export function probeSummary(status) {
  if (!status) return ''
  if (status.state === 'pending') return STATE_TEXT.pending
  const parts = [probeStateText(status)]
  if (status.state === 'down' && status.reason && status.reason !== 'http_status') {
    parts.push(REASON_TEXT[status.reason] || status.reason)
  }
  if (status.status_code) parts.push(`HTTP ${status.status_code}`)
  if (status.state === 'up' && status.duration_ms != null) parts.push(`${status.duration_ms} ms`)
  const time = formatTime(status.checked_at)
  if (time) parts.push(`${time} 检测`)
  return parts.join(' · ')
}
