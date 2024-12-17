<template>
  <div class="tools-container">
    <el-tabs v-model="activeEnv" @tab-click="handleEnvChange">
      <el-tab-pane label="开发环境" name="dev">
        <tool-grid :tools="tools" />
      </el-tab-pane>
      <el-tab-pane label="生产环境" name="prod">
        <tool-grid :tools="tools" />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'
import ToolGrid from '../components/ToolGrid.vue'

const activeEnv = ref('dev')
const tools = ref([])

const fetchTools = async (env) => {
  try {
    console.log('Fetching tools for environment:', env)
    const response = await axios.get(`http://localhost:8080/api/tools?environment=${env}`)
    console.log('Response:', response.data)
    tools.value = response.data
  } catch (error) {
    console.error('Error details:', error.response || error)
    ElMessage.error(`获取工具列表失败: ${error.message}`)
  }
}

const handleEnvChange = () => {
  fetchTools(activeEnv.value)
}

watch(activeEnv, (newEnv) => {
  fetchTools(newEnv)
})

onMounted(() => {
  fetchTools(activeEnv.value)
})
</script>
