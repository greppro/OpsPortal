<template>
  <div class="tool-grid">
    <el-row :gutter="20">
      <el-col :xs="12" :sm="8" :md="6" :lg="4" v-for="tool in tools" :key="tool.id">
        <el-card class="tool-card" shadow="hover">
          <div class="content-wrapper">
            <div class="tool-icon">
              <img 
                :src="getToolIcon(tool.name)" 
                :alt="tool.name"
                @error="handleIconError"
                class="icon-image"
              />
            </div>
            <h3>{{ tool.name }}</h3>
            <p>{{ tool.description }}</p>
          </div>
          <div class="tool-actions">
            <el-button type="primary" size="small" @click="openTool(tool.url)">访问</el-button>
            <el-button v-if="isManagement" type="warning" size="small" @click="$emit('edit', tool)">编辑</el-button>
            <el-button v-if="isManagement" type="danger" size="small" @click="$emit('delete', tool.id)">删除</el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref } from 'vue'

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

const defaultIcon = ref('data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCI+PHBhdGggZmlsbD0iI2NjYyIgZD0iTTEyIDJDNi40OCAyIDIgNi40OCAyIDEyczQuNDggMTAgMTAgMTAgMTAtNC40OCAxMC0xMFMxNy41MiAyIDEyIDJ6bTAgMThjLTQuNDEgMC04LTMuNTktOC04czMuNTktOCA4LTggOCAzLjU5IDggOC0zLjU5IDgtOCA4eiIvPjwvc3ZnPg==')

// 获取工具图标
const getToolIcon = (name) => {
  const toolName = name.toLowerCase().trim()
  // 尝试从 simpleicons.org 获取图标
  return `https://simpleicons.org/icons/${toolName}.svg`
}

// 处理图标加载错误
const handleIconError = (e) => {
  e.target.src = defaultIcon.value
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
  transition: all 0.3s;
}

.tool-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 12px 0 rgba(0,0,0,.1);
}

.content-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  margin-bottom: 8px;
}

.tool-icon {
  margin: 4px 0 8px 0;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 40px;
}

.icon-image {
  width: 32px;
  height: 32px;
  object-fit: contain;
}

h3 {
  margin: 0 0 6px 0;
  font-size: 16px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
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
}

.tool-actions {
  margin-top: 8px;
  padding: 8px 0 4px;
  border-top: 1px solid #f0f0f0;
  display: flex;
  gap: 8px;
  justify-content: center;
}

.el-button--small {
  padding: 5px 12px;
  font-size: 12px;
  height: 26px;
  min-width: 54px;
}
</style>
