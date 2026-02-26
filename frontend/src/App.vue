<script setup>
import { useRouter } from 'vue-router'
import { useAuth } from './composables/useAuth.js'

const router = useRouter()
const { userName, userId, logout } = useAuth()

const handleLogout = () => {
  logout()
  router.push('/login')
}
</script>

<template>
  <div id="app">
    <nav>
      <div class="nav-links">
        <router-link to="/">Home</router-link>
        <router-link to="/workouts">Workouts</router-link>
      </div>
      <div v-if="userId" class="nav-user">
        <span>{{ userName }}</span>
        <button @click="handleLogout">Sign out</button>
      </div>
    </nav>
    <main>
      <RouterView />
    </main>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell,
    sans-serif;
  line-height: 1.6;
  color: #333;
}

#app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

nav {
  background: #2c3e50;
  padding: 1rem 2rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.nav-links {
  display: flex;
  gap: 2rem;
}

nav a {
  color: white;
  text-decoration: none;
  font-weight: 500;
  transition: color 0.3s;
}

nav a:hover,
nav a.router-link-active {
  color: #42b983;
}

.nav-user {
  display: flex;
  align-items: center;
  gap: 1rem;
  color: #ccc;
  font-size: 0.9rem;
}

.nav-user button {
  background: transparent;
  border: 1px solid #ccc;
  color: #ccc;
  padding: 0.25rem 0.75rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.85rem;
  transition: all 0.2s;
}

.nav-user button:hover {
  background: #ccc;
  color: #2c3e50;
}

main {
  flex: 1;
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}
</style>
