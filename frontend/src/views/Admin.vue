<template>
  <div class="admin-container">
    <div class="page-header">
      <el-tabs class="admin-tabs" v-model="activeTab">
        <el-tab-pane label="网址管理" name="sites" />
        <el-tab-pane label="环境管理" name="environments" />
        <el-tab-pane label="项目管理" name="projects" />
      </el-tabs>
    </div>
    <div class="page-content">
      <component :is="currentComponent" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import Manage from './Manage.vue'
import Environment from './EnvironmentManage.vue'
import Project from './Project.vue'

const activeTab = ref('sites')

// 根据当前标签页返回对应的组件
const currentComponent = computed(() => {
  const components = {
    sites: Manage,
    environments: Environment,
    projects: Project
  }
  return components[activeTab.value]
})
</script>

<style scoped>
.admin-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #f0f2f5;
}

.page-header {
  background-color: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
}

.admin-tabs {
  height: 60px;
  padding: 0 20px;
}

.page-content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

:deep(.el-tabs__header) {
  margin: 0;
}

:deep(.el-tabs__nav-wrap) {
  height: 60px;
}

:deep(.el-tabs__nav) {
  height: 60px;
}

:deep(.el-tabs__item) {
  height: 60px;
  line-height: 60px;
  color: #606266;
}

:deep(.el-tabs__item.is-active) {
  color: #409EFF;
}

:deep(.el-tabs__active-bar) {
  bottom: 0;
}
</style> 