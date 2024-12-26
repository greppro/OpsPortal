<template>
  <el-select 
    v-model="selectedEnv" 
    :placeholder="placeholder"
    :disabled="!projectId"
    clearable
    @clear="handleClear"
    @change="handleChange"
  >
    <el-option
      v-for="env in environments"
      :key="env.id"
      :label="env.label"
      :value="env.name"
    />
  </el-select>
</template>

<script setup>
import { ref, watch } from 'vue'
import request from '../utils/request'

const props = defineProps({
  projectId: {
    type: [String, Number],
    default: ''
  },
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: '选择环境'
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

const environments = ref([])
const selectedEnv = ref(props.modelValue)

// 监听项目变化
watch(() => props.projectId, async (newVal) => {
  if (newVal) {
    await fetchEnvironments()
  } else {
    environments.value = []
    selectedEnv.value = ''
    emit('update:modelValue', '')
  }
}, { immediate: true })

// 监听值变化
watch(() => props.modelValue, (newVal) => {
  selectedEnv.value = newVal
})

// 获取环境列表
const fetchEnvironments = async () => {
  try {
    const response = await request.get('/api/environments', {
      params: { project_id: props.projectId }
    })
    environments.value = response.data
  } catch (error) {
    console.error('Fetch environments error:', error)
  }
}

const handleChange = (value) => {
  emit('update:modelValue', value)
  emit('change', value)
}

const handleClear = () => {
  emit('update:modelValue', '')
  emit('change', '')
}
</script> 