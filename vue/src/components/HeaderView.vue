<template>
  <header>
    <h1>KanjiMap</h1>
    <h2>漢字地圖</h2>
    <h2>by Ryan Alex Martin</h2>
    <h2>馬丁瑞安創造</h2>
    <a href="https://ryanalexmartin.com">ryanalexmartin.com</a>
    <br />
    <a href="https://github.com/ryanalexmartin/kanjimap">github.com/ryanalexmartin/kanjimap</a>
    <br />
    <br />
    <div>
      Click on a character to mark it as learned. Click again to mark it as
      unlearned. I recommend something like
      <a href="https://github.com/gkovacs/LiuChanFirefox/tree/firefox">LiuChan</a>
      to help you learn the characters.
      <br />
      點擊一個字符標記為已學。再次點擊標記為未學。
    </div>
    <PassageLearner />
    <nav>
      <a href="#" @click.prevent="logout" v-if="store.isLoggedIn">Logout</a>
      <button @click="logout" v-if="store.isLoggedIn" class="logout-button">Logout</button>
    </nav>
    <p v-if="errorMessage" class="error-message">{{ errorMessage }}</p>
  </header>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useCharacterStore } from '../store'
import PassageLearner from './PassageLearner.vue'
import axios from 'axios'

const router = useRouter()
const store = useCharacterStore()
const errorMessage = ref('')

async function logout() {
  try {
    // Assuming you have a logout API endpoint
    await axios.post('/api/logout')
    store.isLoggedIn = false
    store.username = ''
    router.push('/login')
    errorMessage.value = ''
  } catch (error) {
    console.error('Logout failed:', error)
    errorMessage.value = 'Logout failed. Please try again.'
  }
}
</script>

<style scoped>
.error-message {
  color: red;
  margin-top: 10px;
}
</style>