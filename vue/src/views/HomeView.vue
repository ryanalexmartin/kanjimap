<template>
  <div class="home">
    <h1>Learned {{ store.learnedCount }} out of {{ store.totalCharacters }}</h1>
    <h1>已學 {{ store.learnedCount }} / {{ store.totalCharacters }}</h1>
    
    <button @click="toggleKnownCharactersList">
      {{ showKnownCharactersList ? "Hide" : "Show" }} Known Characters
    </button>
    
    <KnownCharactersList v-if="showKnownCharactersList" :characters="store.characters" />

    <div class="sorting-controls">
      <label for="sort-select">Sort by: </label>
      <select id="sort-select" v-model="sortBy" @change="sortCharacters">
        <option value="default">Default</option>
        <option value="frequency">Frequency</option>
      </select>
    </div>

    <CharacterGrid :characters="sortedCharacters" @toggle-learned="store.toggleCharacterLearned" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useCharacterStore } from '../store'
import KnownCharactersList from '../components/KnownCharactersList.vue'
import CharacterGrid from '../components/CharacterGrid.vue'

const store = useCharacterStore()
const showKnownCharactersList = ref(false)
const sortBy = ref('default')

const toggleKnownCharactersList = () => {
  showKnownCharactersList.value = !showKnownCharactersList.value
}

const sortCharacters = () => {
  // This function will be called when the sorting option changes
  // The actual sorting is done in the computed property 'sortedCharacters'
}

const sortedCharacters = computed(() => {
  if (sortBy.value === 'frequency') {
    return [...store.characters].sort((a, b) => (b.frequency || 0) - (a.frequency || 0))
  } else {
    return store.characters
  }
})
</script>