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
                  @error="handleIconError(tool)"
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
            <div class="env-tags">
              <el-tag
                v-for="env in tool.envs"
                :key="env.id"
                size="small"
                class="env-tag"
                @click.stop="openTool(env.url)"
              >
                {{ env.label }}
              </el-tag>
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

const { favoriteIds, isFavorite, saveFavorites } = useFavorites()

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

const emit = defineEmits(['edit', 'delete'])

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

function handleIconError(tool) {
  const url = firstUrl(tool)
  const img = event?.target
  if (!img?.dataset?.tried) {
    img.dataset.tried = 'true'
    try {
      const urlObj = new URL(url)
      img.src = `https://www.google.com/s2/favicons?domain=${urlObj.hostname}&sz=64`
      img.onerror = () => { img.src = '/default-icon.svg' }
    } catch (e) {
      img.src = '/default-icon.svg'
    }
  }
}

function openTool(url) {
  if (url) window.open(url, '_blank')
}
</script>

<style scoped>
.tool-grid {
  padding: 16px;
}

.tool-card {
  margin-bottom: 24px;
  text-align: center;
  min-height: 200px;
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
}

.tool-card:hover::after {
  opacity: 1;
}

.content-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  margin-bottom: 8px;
  padding: 12px;
}

.card-header {
  position: relative;
  width: 100%;
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
}

.env-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  justify-content: center;
  margin-top: 8px;
}

.env-tag {
  cursor: pointer;
}

.env-tag:hover {
  opacity: 0.9;
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
