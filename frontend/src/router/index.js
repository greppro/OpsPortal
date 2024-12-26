import { createRouter, createWebHistory } from 'vue-router'
import Monitor from '../views/Monitor.vue'
import Admin from '../views/Admin.vue'
import Login from '../views/Login.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/monitor'
    },
    {
      path: '/monitor',
      name: 'Monitor',
      component: Monitor
    },
    {
      path: '/admin',
      name: 'Admin',
      component: Admin,
      meta: { requiresAuth: true }
    },
    {
      path: '/login',
      name: 'Login',
      component: Login
    }
  ]
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const isLoggedIn = !!localStorage.getItem('token')
  
  if (to.meta.requiresAuth && !isLoggedIn) {
    // 保存尝试访问的路径
    localStorage.setItem('redirectPath', to.fullPath)
    next('/login')
  } else {
    next()
  }
})

export default router
