<script setup>
import axios from 'axios'
import { onMounted, ref } from 'vue'

const health = ref(null)
const loading = ref(true)

onMounted(async () => {
  try {
    const response = await axios.get('/health')
    health.value = response.data
  } catch (error) {
    console.error('Failed to fetch health status:', error)
    health.value = { status: 'error' }
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="home">
    <h1>Workout Tracker</h1>
    <p>Welcome to your workout tracking application!</p>

    <div class="status">
      <h2>Backend Status</h2>
      <div v-if="loading">Loading...</div>
      <div v-else-if="health?.status === 'ok'" class="status-ok">
        ✓ Backend is running
      </div>
      <div v-else class="status-error">✗ Backend is not responding</div>
    </div>
  </div>
</template>

<style scoped>
.home {
  text-align: center;
}

h1 {
  color: #2c3e50;
  margin-bottom: 1rem;
  font-size: 2.5rem;
}

p {
  font-size: 1.2rem;
  color: #666;
  margin-bottom: 2rem;
}

.status {
  margin-top: 3rem;
  padding: 2rem;
  background: #f5f5f5;
  border-radius: 8px;
}

.status h2 {
  margin-bottom: 1rem;
  color: #2c3e50;
}

.status-ok {
  color: #42b983;
  font-weight: 600;
  font-size: 1.1rem;
}

.status-error {
  color: #e74c3c;
  font-weight: 600;
  font-size: 1.1rem;
}
</style>
