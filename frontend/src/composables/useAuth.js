import { ref } from 'vue'

// Module-level refs so state is shared across all callers.
const userId = ref(localStorage.getItem('userId') ? Number(localStorage.getItem('userId')) : null)
const userName = ref(localStorage.getItem('userName') || null)

export function useAuth() {
  function login(id, name) {
    userId.value = id
    userName.value = name
    localStorage.setItem('userId', String(id))
    localStorage.setItem('userName', name)
  }

  function logout() {
    userId.value = null
    userName.value = null
    localStorage.removeItem('userId')
    localStorage.removeItem('userName')
  }

  return { userId, userName, login, logout }
}
