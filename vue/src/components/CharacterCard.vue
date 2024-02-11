<template>
  <div @click="markLearned" :class="{ learned: character.learned }">
    {{ character.char }}
  </div>
</template>

<script>
export default {
  props: {
    character: Object
  },
  setup(props, { emit }) {
    
    function markLearned() {
      const newLearnedState = !props.character.learned;
      
      // Emit the change to the parent
      emit('update-learned', { id: props.character.id, learned: newLearnedState });

      // Save to local storage
      if (newLearnedState) {
        localStorage.setItem(props.character.char, 'true');
      } else {
        localStorage.removeItem(props.character.char);
      }
    }

    return {
      markLearned
    };
  }
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
