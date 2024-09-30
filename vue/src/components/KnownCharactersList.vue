<template>
  <div class="known-characters-list">
    <h2>Known Characters</h2>
    <div class="character-grid">
      <div v-for="char in knownCharacters" :key="char.id" class="character-item">
        {{ char.char }}
        <span class="frequency">({{ char.frequency }})</span>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    characters: {
      type: Array,
      required: true
    }
  },
  computed: {
    knownCharacters() {
      return this.characters.filter(char => char.learned).sort((a, b) => b.frequency - a.frequency);
    }
  }
}
</script>

<style scoped>
.known-characters-list {
  margin-top: 20px;
}
.character-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(40px, 1fr));
  gap: 5px;
}
.character-item {
  border: 1px solid #ccc;
  padding: 5px;
  text-align: center;
  background-color: #f0f0f0;
}
.frequency {
  font-size: 0.8em;
  color: #666;
}
</style>
