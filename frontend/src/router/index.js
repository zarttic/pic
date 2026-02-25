import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/gallery',
    name: 'Gallery',
    component: () => import('../views/Gallery.vue')
  },
  {
    path: '/albums',
    name: 'Albums',
    component: () => import('../views/Albums.vue')
  },
  {
    path: '/albums/:id',
    name: 'AlbumDetail',
    component: () => import('../views/AlbumDetail.vue')
  },
  {
    path: '/about',
    name: 'About',
    component: () => import('../views/About.vue')
  },
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: () => import('../views/admin/Login.vue')
  },
  {
    path: '/admin',
    component: () => import('../views/admin/AdminLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        redirect: '/admin/photos'
      },
      {
        path: 'photos',
        name: 'AdminPhotos',
        component: () => import('../views/admin/PhotoManagement.vue')
      },
      {
        path: 'albums',
        name: 'AdminAlbums',
        component: () => import('../views/admin/AlbumManagement.vue')
      },
      {
        path: 'albums/:id/photos',
        name: 'AdminAlbumPhotos',
        component: () => import('../views/admin/AlbumPhotoManagement.vue')
      },
      {
        path: 'statistics',
        name: 'AdminStatistics',
        component: () => import('../views/admin/Statistics.vue')
      },
      {
        path: 'settings',
        name: 'AdminSettings',
        component: () => import('../views/admin/Settings.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.checkAuth()) {
    // 需要认证但未登录，跳转到登录页
    next({
      path: '/admin/login',
      query: { redirect: to.fullPath }
    })
  } else if (to.path === '/admin/login' && authStore.checkAuth()) {
    // 已登录访问登录页，跳转到管理后台
    next('/admin')
  } else {
    next()
  }
})

export default router
