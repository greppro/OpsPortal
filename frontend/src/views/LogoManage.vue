<template>
  <div class="logo-manage">
    <div class="current-logo">
      <h3>当前Logo</h3>
      <div class="logo-preview">
        <img 
          v-if="currentLogo" 
          :src="currentLogo" 
          alt="Current Logo"
          @error="handleImageError"
        >
        <div v-else class="no-logo">暂无Logo</div>
      </div>
    </div>
    
    <div class="logo-actions">
      <el-button type="primary" @click="handleUpload">更换Logo</el-button>
      <el-button 
        type="danger" 
        @click="handleDelete" 
        :disabled="!currentLogo"
      >删除Logo</el-button>
    </div>

    <el-dialog
      title="更换Logo"
      v-model="dialogVisible"
      width="500px"
    >
      <div class="upload-container">
        <el-upload
          class="logo-uploader"
          action="#"
          :headers="uploadHeaders"
          :show-file-list="false"
          :on-success="handleSuccess"
          :before-upload="beforeUpload"
          :on-change="handleChange"
          :auto-upload="false"
          ref="uploadRef"
        >
          <img v-if="previewUrl" :src="previewUrl" class="preview">
          <div v-else class="upload-placeholder">
            <el-icon class="upload-icon"><Plus /></el-icon>
            <div class="upload-text">点击选择图片</div>
          </div>
        </el-upload>
        
        <div class="upload-tip">
          支持 JPG、PNG、SVG 格式，大小不超过 2MB
        </div>
      </div>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="handleCancel">取消</el-button>
          <el-button 
            type="primary" 
            @click="handleSubmit"
            :disabled="!selectedFile"
          >确认上传</el-button>
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

const dialogVisible = ref(false)
const currentLogo = ref('')
const previewUrl = ref('')
const uploadRef = ref(null)
const selectedFile = ref(null)

const uploadUrl = `${import.meta.env.VITE_API_URL}/api/upload/logo`
const uploadHeaders = {
  Authorization: `Bearer ${localStorage.getItem('token')}`
}

// 获取当前 Logo
const fetchLogo = async () => {
  try {
    const res = await request.get('/api/logo')
    currentLogo.value = res.data.url || ''
  } catch (error) {
    console.error('Error fetching logo:', error)
    currentLogo.value = ''
  }
}

const handleUpload = () => {
  dialogVisible.value = true
  previewUrl.value = ''
}

const handleSuccess = (res) => {
  dialogVisible.value = false
  currentLogo.value = res.url
  ElMessage.success('Logo更新成功')
}

const handleSubmit = async () => {
  if (!selectedFile.value) return

  const formData = new FormData()
  formData.append('logo', selectedFile.value)

  try {
    const res = await request.post('/api/upload/logo', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        ...uploadHeaders
      }
    })
    handleSuccess(res.data)
  } catch (error) {
    console.error('Upload error:', error)
    ElMessage.error('上传失败')
  }
}

const handleCancel = () => {
  dialogVisible.value = false
  previewUrl.value = ''
  selectedFile.value = null
}

const handleImageError = (e) => {
  console.error('Logo image load failed:', e)
  currentLogo.value = ''  // 如果图片加载失败，清空URL
  ElMessage.error('Logo 图片加载失败')
}

const beforeUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过 2MB!')
    return false
  }

  // 预览
  previewUrl.value = URL.createObjectURL(file)
  selectedFile.value = file
  return true // 返回 true 允许选择文件
}

const handleChange = (file) => {
  if (file.status === 'ready') {
    // 文件已选择但未上传
    selectedFile.value = file.raw
    previewUrl.value = URL.createObjectURL(file.raw)
  }
}

const handleDelete = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要删除当前Logo吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const res = await request.delete('/api/logo')
    currentLogo.value = ''
    ElMessage.success(res.data.message || 'Logo删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Delete logo error:', error)
      ElMessage.error(error.response?.data?.error || '删除失败')
      // 重新获取 logo 状态
      await fetchLogo()
    }
  }
}

onMounted(() => {
  fetchLogo()
})
</script>

<style scoped>
.logo-manage {
  padding: 20px;
}

.current-logo {
  margin-bottom: 20px;
}

.logo-preview {
  margin-top: 10px;
  width: 200px;
  height: 100px;
  border: 1px dashed #dcdfe6;
  border-radius: 4px;
  display: flex;
  justify-content: center;
  align-items: center;
  background: var(--bg-secondary, #f5f6f8);
}

.logo-preview img {
  max-width: 180px;
  max-height: 80px;
  object-fit: contain;
}

.no-logo {
  color: #909399;
  font-size: 14px;
}

.upload-container {
  text-align: center;
}

.logo-uploader {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: var(--el-transition-duration-fast);
  display: inline-block;
}

.logo-uploader:hover {
  border-color: var(--el-color-primary);
}

.upload-placeholder {
  width: 300px;
  height: 150px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.upload-icon {
  font-size: 28px;
  color: #8c939d;
}

.upload-text {
  color: #8c939d;
  font-size: 14px;
  margin-top: 8px;
}

.preview {
  width: 300px;
  height: 150px;
  display: block;
  object-fit: contain;
}

.upload-tip {
  font-size: 12px;
  color: var(--text-secondary, #64748b);
  margin-top: 10px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.logo-actions {
  display: flex;
  gap: 12px;
}
</style> 