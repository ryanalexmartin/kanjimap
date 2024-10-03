<template>
  <div v-if="store.isLoggedIn" class="learn-from-passage">
    Paste Chinese text here to learn characters automatically:
    <br />
    將中文文本粘貼到此處以自動學習字符：
    <textarea v-model="passage"></textarea>
    <button @click="learnFromPassage">Learn from Passage 從文本學習</button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useCharacterStore } from '../store'

const store = useCharacterStore()
const passage = ref('')

async function learnFromPassage() {
  const characters = passage.value.match(/[\u4e00-\u9fa5]/g)
  if (characters) {
    for (const char of characters) {
      await store.toggleCharacterLearned({
          char,
          pinyin: '',
          learned: true,
          id: 0
      })
    }
  }
  passage.value = ''
}
</script>

<style scoped>
.learn-from-passage {
  margin-top: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.learn-from-passage textarea {
  width: 300px;
  height: 200px;
  font-size: 12px;
  padding: 5px;
  margin-right: 10px;
  border-radius: 5px;
  resize: none;
}
</style>