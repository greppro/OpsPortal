<template>
  <div v-if="showNotice && notice.content" class="notice-bar">
    <div class="notice-content">
      <div class="scroll-wrapper" :class="{ 'scroll': shouldScroll, 'centered': !shouldScroll }">
        <span>{{ notice.content }}</span>
        <span class="close-wrapper">
          <el-icon class="close-icon" @click="closeNotice">
            <Close />
          </el-icon>
        </span>
      </div>
    </div>
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

const handleAnimationEnd = () => {
  const wrapper = document.querySelector('.scroll-wrapper')
  if (wrapper) {
    const rect = wrapper.getBoundingClientRect()
    if (rect.right <= 0) {
      showNotice.value = false
      emit('notice-visible-change', false)
    }
  }
}

onMounted(() => {
  getNotice()
  window.addEventListener('resize', checkOverflow)
  emit('notice-visible-change', true)
  
  const wrapper = document.querySelector('.scroll-wrapper')
  if (wrapper) {
    wrapper.addEventListener('animationiteration', handleAnimationEnd)
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', checkOverflow)
  const wrapper = document.querySelector('.scroll-wrapper')
  if (wrapper) {
    wrapper.removeEventListener('animationiteration', handleAnimationEnd)
  }
})
</script>

<style scoped>
.notice-bar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 40px;
  background: linear-gradient(to right, rgb(147, 38, 185), rgb(218, 117, 185));
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  z-index: 1000;
  color: #fff;
  font-size: 14px;
}

.notice-content {
  flex: 1;
  position: relative;
  overflow: hidden;
  height: 100%;
  margin: 0 auto;
  max-width: 1200px;
  padding: 0 220px;
}

.scroll-wrapper {
  position: absolute;
  white-space: nowrap;
  display: inline-flex;
  align-items: center;
  height: 100%;
  line-height: 40px;
  padding-right: 0;
  left: 220px;
  width: calc(100% - 440px);
}

.close-wrapper {
  display: inline-flex;
  align-items: center;
  margin-left: 20px;
  padding-left: 20px;
  border-left: 2px solid rgba(255, 255, 255, 0.4);
  height: 24px;
}

.close-icon {
  cursor: pointer;
  color: rgba(255, 255, 255, 0.9);
  font-size: 20px;
  font-weight: bold;
  padding: 4px 8px;
  transition: all 0.3s ease;
}

.close-icon:hover {
  color: #ffffff;
  transform: scale(1.1);
}

.scroll {
  justify-content: flex-start;
  left: 220px;
  transform: none;
  animation: scroll 25s linear infinite;
}

.scroll:hover {
  animation-play-state: paused;
}

.centered {
  justify-content: center;
  left: 50%;
  transform: translateX(-50%);
}

@keyframes scroll {
  0% {
    transform: translateX(0);
  }
  100% {
    transform: translateX(calc(-100% - 40px));
  }
}
</style>