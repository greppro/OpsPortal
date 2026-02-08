<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <h3>登录</h3>
      </template>
      <el-form :model="loginForm" @submit.prevent="handleLogin">
        <el-form-item>
          <el-input v-model="loginForm.username" placeholder="用户名" prefix-icon="User" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="loginForm.password" type="password" placeholder="密码" prefix-icon="Lock" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" native-type="submit" :loading="loading" style="width: 100%">登录</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const router = useRouter()
const loginForm = ref({
  username: '',
  password: ''
})
const loading = ref(false)

const handleLogin = async () => {
  if (!loginForm.value.username || !loginForm.value.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  
  loading.value = true
  try {
    const response = await request({
      url: '/api/auth/login',
      method: 'post',
      data: {
        username: loginForm.value.username,
        password: loginForm.value.password
      }
    })
    
    if (response.data && response.data.token) {
      localStorage.setItem('token', response.data.token)
      ElMessage.success('登录成功')
      router.push('/monitor')
    } else {
      throw new Error('登录响应数据格式错误')
    }
  } catch (error) {
    console.error('登录错误:', error)
    ElMessage.error(error.response?.data?.message || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: var(--bg-secondary, #f5f6f8);
}

.login-card {
  width: 400px;
}

.login-card :deep(.el-card__header) {
  text-align: center;
  font-size: 18px;
}

h3 {
  margin: 0;
}
</style>
