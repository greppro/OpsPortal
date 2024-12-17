<template>
  <div class="tool-grid">
    <el-row :gutter="20">
      <el-col :span="6" v-for="tool in tools" :key="tool.id">
        <el-card class="tool-card" shadow="hover">
          <div class="tool-icon">
            <el-icon :size="24"><component :is="tool.icon || 'Link'" /></el-icon>
          </div>
          <h3>{{ tool.name }}</h3>
          <p>{{ tool.description }}</p>
          <div class="tool-actions">
            <el-button type="primary" @click="openTool(tool.url)">访问</el-button>
            <el-button v-if="isManagement" type="warning" @click="$emit('edit', tool)">编辑</el-button>
            <el-button v-if="isManagement" type="danger" @click="$emit('delete', tool.id)">删除</el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue'

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
</script>

<style scoped>
.tool-grid {
  padding: 20px;
}

.tool-card {
  margin-bottom: 20px;
  text-align: center;
  height: 180px;
}

.tool-icon {
  margin: 8px 0;
}

h3 {
  margin: 8px 0;
  font-size: 16px;
}

p {
  margin: 8px 0;
  font-size: 14px;
  color: #666;
}

.tool-actions {
  margin-top: 8px;
}
</style>
