<template>
  <div class="tools-container">
    <div class="filter-header">
      <el-select 
        v-model="activeProject" 
        placeholder="选择项目"
        clearable
        @change="handleProjectChange"
        style="width: 200px"
      >
        <el-option 
          v-for="proj in projects" 
          :key="proj.id" 
          :label="proj.label" 
          :value="proj.name"
        />
      </el-select>
    </div>

    <el-tabs v-model="activeEnv" @tab-click="handleEnvChange">
      <el-tab-pane 
        v-for="env in filteredEnvironments" 
        :key="env.id" 
        :label="env.label" 
        :name="env.name"
      >
        <tool-grid :tools="tools" />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '../utils/request'
import ToolGrid from '../components/ToolGrid.vue'

const activeProject = ref('')
const activeEnv = ref('')
const tools = ref([])
const projects = ref([])
const environments = ref([])

// 根据选中的项目过滤环境
const filteredEnvironments = computed(() => {
  if (!activeProject.value) {
    return environments.value
  }
  return environments.value.filter(env => env.project === activeProject.value)
})

// 获取项目列表
const fetchProjects = async () => {
  try {
    const res = await request.get('/api/projects')
    projects.value = res.data
  } catch (error) {
    console.error('Error fetching projects:', error)
    ElMessage.error('获取项目列表失败')
  }
}

// 获取环境列表
const fetchEnvironments = async () => {
  try {
    const res = await request.get('/api/environments')
    environments.value = res.data
    if (environments.value.length > 0) {
      activeEnv.value = environments.value[0].name
      fetchTools()
    }
  } catch (error) {
    console.error('Error fetching environments:', error)
    ElMessage.error('获取环境列表失败')
  }
}

// 获取工具列表
const fetchTools = async () => {
  try {
    let url = '/api/sites'
    const params = {}
    
    if (activeEnv.value) {
      params.env = activeEnv.value
      console.log('Fetching tools for env:', activeEnv.value)
    }
    if (activeProject.value) {
      params.project = activeProject.value
      console.log('Fetching tools for project:', activeProject.value)
    }

    console.log('Request params:', params)
    const response = await request.get(url, { params })
    tools.value = response.data
  } catch (error) {
    console.error('Error fetching tools:', error)
    ElMessage.error('获取工具列表失败')
  }
}

// 处理项目变更
const handleProjectChange = () => {
  if (filteredEnvironments.value.length > 0) {
    const currentEnvExists = filteredEnvironments.value.some(env => env.name === activeEnv.value)
    if (!currentEnvExists) {
      activeEnv.value = filteredEnvironments.value[0].name
    }
  }
  fetchTools()
}

// 处理环境变更
const handleEnvChange = (tab) => {
  activeEnv.value = tab.paneName
  console.log('Environment changed to:', activeEnv.value)
  fetchTools()
}

onMounted(async () => {
  await fetchProjects()
  await fetchEnvironments()
})
</script>

<style scoped>
.tools-container {
  padding: 20px;
}

.filter-header {
  margin-bottom: 20px;
}

:deep(.el-tabs__header) {
  margin-bottom: 20px;
}

:deep(.el-tabs__nav-wrap) {
  padding: 0 20px;
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