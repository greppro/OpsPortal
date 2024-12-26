<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <h2 class="login-title">登录</h2>
      </template>
      
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
        @keyup.enter="handleSubmit"
      >
        <el-form-item label="用户名" prop="username">
          <el-input 
            v-model="form.username" 
            placeholder="请输入用户名"
            :disabled="loading"
          />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input 
            v-model="form.password" 
            type="password" 
            show-password
            placeholder="请输入密码"
            :disabled="loading"
          />
        </el-form-item>
        <el-form-item>
          <el-button 
            type="primary" 
            @click="handleSubmit"
            :loading="loading"
            style="width: 100%"
          >
            {{ loading ? '登录中...' : '登录' }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '../utils/request'

const router = useRouter()
const loading = ref(false)

const form = ref({
  username: '',
  password: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const formRef = ref(null)

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const response = await request.post('/api/auth/login', form.value)
        localStorage.setItem('token', response.data.token)
        localStorage.setItem('user', response.data.user.username)
        
        ElMessage({
          type: 'success',
          message: '登录成功，正在跳转...',
          duration: 2000
        })
        
        // 获取保存的重定向路径
        const redirectPath = localStorage.getItem('redirectPath') || '/admin'
        localStorage.removeItem('redirectPath') // 清除保存的路径
        
        // 延迟跳转，让用户看到提示
        setTimeout(() => {
          router.push(redirectPath).then(() => {
            // 在路由跳转完成后再刷新页面
            window.location.reload()
          })
        }, 1000)
      } catch (error) {
        ElMessage({
          type: 'error',
          message: error.response?.data?.error || '登录失败',
          duration: 3000
        })
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px;
}

.login-card {
  width: 400px;
}

.login-card :deep(.el-card__header) {
  text-align: center;
  padding: 15px;
}

h2.login-title {
  margin: 0;
  font-size: 18px;
  color: #303133;
}
</style> 