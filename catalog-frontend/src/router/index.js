import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { auth } from '../services/api'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: { title: 'Каталог товаров' }
    },
    {
      path: '/product/:id',
      name: 'product',
      component: () => import('../views/ProductDetail.vue'),
      meta: { title: 'Товар' }
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('../views/AboutView.vue'),
      meta: { title: 'О магазине' }
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('../views/AdminView.vue'),
      meta: { 
        title: 'Админ панель',
        requiresAuth: true 
      }
    },
    {
      path: '/admin/login',
      name: 'login',
      component: () => import('../views/AdminLogin.vue'),
      meta: { title: 'Вход в админку' }
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/'
    }
  ]
})

// Навигационный гард для админки
router.beforeEach((to, from, next) => {
  // Устанавливаем заголовок страницы
  if (to.meta.title) {
    document.title = to.meta.title;
  }

  // Проверка аутентификации для админки
  if (to.meta.requiresAuth && !auth.isAuthenticated()) {
    next('/admin/login');
  } else {
    next();
  }
})

export default router