<template>
  <div class="favorites-container">
    <div class="page-header">
      <h2>我的收藏</h2>
      <span class="subtitle">共 {{ favoriteTools.length }} 个收藏工具</span>
    </div>

    <div class="favorites-content" v-loading="loading">
      <div v-if="favoriteTools.length === 0" class="empty-state">
        <el-empty description="暂无收藏工具，快去添加吧！" />
        <el-button type="primary" @click="$router.push('/monitor')">去浏览工具</el-button>
      </div>
      <tool-grid v-else :tools="favoriteTools" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import request from '../utils/request'
import ToolGrid from '../components/ToolGrid.vue'
import { useFavorites } from '../composables/useFavorites'
import { getEnvLabel } from '../config/environment'

const { favoriteIds, loadFavorites } = useFavorites()
const tools = ref([])
const environments = ref([])
const loading = ref(false)

function getEnvLabelFor(environment, project) {
  const env = environments.value.find(e => e.name === environment && e.project === project)
  return env ? env.label : getEnvLabel(environment)
}

const groupedTools = computed(() => {
  const list = tools.value
  const map = new Map()
  for (const t of list) {
    const key = `${t.project}_${t.name}`
    if (!map.has(key)) {
      map.set(key, {
        key,
        name: t.name,
        description: t.description || '',
        category: t.category || '',
        project: t.project,
        envs: []
      })
    }
    const g = map.get(key)
    g.envs.push({
      id: t.id,
      environment: t.environment,
      url: t.url,
      label: getEnvLabelFor(t.environment, t.project)
    })
  }
  return Array.from(map.values())
})

const favoriteTools = computed(() => {
  return groupedTools.value.filter(g => g.envs.some(e => favoriteIds.value.has(e.id)))
})

const fetchTools = async () => {
  loading.value = true
  try {
    const [sitesRes, envsRes] = await Promise.all([
      request.get('/api/sites'),
      request.get('/api/environments')
    ])
    tools.value = sitesRes.data ?? []
    environments.value = envsRes.data ?? []
  } catch (error) {
    console.error('获取工具列表失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadFavorites()
  fetchTools()
})
</script>

<style scoped>
.favorites-container {
  flex: 1;
  min-height: 0;
  padding: 20px;
  display: flex;
  flex-direction: column;
  background-color: var(--bg-secondary, #f5f6f8);
}

.page-header {
  margin-bottom: 20px;
  display: flex;
  align-items: baseline;
  gap: 12px;
}

.page-header h2 {
  margin: 0;
  font-size: 20px;
  color: var(--text-primary);
}

.subtitle {
  font-size: 14px;
  color: var(--text-secondary);
}

.favorites-content {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
}
</style>
