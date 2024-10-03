<template>
  <div class="character-card" :class="{ learned: character.learned }" @click="toggleLearned">
    <span class="character">{{ character.char }}</span>
    <span class="pinyin">{{ character.pinyin }}</span>
  </div>
</template>

<script setup>
const props = defineProps(['character'])
const emit = defineEmits(['toggle-learned'])

function toggleLearned() {
  emit('toggle-learned', props.character)
}
</script>

<style scoped>
.character-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border: 1px solid #ccc;
  border-radius: 4px;
  padding: 0.5rem;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.character-card:hover {
  background-color: #f0f0f0;
}

.character-card.learned {
  background-color: #e6f7ff;
}

.character {
  font-size: 1.5rem;
  font-weight: bold;
}

.pinyin {
  font-size: 0.8rem;
  color: #666;
}
</style>