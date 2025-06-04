import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LandingPage from '../views/LandingPage.vue'

const routes = [
  { path: '/', name: 'Home', component: HomeView },
  { path: '/landing', name: 'LandingPage', component: LandingPage, meta: { requiresAuth: true } },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

// Decode JWT and check expiry
function isTokenExpired(token) {
  try {
    const payloadBase64 = token.split('.')[1]
    const payloadJson = atob(payloadBase64)
    const payload = JSON.parse(payloadJson)

    const now = Math.floor(Date.now() / 1000) // current time in seconds
    return payload.exp && payload.exp < now
  } catch (err) {
    return true // treat malformed token as expired
  }
}

// Global route guard
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('jwtToken')
  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth)

  if (requiresAuth) {
    if (!token || isTokenExpired(token)) {
      localStorage.removeItem('jwtToken') // Clean up
      return next({ path: '/', query: { sessionExpired: '1' } })
    }
  }

  next()
})

export default router
