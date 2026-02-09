<template>
  <el-aside width="220px" class="sidebar-container">
    <div class="system-title" @click="goHome">{{ siteTitle }}</div>
    
    <!-- 项目选择器 -->
    <div class="project-selector">
      <el-select 
        v-model="activeProject" 
        placeholder="选择项目"
        clearable
        @change="handleProjectChange"
        size="default"
        popper-class="sidebar-project-dropdown"
        class="sidebar-project-select"
      >
        <el-option 
          v-for="proj in projects" 
          :key="proj.name" 
          :label="proj.label" 
          :value="proj.name"
        />
      </el-select>
    </div>
    
    <!-- 分类菜单 -->
    <el-menu
      :default-active="activeCategory"
      class="sidebar-menu"
      @select="handleCategorySelect"
    >
      <el-menu-item index="all">
        <el-icon><List /></el-icon>
        <span>全部工具</span>
      </el-menu-item>

      <el-menu-item index="favorites">
        <el-icon><Star /></el-icon>
        <span>我的收藏</span>
      </el-menu-item>
      
      <el-menu-item 
        v-for="cat in categories" 
        :key="cat"
        :index="cat"
      >
        <el-icon>
          <component :is="iconMap[categoryIconMap[cat] || 'Grid']" />
        </el-icon>
        <span>{{ cat }}</span>
      </el-menu-item>
    </el-menu>
  </el-aside>
</template>

<script setup>
import { ref, onMounted, inject } from 'vue'
import { useRouter } from 'vue-router'
import {
  Monitor,
  List,
  Grid,
  ChromeFilled,
  Operation,
  Star,
  DataLine,
  DataAnalysis,
  Cpu,
  Connection,
  Document,
  Folder,
  Setting,
  SetUp,
  Service,
  Tools,
  Bell,
  PieChart,
} from '@element-plus/icons-vue'
import request from '../../utils/request'

const router = useRouter()
const siteTitle = inject('siteTitle', ref('OpsPortal运维导航'))
const activeProject = ref('')
const projects = ref([])
const activeCategory = ref('all')
const categories = ref([])
const categoryIconMap = ref({}) // 分类名 -> 图标组件名（来自 API）
const iconMap = {
  Monitor,
  Operation,
  ChromeFilled,
  Grid,
  List,
  DataLine,
  DataAnalysis,
  Cpu,
  Connection,
  Document,
  Folder,
  Setting,
  SetUp,
  Service,
  Tools,
  Bell,
  PieChart,
}

// 定义 props 接收父组件传递的数据（如果有必要）
// 这里改为自包含逻辑，但通过 emit 通知父组件或 RouterView
const emit = defineEmits(['project-change', 'category-change'])

const goHome = () => {
  // 保留当前选中的项目；若未选择则使用默认项目
  if (!activeProject.value && projects.value.length > 0) {
    const defaultProject = projects.value.find(p => p.is_default)
    if (defaultProject) {
      activeProject.value = defaultProject.name
      emit('project-change', defaultProject.name)
    }
  }
  activeCategory.value = 'all'
  emit('category-change', 'all')
  router.push('/monitor')
}

const fetchProjects = async () => {
  try {
    const res = await request.get('/api/projects')
    projects.value = res.data
    
    // 如果当前没有选中项目，且有默认项目，则自动选中默认项目
    if (!activeProject.value && projects.value.length > 0) {
      const defaultProject = projects.value.find(p => p.is_default)
      if (defaultProject) {
        activeProject.value = defaultProject.name
        emit('project-change', defaultProject.name)
      }
    }
  } catch (error) {
    console.error('Error fetching projects:', error)
  }
}

function getDefaultIcon(name) {
  if (!name) return 'Grid'
  const n = name.trim()
  if (n.includes('监控') || n.includes('可观测')) return 'Monitor'
  if (n.includes('DevOps') || n.includes('运维')) return 'Operation'
  if (n.includes('云')) return 'ChromeFilled'
  return 'Grid'
}

const fetchCategories = async () => {
  try {
    const [catRes, sitesRes] = await Promise.all([
      request.get('/api/categories').catch(() => ({ data: [] })),
      request.get('/api/sites')
    ])
    const uniqueCategories = new Set()
    const iconByName = {}
    const dbCats = Array.isArray(catRes.data) ? catRes.data : []
    dbCats.forEach(c => {
      if (c.name && c.name.trim()) {
        uniqueCategories.add(c.name.trim())
        iconByName[c.name.trim()] = (c.icon && c.icon.trim()) ? c.icon.trim() : getDefaultIcon(c.name)
      }
    })
    const tools = sitesRes.data ?? []
    tools.forEach(tool => {
      if (tool.category && tool.category.trim()) {
        uniqueCategories.add(tool.category.trim())
        if (!(tool.category.trim() in iconByName)) {
          iconByName[tool.category.trim()] = getDefaultIcon(tool.category)
        }
      }
    })
    uniqueCategories.add('其它')
    if (!('其它' in iconByName)) iconByName['其它'] = getDefaultIcon('其它')
    categories.value = Array.from(uniqueCategories).sort()
    categoryIconMap.value = iconByName
  } catch (error) {
    console.error('Error fetching categories:', error)
  }
}

const handleProjectChange = (val) => {
  emit('project-change', val)
}

const handleCategorySelect = (val) => {
  if (val === 'favorites') {
    router.push('/favorites')
    activeCategory.value = 'favorites'
  } else {
    // 如果之前是收藏页，现在切回其他分类，需要跳转回 monitor
    if (activeCategory.value === 'favorites' || router.currentRoute.value.path === '/favorites') {
      router.push('/monitor')
    }
    activeCategory.value = val
    emit('category-change', val)
  }
}

onMounted(() => {
  fetchProjects()
  fetchCategories()
})
</script>

<style scoped>
.sidebar-container {
  background-color: var(--bg-sidebar, #f9fafb);
  color: var(--text-sidebar, #374151);
  height: 100%;
  overflow-y: auto;
  border-right: 1px solid var(--border-color, #e5e7eb);
  .el-menu {
    border-right: none;
    background-color: transparent;
  }
}

.system-title {
  height: 64px;
  line-height: 64px;
  text-align: center;
  font-size: 18px;
  font-weight: 700;
  color: var(--primary, #615ced);
  background: transparent;
  cursor: pointer;
  letter-spacing: -0.5px;
}

.system-title:hover {
  background-color: transparent;
  opacity: 0.8;
}

.system-title:focus,
.system-title:focus-visible {
  outline: none;
  box-shadow: none;
}

/* 选择项目区域与侧栏背景一致，更干净 */
.project-selector {
  padding: 16px 16px;
  background-color: transparent;
  border-bottom: none;
}

.sidebar-project-select :deep(.el-input__wrapper) {
  background-color: var(--bg-primary, #fff);
  box-shadow: 0 1px 2px rgba(0,0,0,0.05);
  border: 1px solid var(--border-color, #e5e7eb);
  border-radius: 99px; /* Pill shape */
  padding-left: 16px;
  transition: all 0.2s;
}

.sidebar-project-select :deep(.el-input__wrapper:hover),
.sidebar-project-select :deep(.el-input.is-focus .el-input__wrapper) {
  border-color: var(--primary, #615ced);
  box-shadow: 0 0 0 2px var(--primary-light);
}

.sidebar-project-select :deep(.el-input__inner) {
  color: var(--text-primary, #111827);
  font-weight: 500;
}

.sidebar-project-select :deep(.el-input__inner::placeholder) {
  color: var(--text-secondary, #9ca3af);
}

.sidebar-project-select :deep(.el-select__caret) {
  color: var(--text-secondary, #9ca3af);
}

:deep(.el-menu-item) {
  color: var(--text-sidebar, #374151);
  margin: 4px 12px;
  border-radius: 8px;
  height: 44px;
  line-height: 44px;
}

:deep(.el-menu-item.is-active) {
  color: var(--primary, #615ced);
  background-color: var(--primary-light, #eef2ff);
  font-weight: 600;
}

:deep(.el-menu-item:hover) {
  background-color: rgba(0, 0, 0, 0.04);
  color: var(--text-primary);
}
</style>

<!-- 下拉框与主内容区底色一致 -->
<style>
.sidebar-project-dropdown.el-select-dropdown {
  background-color: var(--bg-primary, #fff);
  border: 1px solid var(--border-color, #e5e7eb);
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.1);
  padding: 8px;
}
.sidebar-project-dropdown .el-select-dropdown__item {
  color: var(--text-primary, #111827);
  border-radius: 8px;
  margin: 2px 0;
  height: 36px;
  line-height: 36px;
}
.sidebar-project-dropdown .el-select-dropdown__item.hover,
.sidebar-project-dropdown .el-select-dropdown__item:hover {
  background-color: var(--border-light, #f3f4f6);
  color: var(--text-primary);
}
.sidebar-project-dropdown .el-select-dropdown__item.selected {
  background-color: var(--primary-light, #eef2ff);
  color: var(--primary, #615ced);
  font-weight: 600;
}
.sidebar-project-dropdown .el-popper__arrow::before {
  background: var(--bg-primary, #fff);
  border: 1px solid var(--border-color, #e5e7eb);
}
</style>
