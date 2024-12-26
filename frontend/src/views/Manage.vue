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
          @change="handleEnvironmentChange"
          style="width: 200px"
        />
      </div>
    </div>

    <el-table :data="filteredSites" style="width: 100%">
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="url" label="地址" />
      <el-table-column prop="description" label="描述" />
      <el-table-column prop="environment" label="环境">
        <template #default="{ row }">
          <el-tag 
            size="small" 
            :type="getEnvTagType(row.environment)"
          >{{ getEnvLabel(row.environment) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="project" label="项目" />
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button 
            type="primary" 
            link
            @click="handleEdit(row)"
          >编辑</el-button>
          <el-button 
            type="danger" 
            link
            @click="handleDelete(row)"
          >删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="地址" prop="url">
          <el-input v-model="form.url" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" />
        </el-form-item>
        <el-form-item label="环境" prop="environment">
          <environment-select
            v-model="form.environment"
            :project-id="getProjectId(form.project)"
          />
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

const siteList = ref([])
const projectList = ref([])
const environmentList = ref([])
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
  project: ''
})

const rules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  url: [{ required: true, message: '请输入地址', trigger: 'blur' }],
  environment: [{ required: true, message: '请选择环境', trigger: 'change' }],
  project: [{ required: true, message: '请选择项目', trigger: 'change' }]
}

// 筛选网址
const filteredSites = computed(() => {
  return siteList.value.filter(site => {
    const matchProject = !activeProject.value || site.project === activeProject.value
    const matchEnv = !activeEnvironment.value || site.environment === activeEnvironment.value
    return matchProject && matchEnv
  })
})

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

// 清除筛选
const handleProjectClear = () => {
  activeProject.value = ''
}

const handleEnvironmentClear = () => {
  activeEnvironment.value = ''
}

// 添加网址
const handleAdd = () => {
  dialogTitle.value = '添加网址'
  form.value = {
    id: null,
    name: '',
    url: '',
    description: '',
    environment: '',
    project: ''
  }
  dialogVisible.value = true
}

// 编辑网址
const handleEdit = (row) => {
  dialogTitle.value = '编辑网址'
  form.value = { ...row }
  dialogVisible.value = true
}

// 删除网址
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该网址吗？', '提示', {
      type: 'warning'
    })
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

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (form.value.id) {
          await request.put(`/api/sites/${form.value.id}`, form.value)
          ElMessage.success('更新成功')
        } else {
          await request.post('/api/sites', form.value)
          ElMessage.success('添加成功')
        }
        dialogVisible.value = false
        fetchData()
      } catch (error) {
        console.error('Submit error:', error)
        ElMessage.error('操作失败')
      }
    }
  })
}

// 获取项目ID
const getProjectId = (projectName) => {
  const project = projectList.value.find(p => p.name === projectName)
  return project?.id || ''
}

// 处理项目变化
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
  background-color: #f5f7fa;
  border: 1px solid #dcdfe6;
  box-shadow: none !important;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style> 