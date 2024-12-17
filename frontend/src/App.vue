<template>
  <el-container class="layout-container">
    <el-aside width="200px">
      <el-menu
        :router="true"
        class="el-menu-vertical"
      >
        <el-menu-item index="/">
          <el-icon><Monitor /></el-icon>
          <span>工具导航</span>
        </el-menu-item>
        <el-menu-item index="/management">
          <el-icon><Setting /></el-icon>
          <span>工具管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    
    <el-container>
      <el-header>
        <h2>OpsPortal 运维导航平台</h2>
      </el-header>
      <el-main>
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { Monitor, Setting } from '@element-plus/icons-vue'
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.el-header {
  background-color: #409EFF;
  color: white;
  line-height: 60px;
}

.el-aside {
  background-color: #545c64;
}
</style>
