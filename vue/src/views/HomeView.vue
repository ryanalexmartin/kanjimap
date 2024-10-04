<template>
  <div class="home">
    <h1>Learned {{ learnedCount }} out of {{ filteredCharacters.length }}</h1>
    <h1>已學 {{ learnedCount }} / {{ filteredCharacters.length }}</h1>
    
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
import { ref, computed, onMounted, watch } from 'vue'
import { useCharacterStore } from '../store'
import KnownCharactersList from '../components/KnownCharactersList.vue'
import CharacterGrid from '../components/CharacterGrid.vue'

const store = useCharacterStore()
const showKnownCharactersList = ref(false)
const sortBy = ref('default')

const learnedCount = computed(() => filteredCharacters.value.filter(c => c.learned).length)

const toggleKnownCharactersList = () => {
  showKnownCharactersList.value = !showKnownCharactersList.value
}

const sortCharacters = () => {
  // This function will be called when the sorting option changes
  // The actual sorting is done in the computed property 'sortedCharacters'
}

const filteredCharacters = computed(() => {
  return store.characters.filter(character => character.id.startsWith('A'))
})

const sortedCharacters = computed(() => {
  const charactersToSort = filteredCharacters.value
  if (sortBy.value === 'frequency') {
    return [...charactersToSort].sort((a, b) => (b.frequency || 0) - (a.frequency || 0))
  } else {
    return charactersToSort
  }
})

onMounted(async () => {
  await store.loadCharacters()
  console.log('Characters loaded:', store.characters)
})

watch(() => store.characters, () => {
  console.log('Characters updated, new learned count:', learnedCount.value)
}, { deep: true })
</script>