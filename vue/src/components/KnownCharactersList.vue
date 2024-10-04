<template>
  <div class="known-characters-list">
    <h2>Known Characters</h2>
      <ins v-for="char in knownCharacters" :key="char.id">
        {{ char.char }}
      </ins>
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
      if (!Array.isArray(this.characters)) {
        console.error('characters prop is not an array:', this.characters);
        return [];
      }
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
