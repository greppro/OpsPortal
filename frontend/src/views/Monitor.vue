<template>
  <div class="tools-container">
    <div class="tools-header">
      <div class="filter-header">
        <el-select
          v-model="activeProject"
          placeholder="选择项目"
          clearable
          @change="handleProjectChange"
        >
          <el-option
            v-for="project in projects"
            :key="project.id"
            :label="project.label"
            :value="project.name"
          />
        </el-select>
      </div>

      <div class="tabs-header" v-if="activeProject">
        <el-tabs v-model="activeEnv" @tab-change="handleEnvChange">
          <el-tab-pane
            v-for="env in environments"
            :key="env.id"
            :label="env.label"
            :name="env.name"
          />
        </el-tabs>
      </div>
    </div>

    <div class="tools-content">
      <tool-grid :tools="tools" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import ToolGrid from '../components/ToolGrid.vue'
import request from '../utils/request'

const projects = ref([])
const environments = ref([])
const tools = ref([])
const activeProject = ref('')
const activeEnv = ref('')

// 获取项目列表
const fetchProjects = async () => {
  try {
    const response = await request.get('/api/projects')
    projects.value = response.data
  } catch (error) {
    console.error('Error fetching projects:', error)
    ElMessage.error('获取项目列表失败')
  }
}

// 获取环境列表
const fetchEnvironments = async () => {
  try {
    const response = await request.get('/api/environments', {
      params: { project_id: getProjectId(activeProject.value) }
    })
    environments.value = response.data
  } catch (error) {
    console.error('Error fetching environments:', error)
    ElMessage.error('获取环境列表失败')
  }
}

// 获取工具列表
const fetchTools = async () => {
  try {
    const params = {}
    
    if (activeProject.value) {
      params.project = activeProject.value
    }
    if (activeEnv.value) {
      params.env = activeEnv.value
    }

    console.log('Fetching tools with params:', params) // 添加调试日志
    const response = await request.get('/api/sites', { params })
    tools.value = response.data
  } catch (error) {
    console.error('Error fetching tools:', error)
    ElMessage.error('获取工具列表失败')
  }
}

// 获取项目ID
const getProjectId = (projectName) => {
  const project = projects.value.find(p => p.name === projectName)
  return project?.id || ''
}

// 处理项目变更
const handleProjectChange = async () => {
  activeEnv.value = '' // 重置环境选择
  if (activeProject.value) {
    await fetchEnvironments() // 重新获取环境列表
  } else {
    environments.value = [] // 清空环境列表
  }
  fetchTools()
}

// 处理环境变更
const handleEnvChange = (tabName) => {
  console.log('Environment changed to:', tabName) // 添加调试日志
  activeEnv.value = tabName
  fetchTools()
}

onMounted(async () => {
  await fetchProjects()
  if (activeProject.value) {
    await fetchEnvironments()
  }
  fetchTools()
})
</script>

<style scoped>
.tools-container {
  height: 100%;
  padding: 20px;
  display: flex;
  flex-direction: column;
}

.tools-header {
  flex-shrink: 0;
}

.filter-header {
  margin-bottom: 20px;
}

.tabs-header {
  margin-bottom: 20px;
}

.tools-content {
  flex: 1;
  overflow-y: auto;
}

:deep(.el-tabs__header) {
  margin-bottom: 0;
}

:deep(.el-tabs__nav-wrap) {
  padding: 0;
}

:deep(.el-tabs__item) {
  font-size: 14px;
  height: 40px;
  line-height: 40px;
}

:deep(.el-select) {
  width: 200px;
}
</style> 