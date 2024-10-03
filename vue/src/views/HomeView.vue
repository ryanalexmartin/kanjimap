<template>
  <div class="home">
    <h1>Learned {{ learnedCount }} out of {{ totalCharacters }}</h1>
    <CharacterGrid :characters="characters" @toggle-learned="toggleCharacterLearned" />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useCharacterStore } from '@/store'
import CharacterGrid from '@/components/CharacterGrid.vue'

const store = useCharacterStore()

onMounted(async () => {
  await store.loadCharacters()
})

const { characters, learnedCount, totalCharacters, toggleCharacterLearned } = store
</script>