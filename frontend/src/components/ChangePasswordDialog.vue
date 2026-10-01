<template>
  <el-dialog
    title="修改密码"
    v-model="visible"
    width="400px"
    @close="handleClose"
  >
    <el-alert
      v-if="insecure"
      type="warning"
      :closable="false"
      show-icon
      class="insecure-alert"
      title="你正在使用默认密码，任何人都可以用它登录后台，请立即修改。"
    />
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item label="旧密码" prop="oldPassword">
        <el-input
          v-model="form.oldPassword"
          type="password"
          show-password
          placeholder="请输入旧密码"
        />
      </el-form-item>
      <el-form-item label="新密码" prop="newPassword">
        <el-input
          v-model="form.newPassword"
          type="password"
          show-password
          placeholder="至少 8 位"
        />
      </el-form-item>
      <el-form-item label="确认新密码" prop="confirmPassword">
        <el-input
          v-model="form.confirmPassword"
          type="password"
          show-password
          placeholder="请再次输入新密码"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="visible = false">{{ insecure ? '稍后修改' : '取消' }}</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import request from '../utils/request'
import { DEFAULT_PASSWORD_FLAG } from '../utils/auth'

const visible = ref(false)
const insecure = ref(false)
const submitting = ref(false)
const formRef = ref(null)

const form = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== form.value.newPassword) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const validateNewPassword = (rule, value, callback) => {
  if (value && value === form.value.oldPassword) {
    callback(new Error('新密码不能与旧密码相同'))
  } else {
    callback()
  }
}

const rules = {
  oldPassword: [
    { required: true, message: '请输入旧密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 8, message: '密码长度不能小于8位', trigger: 'blur' },
    { validator: validateNewPassword, trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      const res = await request.post('/api/auth/change-password', {
        oldPassword: form.value.oldPassword,
        newPassword: form.value.newPassword
      })
      // 修改密码后旧 token 全部失效，换成后端返回的新 token
      if (res.data?.token) {
        localStorage.setItem('token', res.data.token)
      }
      localStorage.removeItem(DEFAULT_PASSWORD_FLAG)
      insecure.value = false
      ElMessage.success('密码修改成功')
      visible.value = false
      handleClose()
    } catch (error) {
      // 错误提示由请求拦截器统一弹出
      console.error('Change password error:', error)
    } finally {
      submitting.value = false
    }
  })
}

const handleClose = () => {
  formRef.value?.resetFields()
  form.value = {
    oldPassword: '',
    newPassword: '',
    confirmPassword: ''
  }
}

// 暴露方法给父组件；insecure 为 true 时显示默认密码警告
defineExpose({
  show: (options = {}) => {
    insecure.value = !!options.insecure
    visible.value = true
  }
})
</script>

<style scoped>
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.insecure-alert {
  margin-bottom: 16px;
}
</style>
