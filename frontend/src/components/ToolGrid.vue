<template>
  <div class="tool-grid">
    <el-row :gutter="20">
      <el-col :xs="12" :sm="8" :md="6" :lg="4" v-for="tool in tools" :key="tool.id">
        <el-card 
          class="tool-card" 
          shadow="hover" 
          @click="openTool(tool.url)"
        >
          <div class="content-wrapper">
            <div class="tool-icon">
              <img 
                :src="getFavicon(tool.url)"
                :alt="tool.name"
                class="icon-image"
                @error="handleIconError(tool)"
              />
            </div>
            <h3>{{ tool.name }}</h3>
            <p>{{ tool.description }}</p>
          </div>
          <div v-if="isManagement" class="tool-actions">
            <el-button type="warning" size="small" @click.stop="$emit('edit', tool)">编辑</el-button>
            <el-button type="danger" size="small" @click.stop="$emit('delete', tool.id)">删除</el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
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
  padding: 16px;
}

.tool-card {
  margin-bottom: 20px;
  text-align: center;
  height: 180px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
  position: relative;
  overflow: hidden;
}

.tool-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 6px 16px rgba(0,0,0,.1);
}

.tool-card::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle, rgba(255,255,255,0.8) 0%, rgba(255,255,255,0) 100%);
  opacity: 0;
  transition: opacity 0.3s;
}

.tool-card:active::after {
  opacity: 1;
}

.content-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  margin-bottom: 8px;
  padding: 12px;
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
  color: #303133;
  transition: color 0.3s;
}

.tool-card:hover h3 {
  color: #409EFF;
}

p {
  margin: 0;
  font-size: 14px;
  color: #666;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
  min-height: 38px;
  line-height: 1.5;
}

.tool-actions {
  margin-top: 8px;
  padding: 8px 0 4px;
  border-top: 1px solid #f0f0f0;
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
