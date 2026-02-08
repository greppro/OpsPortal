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
          <Monitor v-if="cat === '可观测性' || cat === '监控'" />
          <Operation v-else-if="cat === 'DevOps' || cat === 'devops'" />
          <ChromeFilled v-else-if="cat === '云平台'" />
          <Grid v-else />
        </el-icon>
        <span>{{ cat }}</span>
      </el-menu-item>
    </el-menu>
  </el-aside>
</template>

<script setup>
import { ref, onMounted, inject } from 'vue'
import { useRouter } from 'vue-router'
import { Monitor, List, Grid, ChromeFilled, Operation, Star } from '@element-plus/icons-vue'
import request from '../../utils/request'

const router = useRouter()
const siteTitle = inject('siteTitle', ref('OpsPortal运维导航'))
const activeProject = ref('')
const projects = ref([])
const activeCategory = ref('all')
const categories = ref([])

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
      // #region agent log
      fetch('http://127.0.0.1:7242/ingest/3c90f934-050e-4fa8-bc2b-4f202bd091da',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'Sidebar.vue:fetchProjects',message:'default project check',data:{projectsCount:projects.value.length,defaultProject:defaultProject?{name:defaultProject.name,label:defaultProject.label}:null,is_defaults:projects.value.map(p=>p.is_default)},timestamp:Date.now(),hypothesisId:'H3,H5'})}).catch(()=>{});
      // #endregion
      if (defaultProject) {
        activeProject.value = defaultProject.name
        emit('project-change', defaultProject.name)
      }
    }
  } catch (error) {
    console.error('Error fetching projects:', error)
  }
}

const fetchCategories = async () => {
  try {
    const [catRes, sitesRes] = await Promise.all([
      request.get('/api/categories').catch(() => ({ data: [] })),
      request.get('/api/sites')
    ])
    const uniqueCategories = new Set()
    const dbCats = Array.isArray(catRes.data) ? catRes.data : []
    dbCats.forEach(c => {
      if (c.name && c.name.trim()) uniqueCategories.add(c.name.trim())
    })
    const tools = sitesRes.data ?? []
    tools.forEach(tool => {
      if (tool.category && tool.category.trim()) {
        uniqueCategories.add(tool.category.trim())
      }
    })
    uniqueCategories.add('其它')
    categories.value = Array.from(uniqueCategories).sort()
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
