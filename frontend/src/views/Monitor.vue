<template>
  <div class="tools-container">
    <div class="tools-toolbar">
      <el-input
        v-model="searchQuery"
        placeholder="搜索工具、描述、URL、项目、环境"
        clearable
        class="search-input"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
    </div>

    <div class="tools-content">
      <section v-if="showRecentTools" class="recent-section">
        <div class="section-header">
          <h2>最近访问</h2>
          <span>最多显示 8 个工具</span>
        </div>
        <tool-grid :tools="recentTools" @record-recent="recordRecentTool" />
      </section>

      <div class="section-header">
        <h2>{{ searchQuery.trim() ? '搜索结果' : '全部工具' }}</h2>
        <span>共 {{ filteredTools.length }} 个工具</span>
      </div>
      <tool-grid :tools="filteredTools" @record-recent="recordRecentTool" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import request from '../utils/request'
import ToolGrid from '../components/ToolGrid.vue'
import { getEnvLabel } from '../config/environment'

const RECENT_TOOLS_KEY = 'opsportal-recent-tools'
const RECENT_TOOLS_LIMIT = 8

const props = defineProps({
  activeProject: {
    type: String,
    default: ''
  },
  activeCategory: {
    type: String,
    default: 'all'
  }
})

const tools = ref([])
const projects = ref([])
const environments = ref([])
const searchQuery = ref('')
const recentTools = ref([])

const currentProject = computed(() => props.activeProject)
const currentCategory = computed(() => props.activeCategory)

// 当前项目下的环境列表（用于取 label）
const projectEnvironments = computed(() => {
  if (!currentProject.value) return []
  return environments.value.filter(env => env.project === currentProject.value)
})

function getEnvLabelFor(environment) {
  const env = projectEnvironments.value.find(e => e.name === environment)
  return env ? env.label : getEnvLabel(environment)
}

// 按 (project, name) 分组，得到逻辑工具列表（同一工具任一条有分类即采用，避免被空分类覆盖）
const groupedTools = computed(() => {
  const list = tools.value
  const map = new Map()
  for (const t of list) {
    const key = `${t.project}_${t.name}`
    const cat = (t.category != null && t.category !== undefined) ? String(t.category).trim() : ''
    if (!map.has(key)) {
      map.set(key, {
        key,
        name: t.name,
        description: t.description || '',
        category: cat,
        project: t.project,
        envs: []
      })
    }
    const g = map.get(key)
    if (cat && !(g.category || '').trim()) g.category = cat
    g.envs.push({
      id: t.id,
      environment: t.environment,
      url: t.url,
      label: getEnvLabelFor(t.environment)
    })
  }
  return Array.from(map.values())
})

// 按分类筛选后的分组列表（比较时 trim，避免空格导致不匹配）
const filteredTools = computed(() => {
  let result = groupedTools.value
  const cat = (currentCategory.value || '').trim()
  if (cat && cat !== 'all') {
    result = result.filter(g => (g.category || '').trim() === cat)
  }
  const query = searchQuery.value.trim().toLowerCase()
  if (query) {
    result = result.filter(tool => toolMatchesSearch(tool, query))
  }
  return result
})

const showRecentTools = computed(() => {
  const cat = (currentCategory.value || '').trim()
  return recentTools.value.length > 0 && !searchQuery.value.trim() && (!cat || cat === 'all')
})

function toolMatchesSearch(tool, query) {
  const fields = [
    tool.name,
    tool.description,
    tool.project,
    tool.category,
    ...(tool.envs || []).flatMap(env => [
      env.environment,
      env.label,
      env.url
    ])
  ]
  return fields.some(value => String(value || '').toLowerCase().includes(query))
}

function normalizeRecentTool(tool) {
  return {
    key: tool.key,
    name: tool.name,
    description: tool.description || '',
    category: tool.category || '',
    project: tool.project || '',
    envs: (tool.envs || []).map(env => ({
      id: env.id,
      environment: env.environment,
      label: env.label,
      url: env.url
    })),
    visitedAt: new Date().toISOString()
  }
}

function loadRecentTools() {
  try {
    const stored = JSON.parse(localStorage.getItem(RECENT_TOOLS_KEY) || '[]')
    recentTools.value = Array.isArray(stored) ? stored.slice(0, RECENT_TOOLS_LIMIT) : []
  } catch (error) {
    console.error('Failed to load recent tools:', error)
    recentTools.value = []
  }
}

function recordRecentTool(tool) {
  if (!tool?.key) return
  const current = recentTools.value.filter(item => item.key !== tool.key)
  recentTools.value = [normalizeRecentTool(tool), ...current].slice(0, RECENT_TOOLS_LIMIT)
  localStorage.setItem(RECENT_TOOLS_KEY, JSON.stringify(recentTools.value))
}

const fetchProjects = async () => {
  try {
    const res = await request.get('/api/projects')
    projects.value = res.data
  } catch (error) {
    console.error('Error fetching projects:', error)
    ElMessage.error('获取项目列表失败')
  }
}

const fetchEnvironments = async () => {
  try {
    const res = await request.get('/api/environments')
    environments.value = res.data
  } catch (error) {
    console.error('Error fetching environments:', error)
    ElMessage.error('获取环境列表失败')
  }
}

let fetchToolsId = 0
const fetchTools = async () => {
  const thisId = ++fetchToolsId
  try {
    const params = {}
    if (currentProject.value) {
      params.project = currentProject.value
    }
    const response = await request.get('/api/sites', { params })
    if (thisId !== fetchToolsId) return
    const raw = response.data ?? []
    tools.value = raw
  } catch (error) {
    if (thisId !== fetchToolsId) return
    console.error('Error fetching tools:', error)
    ElMessage.error('获取工具列表失败')
  }
}

watch(() => props.activeProject, () => {
  fetchTools()
})

onMounted(async () => {
  loadRecentTools()
  await fetchProjects()
  await fetchEnvironments()
  fetchTools()
})
</script>

<style scoped>
.tools-container {
  flex: 1;
  min-height: 0;
  padding: 20px;
  display: flex;
  flex-direction: column;
  background-color: var(--bg-secondary, #f5f6f8);
}

.tools-toolbar {
  flex-shrink: 0;
  display: flex;
  justify-content: flex-start;
  padding-bottom: 16px;
}

.search-input {
  width: min(520px, 100%);
}

.tools-content {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}

.recent-section {
  margin-bottom: 8px;
}

.section-header {
  display: flex;
  align-items: baseline;
  gap: 10px;
  padding: 0 16px;
  margin: 4px 0 8px;
}

.section-header h2 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary, #1e293b);
}

.section-header span {
  font-size: 13px;
  color: var(--text-secondary, #64748b);
}
</style>
