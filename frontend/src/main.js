import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import { useAuth } from './composables/useAuth.js'
import Home from './views/Home.vue'
import Login from './views/Login.vue'
import Workouts from './views/Workouts.vue'

const routes = [
  { path: '/', component: Home },
  { path: '/login', component: Login },
  { path: '/workouts', component: Workouts },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Redirect to /login if the user isn't authenticated.
router.beforeEach((to) => {
  const { userId } = useAuth()
  if (to.path !== '/login' && userId.value === null) {
    return '/login'
  }
})

const app = createApp(App)
app.use(router)
app.mount('#app')
