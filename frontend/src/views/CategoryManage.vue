<template>
  <div class="category-manage">
    <div class="toolbar">
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>
        添加分类
      </el-button>
    </div>

    <el-table :data="categories" border stripe style="width: 100%">
      <el-table-column prop="name" label="分类名称" min-width="150" />
      <el-table-column prop="description" label="描述" min-width="200" />
      <el-table-column prop="icon" label="图标" width="80" align="center">
        <template #default="scope">
          <el-icon v-if="scope.row.icon" :size="24">
            <component :is="scope.row.icon" />
          </el-icon>
        </template>
      </el-table-column>
      <el-table-column prop="tool_count" label="工具数量" width="100" align="center" />
      <el-table-column label="操作" width="200" align="center">
        <template #default="scope">
          <el-button size="small" @click="handleEdit(scope.row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加/编辑分类对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="分类名称" prop="name">
          <el-input v-model="form.name" placeholder="例如：可观测性、DevOps" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input 
            v-model="form.description" 
            type="textarea" 
            :rows="3"
            placeholder="分类的简要描述"
          />
        </el-form-item>
        <el-form-item label="图标" prop="icon">
          <el-select v-model="form.icon" placeholder="选择图标">
            <el-option label="监控 (Monitor)" value="Monitor" />
            <el-option label="操作 (Operation)" value="Operation" />
            <el-option label="云平台 (ChromeFilled)" value="ChromeFilled" />
            <el-option label="网格 (Grid)" value="Grid" />
            <el-option label="列表 (List)" value="List" />
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
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import request from '../utils/request'

const categories = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref('添加分类')
const formRef = ref(null)
const editingId = ref(null)       // 编辑时的原分类名（用于查找要更新的工具）
const editingCategoryId = ref(null) // 编辑时的分类 ID（来自 DB 时用于 PUT）

const form = ref({
  name: '',
  description: '',
  icon: ''
})

const rules = {
  name: [
    { required: true, message: '请输入分类名称', trigger: 'blur' }
  ]
}

// 获取分类列表：后端存储的分类 + 工具中的分类（合并并带工具数）
const fetchCategories = async () => {
  try {
    const [catRes, sitesRes] = await Promise.all([
      request.get('/api/categories').catch(() => ({ data: [] })),
      request.get('/api/sites')
    ])
    const dbCats = Array.isArray(catRes.data) ? catRes.data : []
    const tools = sitesRes.data ?? []
    // #region agent log
    fetch('http://127.0.0.1:7242/ingest/3c90f934-050e-4fa8-bc2b-4f202bd091da',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'CategoryManage.vue:fetchCategories',message:'fetchCategories called',data:{dbCategoriesLength:dbCats.length,toolsLength:tools.length,runId:'post-fix'},timestamp:Date.now(),hypothesisId:'H1'})}).catch(()=>{});
    // #endregion

    const toolCountByCategory = new Map()
    tools.forEach(tool => {
      if (tool.category && tool.category.trim()) {
        const cat = tool.category.trim()
        toolCountByCategory.set(cat, (toolCountByCategory.get(cat) || 0) + 1)
      }
    })

    const byName = new Map()
    dbCats.forEach(c => {
      byName.set(c.name, {
        id: c.id,
        name: c.name,
        description: c.description || '',
        icon: c.icon || getDefaultIcon(c.name),
        tool_count: toolCountByCategory.get(c.name) || 0
      })
    })
    toolCountByCategory.forEach((count, name) => {
      if (!byName.has(name)) {
        byName.set(name, {
          id: null,
          name,
          description: '',
          icon: getDefaultIcon(name),
          tool_count: count
        })
      }
    })
    if (!byName.has('其它')) {
      byName.set('其它', {
        id: null,
        name: '其它',
        description: '未归类的工具',
        icon: 'Grid',
        tool_count: toolCountByCategory.get('其它') || 0
      })
    }
    const list = Array.from(byName.values()).sort((a, b) => a.name.localeCompare(b.name))
    categories.value = list
    // #region agent log
    fetch('http://127.0.0.1:7242/ingest/3c90f934-050e-4fa8-bc2b-4f202bd091da',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'CategoryManage.vue:fetchCategories',message:'fetchCategories result',data:{categoryNames:list.map(c=>c.name),listLength:list.length,runId:'post-fix'},timestamp:Date.now(),hypothesisId:'H1'})}).catch(()=>{});
    // #endregion
  } catch (error) {
    // #region agent log
    fetch('http://127.0.0.1:7242/ingest/3c90f934-050e-4fa8-bc2b-4f202bd091da',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'CategoryManage.vue:fetchCategories',message:'fetchCategories error',data:{err:String(error)},timestamp:Date.now(),hypothesisId:'H2'})}).catch(()=>{});
    // #endregion
    console.error('Error fetching categories:', error)
    ElMessage.error('获取分类列表失败')
  }
}

// 根据分类名称返回默认图标
const getDefaultIcon = (name) => {
  if (name.includes('监控') || name.includes('可观测')) return 'Monitor'
  if (name.includes('DevOps') || name.includes('运维')) return 'Operation'
  if (name.includes('云')) return 'ChromeFilled'
  return 'Grid'
}

// 添加分类
const handleAdd = () => {
  dialogTitle.value = '添加分类'
  editingId.value = null
  editingCategoryId.value = null
  form.value = {
    name: '',
    description: '',
    icon: ''
  }
  dialogVisible.value = true
}

// 编辑分类
const handleEdit = (row) => {
  dialogTitle.value = '编辑分类'
  editingId.value = row.name
  editingCategoryId.value = row.id != null ? row.id : null
  form.value = {
    name: row.name,
    description: row.description,
    icon: row.icon
  }
  dialogVisible.value = true
}

// 删除分类
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除分类"${row.name}"吗？该操作会将所有属于该分类的工具的分类字段清空。`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    if (row.id != null) {
      await request.delete(`/api/categories/${row.id}`)
    }
    const res = await request.get('/api/sites')
    const toolsToUpdate = res.data.filter(tool => tool.category === row.name)
    for (const tool of toolsToUpdate) {
      await request.put(`/api/sites/${tool.id}`, { ...tool, category: '' })
    }
    ElMessage.success('分类删除成功')
    fetchCategories()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Error deleting category:', error)
      ElMessage.error('删除分类失败')
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    // #region agent log
    fetch('http://127.0.0.1:7242/ingest/3c90f934-050e-4fa8-bc2b-4f202bd091da',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'CategoryManage.vue:handleSubmit',message:'submit branch',data:{editingId:editingId.value,formName:form.value.name,isAdd:!editingId.value},timestamp:Date.now(),hypothesisId:'H3'})}).catch(()=>{});
    // #endregion

    if (editingId.value) {
      if (editingCategoryId.value != null) {
        await request.put(`/api/categories/${editingCategoryId.value}`, {
          name: form.value.name,
          description: form.value.description || '',
          icon: form.value.icon || ''
        })
      }
      const res = await request.get('/api/sites')
      const toolsToUpdate = res.data.filter(tool => tool.category === editingId.value)
      for (const tool of toolsToUpdate) {
        await request.put(`/api/sites/${tool.id}`, { ...tool, category: form.value.name })
      }
      ElMessage.success('分类更新成功')
    } else {
      await request.post('/api/categories', {
        name: form.value.name,
        description: form.value.description || '',
        icon: form.value.icon || ''
      })
      // #region agent log
      fetch('http://127.0.0.1:7242/ingest/3c90f934-050e-4fa8-bc2b-4f202bd091da',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'CategoryManage.vue:handleSubmit',message:'add branch persisted',data:{formName:form.value.name},timestamp:Date.now(),hypothesisId:'H1',runId:'post-fix'})}).catch(()=>{});
      // #endregion
      ElMessage.success('分类创建成功，请在创建工具时使用该分类')
    }

    dialogVisible.value = false
    fetchCategories()
  } catch (error) {
    // #region agent log
    fetch('http://127.0.0.1:7242/ingest/3c90f934-050e-4fa8-bc2b-4f202bd091da',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'CategoryManage.vue:handleSubmit',message:'submit error',data:{err:String(error)},timestamp:Date.now(),hypothesisId:'H2'})}).catch(()=>{});
    // #endregion
    console.error('Error submitting form:', error)
  }
}

onMounted(() => {
  fetchCategories()
})
</script>

<style scoped>
.category-manage {
  background-color: var(--bg-primary, #fff);
  padding: 20px;
  border-radius: var(--radius-md, 8px);
}

.toolbar {
  margin-bottom: 20px;
}
</style>
