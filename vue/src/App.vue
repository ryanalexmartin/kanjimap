<template>
  <div id="app">
    <div v-if="state.isLoggedIn">
      <p>Welcome, {{ state.username }}!</p>
      <button @click="logout">Logout</button>
      <h1>Learned {{ learnedCount }} out of {{ totalCharacters }}</h1>
      <RecycleScroller class="scroller" :items="characters" :item-size="50" :gridItems="20" key-field="id"
        v-slot="{ item }">
        <CharacterCard :character="item" @update-learned="updateCharacterLearnedState" />
      </RecycleScroller>
    </div>
    <div v-else>
      <LoginView v-if="!state.isRegistering" @login="handleLogin" @register="showRegisterView" />
      <RegisterView v-else @registered="showLoginView" />
    </div>
  </div>
</template>

<script>
import { ref, onMounted, reactive } from 'vue';
import { RecycleScroller } from 'vue-virtual-scroller';
import 'vue-virtual-scroller/dist/vue-virtual-scroller.css';
import CharacterCard from './components/CharacterCard.vue';
import charactersData from './data/variant-WordData.json';
import LoginView from './components/LoginView.vue';
import RegisterView from './components/RegisterView.vue';


export default {
  components: {
    CharacterCard,
    RecycleScroller,
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

    const characters = ref([]);
    const learnedCount = ref(0);
    // const refreshKey = ref(0);
    const totalCharacters = charactersData.filter(c => c.serial.includes('A')).length;

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

    // function forceRefresh() {
    //   refreshKey.value++;
    // }

    function updateCharacterLearnedState(character, learned) {
      // Update the state in the client
      character.learned = learned;
      this.learnedCount = this.characters.filter(c => c.learned).length;

      // Send an update to the server
      fetch('http://localhost:8081/learn', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          username: this.username,
          character: character.character,
          learned: learned,
        }),
      })
        .then(response => {
          if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
          }
          return response.text();
        })
        .then(message => {
          console.log(message);
        })
        .catch(error => {
          console.log('There was a problem with the learn request.', error);
        });
    }

    return {
      state,
      showRegisterView,
      showLoginView,
      characters,
      learnedCount,
      totalCharacters,
      updateCharacterLearnedState,
      handleLogin,
      logout
    };
  }
}
</script>

<style>
.vue-recycle-scroller.direction-vertical.scroller {
  height: 500px;
  /* adjust this value based on your design */
  overflow-y: auto;
}

.character-card {
  height: 32%;
  padding: 0 12px;
  display: flex;
  align-items: center;
}
</style>
