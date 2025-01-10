<template>
    <div v-if="showNotice && notice.content" class="notice-bar">
      <div class="notice-content" :class="{ 'scroll': shouldScroll }">
        <span>{{ notice.content }}</span>
      </div>
      <el-icon class="close-icon" @click="closeNotice">
        <Close />
      </el-icon>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted, onUnmounted, nextTick } from 'vue'
  import { Close } from '@element-plus/icons-vue'
  import { getActiveNotice } from '../api'
  
  const emit = defineEmits(['notice-visible-change'])
  const notice = ref({ content: '' })
  const showNotice = ref(true)
  const shouldScroll = ref(false)
  
  const getNotice = async () => {
    try {
      const res = await getActiveNotice()
      notice.value = res.data
      if (notice.value.content) {
        checkOverflow()
      }
    } catch (error) {
      console.error('获取公告失败:', error)
    }
  }
  
  const checkOverflow = async () => {
    await nextTick()
    const content = document.querySelector('.notice-content')
    if (content) {
      shouldScroll.value = content.scrollWidth > content.clientWidth
    }
  }
  
  const closeNotice = () => {
    showNotice.value = false
    emit('notice-visible-change', false)
  }
  
  onMounted(() => {
    getNotice()
    window.addEventListener('resize', checkOverflow)
    emit('notice-visible-change', true)
  })
  
  onUnmounted(() => {
    window.removeEventListener('resize', checkOverflow)
  })
  </script>
  
  <style scoped>
  .notice-bar {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    height: 40px;
    background: linear-gradient(to right, #9c27b0, #673ab7);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0 20px 0 220px;
    z-index: 1000;
    color: #fff;
  }
  
  .notice-content {
    flex: 1;
    overflow: hidden;
    white-space: nowrap;
    margin-right: 20px;
    text-align: center;
  }
  
  .scroll {
    animation: scroll-left 30s linear infinite;
  }
  
  .scroll:hover {
    animation-play-state: paused;
  }
  
  @keyframes scroll-left {
    0% {
      transform: translateX(100%);
    }
    100% {
      transform: translateX(-100%);
    }
  }
  
  .close-icon {
    cursor: pointer;
    color: rgba(255, 255, 255, 0.8);
  }
  
  .close-icon:hover {
    color: #ffffff;
  }
  </style>