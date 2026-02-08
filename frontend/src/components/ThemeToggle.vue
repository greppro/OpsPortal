<template>
  <div 
    class="theme-toggle-btn"
    @click="toggleTheme"
    :title="currentTheme === 'light' ? '切换到夜间模式' : '切换到日间模式'"
  >
    <el-icon size="20">
      <Sunny v-if="currentTheme === 'light'" />
      <Moon v-else />
    </el-icon>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { Sunny, Moon } from '@element-plus/icons-vue'
import { useTheme } from '../composables/useTheme'

const { theme, toggleTheme } = useTheme()
const currentTheme = computed(() => theme.value)
</script>

<style scoped>
.theme-toggle-btn {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s;
  color: var(--text-primary, #606266);
  border: 1px solid var(--border-color, #dcdfe6);
  background-color: transparent;
}

.theme-toggle-btn:hover {
  background-color: var(--primary-light, #e0f2fe);
  color: var(--primary, #0ea5e9);
  border-color: var(--primary);
  transform: rotate(20deg);
}
</style>
