import { ref, onMounted, onUnmounted } from 'vue'
import request from '../utils/request'

// 多个组件共用同一份状态和同一个轮询定时器
const enabled = ref(false)
const results = ref({})
const intervalSeconds = ref(60)

const MIN_POLL_SECONDS = 15
const MAX_POLL_SECONDS = 300

let subscribers = 0
let timer = null
let inflight = null

function pollDelayMs() {
  const s = Math.min(Math.max(intervalSeconds.value || 60, MIN_POLL_SECONDS), MAX_POLL_SECONDS)
  return s * 1000
}

async function refresh() {
  if (inflight) return inflight
  inflight = request
    .get('/api/probe/status', { silent: true })
    .then(res => {
      enabled.value = !!res.data?.enabled
      results.value = res.data?.results || {}
      if (res.data?.interval_seconds) intervalSeconds.value = res.data.interval_seconds
    })
    .catch(() => {
      // 静默失败：保留上一次的结果
    })
    .finally(() => {
      inflight = null
    })
  return inflight
}

function schedule() {
  clearTimeout(timer)
  timer = setTimeout(async () => {
    if (!document.hidden) await refresh()
    schedule()
  }, pollDelayMs())
}

function onVisibilityChange() {
  if (!document.hidden) refresh()
}

function start() {
  refresh()
  schedule()
  document.addEventListener('visibilitychange', onVisibilityChange)
}

function stop() {
  clearTimeout(timer)
  timer = null
  document.removeEventListener('visibilitychange', onVisibilityChange)
}

// 管理页"立即检测"后直接写入结果，不用等下一次轮询
function setResult(id, status) {
  results.value = { ...results.value, [String(id)]: status }
}

function statusOf(id) {
  if (!enabled.value || id == null) return null
  return results.value[String(id)] || null
}

export function useProbeStatus() {
  onMounted(() => {
    subscribers += 1
    if (subscribers === 1) start()
  })
  onUnmounted(() => {
    subscribers -= 1
    if (subscribers === 0) stop()
  })
  return { enabled, results, statusOf, refresh, setResult }
}
