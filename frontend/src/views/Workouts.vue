<script setup>
import useSWRV from 'swrv'
import { ref } from 'vue'
import { api } from '@/api/client'
import { useAuth } from '@/composables/useAuth.js'

const { userId } = useAuth()

const fetcher = () => api.GET('/api/v1/workouts').then(({ data }) => data)
const { data: workouts, isValidating: loading, mutate } = useSWRV('/api/v1/workouts', fetcher)

const newWorkout = ref({ name: '', description: '', duration_minutes: 0 })
const error = ref('')

const createWorkout = async () => {
  error.value = ''
  const { error: apiErr } = await api.POST('/api/v1/workouts', {
    body: { ...newWorkout.value, user_id: userId.value },
  })
  if (apiErr) {
    error.value = apiErr.detail ?? 'Failed to create workout.'
    return
  }
  newWorkout.value = { name: '', description: '', duration_minutes: 0 }
  mutate()
}
</script>

<template>
  <div class="workouts">
    <h1>My Workouts</h1>

    <div class="create-workout">
      <h2>Create New Workout</h2>
      <form @submit.prevent="createWorkout">
        <input v-model="newWorkout.name" placeholder="Workout name" required />
        <textarea v-model="newWorkout.description" placeholder="Description" />
        <input
          v-model.number="newWorkout.duration_minutes"
          type="number"
          placeholder="Duration (minutes)"
          min="0"
        />
        <p v-if="error" class="error">{{ error }}</p>
        <button type="submit">Create Workout</button>
      </form>
    </div>

    <div class="workout-list">
      <h2>Your Workouts</h2>
      <div v-if="loading">Loading workouts...</div>
      <div v-else-if="!workouts?.length" class="empty">
        No workouts yet. Create your first workout above!
      </div>
      <div v-else class="workouts-grid">
        <div v-for="workout in workouts" :key="workout.id" class="workout-card">
          <h3>{{ workout.name }}</h3>
          <p v-if="workout.description">{{ workout.description }}</p>
          <p class="duration">{{ workout.duration_minutes }} min</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.workouts {
  max-width: 800px;
  margin: 0 auto;
}

h1 {
  color: #2c3e50;
  margin-bottom: 2rem;
}

h2 {
  color: #2c3e50;
  margin-bottom: 1rem;
  font-size: 1.5rem;
}

.create-workout {
  background: #f5f5f5;
  padding: 2rem;
  border-radius: 8px;
  margin-bottom: 3rem;
}

form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

input,
textarea {
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  font-family: inherit;
}

textarea {
  min-height: 100px;
  resize: vertical;
}

button {
  background: #42b983;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background 0.3s;
}

button:hover {
  background: #359268;
}

.error {
  color: #e74c3c;
  margin: 0;
}

.workout-list {
  margin-top: 2rem;
}

.empty {
  text-align: center;
  color: #666;
  padding: 2rem;
  background: #f5f5f5;
  border-radius: 8px;
}

.workouts-grid {
  display: grid;
  gap: 1rem;
}

.workout-card {
  background: white;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 1.5rem;
  transition: box-shadow 0.3s;
}

.workout-card:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.workout-card h3 {
  color: #2c3e50;
  margin-bottom: 0.5rem;
}

.workout-card p {
  color: #666;
}

.duration {
  margin-top: 0.5rem;
  font-size: 0.9rem;
  color: #42b983;
  font-weight: 500;
}
</style>
