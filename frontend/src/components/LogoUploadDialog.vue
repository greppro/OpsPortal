<template>
  <el-dialog
    title="更换Logo"
    v-model="visible"
    width="400px"
    class="upload-dialog"
    destroy-on-close
  >
    <el-upload
      class="logo-uploader"
      :action="uploadUrl"
      :headers="uploadHeaders"
      :show-file-list="false"
      :on-success="handleSuccess"
      :before-upload="beforeUpload"
    >
      <img v-if="previewUrl" :src="previewUrl" class="preview">
      <el-icon v-else class="upload-icon"><Plus /></el-icon>
    </el-upload>
    <div class="upload-tip">
      支持 JPG、PNG、SVG 格式，大小不超过 2MB
    </div>
  </el-dialog>
</template>

<script setup>
import { ref } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const visible = ref(false)
const previewUrl = ref('')

const uploadUrl = `${import.meta.env.VITE_API_URL}/api/upload/logo`
const uploadHeaders = {
  Authorization: `Bearer ${localStorage.getItem('token')}`
}

const emit = defineEmits(['success'])

const handleSuccess = (res) => {
  visible.value = false
  emit('success', res)
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
  return true
}

const show = () => {
  visible.value = true
  previewUrl.value = ''
}

defineExpose({ show })
</script>

<style scoped>
.logo-uploader {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: var(--el-transition-duration-fast);
}

.logo-uploader:hover {
  border-color: var(--el-color-primary);
}

.upload-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  text-align: center;
  line-height: 178px;
}

.preview {
  width: 178px;
  height: 178px;
  display: block;
  object-fit: contain;
}

.upload-tip {
  font-size: 12px;
  color: #606266;
  margin-top: 10px;
  text-align: center;
}
</style> 