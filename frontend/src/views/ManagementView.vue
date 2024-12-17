<template>
  <div class="management-container">
    <div class="actions">
      <el-button type="primary" @click="showToolDialog()">添加工具</el-button>
    </div>
    
    <el-tabs v-model="activeEnv" @tab-click="handleEnvChange">
      <el-tab-pane label="开发环境" name="dev">
        <tool-grid 
          :tools="tools" 
          :is-management="true"
          @edit="showToolDialog"
          @delete="handleDelete"
        />
      </el-tab-pane>
      <el-tab-pane label="生产环境" name="prod">
        <tool-grid 
          :tools="tools" 
          :is-management="true"
          @edit="showToolDialog"
          @delete="handleDelete"
        />
      </el-tab-pane>
    </el-tabs>

    <el-dialog 
      :title="editingTool.id ? '编辑工具' : '添加工具'" 
      v-model="dialogVisible"
    >
      <el-form :model="editingTool" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="editingTool.name" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="editingTool.description" type="textarea" />
        </el-form-item>
        <el-form-item label="URL">
          <el-input v-model="editingTool.url" />
        </el-form-item>
        <el-form-item label="环境">
          <el-select v-model="editingTool.environment">
            <el-option label="开发环境" value="dev" />
            <el-option label="生产环境" value="prod" />
          </el-select>
        </el-form-item>
        <el-form-item label="分类">
          <el-input v-model="editingTool.category" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveTool">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'
import ToolGrid from '../components/ToolGrid.vue'

const activeEnv = ref('dev')
const tools = ref([])
const dialogVisible = ref(false)
const editingTool = ref({
  name: '',
  description: '',
  url: '',
  environment: 'dev',
  category: ''
})

const showToolDialog = (tool = null) => {
  if (tool) {
    editingTool.value = { ...tool }
  } else {
    editingTool.value = {
      name: '',
      description: '',
      url: '',
      environment: activeEnv.value,
      category: ''
    }
  }
  dialogVisible.value = true
}

const saveTool = async () => {
  try {
    console.log('Saving tool:', editingTool.value)
    if (editingTool.value.id) {
      await axios.put(`http://localhost:8080/api/tools/${editingTool.value.id}`, editingTool.value)
    } else {
      await axios.post('http://localhost:8080/api/tools', editingTool.value)
    }
    dialogVisible.value = false
    fetchTools(activeEnv.value)
    ElMessage.success('保存成功')
  } catch (error) {
    console.error('Error saving tool:', error.response || error)
    ElMessage.error('保存失败')
  }
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除这个工具吗？', '提示', {
      type: 'warning'
    })
    await axios.delete(`http://localhost:8080/api/tools/${id}`)
    fetchTools(activeEnv.value)
    ElMessage.success('删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Error deleting tool:', error)
      ElMessage.error('删除失败')
    }
  }
}

const fetchTools = async (env) => {
  try {
    const response = await axios.get(`http://localhost:8080/api/tools?environment=${env}`)
    tools.value = response.data
  } catch (error) {
    console.error('Error fetching tools:', error)
    ElMessage.error('获取工具列表失败')
  }
}

watch(activeEnv, (newEnv) => {
  fetchTools(newEnv)
})

onMounted(() => {
  fetchTools(activeEnv.value)
})
</script>
