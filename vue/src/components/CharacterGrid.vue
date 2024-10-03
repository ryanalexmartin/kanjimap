<template>
  <div class="character-grid">
    <CharacterCard
      v-for="character in characters"
      :key="character.id"
      :character="character"
      @toggle-learned="emitToggleLearned(character)"
    />
  </div>
</template>

<script setup lang="ts">
import { Character } from '../types/types'
import CharacterCard from '@/components/CharacterCard.vue'

const props = defineProps<{
  characters: Character[]
}>()

const emit = defineEmits<{
  (e: 'toggle-learned', character: Character): void
}>()

const emitToggleLearned = (character: Character) => {
  emit('toggle-learned', character)
}
</script>

<style scoped>
.character-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(40px, 1fr));
  gap: 5px;
}
</style>