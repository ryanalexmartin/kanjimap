import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { Character } from '@/types/types'
import { fetchCharacters, updateCharacterLearned } from '@/api/characters'
import characterData from '@/data/variant-WordData.json'

export const useCharacterStore = defineStore('characters', () => {
  const characters = ref<Character[]>([])
  const isLoggedIn = ref(false)
  const username = ref('insomagent')

  const learnedCount = computed(() => Array.isArray(characters.value) ? characters.value.filter(c => c.learned).length : 0)
  const totalCharacters = computed(() => characters.value.length)

  async function loadCharacters() {
    try {
      const fetchedCharacters = await fetchCharacters(username.value)
      console.log('Fetched characters:', fetchedCharacters)
      
      const charactersArray = Array.isArray(fetchedCharacters) ? fetchedCharacters : [];
      const characterDataArray = Array.isArray(characterData) ? characterData : [];

      // Load learned characters from localStorage
      const learnedCharacterIds = JSON.parse(localStorage.getItem('learnedCharacters') || '[]')

      const mergedCharacters = characterDataArray.map(char => {
        const fetchedChar = charactersArray.find(c => c.char === char.word)
        return {
          id: char.serial,
          char: char.word,
          pinyin: char.meanings[0].pinyin,
          learned: fetchedChar ? fetchedChar.learned : learnedCharacterIds.includes(char.serial),
          frequency: 0,
        }
      })
      
      characters.value = mergedCharacters
    } catch (error) {
      console.error('Error loading characters:', error)
    }
  }

  async function toggleCharacterLearned(character: Character) {
    try {
      const updatedCharacter = await updateCharacterLearned(character)
      const index = characters.value.findIndex(c => c.id === updatedCharacter.id)
      if (index !== -1) {
        characters.value[index] = { ...characters.value[index], ...updatedCharacter }
      } else {
        characters.value.push(updatedCharacter)
      }
      // Force reactivity update
      characters.value = [...characters.value]
      // Update localStorage
      localStorage.setItem('learnedCharacters', JSON.stringify(characters.value.filter(c => c.learned).map(c => c.id)))
    } catch (error) {
      console.error('Error toggling character learned status:', error)
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