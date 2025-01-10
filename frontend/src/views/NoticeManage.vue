<template>
  <div class="notice-manage">
    <div class="header">
      <el-button type="primary" @click="handleAdd">添加公告</el-button>
    </div>
    
    <el-table :data="noticeList" style="width: 100%">
      <el-table-column prop="content" label="公告内容" />
      <el-table-column prop="active" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.active ? 'success' : 'info'">
            {{ row.active ? '激活' : '未激活' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
          <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
      :title="dialogType === 'add' ? '添加公告' : '编辑公告'"
      v-model="dialogVisible"
      width="500px"
      @close="handleClose"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="内容" prop="content">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="3"
            placeholder="请输入公告内容"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="form.active" />
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
import request from '../utils/request'
import { getNotices, createNotice, updateNotice, deleteNotice } from '../api'

const noticeList = ref([])
const dialogVisible = ref(false)
const dialogType = ref('add')
const formRef = ref(null)

const form = ref({
  content: '',
  active: true
})

const rules = {
  content: [
    { required: true, message: '请输入公告内容', trigger: 'blur' },
    { min: 1, max: 200, message: '长度在 1 到 200 个字符', trigger: 'blur' }
  ]
}

const fetchNoticeList = async () => {
  try {
    const res = await getNotices()
    noticeList.value = res.data
  } catch (error) {
    console.error('获取公告列表失败:', error)
    ElMessage.error('获取公告列表失败')
  }
}

const handleAdd = () => {
  dialogType.value = 'add'
  form.value = {
    content: '',
    active: true
  }
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogType.value = 'edit'
  form.value = {
    id: row.id,
    content: row.content,
    active: row.active
  }
  dialogVisible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该公告吗？', '提示', {
      type: 'warning'
    })
    await deleteNotice(row.id)
    ElMessage.success('删除成功')
    fetchNoticeList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除公告失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (dialogType.value === 'add') {
          await createNotice(form.value)
          ElMessage.success('添加成功')
        } else {
          const { id, ...data } = form.value
          await updateNotice(id, data)
          ElMessage.success('更新成功')
        }
        dialogVisible.value = false
        fetchNoticeList()
      } catch (error) {
        console.error('提交失败:', error)
        ElMessage.error('操作失败')
      }
    }
  })
}

const resetForm = () => {
  form.value = {
    content: '',
    active: true
  }
}

const handleClose = () => {
  dialogVisible.value = false
  resetForm()
}

onMounted(() => {
  fetchNoticeList()
})
</script>

<style scoped>
.notice-manage {
  padding: 20px;
}

.header {
  margin-bottom: 20px;
}
</style> 