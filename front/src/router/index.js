import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LandingPage from '../views/LandingPage.vue'

const routes = [
  { path: '/', name: 'Home', component: HomeView },
  { path: '/landing', name: 'LandingPage', component: LandingPage },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

export default router
