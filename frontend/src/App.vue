<template>
  <el-container class="layout-container">
    <el-header>
      <div class="header-content">
        <h2>OpsPortal 运维导航平台</h2>
      </div>
    </el-header>
    <el-container class="main-container">
      <el-aside width="200px">
        <el-menu
          :router="true"
          class="el-menu-vertical"
          :default-active="currentRoute"
          background-color="#545c64"
          text-color="#fff"
          active-text-color="#ffd04b"
        >
          <el-menu-item index="/">
            <el-icon><Monitor /></el-icon>
            <span>网址导航</span>
          </el-menu-item>
          <el-menu-item index="/management">
            <el-icon><Setting /></el-icon>
            <span>网址管理</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      
      <el-main>
        <router-view @showChangePassword="showChangePassword = true"></router-view>
      </el-main>
    </el-container>
  </el-container>

  <el-dialog
    v-model="showChangePassword"
    title="修改密码"
    width="400px"
  >
    <el-form
      ref="passwordFormRef"
      :model="passwordForm"
      :rules="passwordRules"
      label-width="100px"
    >
      <el-form-item label="原密码" prop="oldPassword">
        <el-input
          v-model="passwordForm.oldPassword"
          type="password"
          show-password
          placeholder="请输入原密码"
        />
      </el-form-item>
      <el-form-item label="新密码" prop="newPassword">
        <el-input
          v-model="passwordForm.newPassword"
          type="password"
          show-password
          placeholder="请输入新密码"
        />
      </el-form-item>
      <el-form-item label="确认密码" prop="confirmPassword">
        <el-input
          v-model="passwordForm.confirmPassword"
          type="password"
          show-password
          placeholder="请再次输入新密码"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="showChangePassword = false">取消</el-button>
        <el-button type="primary" @click="handleChangePassword" :loading="changing">
          确认
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Monitor, Setting } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

const router = useRouter()
const route = useRoute()
const user = ref(null)
const passwordFormRef = ref(null)
const showChangePassword = ref(false)
const changing = ref(false)

const currentRoute = computed(() => route.path)

const showUserInfo = computed(() => {
  return route.path === '/management' && user.value
})

watch(
  () => route.path,
  () => {
    const userStr = localStorage.getItem('user')
    if (userStr) {
      user.value = JSON.parse(userStr)
    }
  },
  { immediate: true }
)

const handleLogout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  user.value = null
  ElMessage.success('已退出登录')
  router.push('/')
}

const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const passwordRules = {
  oldPassword: [
    { required: true, message: '请输入原密码', trigger: 'blur' },
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== passwordForm.value.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

const handleChangePassword = async () => {
  if (!passwordFormRef.value) return
  
  await passwordFormRef.value.validate(async (valid) => {
    if (valid) {
      changing.value = true
      try {
        const response = await axios.post('http://localhost:8080/api/auth/change-password', {
          username: user.value.username,
          oldPassword: passwordForm.value.oldPassword,
          newPassword: passwordForm.value.newPassword
        })
        
        ElMessage.success('密码修改成功')
        showChangePassword.value = false
        passwordForm.value = {
          oldPassword: '',
          newPassword: '',
          confirmPassword: ''
        }
      } catch (error) {
        ElMessage.error(error.response?.data?.error || '修改密码失败')
      } finally {
        changing.value = false
      }
    }
  })
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-container {
  flex: 1;
  overflow: hidden;
}

.el-header {
  background-color: #409EFF;
  color: white;
  line-height: 80px;
  height: 80px !important;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.12);
  padding: 0 20px;
  z-index: 1000;
}

.el-aside {
  background-color: #545c64;
  border-right: solid 1px #434a50;
  height: 100%;
}

.el-main {
  background-color: #f5f7fa;
  padding: 20px;
}

.el-menu {
  border-right: none;
}

.el-menu-item {
  border-left: 3px solid transparent;
}

.el-menu-item.is-active {
  border-left: 3px solid #ffd04b;
  background-color: #434a50 !important;
}

.el-menu-item:hover {
  background-color: #434a50 !important;
}

h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1400px;
  margin: 0 auto;
  height: 100%;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-info .el-button {
  color: white;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
