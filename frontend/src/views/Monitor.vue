<template>
  <div class="monitor-container">
    <div class="page-header">
      <div class="filter-container">
        <el-select 
          v-model="activeProject" 
          placeholder="选择项目"
          clearable
          @clear="handleProjectClear"
          style="width: 200px"
          class="project-select"
          size="large"
        >
          <el-option 
            label="全部项目"
            value=""
          />
          <el-option 
            v-for="proj in projectList" 
            :key="proj.name" 
            :label="proj.label" 
            :value="proj.name"
          />
        </el-select>
        
        <el-tabs v-model="activeTab" class="flex-grow env-tabs">
          <el-tab-pane 
            label="全部环境"
            name="all"
          />
          <el-tab-pane 
            v-for="env in environmentList" 
            :key="env.name" 
            :label="env.label" 
            :name="env.name"
          />
        </el-tabs>
      </div>
    </div>

    <div class="page-content">
      <div class="card-container">
        <el-card 
          v-for="site in filteredSites" 
          :key="site.id" 
          class="monitor-card"
        >
          <div class="card-content">
            <div class="icon">
              <el-icon><component :is="site.icon || 'Link'" /></el-icon>
            </div>
            <h3>{{ site.name }}</h3>
            <p>{{ site.description }}</p>
            <div class="tags">
              <el-tag 
                size="small" 
                :type="getEnvTagType(site.environment)"
              >{{ getEnvLabel(site.environment) }}</el-tag>
              <el-tag size="small" type="success">{{ site.project }}</el-tag>
            </div>
            <div class="button-group">
              <el-button 
                type="primary" 
                size="small"
                @click="handleVisit(site.url)"
              >访问</el-button>
            </div>
          </div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { TrendCharts, List, PieChart, Link } from '@element-plus/icons-vue'
import request from '../utils/request'
import { getEnvTagType, getEnvLabel } from '../config/environment'

const siteList = ref([])
const projectList = ref([])
const environmentList = ref([])

const activeProject = ref('')
const activeTab = ref('all')

// 获取数据
const fetchData = async () => {
  try {
    const [sitesRes, projectsRes, envsRes] = await Promise.all([
      request.get('/api/sites'),
      request.get('/api/projects'),
      request.get('/api/environments')
    ])
    siteList.value = sitesRes.data
    projectList.value = projectsRes.data
    environmentList.value = envsRes.data
  } catch (error) {
    console.error('Error fetching data:', error)
  }
}

// 筛选网址
const filteredSites = computed(() => {
  return siteList.value.filter(site => {
    const matchProject = !activeProject.value || site.project === activeProject.value
    const matchEnv = activeTab.value === 'all' || site.environment === activeTab.value
    return matchProject && matchEnv
  })
})

// 清除项目筛选
const handleProjectClear = () => {
  activeProject.value = ''
}

// 访问网址
const handleVisit = (url) => {
  window.open(url, '_blank')
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.monitor-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #f0f2f5;
}

.page-header {
  background-color: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
}

.filter-container {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 0 20px;
  height: 60px;
  width: 100%;
}

.page-content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

.card-container {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
}

.monitor-card {
  width: 300px;
  text-align: center;
}

.card-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

.icon {
  font-size: 40px;
  color: #409EFF;
}

.tags {
  display: flex;
  gap: 8px;
}

.button-group {
  display: flex;
  gap: 10px;
  margin-top: 10px;
}

h3 {
  margin: 10px 0;
  color: #303133;
}

p {
  color: #909399;
  font-size: 14px;
  margin: 0;
}

.project-select :deep(.el-input__wrapper) {
  background-color: #f5f7fa;
  border: 1px solid #dcdfe6;
  box-shadow: none !important;
  height: 40px;
}

.env-tabs :deep(.el-tabs__header) {
  margin: 0;
  border: none;
  height: 60px;
  line-height: 60px;
}

.env-tabs :deep(.el-tabs__nav-wrap) {
  height: 60px;
}

.env-tabs :deep(.el-tabs__nav) {
  border: none;
  height: 60px;
  display: flex;
  align-items: center;
}

.env-tabs :deep(.el-tabs__item) {
  padding: 0 20px;
  height: 40px;
  line-height: 40px;
  border: none;
  transition: all 0.3s;
  color: #606266;
}

.env-tabs :deep(.el-tabs__item.is-active) {
  color: #409EFF;
  background-color: #ecf5ff;
  border-radius: 4px;
}

.env-tabs :deep(.el-tabs__item:hover) {
  color: #409EFF;
}
</style> 