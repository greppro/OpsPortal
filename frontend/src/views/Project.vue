<template>
  <div class="project-content">
    <div class="header-actions">
      <el-button type="primary" @click="showAddDialog = true">
        <el-icon><Plus /></el-icon>添加项目
      </el-button>
    </div>

    <el-table :data="projectList" style="width: 100%" v-loading="loading">
      <el-table-column prop="name" label="项目标识" />
      <el-table-column prop="label" label="项目名称" />
      <el-table-column label="默认项目" width="100" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.is_default" type="success" size="small">默认</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button size="small" type="primary" @click="handleEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      :title="dialogType === 'add' ? '添加项目' : '编辑项目'"
      v-model="showAddDialog"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="项目标识" prop="name">
          <el-input 
            v-model="form.name" 
            :disabled="dialogType === 'edit'"
            placeholder="如：project1, project2"
          />
        </el-form-item>
        <el-form-item label="项目名称" prop="label">
          <el-input 
            v-model="form.label"
            placeholder="如：项目一、项目二"
          />
        </el-form-item>
        <el-form-item label="设为默认">
          <el-switch v-model="form.is_default" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showAddDialog = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '../utils/request'

const loading = ref(false)
const showAddDialog = ref(false)
const dialogType = ref('add')
const projectList = ref([])
const formRef = ref(null)

const form = ref({
  name: '',
  label: '',
  is_default: false
})

const rules = {
  name: [
    { required: true, message: '请输入项目标识', trigger: 'blur' },
    { pattern: /^[a-z0-9-]+$/, message: '项目标识只能包含小写字母、数字和横线', trigger: 'blur' }
  ],
  label: [{ required: true, message: '请输入项目名称', trigger: 'blur' }]
}

// 获取项目列表
const fetchProjectList = async () => {
  loading.value = true
  try {
    const response = await request.get('/api/projects')
    projectList.value = response.data
  } catch (error) {
    console.error('Fetch projects error:', error)
    ElMessage.error(error.response?.data?.error || '获取项目列表失败')
  } finally {
    loading.value = false
  }
}

// 处理添加/编辑提交
const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (dialogType.value === 'add') {
          await request.post('/api/projects', form.value)
          ElMessage.success('添加成功')
        } else {
          await request.put(`/api/projects/${form.value.id}`, form.value)
          ElMessage.success('更新成功')
        }
        showAddDialog.value = false
        fetchProjectList()
      } catch (error) {
        console.error('Submit error:', error)
        ElMessage.error(error.response?.data?.error || (dialogType.value === 'add' ? '添加失败' : '更新失败'))
      }
    }
  })
}

// 处理编辑
const handleEdit = (row) => {
  dialogType.value = 'edit'
  form.value = { ...row }
  showAddDialog.value = true
}

// 处理删除
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该项目吗？删除后无法恢复。', '提示', {
      type: 'warning'
    })
    await request.delete(`/api/projects/${row.id}`)
    ElMessage.success('删除成功')
    fetchProjectList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Delete error:', error)
      ElMessage.error(error.response?.data?.error || '删除失败')
    }
  }
}

// 重置表单
const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  form.value = {
    name: '',
    label: '',
    is_default: false
  }
}

// 监听对话框关闭
const handleDialogClose = () => {
  resetForm()
  dialogType.value = 'add'
}

onMounted(() => {
  fetchProjectList()
})
</script>

<style scoped>
.project-content {
  padding: 20px;
}

.header-actions {
  margin-bottom: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style> 