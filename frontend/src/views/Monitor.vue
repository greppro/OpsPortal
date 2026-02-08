<template>
  <div class="tools-container">
    <div class="tools-content">
      <tool-grid :tools="filteredTools" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import request from '../utils/request'
import ToolGrid from '../components/ToolGrid.vue'
import { getEnvLabel } from '../config/environment'

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
    // #region agent log
    const groupCats = groupedTools.value.map(g => ({ name: g.name, category: g.category, match: (g.category || '').trim() === cat }))
    fetch('http://127.0.0.1:7242/ingest/3c90f934-050e-4fa8-bc2b-4f202bd091da',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'Monitor.vue:filteredTools',message:'category filter',data:{activeCategory:cat,filteredCount:result.length,totalGroups:groupedTools.value.length,sampleGroupCategories:groupCats.slice(0,20),runId:'post-fix'},timestamp:Date.now(),hypothesisId:'H2,H3,H4'})}).catch(()=>{});
    // #endregion
  }
  return result
})

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
    // #region agent log
    const sample = raw.slice(0, 15).map(t => ({ name: t.name, category: t.category, hasCategory: 'category' in t, url: t.url, hasUrl: !!t.url }))
    const uniq = [...new Set(raw.map(t => (t && t.category != null ? String(t.category).trim() : '')))].filter(Boolean)
    fetch('http://127.0.0.1:7242/ingest/3c90f934-050e-4fa8-bc2b-4f202bd091da',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'Monitor.vue:fetchTools',message:'sites loaded',data:{toolsLength:raw.length,sample,uniqueCategoriesFromApi:uniq},timestamp:Date.now(),hypothesisId:'H1,H2,H5'})}).catch(()=>{});
    // #endregion
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

.tools-content {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}
</style>
