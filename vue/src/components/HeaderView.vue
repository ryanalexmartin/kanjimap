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
      <button @click="logout" v-if="isLoggedIn" class="logout-button">Logout</button>
    </nav>
    <p v-if="errorMessage" class="error-message">{{ errorMessage }}</p>
    {{ console.log('store.isLoggedIn:', store.isLoggedIn) }}
  </header>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useCharacterStore } from '../store'
import PassageLearner from './PassageLearner.vue'
import { computed } from 'vue'
import { watch } from 'vue'
import { onMounted } from 'vue'

const router = useRouter()
const store = useCharacterStore()
const errorMessage = ref('')

const isLoggedIn = computed(() => store.isLoggedIn)

async function logout() {
  try {
    // await axios.post('/logout') // TODO: implement logout on backend
    store.setIsLoggedIn(false)
    store.setUsername('')
    router.push('/login')
    errorMessage.value = ''
    localStorage.removeItem('token')
  } catch (error) {
    console.error('Logout failed:', error)
    errorMessage.value = 'Logout failed. Please try again.'
  }
}

watch(() => store.isLoggedIn, (newValue) => {
  console.log('isLoggedIn changed:', newValue)
})

onMounted(() => {
  console.log('HeaderView mounted, store.isLoggedIn:', store.isLoggedIn)
})
</script>

<style scoped>
.error-message {
  color: red;
  margin-top: 10px;
}

.logout-button {
  padding: 8px 16px;
  background-color: #f44336;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.logout-button:hover {
  background-color: #d32f2f;
}
</style>