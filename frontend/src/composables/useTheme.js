import { ref, watch } from 'vue'

const THEME_KEY = 'opsportal-theme'

// 全局共享的响应式主题状态
const theme = ref('light')

// 应用主题到DOM
const applyTheme = () => {
    document.documentElement.setAttribute('data-theme', theme.value)
    if (theme.value === 'dark') {
        document.documentElement.classList.add('dark')
    } else {
        document.documentElement.classList.remove('dark')
    }
}

// 初始化主题
const initTheme = () => {
    const saved = localStorage.getItem(THEME_KEY)
    if (saved) {
        theme.value = saved
    } else {
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
        theme.value = prefersDark ? 'dark' : 'light'
    }
    applyTheme()
}

// 切换主题
const toggleTheme = () => {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
    localStorage.setItem(THEME_KEY, theme.value)
    applyTheme()
}

// 监听主题变化
watch(theme, () => {
    applyTheme()
})

// 导出composable
export function useTheme() {
    if (!document.documentElement.hasAttribute('data-theme')) {
        initTheme()
    }
    return {
        theme,
        toggleTheme,
        initTheme
    }
}
