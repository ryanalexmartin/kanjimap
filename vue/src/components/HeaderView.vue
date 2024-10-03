<template>
  <header>
    <nav>
      <router-link to="/">Home</router-link> |
      <router-link to="/login" v-if="!isLoggedIn">Login</router-link>
      <a href="#" @click.prevent="logout" v-else>Logout</a>
    </nav>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useCharacterStore } from '@/store'

const store = useCharacterStore()
const router = useRouter()

const isLoggedIn = computed(() => store.isLoggedIn)

function logout() {
  store.isLoggedIn = false
  store.username = ''
  localStorage.removeItem('token')
  localStorage.removeItem('username')
  router.push('/login')
}
</script>

<style scoped>
header {
  padding: 1rem;
  background-color: #f0f0f0;
}

nav a {
  margin-right: 1rem;
  text-decoration: none;
  color: #333;
}
</style>