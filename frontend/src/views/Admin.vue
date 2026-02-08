<template>
  <div class="admin-container">
    <aside class="admin-sidebar">
      <nav class="admin-nav">
        <div class="nav-section">
          <span class="nav-group-label">业务管理</span>
          <button
            v-for="item in businessItems"
            :key="item.name"
            type="button"
            class="nav-item"
            :class="{ active: activeTab === item.name }"
            @click="activeTab = item.name"
          >
            {{ item.label }}
          </button>
        </div>
        <div class="nav-divider" aria-hidden="true" />
        <div class="nav-section">
          <span class="nav-group-label">系统与外观</span>
          <button
            v-for="item in systemItems"
            :key="item.name"
            type="button"
            class="nav-item"
            :class="{ active: activeTab === item.name }"
            @click="activeTab = item.name"
          >
            {{ item.label }}
          </button>
        </div>
      </nav>
    </aside>
    <main class="page-content">
      <component :is="currentComponent" />
    </main>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import Manage from './Manage.vue'
import CategoryManage from './CategoryManage.vue'
import Environment from './EnvironmentManage.vue'
import Project from './Project.vue'
import NoticeManage from './NoticeManage.vue'
import LogoManage from './LogoManage.vue'
import SettingsManage from './SettingsManage.vue'

const businessItems = [
  { name: 'sites', label: '网址管理' },
  { name: 'categories', label: '分类管理' },
  { name: 'environments', label: '环境管理' },
  { name: 'projects', label: '项目管理' },
  { name: 'notices', label: '公告管理' }
]

const systemItems = [
  { name: 'logo', label: 'Logo管理' },
  { name: 'settings', label: '系统设置' }
]

const activeTab = ref('sites')

const componentMap = {
  sites: Manage,
  categories: CategoryManage,
  environments: Environment,
  projects: Project,
  notices: NoticeManage,
  logo: LogoManage,
  settings: SettingsManage
}

const currentComponent = computed(() => componentMap[activeTab.value])
</script>

<style scoped>
.admin-container {
  flex: 1;
  min-height: 0;
  width: 100%;
  display: flex;
  background-color: var(--bg-secondary);
}

.admin-sidebar {
  width: 220px;
  flex-shrink: 0;
  background-color: var(--bg-primary);
  border-right: 1px solid var(--border-color);
  overflow-y: auto;
}

.admin-nav {
  display: flex;
  flex-direction: column;
  padding: 20px 12px;
  gap: 8px;
}

.nav-section {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.nav-group-label {
  font-size: 11px;
  font-weight: 600;
  color: var(--text-sidebar-muted, var(--text-secondary));
  text-transform: uppercase;
  letter-spacing: 0.06em;
  padding: 12px 12px 6px;
  margin: 0 -4px 0 0;
}

.nav-item {
  padding: 10px 14px;
  font-size: 14px;
  color: var(--text-secondary);
  background: none;
  border: none;
  border-radius: var(--radius-md, 8px);
  cursor: pointer;
  text-align: left;
  width: 100%;
  transition: color 0.15s, background-color 0.15s;
  outline: none;
  box-shadow: none;
}

.nav-item:focus {
  outline: none;
  box-shadow: none;
}

.nav-item:hover {
  color: var(--text-primary);
  background-color: var(--bg-secondary);
}

.nav-item.active {
  color: var(--primary);
  font-weight: 500;
}

.nav-divider {
  height: 1px;
  margin: 12px 8px;
  background-color: var(--border-color);
}

.page-content {
  flex: 1;
  min-width: 0;
  padding: 24px;
  overflow-y: auto;
  background-color: var(--bg-secondary);
}
</style> 