<template>
  <div class="register">
    <h2>Register</h2>
    <form @submit.prevent="register">
      <input v-model="username" type="text" placeholder="Username" required>
      <input v-model="password" type="password" placeholder="Password" required>
      <input v-model="confirmPassword" type="password" placeholder="Confirm Password" required>
      <button type="submit">Register</button>
    </form>
    <p>Already have an account? <router-link to="/login">Login</router-link></p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()

const username = ref('')
const password = ref('')
const confirmPassword = ref('')

async function register() {
  if (password.value !== confirmPassword.value) {
    alert("Passwords don't match")
    return
  }

  try {
    await axios.post(`${process.env.VUE_APP_API_URL}/register`, {
      username: username.value,
      password: password.value,
    })
    
    alert('Registration successful! Please login.')
    router.push('/login')
  } catch (error) {
    console.error('Registration failed:', error)
    // Handle registration error (e.g., show error message to user)
  }
}
</script>

<style scoped>
.register {
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