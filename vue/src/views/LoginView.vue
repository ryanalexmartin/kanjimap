<template>
  <div class="login">
    <h2>Login</h2>
    <form @submit.prevent="login">
      <input v-model="username" type="text" placeholder="Username" required>
      <input v-model="password" type="password" placeholder="Password" required>
      <button type="submit">Login</button>
    </form>
    <p>Don't have an account? <router-link to="/register">Register</router-link></p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useCharacterStore } from '@/store'
import axios from 'axios'


const router = useRouter()
const store = useCharacterStore()

const username = ref('')
const password = ref('')

async function login() {
  console.log('Logging in user:', username.value);
  try {
    const apiUrl = import.meta.env.VITE_API_BASE_URL
    if (!apiUrl) {
      throw new Error('VITE_API_BASE_URL is not defined')
    }
    const response = await axios.post(`${apiUrl}/login`, {
      username: username.value,
      password: password.value,
    })
    const { token } = response.data
    localStorage.setItem('token', token)
    localStorage.setItem('username', username.value)
    
    store.isLoggedIn = true
    store.username = username.value
    await store.loadCharacters()
    
    router.push('/')
  } catch (error) {
    console.error('Login failed:', error)
    // Handle login error (e.g., show error message to user)
  }
}
</script>

<style scoped>
.login {
  max-width: 300px;
  margin: 0 auto;
  padding: 2rem;
}

form {
  display: flex;
  flex-direction: column;
}

input {
  margin-bottom: 1rem;
  padding: 0.5rem;
}

button {
  padding: 0.5rem;
  background-color: #4CAF50;
  color: white;
  border: none;
  cursor: pointer;
}

button:hover {
  background-color: #45a049;
}
</style>