<template>
  <Notice @notice-visible-change="handleNoticeVisibleChange" />
  <el-container class="layout-container" :class="{ 'no-notice': !showNoticeBar }">
    <el-aside width="200px" class="aside">
      <div class="system-title">OpsPortal运维导航</div>
      <el-menu
        :default-active="activeMenu"
        class="menu"
        router
      >
        <el-menu-item index="/monitor">
          <el-icon><Monitor /></el-icon>
          <span>网址导航</span>
        </el-menu-item>
        <el-menu-item index="/admin">
          <el-icon><Setting /></el-icon>
          <span>后台管理</span>
          <el-tag size="small" type="warning" class="auth-tag" v-if="!isLoggedIn">
            需要登录
          </el-tag>
        </el-menu-item>
      </el-menu>
    </el-aside>
    
    <el-container class="main-container">
      <el-header class="header">
        <div class="header-logo">
          <img 
            :src="logoUrl" 
            alt="Logo" 
            class="site-logo" 
            v-if="logoUrl"
            @error="handleLogoError"
          >
        </div>
        <div class="header-right" v-if="showUserInfo">
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-icon><User /></el-icon>
              {{ username }}
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="changePassword">修改密码</el-dropdown-item>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      
      <el-main class="main-content">
        <router-view></router-view>
      </el-main>
    </el-container>
    <ChangePasswordDialog ref="changePasswordDialogRef" />
    <LogoUploadDialog 
      ref="logoUploadDialogRef" 
      @success="handleLogoSuccess"
    />
  </el-container>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Monitor, Setting, User } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import ChangePasswordDialog from './components/ChangePasswordDialog.vue'
import Notice from './components/Notice.vue'
import request from './utils/request'
import LogoUploadDialog from './components/LogoUploadDialog.vue'

const router = useRouter()
const route = useRoute()

const isLoggedIn = computed(() => !!localStorage.getItem('token'))
const username = computed(() => localStorage.getItem('user') || '')
const activeMenu = computed(() => route.path)

// 只在管理页面显示用户信息
const showUserInfo = computed(() => isLoggedIn.value && route.path === '/admin')

const changePasswordDialogRef = ref(null)
const logoUploadDialogRef = ref(null)

const handleCommand = async (command) => {
  if (command === 'logout') {
    // 显示确认对话框
    try {
      await ElMessageBox.confirm(
        '确定要退出登录吗？',
        '提示',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      )
      
      // 清除用户信息
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      
      ElMessage({
        type: 'success',
        message: '已退出登录，正在跳转...',
        duration: 2000
      })
      
      // 延迟跳转，让用户看到提示
      setTimeout(() => {
        router.push('/monitor')
        // 触发页面刷新
        window.location.reload()
      }, 1000)
    } catch {
      // 用户取消退出
    }
  } else if (command === 'changePassword') {
    changePasswordDialogRef.value?.show()
  } else if (command === 'uploadLogo') {
    logoUploadDialogRef.value?.show()
  }
}

const showNoticeBar = ref(true)

const handleNoticeVisibleChange = (visible) => {
  showNoticeBar.value = visible
}

const logoUrl = ref('')
const isAdmin = computed(() => isLoggedIn.value)

// 获取 Logo
const fetchLogo = async () => {
  try {
    const res = await request.get('/api/logo')
    logoUrl.value = res.data.url
  } catch (error) {
    console.error('Error fetching logo:', error)
  }
}

// 上传相关
const uploadHeaders = {
  Authorization: `Bearer ${localStorage.getItem('token')}`
}

const handleLogoSuccess = (res) => {
  logoUrl.value = res.url
  ElMessage.success('Logo更新成功')
}

const beforeLogoUpload = (file) => {
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
  return true
}

const uploadUrl = `${import.meta.env.VITE_API_URL}/api/upload/logo`

const handleLogoError = (e) => {
  console.error('Logo image load failed:', e)
  logoUrl.value = ''  // 如果图片加载失败，清空URL
}

onMounted(() => {
  fetchLogo()
})
</script>

<style scoped>
.layout-container {
  height: 100vh;
  padding-top: 40px;
  overflow: hidden;
  transition: padding-top 0.3s;
}

.layout-container.no-notice {
  padding-top: 0;
}

.aside {
  background-color: #304156;
  color: #fff;
  height: 100%;
  overflow-y: hidden;
}

.main-container {
  height: 100%;
  overflow: hidden;
}

.main-content {
  padding: 0;
  height: calc(100% - 64px);
  overflow-y: auto;
}

.logo {
  height: 60px;
  line-height: 60px;
  text-align: center;
  font-size: 18px;
  font-weight: bold;
  color: #fff;
  background-color: #2b2f3a;
}

.menu {
  border-right: none;
  background-color: #304156;
}

.header {
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  height: 64px;
}

.header-logo {
  position: relative;
  flex: 1;
  height: 100%;
}

.site-logo {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  max-height: 40px;
  max-width: 180px;
  object-fit: contain;
  display: block;
}

.system-title {
  height: 60px;
  line-height: 60px;
  text-align: center;
  font-size: 18px;
  font-weight: bold;
  color: #fff;
  background-color: #2b2f3a;
}

.header-right {
  display: flex;
  align-items: center;
  margin-left: auto;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  color: #606266;
}

.user-info .el-icon {
  margin-right: 8px;
}

:deep(.el-menu) {
  border-right: none;
}

:deep(.el-menu-item) {
  color: #bfcbd9;
}

:deep(.el-menu-item.is-active) {
  color: #409EFF;
  background-color: #263445;
}

:deep(.el-menu-item:hover) {
  background-color: #263445;
}

.auth-tag {
  margin-left: 8px;
}

.app-container {
  margin-top: 40px;
}

.logo-area {
  padding: 20px;
  text-align: center;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.site-logo {
  max-width: 160px;
  max-height: 60px;
  margin-bottom: 10px;
}

.logo-placeholder {
  font-size: 24px;
  color: #fff;
  margin-bottom: 10px;
}

.logo-upload {
  margin-top: 10px;
}

.logo-upload :deep(.el-upload) {
  display: block;
}

.logo-upload :deep(.el-button) {
  width: 100%;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  color: #fff;
}

.logo-upload :deep(.el-button:hover) {
  background: rgba(255, 255, 255, 0.2);
}

.upload-dialog {
  .el-upload {
    width: 100%;
    text-align: center;
  }
}
</style>
