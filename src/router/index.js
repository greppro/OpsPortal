import { createRouter, createWebHistory } from 'vue-router'
import ToolsView from '../views/ToolsView.vue'
import ManagementView from '../views/ManagementView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'tools',
      component: ToolsView
    },
    {
      path: '/management',
      name: 'management',
      component: ManagementView
    }
  ]
})

export default router
