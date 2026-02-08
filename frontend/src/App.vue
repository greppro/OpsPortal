<template>
  <el-config-provider>
    <router-view></router-view>
  </el-config-provider>
</template>

<script setup>
import { onMounted } from 'vue'
import { useTheme } from './composables/useTheme'

const { initTheme } = useTheme()

onMounted(() => {
  initTheme()
})
</script>

<style>
/* 全局样式 */
html, body {
  margin: 0;
  padding: 0;
  height: 100vh;
  overflow: hidden;
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
}

#app {
  height: 100vh;
}
</style>
