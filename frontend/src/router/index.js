import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/monitor'
    },
    {
      path: '/monitor',
      name: 'monitor',
      component: () => import('../views/Monitor.vue')
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('../views/Admin.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/Login.vue')
    }
  ]
})

// 路由守卫
router.beforeEach((to, from, next) => {
  if (to.meta.requiresAuth && !localStorage.getItem('token')) {
    localStorage.setItem('redirectPath', to.fullPath)
    next('/login')
  } else {
    next()
  }
})

export default router
