<template>
  <div class="main-layout">
    <Notice v-if="!isAdminRoute" @notice-visible-change="handleNoticeVisibleChange" />
    <div v-if="isAdminRoute" class="admin-fullscreen">
      <Header />
      <router-view />
    </div>
    <el-container v-else class="layout-container" :class="{ 'has-notice': showNoticeBar }">
      <Sidebar 
        @project-change="handleProjectChange"
        @category-change="handleCategorySelect"
      />
      <el-container class="main-container">
        <Header />
        <el-main class="main-content">
          <router-view 
            :active-project="activeProject"
            :active-category="activeCategory"
          ></router-view>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { ref, onMounted, provide, computed } from 'vue'
import { useRoute } from 'vue-router'
import Header from './components/Header.vue'
import Sidebar from './components/Sidebar.vue'
import Notice from '../components/Notice.vue'
import request from '../utils/request'

const route = useRoute()
const isAdminRoute = computed(() => route.name === 'admin')

const showNoticeBar = ref(true)
const activeProject = ref('')
const activeCategory = ref('all')
const siteTitle = ref('OpsPortal运维导航')

const refreshSiteTitle = async () => {
  try {
    const res = await request.get('/api/config/site-title')
    siteTitle.value = res.data?.site_title ?? 'OpsPortal运维导航'
  } catch (_) {
    siteTitle.value = 'OpsPortal运维导航'
  }
}

onMounted(() => {
  refreshSiteTitle()
})
provide('siteTitle', siteTitle)
provide('refreshSiteTitle', refreshSiteTitle)

const handleNoticeVisibleChange = (visible) => {
  showNoticeBar.value = visible
}

const handleProjectChange = (val) => {
  activeProject.value = val
}

const handleCategorySelect = (val) => {
  activeCategory.value = val
}
</script>

<style scoped>
.main-layout {
  height: 100%;
  min-height: 100vh;
}

.layout-container {
  height: 100vh;
  padding-top: 0;
  overflow: hidden;
  transition: padding-top 0.3s;
  background-color: var(--bg-secondary, #f5f6f8);
}

.layout-container.has-notice {
  padding-top: 40px;
}

.main-container {
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.main-content {
  padding: 0;
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  background-color: var(--bg-secondary, #f5f6f8);
  display: flex;
  flex-direction: column;
}

.admin-fullscreen {
  height: 100vh;
  width: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  background-color: var(--bg-secondary, #f5f6f8);
}
</style>
