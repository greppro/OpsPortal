<template>
  <el-container class="layout-container">
    <el-aside width="200px" class="aside">
      <div class="logo">OpsPortal运维导航</div>
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
    
    <el-container>
      <el-header class="header" v-if="showUserInfo">
        <div class="header-right">
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
      
      <el-main>
        <router-view></router-view>
      </el-main>
    </el-container>
    <ChangePasswordDialog ref="changePasswordDialogRef" />
  </el-container>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Monitor, Setting, User } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import ChangePasswordDialog from './components/ChangePasswordDialog.vue'

const router = useRouter()
const route = useRoute()

const isLoggedIn = computed(() => !!localStorage.getItem('token'))
const username = computed(() => localStorage.getItem('user') || '')
const activeMenu = computed(() => route.path)

// 只在管理页面显示用户信息
const showUserInfo = computed(() => isLoggedIn.value && route.path === '/admin')

const changePasswordDialogRef = ref(null)

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
  }
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.aside {
  background-color: #304156;
  color: #fff;
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
  justify-content: flex-end;
  padding: 0 20px;
}

.header-right {
  display: flex;
  align-items: center;
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
</style>
