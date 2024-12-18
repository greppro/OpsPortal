<template>
  <div class="tool-grid">
    <el-row :gutter="16">
      <el-col :xs="12" :sm="8" :md="6" :lg="4" v-for="tool in tools" :key="tool.id">
        <el-card class="tool-card" shadow="hover">
          <div class="content-wrapper">
            <div class="tool-icon">
              <el-icon :size="28" :color="getIconColor(tool.category)">
                <component :is="getToolIcon(tool.name)" />
              </el-icon>
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
import { defineProps, defineEmits } from 'vue'
import { Monitor, DataLine, Connection, List, Link } from '@element-plus/icons-vue'

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

const openTool = (url) => {
  window.open(url, '_blank')
}

const getToolIcon = (name) => {
  const lowerName = name.toLowerCase()
  if (lowerName.includes('grafana')) return Monitor
  if (lowerName.includes('prometheus')) return DataLine
  if (lowerName.includes('jenkins')) return Connection
  if (lowerName.includes('gitlab')) return List
  return Link
}

const getIconColor = (category) => {
  switch (category?.toLowerCase()) {
    case '监控': return '#409EFF'
    case '部署': return '#67C23A'
    case '日志': return '#E6A23C'
    default: return '#909399'
  }
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
  display: flex;
  flex-direction: column;
}

.tool-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 12px 0 rgba(0,0,0,.1);
}

.el-card__body {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 12px !important;
  justify-content: space-between;
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
