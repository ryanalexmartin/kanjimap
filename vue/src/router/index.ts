import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'

const routes = [
  { path: '/', component: HomeView },
  { path: '/login', component: LoginView },
  { path: '/register', component: RegisterView },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Add navigation guard for debugging
router.beforeEach((to, from, next) => {
  console.log('Navigating to:', to.path)
  next()
})

// Add navigation guard for handling redirection
router.beforeEach((to, from, next) => {
  const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true' && !!localStorage.getItem('token')
  if (to.path === '/' && !isLoggedIn) {
    next('/login')
  } else {
    next()
  }
})

export default router