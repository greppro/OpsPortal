<template>
  <div class="manage-container">
    <div class="filter-header">
      <el-button type="primary" @click="handleAdd">添加网址</el-button>
      <div class="filters">
        <el-select
          v-model="activeProject"
          placeholder="选择项目"
          clearable
          @clear="handleProjectClear"
          @change="handleProjectChange"
          style="width: 200px"
        >
          <el-option
            v-for="proj in projectList"
            :key="proj.name"
            :label="proj.label"
            :value="proj.name"
          />
        </el-select>

        <environment-select
          v-model="activeEnvironment"
          :project-id="getProjectId(activeProject)"
          style="width: 200px"
        />
      </div>
    </div>

    <el-table :data="groupedSites" row-key="key" style="width: 100%">
      <el-table-column type="expand">
        <template #default="{ row }">
          <el-table :data="row.envs" size="small" border>
            <el-table-column prop="environment" label="环境" width="120">
              <template #default="{ row: envRow }">
                <el-tag size="small" :type="getEnvTagType(envRow.environment)">
                  {{ getEnvLabel(envRow.environment) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="url" label="地址" min-width="200" show-overflow-tooltip />
            <el-table-column v-if="probeEnabled" label="可用性" min-width="240">
              <template #default="{ row: envRow }">
                <span v-if="envRow.probe_disabled" class="probe-cell probe-cell--off">未开启检测</span>
                <span v-else class="probe-cell">
                  <span
                    class="probe-dot"
                    :class="'probe-dot--' + (statusOf(envRow.id)?.state || 'pending')"
                    aria-hidden="true"
                  ></span>
                  {{ probeSummary(statusOf(envRow.id) || { state: 'pending' }) }}
                </span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200" align="right">
              <template #default="{ row: envRow }">
                <el-button
                  v-if="probeEnabled && !envRow.probe_disabled"
                  type="primary"
                  link
                  :loading="checkingId === envRow.id"
                  @click="handleCheck(envRow)"
                >检测</el-button>
                <el-button type="primary" link @click="handleEdit(envRow)">编辑</el-button>
                <el-button type="danger" link @click="handleDelete(envRow)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </template>
      </el-table-column>
      <el-table-column prop="name" label="名称" width="160" />
      <el-table-column prop="project" label="项目" width="120" />
      <el-table-column prop="description" label="描述" min-width="180" show-overflow-tooltip />
      <el-table-column label="环境数" width="90" align="center">
        <template #default="{ row }">
          {{ row.envs.length }}
        </template>
      </el-table-column>
    </el-table>

    <el-dialog :title="dialogTitle" v-model="dialogVisible" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="地址" prop="url">
          <el-input v-model="form.url" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" />
        </el-form-item>
        <el-form-item label="项目" prop="project">
          <el-select v-model="form.project" style="width: 100%">
            <el-option
              v-for="proj in projectList"
              :key="proj.name"
              :label="proj.label"
              :value="proj.name"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="环境" prop="environment">
          <environment-select
            v-model="form.environment"
            :project-id="getProjectId(form.project)"
          />
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-select v-model="form.category" placeholder="不选则归为其它" style="width: 100%">
            <el-option
              v-for="cat in categoryOptions"
              :key="cat"
              :label="cat"
              :value="cat"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="可用性检测">
          <div class="probe-switch">
            <el-switch v-model="form.probe_enabled" />
            <span class="form-hint">定时检测这个地址能否访问。关闭后不检测，也不会出现在 Prometheus 服务发现里</span>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '../utils/request'
import { getEnvTagType, getEnvLabel } from '../config/environment'
import EnvironmentSelect from '../components/EnvironmentSelect.vue'
import { useProbeStatus } from '../composables/useProbeStatus'
import { probeSummary } from '../utils/probe'

const DEFAULT_CATEGORY = '其它'

const {
  enabled: probeEnabled,
  statusOf,
  refresh: refreshProbe,
  setResult: setProbeResult
} = useProbeStatus()

const siteList = ref([])
const projectList = ref([])
const environmentList = ref([])
const categoryOptions = ref([DEFAULT_CATEGORY])
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref(null)

const activeProject = ref('')
const activeEnvironment = ref('')

const form = ref({
  id: null,
  name: '',
  url: '',
  description: '',
  environment: '',
  project: '',
  category: DEFAULT_CATEGORY,
  probe_enabled: true
})

const rules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  url: [{ required: true, message: '请输入地址', trigger: 'blur' }],
  project: [{ required: true, message: '请选择项目', trigger: 'change' }],
  environment: [{ required: true, message: '请选择环境', trigger: 'change' }]
}

const filteredSites = computed(() => {
  return siteList.value.filter(site => {
    const matchProject = !activeProject.value || site.project === activeProject.value
    const matchEnv = !activeEnvironment.value || site.environment === activeEnvironment.value
    return matchProject && matchEnv
  })
})

const groupedSites = computed(() => {
  const list = filteredSites.value
  const map = new Map()
  for (const t of list) {
    const key = `${t.project}_${t.name}`
    if (!map.has(key)) {
      map.set(key, {
        key,
        name: t.name,
        description: t.description || '',
        project: t.project,
        envs: []
      })
    }
    map.get(key).envs.push({ ...t })
  }
  return Array.from(map.values())
})

const fetchData = async () => {
  try {
    const [sitesRes, projectsRes, envsRes, catRes] = await Promise.all([
      request.get('/api/sites'),
      request.get('/api/projects'),
      request.get('/api/environments'),
      request.get('/api/categories').catch(() => ({ data: [] }))
    ])
    siteList.value = sitesRes.data ?? []
    projectList.value = projectsRes.data ?? []
    environmentList.value = envsRes.data ?? []
    const cats = Array.isArray(catRes.data) ? catRes.data : []
    const names = cats.map(c => c.name).filter(Boolean)
    if (!names.includes(DEFAULT_CATEGORY)) names.unshift(DEFAULT_CATEGORY)
    categoryOptions.value = names.sort((a, b) => a.localeCompare(b))
  } catch (error) {
    console.error('Error fetching data:', error)
  }
}

const handleProjectClear = () => {
  activeProject.value = ''
}

const handleAdd = () => {
  dialogTitle.value = '添加网址'
  form.value = {
    id: null,
    name: '',
    url: '',
    description: '',
    environment: '',
    project: '',
    category: DEFAULT_CATEGORY,
    probe_enabled: true
  }
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑网址'
  const { probe_disabled, ...rest } = row
  form.value = {
    ...rest,
    category: row.category && row.category.trim() ? row.category.trim() : DEFAULT_CATEGORY,
    probe_enabled: !probe_disabled
  }
  dialogVisible.value = true
}

// 立即检测一个地址，结果同步到首页状态
const checkingId = ref(null)
const handleCheck = async (row) => {
  checkingId.value = row.id
  try {
    const res = await request.post(`/api/probe/check/${row.id}`, null, { timeout: 30000 })
    setProbeResult(row.id, res.data)
    const text = probeSummary(res.data)
    if (res.data?.state === 'up') {
      ElMessage.success(text)
    } else {
      ElMessage.warning(text)
    }
  } catch (error) {
    console.error('Check error:', error)
  } finally {
    checkingId.value = null
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该网址吗？', '提示', { type: 'warning' })
    await request.delete(`/api/sites/${row.id}`)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Delete error:', error)
      ElMessage.error('删除失败')
    }
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        const { probe_enabled, ...rest } = form.value
        const payload = { ...rest, probe_disabled: !probe_enabled }
        if (!payload.category || !payload.category.trim()) {
          payload.category = DEFAULT_CATEGORY
        }
        if (form.value.id) {
          await request.put(`/api/sites/${form.value.id}`, payload)
          ElMessage.success('更新成功')
        } else {
          await request.post('/api/sites', payload)
          ElMessage.success('添加成功')
        }
        dialogVisible.value = false
        fetchData()
        // 后端保存后会立即检测一次，稍后刷新状态
        if (probe_enabled) setTimeout(refreshProbe, 2000)
      } catch (error) {
        console.error('Submit error:', error)
        ElMessage.error('操作失败')
      }
    }
  })
}

const getProjectId = (projectName) => {
  const project = projectList.value.find(p => p.name === projectName)
  return project?.id || ''
}

const handleProjectChange = () => {
  activeEnvironment.value = ''
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.manage-container {
  padding: 20px;
}

.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.filters {
  display: flex;
  gap: 20px;
}

:deep(.el-input__wrapper) {
  background-color: var(--bg-secondary, #f5f6f8);
  border: 1px solid #dcdfe6;
  box-shadow: none !important;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.probe-switch {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

.form-hint {
  font-size: 12px;
  line-height: 1.5;
  color: var(--text-secondary, #6b7280);
}

.probe-cell {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--text-primary, #111827);
}

.probe-cell--off {
  color: var(--text-secondary, #6b7280);
}

.probe-dot {
  flex-shrink: 0;
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.probe-dot--up {
  background: #22c55e;
}

.probe-dot--down {
  background: #ef4444;
}

.probe-dot--pending {
  box-shadow: inset 0 0 0 1.5px var(--text-secondary, #9ca3af);
}
</style>
