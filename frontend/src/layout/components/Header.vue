<template>
  <el-header class="app-header">
    <div class="header-left">
      <button
        v-if="isAdminRoute"
        type="button"
        class="site-title-btn"
        @click="goHome"
      >
        {{ siteTitle }}
      </button>
      <div class="header-logo">
        <img 
          :src="logoUrl" 
          alt="Logo" 
          class="site-logo" 
          v-if="logoUrl"
          @error="handleLogoError"
        >
      </div>
    </div>
    <div class="header-right">
      <ThemeToggle />
      
      <el-tooltip content="后台管理" placement="bottom">
        <el-button
          circle
          @click="goToAdmin"
          class="admin-btn"
        >
          <el-icon><Setting /></el-icon>
        </el-button>
      </el-tooltip>
      
      <el-dropdown v-if="showUserInfo" @command="handleCommand">
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
    
    <ChangePasswordDialog ref="changePasswordDialogRef" />
  </el-header>
</template>

<script setup>
import { ref, computed, onMounted, inject } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Setting, User } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import ThemeToggle from '../../components/ThemeToggle.vue'
import ChangePasswordDialog from '../../components/ChangePasswordDialog.vue'
import request from '../../utils/request'

const router = useRouter()
const route = useRoute()
const siteTitle = inject('siteTitle', ref('OpsPortal运维导航'))

const logoUrl = ref('')
const changePasswordDialogRef = ref(null)

const isLoggedIn = computed(() => !!localStorage.getItem('token'))
const username = computed(() => localStorage.getItem('user') || '')
const isAdminRoute = computed(() => route.path === '/admin')
const showUserInfo = computed(() => isLoggedIn.value && isAdminRoute.value)

const fetchLogo = async () => {
  try {
    const res = await request.get('/api/logo')
    logoUrl.value = res.data.url
  } catch (error) {
    console.error('Error fetching logo:', error)
  }
}

const handleLogoError = () => {
  logoUrl.value = ''
}

const goHome = () => {
  router.push('/monitor')
}

const goToAdmin = () => {
  if (!isLoggedIn.value) {
    router.push('/login')
  } else {
    router.push('/admin')
  }
}

const handleCommand = async (command) => {
    if (command === 'logout') {
        try {
            await ElMessageBox.confirm('确定要退出登录吗？', '提示', { type: 'warning' })
            localStorage.removeItem('token')
            localStorage.removeItem('user')
            ElMessage.success('已退出登录')
            setTimeout(() => {
                router.push('/monitor')
                window.location.reload()
            }, 1000)
        } catch {}
    } else if (command === 'changePassword') {
        changePasswordDialogRef.value?.show()
    }
}

onMounted(() => {
    fetchLogo()
})
</script>

<style scoped>
.app-header {
  background-color: var(--bg-primary, #fff);
  border-bottom: 1px solid var(--border-color, #e2e8f0);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  height: 64px;
  position: relative;
  z-index: 1000;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.site-title-btn {
  font-size: 18px;
  font-weight: 700;
  color: var(--primary, #615ced);
  background: none;
  border: none;
  cursor: pointer;
  letter-spacing: -0.5px;
  padding: 0;
  white-space: nowrap;
  outline: none;
  box-shadow: none;
}

.site-title-btn:hover {
  opacity: 0.85;
}

.site-title-btn:focus {
  outline: none;
  box-shadow: none;
}

.header-logo {
  height: 100%;
  padding-right: 20px;
  display: flex;
  align-items: center;
}

.site-logo {
  max-height: 40px;
  max-width: 180px;
  object-fit: contain;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-left: auto;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  color: var(--text-primary, #1e293b);
}

.admin-btn:hover {
  transform: rotate(90deg);
  transition: transform 0.3s;
}
</style>
