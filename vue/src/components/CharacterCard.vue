<template>
  <div @click="updateCharacterLearned(character, !character.learned)" :class="{ learned: character.learned }">
    {{ character.char }}
  </div>
</template>

<script>
import { mapActions } from 'vuex';

export default {
  methods: {
    ...mapActions(['updateCharacterLearned']),
  },
  props: {
    character: Object,
  },

  setup() {
  },
}
</script>

<style scoped>
/* Style for your card here, e.g. */
div {
  border: 1px solid black;
  padding: 10px;
  margin: 5px;
  cursor: pointer;
}

div.learned {
  background-color: lightgreen;
}
</style>
