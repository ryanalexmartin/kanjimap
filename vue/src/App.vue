<template>
  <div id="app">
    <Login />
    <Register />
    <h1>Learned {{ learnedCount }} out of {{ totalCharacters }}</h1>
    <RecycleScroller 
      class="scroller"
      :items="characters"
      :item-size="50"
      :gridItems="20"
      key-field="id"
      v-slot="{ item }"
    >
      <CharacterCard :character="item" @update-learned="updateCharacterLearnedState"/>
    </RecycleScroller>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import { RecycleScroller } from 'vue-virtual-scroller';
import 'vue-virtual-scroller/dist/vue-virtual-scroller.css';
import CharacterCard from './components/CharacterCard.vue';
import charactersData from './data/variant-WordData.json';
import Login from './components/Login.vue'
import Register from './components/Register.vue'

export default {
  components: {
    CharacterCard,
    RecycleScroller,
    Login,
    Register
  },
  setup() {
    const characters = ref([]);
    const learnedCount = ref(0);
    const refreshKey = ref(0);

    onMounted(() => {
      characters.value = charactersData
      .filter(char => char.serial.includes('A'))
      .map((char) => {
        return {
          id: char.serial,  
          char: char.word,
          learned: !!localStorage.getItem(char.word)
        };
      });

      learnedCount.value = characters.value.filter(c => c.learned).length;
    });

    const totalCharacters = charactersData.filter(c => c.serial.includes('A')).length;

    function forceRefresh() {
      refreshKey.value++;
    }

    function updateCharacterLearnedState({ id, learned }) {
      const characterToUpdate = characters.value.find(char => char.id === id);
      if (characterToUpdate) {
        characterToUpdate.learned = learned;
        if (learned) {
          learnedCount.value++;
        } else {
          learnedCount.value--;
        }
        forceRefresh();
      }
    }

    return {
      characters,
      learnedCount,
      totalCharacters,
      updateCharacterLearnedState
    };
  }
}
</script>

<style>
.vue-recycle-scroller.direction-vertical.scroller {
  height: 500px;  /* adjust this value based on your design */
  overflow-y: auto;
}
.character-card{
  height: 32%;
  padding: 0 12px;
  display: flex;
  align-items: center;
}
</style>
