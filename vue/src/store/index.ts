import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { Character } from '@/types/types'
import { fetchCharacters, updateCharacterLearned } from '@/api/characters'

export const useCharacterStore = defineStore('characters', () => {
  const characters = ref<Character[]>([])
  const isLoggedIn = ref(false)
  const username = ref('')

  const learnedCount = computed(() => characters.value.filter(c => c.learned).length)
  const totalCharacters = computed(() => characters.value.length)

  async function loadCharacters() {
    characters.value = await fetchCharacters(username.value)
  }

  async function toggleCharacterLearned(character: Character) {
    const updatedCharacter = await updateCharacterLearned(character)
    const index = characters.value.findIndex(c => c.id === updatedCharacter.id)
    if (index !== -1) {
      characters.value[index] = updatedCharacter
    }
  }

  return {
    characters,
    isLoggedIn,
    username,
    learnedCount,
    totalCharacters,
    loadCharacters,
    toggleCharacterLearned,
  }
})