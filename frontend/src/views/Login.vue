<template>
  <div class="login-page">
    <div class="login-container">
      <div class="background-pattern"></div>
      <div class="brand-title">OpsPortal</div>
      <el-card class="login-card">
        <template #header>
          <h3 class="login-title">登录</h3>
        </template>
        
        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          label-width="0"
          @keyup.enter="handleSubmit"
        >
          <el-form-item prop="username">
            <el-input 
              v-model="form.username" 
              placeholder="请输入用户名"
              :prefix-icon="User"
              :disabled="loading"
            />
          </el-form-item>
          <el-form-item prop="password">
            <el-input 
              v-model="form.password" 
              type="password" 
              show-password
              placeholder="请输入密码"
              :prefix-icon="Lock"
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
.login-page {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 200px; /* 左侧菜单宽度 */
  overflow: hidden;
}

.login-container {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
  position: relative;
  overflow: hidden;
}

.background-pattern {
  position: absolute; /* 改回 absolute */
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  opacity: 0.4;
  background-color: #f5f7fa;
  background-image: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%239C92AC' fill-opacity='0.15'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
  background-size: 60px 60px;
  animation: moveBackground 60s linear infinite;
}

.login-card {
  width: 400px;
  margin: 0;
  position: relative;
  z-index: 1;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border: none;
}

.login-title {
  margin: 0;
  text-align: center;
  font-size: 20px;
  color: #303133;
}

:deep(.el-card__header) {
  padding: 20px;
  border-bottom: 1px solid #ebeef5;
  background: rgba(255, 255, 255, 0.8);
}

:deep(.el-card__body) {
  padding: 30px;
}

:deep(.el-form-item:last-child) {
  margin-bottom: 0;
}

:deep(.el-input__wrapper) {
  box-shadow: 0 0 0 1px #dcdfe6 inset;
  background: rgba(255, 255, 255, 0.9);
}

:deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #c0c4cc inset;
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px #409eff inset;
}

@keyframes moveBackground {
  from {
    background-position: 0 0;
  }
  to {
    background-position: 100% 100%;
  }
}

/* 添加一些渐变光效 */
.login-container::before {
  position: absolute; /* 改回 absolute */
  content: '';
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0) 70%);
  animation: rotate 30s linear infinite;
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.brand-title {
  position: absolute; /* 改回 absolute */
  top: 10%;
  left: 50%;
  transform: translateX(-50%);
  font-size: 48px;
  font-weight: 200;
  color: rgba(64, 158, 255, 0.15);
  letter-spacing: 4px;
  text-transform: uppercase;
  pointer-events: none;
  font-family: "Helvetica Neue", Helvetica, "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", "微软雅黑", Arial, sans-serif;
  z-index: 1;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  animation: fadeIn 1.5s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translate(-50%, -20px);
  }
  to {
    opacity: 1;
    transform: translate(-50%, 0);
  }
}
</style> 