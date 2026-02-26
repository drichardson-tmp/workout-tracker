<script setup>
import axios from 'axios'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth.js'

const router = useRouter()
const { login } = useAuth()

const name = ref('')
const email = ref('')
const error = ref('')
const loading = ref(false)

const submit = async () => {
  loading.value = true
  error.value = ''
  try {
    // Look up existing user by email first.
    const { data: found } = await axios.get(`/api/users?email=${encodeURIComponent(email.value)}`)
    if (found.length > 0) {
      login(found[0].id, found[0].name)
      router.push('/workouts')
      return
    }
    // No account yet — create one.
    if (!name.value.trim()) {
      error.value = 'Name is required to create an account.'
      return
    }
    const { data: created } = await axios.post('/api/users', {
      email: email.value,
      name: name.value.trim(),
      password: 'dev-placeholder',
    })
    login(created.id, created.name)
    router.push('/workouts')
  } catch (e) {
    error.value = e.response?.data?.detail ?? 'Something went wrong.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login">
    <h1>Welcome</h1>
    <p>Enter your email to sign in or create an account.</p>

    <form @submit.prevent="submit">
      <label>
        Email
        <input v-model="email" type="email" placeholder="you@example.com" required autofocus />
      </label>
      <label>
        Name <span class="hint">(required for new accounts)</span>
        <input v-model="name" type="text" placeholder="Your name" />
      </label>
      <p v-if="error" class="error">{{ error }}</p>
      <button type="submit" :disabled="loading">
        {{ loading ? 'Signing in…' : 'Continue' }}
      </button>
    </form>
  </div>
</template>

<style scoped>
.login {
  max-width: 400px;
  margin: 4rem auto;
}

h1 {
  color: #2c3e50;
  margin-bottom: 0.5rem;
}

p {
  color: #666;
  margin-bottom: 2rem;
}

form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

label {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
  font-weight: 500;
  color: #2c3e50;
}

.hint {
  font-weight: 400;
  font-size: 0.85rem;
  color: #999;
}

input {
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  font-family: inherit;
}

button {
  background: #42b983;
  color: white;
  border: none;
  padding: 0.75rem;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background 0.2s;
}

button:hover:not(:disabled) {
  background: #359268;
}

button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.error {
  color: #e74c3c;
  margin: 0;
}
</style>
