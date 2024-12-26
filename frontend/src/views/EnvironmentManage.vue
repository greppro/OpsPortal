<template>
  <div class="manage-container">
    <div class="filter-header">
      <el-button type="primary" @click="handleAdd">添加环境</el-button>
      <div class="filters">
        <el-select 
          v-model="activeProject" 
          placeholder="选择项目"
          clearable
          @clear="handleProjectClear"
          style="width: 200px"
        >
          <el-option 
            v-for="proj in projectList" 
            :key="proj.id" 
            :label="proj.label" 
            :value="proj.id"
          />
        </el-select>
      </div>
    </div>

    <el-table :data="filteredEnvironments" style="width: 100%">
      <el-table-column prop="name" label="环境标识" />
      <el-table-column prop="label" label="环境名称" />
      <el-table-column label="所属项目">
        <template #default="{ row }">
          <el-tag size="small">{{ getProjectLabel(row.project_id) }}</el-tag>
          <el-tag size="small" type="info" class="ml-2">{{ row.project }}</el-tag>
        </template>
      </el-table-column>
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
      @closed="handleDialogClose"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <!-- 项目选择 -->
        <el-form-item label="所属项目" prop="project_id">
          <el-select 
            v-model="form.project_id" 
            style="width: 100%"
            placeholder="请选择项目"
            @change="handleProjectChange"
            filterable
          >
            <el-option
              v-for="proj in projectList"
              :key="proj.id"
              :label="proj.label"
              :value="proj.id"
            >
              <span>{{ proj.label }}</span>
              <span class="project-id">{{ proj.name }}</span>
            </el-option>
          </el-select>
        </el-form-item>

        <!-- 环境标识 -->
        <el-form-item label="环境标识" prop="name">
          <el-input 
            v-model="form.name" 
            placeholder="如: dev, test, prod"
          />
        </el-form-item>

        <!-- 环境名称 -->
        <el-form-item label="环境名称" prop="label">
          <el-input 
            v-model="form.label" 
            placeholder="如: 开发环境, 测试环境"
          />
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
import { ref, computed, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '../utils/request'

const environmentList = ref([])
const projectList = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref(null)
const activeProject = ref('')

const form = ref({
  id: null,
  name: '',
  label: '',
  project_id: undefined,
  project: ''
})

const rules = {
  project_id: [
    { required: true, message: '请选择项目', trigger: 'change' }
  ],
  name: [
    { required: true, message: '请输入环境标识', trigger: 'blur' }
  ],
  label: [
    { required: true, message: '请输入环境名称', trigger: 'blur' }
  ]
}

// 筛选环境
const filteredEnvironments = computed(() => {
  if (!activeProject.value) return environmentList.value
  return environmentList.value.filter(env => env.project_id === activeProject.value)
})

// 获取项目名称
const getProjectLabel = (projectId) => {
  const project = projectList.value.find(p => p.id === projectId)
  return project?.label || '-'
}

// 获取数据
const fetchData = async () => {
  try {
    const [envsRes, projsRes] = await Promise.all([
      request.get('/api/environments'),
      request.get('/api/projects')
    ])
    console.log('Projects data:', projsRes.data)
    environmentList.value = envsRes.data
    projectList.value = projsRes.data
  } catch (error) {
    console.error('Error fetching data:', error)
    ElMessage.error('获取数据失败')
  }
}

// 清除筛选
const handleProjectClear = () => {
  activeProject.value = ''
}

// 添加环境
const handleAdd = () => {
  dialogTitle.value = '添加环境'
  form.value = {
    id: null,
    name: '',
    label: '',
    project_id: undefined,
    project: ''
  }
  dialogVisible.value = true
}

// 编辑环境
const handleEdit = (row) => {
  dialogTitle.value = '编辑环境'
  nextTick(() => {
    form.value = { 
      ...row,
      project_id: row.project_id,
      project: row.project
    }
  })
  dialogVisible.value = true
}

// 删除环境
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该环境吗？', '提示', {
      type: 'warning'
    })
    await request.delete(`/api/environments/${row.id}`)
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
        const submitData = {
          ...form.value,
          project_id: Number(form.value.project_id)
        }
        
        if (form.value.id) {
          await request.put(`/api/environments/${form.value.id}`, submitData)
          ElMessage.success('更新成功')
        } else {
          await request.post('/api/environments', submitData)
          ElMessage.success('添加成功')
        }
        dialogVisible.value = false
        fetchData()
      } catch (error) {
        console.error('Submit error:', error)
        ElMessage.error(error.response?.data?.error || '操作失败')
      }
    }
  })
}

// 处理项目变化
const handleProjectChange = (projectId) => {
  const project = projectList.value.find(p => p.id === projectId)
  if (project) {
    form.value.project = project.name
    console.log('Selected project:', project)
  }
}

// 判断项目是否应该禁用
const isProjectDisabled = (projectId) => {
  // 如果是编辑模式，且环境已经被使用，则只允许选择当前项目
  if (form.value.id && form.value.project) {
    const hasTools = true // 这里可以添加检查逻辑
    if (hasTools) {
      return projectId !== form.value.project_id
    }
  }
  return false
}

// 添加对话框关闭处理
const handleDialogClose = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  form.value = {
    id: null,
    name: '',
    label: '',
    project_id: undefined,
    project: ''
  }
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

.ml-2 {
  margin-left: 8px;
}

.project-id {
  color: #909399;
  font-size: 13px;
  margin-left: 12px;
}

:deep(.el-select-dropdown__item) {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 12px;
}

:deep(.el-form-item) {
  margin-bottom: 20px;
}

:deep(.el-select) {
  width: 100%;
}
</style> 