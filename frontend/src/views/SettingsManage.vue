<template>
  <div class="settings-container">
    <el-card class="settings-card">
      <template #header>
        <span>基础设置</span>
      </template>
      <el-form label-width="100px" :model="form">
        <el-form-item label="站点标题">
          <el-input v-model="form.site_title" placeholder="如：OpsPortal运维导航" style="max-width: 400px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, inject } from 'vue'
import { ElMessage } from 'element-plus'
import request from '../utils/request'

const form = ref({ site_title: 'OpsPortal运维导航' })
const saving = ref(false)
const refreshSiteTitle = inject('refreshSiteTitle', null)

const fetchSettings = async () => {
  try {
    const res = await request.get('/api/config/site-title')
    form.value.site_title = res.data?.site_title ?? 'OpsPortal运维导航'
  } catch (e) {
    form.value.site_title = 'OpsPortal运维导航'
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    await request.put('/api/config/site-title', { site_title: form.value.site_title })
    ElMessage.success('保存成功')
    if (typeof refreshSiteTitle === 'function') {
      await refreshSiteTitle()
    }
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  fetchSettings()
})
</script>

<style scoped>
.settings-container {
  padding: 20px;
}
.settings-card {
  max-width: 600px;
}
</style>
