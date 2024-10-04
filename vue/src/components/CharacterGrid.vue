<template>
    <div>
      <div class="sorting-controls">
        <label for="sort-select">Sort by: </label>
        <select id="sort-select" v-model="sortBy" @change="sortCharacters">
          <option value="default">Default</option>
          <option value="frequency">Frequency</option>
        </select>
      </div>
  
      <div class="big-table">
        <CharacterCard
          v-for="character in sortedCharacters"
          :key="character.id"
          :character="character"
          @toggle-learned="toggleCharacterLearned"
        />
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, computed, onMounted, watch } from 'vue'
  import { Character } from '../types/types'
  import CharacterCard from '@/components/CharacterCard.vue'
  import { useCharacterStore } from '@/store'
  
  const props = defineProps<{
    characters: Character[]
  }>()
  
  const store = useCharacterStore()
  
  const sortBy = ref('default')
  
  const sortCharacters = () => {
    // This function will be called when the sorting option changes
    // The actual sorting is done in the computed property 'sortedCharacters'
  }
  
  const sortedCharacters = computed(() => {
    console.log('Characters in sortedCharacters:', props.characters)
    if (sortBy.value === 'frequency') {
      return [...props.characters].sort((a, b) => (b.frequency || 0) - (a.frequency || 0))
    } else {
      // Default sorting (by id)
      return [...props.characters].sort((a, b) => b.id - a.id)
    }
  })
  
  const toggleCharacterLearned = async (character: Character) => {
    await store.toggleCharacterLearned(character)
  }

  onMounted(() => {
    console.log('CharacterGrid mounted. Characters:', props.characters)
  })

  watch(() => props.characters, (newCharacters) => {
    console.log('Characters updated:', newCharacters)
  }, { deep: true })
  </script>



  <style scoped>
  .big-table {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(40px, 1fr));
    gap: 5px;
    border: 1px solid #ccc;
    border-radius: 5px;
    margin: 10px auto;
    padding: 10px;
  }
  
  .sorting-controls {
    margin-bottom: 20px;
  }
  
  .sorting-controls select {
    margin-left: 10px;
    padding: 5px;
    font-size: 16px;
  }
  </style>)