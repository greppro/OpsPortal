import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

const FAVORITES_KEY = 'opsportal-favorites'
const favoriteIds = ref(new Set())

export function useFavorites() {
  // 初始化加载
  const loadFavorites = () => {
    try {
      const stored = localStorage.getItem(FAVORITES_KEY)
      if (stored) {
        favoriteIds.value = new Set(JSON.parse(stored))
      }
    } catch (e) {
      console.error('Failed to load favorites', e)
    }
  }

  // 切换收藏状态
  const toggleFavorite = (toolId) => {
    if (favoriteIds.value.has(toolId)) {
      favoriteIds.value.delete(toolId)
      ElMessage.success('已取消收藏')
    } else {
      favoriteIds.value.add(toolId)
      ElMessage.success('已添加收藏')
    }
    saveFavorites()
  }

  // 检查是否收藏
  const isFavorite = (toolId) => {
    return favoriteIds.value.has(toolId)
  }

  // 保存到本地
  const saveFavorites = () => {
    localStorage.setItem(FAVORITES_KEY, JSON.stringify([...favoriteIds.value]))
  }

  onMounted(() => {
    loadFavorites()
  })

  return {
    favoriteIds,
    toggleFavorite,
    isFavorite,
    loadFavorites,
    saveFavorites
  }
}
