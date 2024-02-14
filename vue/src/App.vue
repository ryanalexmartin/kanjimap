<template>
  <div id="app">
    <div v-if="state.isLoggedIn">
      <p>Welcome, {{ state.username }}!</p>
      <button @click="logout">Logout</button>
      <h1>Learned {{ learnedCount }} out of {{ totalCharacters }}</h1>
      <div class="big-table">
      <div v-for="item in characters" :key="item.id">
        <div class="character-card" :class="{ learned: item.learned }" @click="updateCharacterLearned(item, !item.learned)">
          {{ item.char }}
        </div>
      </div>
      </div>
    </div>
    <div v-else>
      <LoginView v-if="!state.isRegistering" @login="handleLogin" @register="showRegisterView" />
      <RegisterView v-else @registered="showLoginView" />
    </div>
  </div>
</template>

<script>
import { ref, onMounted, reactive } from 'vue';
import 'vue-virtual-scroller/dist/vue-virtual-scroller.css';
import charactersData from './data/variant-WordData.json';
import LoginView from './components/LoginView.vue';
import RegisterView from './components/RegisterView.vue';
import store from './store'; 

export default {
  components: {
    LoginView,
    RegisterView
  },
  setup() {
    const state = reactive({
      isRegistering: false,
      isLoggedIn: false,
      username: ''
    });
    const handleLogin = (username) => {
      state.isLoggedIn = true;
      state.username = username;
     };
    const logout = ({ commit }) => {
      commit('setLoggedIn', false);
      commit('setUsername', '');
    };
    const showRegisterView = () => {
      state.isRegistering = true;
    };
    const showLoginView = () => {
      state.isRegistering = false;
    };

    // const characters = ref(0);
    const characters = ref(0);
    const learnedCount = ref(0);
    const totalCharacters = charactersData.filter(c => c.serial.includes('A')).length;

    const updateCharacterLearned = (character, learned) => {
      character.learned = learned;
      localStorage.setItem(character.char, learned);
      learnedCount.value = characters.value.filter(c => c.learned).length;
      // Send a request to the server to update the learned status (JSON)
      let response = fetch('http://localhost:8081/learn-character', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          username: state.username,
          chinese_character: character.char,
          characterId: character.id,
          learned: learned,
        }),
      });

      console.log(response);

    };

    onMounted(() => {
      characters.value = charactersData
        .filter(char => char.serial.includes('A'))
        .map((char) => {
          return {
            id: char.serial,
            char: char.word,
            learned: localStorage.getItem(char.word) === 'true',
          };
        });
      learnedCount.value = characters.value.filter(c => c.learned).length;
    });

    return {
      state,
      showRegisterView,
      showLoginView,
      characters,
      learnedCount,
      totalCharacters,
      updateCharacterLearned,
      handleLogin,
      logout,
      store
    };
  }
}
</script>

<style>
.character-card {
  display: flex;
  align-items: center;
  justify-content: center; /* Add this line to center the text horizontally */
  border: 1px solid black;
  padding: 5px;
  margin: 1px;
  cursor: pointer;
}
div.learned {
  background-color: lightgreen;
}

/* For debugging */
/* div {
  border: 1px solid black;
}  */

/* Display a dynamic grid of character cards based on browser width */
.big-table{
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(40px, 1fr));
  gap: 5px;
  border: 1px solid #ccc;
  border-radius: 5px;
  margin: 10 auto;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
}

</style>
