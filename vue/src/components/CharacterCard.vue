<template>
  <div :class="{ learned: character.learned }">
    {{ character.char }}
  </div>
</template>

<script>
export default {
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
