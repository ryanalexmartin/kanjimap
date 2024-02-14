<template>
  <header>
    <h1>HanziMap</h1>
    <h2>by Ryan Alex Martin</h2>
    <a href="https://ryanalexmartin.com">ryanalexmartin.com</a>
    <br />
    <a href="https://github.com/ryanalexmartin/kanjimap">github.com/ryanalexmartin/kanjimap</a>
    <br />
    <br />
    <div>
      Click on a character to mark it as learned.  Click again to mark it as unlearned. I recommend something like <a href="https://github.com/gkovacs/LiuChanFirefox/tree/firefox">LiuChan</a> to help you learn the characters.  
    </div>
  </header>
  <div id="app">
    <div v-if="state.isLoggedIn">
      <p>Welcome, {{ state.username }}!</p>
      <button @click="logout">Logout</button>
      <h1>Learned {{ learnedCount }} out of {{ totalCharacters }}</h1>
      <div class="big-table">
        <div v-for="item in characters" :key="item.id">
          <div class="character-card" :class="{ learned: item.learned }"
            @click="updateCharacterLearned(item, !item.learned)">
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
import charactersData from './data/variant-WordData.json';
import LoginView from './components/LoginView.vue';
import RegisterView from './components/RegisterView.vue';

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

    const handleLogin = async (username) => { // Add the missing 'commit' parameter
      state.isLoggedIn = true;
      state.username = username;
      const response = await fetch(
        `http://localhost:8081/fetch-characters?username=${username}`
      );
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      } else {
        const characterCards = await response.json();
        for (const card of characterCards) {
          if (card.learned) {
            localStorage.setItem(card.character, card.learned);
          } else {
            localStorage.removeItem(card.character);
          }
        }
      }
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
    };
    const logout = () => {
      state.isLoggedIn = false;
      state.username = '';
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
    };
  }
}
</script>

<style>
.character-card {
  display: flex;
  align-items: center;
  justify-content: center;
  /* Add this line to center the text horizontally */
  border: 1px solid black;
  padding: 5px;
  margin: 1px;
  cursor: pointer;
  background-color: #f0f0f0;
}

div.learned {
  background-color: lightgreen;
}

/* For debugging */
/* div {
  border: 1px solid black;
}  */

/* Display a dynamic grid of character cards based on browser width */
.big-table {
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

/* Make everything look nice */
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
  background-color: lightgray;
}

html {
  background-color: gray;
}
</style>
