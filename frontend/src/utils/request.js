import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '../router'

const request = axios.create({
    baseURL: import.meta.env.VITE_API_URL,
    timeout: 5000
})

// 请求拦截器
request.interceptors.request.use(
    config => {
        const token = localStorage.getItem('token')
        if (token) {
            config.headers['Authorization'] = `Bearer ${token}`
        }
        return config
    },
    error => {
        return Promise.reject(error)
    }
)

// 响应拦截器
request.interceptors.response.use(
    response => response,
    error => {
        const status = error.response?.status
        const data = error.response?.data
        const msg = data?.message ?? data?.error

        // 401 且为 token 无效/过期：清除 token 并跳转登录
        if (status === 401 && (data?.error === 'token无效或已过期' || data?.error === '未授权访问')) {
            localStorage.removeItem('token')
            ElMessage.error(msg || '登录已过期，请重新登录')
            router.push('/login')
            return Promise.reject(error)
        }

        ElMessage.error(msg || '请求失败')
        return Promise.reject(error)
    }
)

export default request 