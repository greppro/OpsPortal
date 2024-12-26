import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '../router'

const request = axios.create({
  baseURL: 'http://localhost:8080',
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
  response => {
    return response
  },
  error => {
    if (error.response?.status === 401) {
      // 清除登录状态
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      
      ElMessage({
        type: 'warning',
        message: '请先登录',
        duration: 2000
      })

      // 保存当前路径
      localStorage.setItem('redirectPath', router.currentRoute.value.fullPath)
      
      // 跳转到登录页
      router.push('/login')
    } else {
      // 其他错误显示错误信息
      ElMessage({
        type: 'error',
        message: error.response?.data?.error || '请求失败',
        duration: 3000
      })
    }
    return Promise.reject(error)
  }
)

export default request 