<template>
  <div class="tool-grid">
    <div class="grid-container">
      <el-row :gutter="20">
        <el-col 
          v-for="tool in tools" 
          :key="tool.id"
          :span="4.8"
          class="tool-col"
        >
          <div class="tool-card">
            <div class="content-wrapper">
              <div class="tool-header">
                <img 
                  :src="getFavicon(tool.url)"
                  :alt="tool.name"
                  class="tool-icon"
                  @error="handleIconError(tool)"
                />
                <h3 class="tool-name">{{ tool.name }}</h3>
              </div>
              <p class="tool-description">{{ tool.description || '暂无描述' }}</p>
              <div class="tool-tags">
                <div class="tag-item tag-gray">{{ tool.project }}</div>
                <div class="tag-item tag-green">{{ tool.environment }}</div>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup>
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

// 获取网站favicon
const getFavicon = (url) => {
  try {
    const urlObj = new URL(url)
    // 优先尝试获取高清图标
    return `${urlObj.protocol}//${urlObj.hostname}/favicon.ico`
  } catch (error) {
    return '/default-icon.svg'
  }
}

// 处理图标加载失败
const handleIconError = (tool) => {
  const img = event.target
  if (!img.dataset.tried) {
    img.dataset.tried = 'true'
    try {
      const urlObj = new URL(tool.url)
      // 尝试 Google 的服务
      img.src = `https://www.google.com/s2/favicons?domain=${urlObj.hostname}&sz=64`
      img.onerror = () => {
        // 如果 Google 服务失败，使用默认图标
        img.src = '/default-icon.svg'
      }
    } catch (error) {
      img.src = '/default-icon.svg'
    }
  }
}

const openTool = (url) => {
  window.open(url, '_blank')
}
</script>

<style scoped>
.tool-grid {
  padding: 20px 40px;
  width: 100%;
  box-sizing: border-box;
}

.grid-container {
  max-width: 1720px;
  margin: 0 auto;
  width: 100%;
}

.tool-col {
  margin-bottom: 20px;
  display: flex;
  justify-content: center;
}

.tool-card {
  width: 308px;
  height: 120px;
  cursor: pointer;
  transition: all 0.2s ease-in-out;
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.04);
}

.tool-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.content-wrapper {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 16px;
}

.tool-header {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.tool-icon {
  width: 24px;
  height: 24px;
  margin-right: 8px;
  object-fit: contain;
  border-radius: 4px;
}

.tool-name {
  margin: 0;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tool-description {
  margin: 0 0 6px 0;
  font-size: 12px;
  color: #666;
  line-height: 1.5;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  max-height: 36px;
}

.tool-tags {
  display: flex;
  gap: 6px;
}

.tag-item {
  display: inline-block;
  padding: 0 8px;
  height: 22px;
  line-height: 22px;
  font-size: 12px;
  border-radius: 4px;
  white-space: nowrap;
  box-sizing: border-box;
}

.tag-gray {
  background-color: #f4f4f5;
  color: #909399;
}

.tag-green {
  background-color: #f0f9eb;
  color: #67c23a;
}

/* 响应式布局 */
@media screen and (max-width: 1920px) {
  :deep(.el-col) {
    width: 20% !important;
  }
}

@media screen and (max-width: 1600px) {
  .tool-grid {
    padding: 20px 32px;
  }
  .grid-container {
    max-width: 1520px;
  }
}

@media screen and (max-width: 1400px) {
  :deep(.el-col) {
    width: 25% !important;
  }
  .tool-grid {
    padding: 20px 24px;
  }
  .grid-container {
    max-width: 1280px;
  }
}

@media screen and (max-width: 1200px) {
  :deep(.el-col) {
    width: 33.333% !important;
  }
  .tool-grid {
    padding: 20px;
  }
  .grid-container {
    max-width: 1080px;
  }
}

@media screen and (max-width: 768px) {
  :deep(.el-col) {
    width: 50% !important;
  }
  .tool-grid {
    padding: 16px;
  }
}

:deep(.el-card__body) {
  overflow: visible;
  height: 100%;
}
</style>
