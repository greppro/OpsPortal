<template>
  <div class="tool-grid">
    <el-row :gutter="20">
      <el-col :xs="12" :sm="8" :md="6" :lg="4" v-for="tool in tools" :key="tool.key">
        <el-card class="tool-card" shadow="hover">
          <div class="content-wrapper">
            <div class="card-header">
              <div class="tool-icon">
                <img
                  :src="getFavicon(firstUrl(tool))"
                  :alt="tool.name"
                  class="icon-image"
                  @error="handleIconError"
                />
              </div>
              <div class="favorite-btn" @click.stop="toggleGroupFavorite(tool)">
                <el-icon :class="{ 'is-active': isGroupFavorite(tool) }" :size="20">
                  <StarFilled v-if="isGroupFavorite(tool)" />
                  <Star v-else />
                </el-icon>
              </div>
            </div>
            <h3>{{ tool.name }}</h3>
            <p>{{ tool.description }}</p>
            <div
              class="env-tags"
              :class="'env-tags--count-' + Math.min((tool.envs || []).length, 4)"
            >
              <template v-for="env in tool.envs" :key="env.id">
                <a
                  v-if="getEnvHref(env) !== '#'"
                  :href="getEnvHref(env)"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="env-tag"
                  :class="{ 'env-tag--down': statusOf(env.id)?.state === 'down' }"
                  :title="envTitle(env)"
                  :aria-label="envAriaLabel(env)"
                  @click="recordRecent(tool)"
                >
                  <span
                    v-if="statusOf(env.id)"
                    class="probe-dot"
                    :class="'probe-dot--' + statusOf(env.id).state"
                    aria-hidden="true"
                  ></span>
                  <span class="env-tag__text">{{ env.environment ?? env.label }}</span>
                </a>
                <span v-else class="env-tag env-tag--disabled">{{ env.environment ?? env.label }}</span>
              </template>
            </div>
          </div>
          <div v-if="isManagement" class="tool-actions">
            <el-button type="warning" size="small" @click.stop="$emit('edit', tool)">编辑</el-button>
            <el-button type="danger" size="small" @click.stop="$emit('delete', tool)">删除</el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { Star, StarFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useFavorites } from '../composables/useFavorites'
import { useProbeStatus } from '../composables/useProbeStatus'
import { probeStateText, probeSummary } from '../utils/probe'

const { favoriteIds, isFavorite, saveFavorites } = useFavorites()
const { statusOf } = useProbeStatus()

const props = defineProps({
  tools: {
    type: Array,
    required: true
  },
  isManagement: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['edit', 'delete', 'record-recent'])

function firstUrl(tool) {
  return tool.envs?.[0]?.url ?? ''
}

function isGroupFavorite(tool) {
  if (!tool.envs?.length) return false
  return tool.envs.some(env => isFavorite(env.id))
}

function toggleGroupFavorite(tool) {
  if (!tool.envs?.length) return
  const ids = tool.envs.map(e => e.id)
  const anyFavorite = ids.some(id => favoriteIds.value.has(id))
  if (anyFavorite) {
    ids.forEach(id => favoriteIds.value.delete(id))
    ElMessage.success('已取消收藏')
  } else {
    ids.forEach(id => favoriteIds.value.add(id))
    ElMessage.success('已添加收藏')
  }
  favoriteIds.value = new Set(favoriteIds.value)
  saveFavorites()
}

function getFavicon(url) {
  if (!url) return '/default-icon.svg'
  try {
    const urlObj = new URL(url)
    return `${urlObj.protocol}//${urlObj.hostname}/favicon.ico`
  } catch (error) {
    return '/default-icon.svg'
  }
}

// 直连站点 favicon 失败时用本地默认图标，不把内网域名发给第三方图标服务
function handleIconError(event) {
  const img = event?.target
  if (img && !img.dataset.fallback) {
    img.dataset.fallback = 'true'
    img.src = '/default-icon.svg'
  }
}

function envText(env) {
  return env.environment ?? env.label
}

function envTitle(env) {
  const status = statusOf(env.id)
  if (!status) return undefined
  return `${env.label || envText(env)}：${probeSummary(status)}`
}

function envAriaLabel(env) {
  const status = statusOf(env.id)
  if (!status) return undefined
  return `${envText(env)}（${probeStateText(status)}）`
}

function getEnvHref(env) {
  const u = env?.url
  if (!u || typeof u !== 'string') return '#'
  let t = u.trim()
  if (!t) return '#'
  if (/^https?\/+/i.test(t)) t = t.replace(/^https?\/+:?/i, 'https://')
  return /^https?:\/\//i.test(t) ? t : `https://${t}`
}

function recordRecent(tool) {
  emit('record-recent', tool)
}
</script>

<style scoped>
.tool-grid {
  padding: 16px;
}

.tool-card {
  margin-bottom: 24px;
  text-align: center;
  height: 248px;
  min-height: 248px;
  display: flex;
  flex-direction: column;
  transition: all 0.3s ease;
  cursor: default;
  position: relative;
  overflow: hidden;
  background-color: var(--bg-primary, #fff);
  border: 1px solid var(--border-color, #e5e7eb);
  border-radius: 16px;
  box-shadow: var(--shadow-light);
}

.tool-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-medium);
  border-color: var(--primary-light, #eef2ff);
}

.tool-card::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle at 50% 0%, var(--primary-light, #eef2ff) 0%, transparent 60%);
  opacity: 0;
  transition: opacity 0.3s;
  pointer-events: none;
}

.tool-card:hover::after {
  opacity: 1;
}

.content-wrapper {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 8px;
  padding: 12px;
}

.card-header {
  position: relative;
  width: 100%;
  display: flex;
  justify-content: center;
}

.favorite-btn {
  position: absolute;
  top: -8px;
  right: -8px;
  padding: 8px;
  cursor: pointer;
  color: var(--text-secondary, #9ca3af);
  transition: all 0.2s;
  opacity: 0;
  z-index: 2;
}

.tool-card:hover .favorite-btn,
.favorite-btn:has(.is-active) {
  opacity: 1;
}

.favorite-btn .is-active {
  color: #f59e0b;
}

.favorite-btn:hover {
  transform: scale(1.1);
  color: #f59e0b;
}

.tool-icon {
  margin: 4px 0 8px 0;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 40px;
  min-height: 40px;
}

.icon-image {
  width: 32px;
  height: 32px;
  object-fit: contain;
  transition: transform 0.3s;
  border-radius: 4px;
  filter: grayscale(0.2);
}

.tool-card:hover .icon-image {
  transform: scale(1.1);
  filter: grayscale(0);
}

h3 {
  margin: 0 0 6px 0;
  font-size: 16px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--text-primary, #1e293b);
  transition: color 0.3s;
  width: 100%;
  text-align: center;
}

.tool-card:hover h3 {
  color: var(--primary, #0ea5e9);
}

p {
  margin: 0;
  font-size: 14px;
  color: var(--text-secondary, #666);
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
  min-height: 38px;
  line-height: 1.5;
  width: 100%;
  text-align: center;
}

.env-tags {
  gap: 6px;
  margin-top: 8px;
  min-height: 72px;
  width: 100%;
  box-sizing: border-box;
}

/* 1 个环境：只占第一行并拉长（跨两格），第二行空着 */
.env-tags--count-1 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: auto auto;
  gap: 8px;
  max-width: 256px;
  margin-left: auto;
  margin-right: auto;
  min-height: 72px;
  align-content: start;
}
.env-tags--count-1 .env-tag {
  grid-column: 1 / -1;
  grid-row: 1;
  width: 100%;
  min-width: 0;
  text-align: center;
}

/* 2 个环境：图中示意 - 两行等大，纵向居中排列，整体水平居中 */
.env-tags--count-2 {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  margin-left: auto;
  margin-right: auto;
}
/* 每行只有一个环境时拉长：两行各一个，适当加宽 */
.env-tags--count-2 .env-tag {
  width: 100%;
  max-width: 200px;
  min-width: 100px;
  text-align: center;
}

/* 3 个环境：图中示意 - 左列两个上下堆叠，右列一个与左上对齐 */
.env-tags--count-3 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr 1fr;
  gap: 8px;
  max-width: 256px;
  margin-left: auto;
  margin-right: auto;
  min-height: 72px;
}
.env-tags--count-3 .env-tag:nth-child(1) {
  grid-column: 1;
  grid-row: 1;
}
.env-tags--count-3 .env-tag:nth-child(2) {
  grid-column: 1;
  grid-row: 2;
}
.env-tags--count-3 .env-tag:nth-child(3) {
  grid-column: 2;
  grid-row: 1;
}
.env-tags--count-3 .env-tag {
  width: 100%;
  min-width: 0;
  text-align: center;
}

/* 4 个环境：图中示意 - 2×2 网格，等大等间距，整体居中 */
.env-tags--count-4 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr 1fr;
  gap: 8px;
  max-width: 256px;
  margin-left: auto;
  margin-right: auto;
  min-height: 72px;
}
.env-tags--count-4 .env-tag {
  width: 100%;
  min-width: 0;
  text-align: center;
}

.env-tag {
  cursor: pointer;
  text-decoration: none;
  color: var(--primary, #615ced);
  min-height: 28px;
  height: 28px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0 10px;
  border-radius: 8px;
  font-size: 12px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  border: 1px solid transparent;
  box-sizing: border-box;
}

/* 环境标签交替底色，与主题协调 */
.env-tag:nth-child(4n + 1) {
  background: rgba(97, 92, 237, 0.12);
  border-color: rgba(97, 92, 237, 0.2);
}
.env-tag:nth-child(4n + 2) {
  background: rgba(14, 165, 233, 0.12);
  border-color: rgba(14, 165, 233, 0.2);
}
.env-tag:nth-child(4n + 3) {
  background: rgba(34, 197, 94, 0.12);
  border-color: rgba(34, 197, 94, 0.2);
}
.env-tag:nth-child(4n + 4) {
  background: rgba(245, 158, 11, 0.12);
  border-color: rgba(245, 158, 11, 0.2);
}

.env-tag:hover {
  opacity: 0.9;
}

.env-tag {
  gap: 6px;
}

.env-tag__text {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 可用性检测状态点 */
.probe-dot {
  position: relative;
  flex-shrink: 0;
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--text-secondary, #9ca3af);
}

.probe-dot--up {
  background: #22c55e;
}

.probe-dot--down {
  background: #ef4444;
}

.probe-dot--down::after {
  content: '';
  position: absolute;
  inset: -3px;
  border-radius: 50%;
  border: 1px solid #ef4444;
  opacity: 0;
  animation: probe-pulse 1.8s ease-out infinite;
}

.probe-dot--pending {
  background: transparent;
  box-shadow: inset 0 0 0 1.5px var(--text-secondary, #9ca3af);
}

.env-tag.env-tag--down {
  border-color: rgba(239, 68, 68, 0.45);
}

@keyframes probe-pulse {
  0% {
    transform: scale(0.6);
    opacity: 0.8;
  }
  100% {
    transform: scale(1.6);
    opacity: 0;
  }
}

@media (prefers-reduced-motion: reduce) {
  .probe-dot--down::after {
    animation: none;
  }
}

.env-tag--disabled {
  cursor: default;
  opacity: 0.7;
  background: var(--bg-secondary, #f3f4f6) !important;
  border-color: var(--border-color, #e5e7eb) !important;
}

.tool-actions {
  margin-top: 8px;
  padding: 8px 0 4px;
  border-top: 1px solid var(--border-color, #f0f0f0);
  display: flex;
  gap: 8px;
  justify-content: center;
  position: relative;
  z-index: 1;
}

.el-button--small {
  padding: 5px 12px;
  font-size: 12px;
  height: 26px;
  min-width: 54px;
}
</style>
