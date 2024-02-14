<template>
  <div id="app">
    <div v-if="state.isLoggedIn">
      <p>Welcome, {{ state.username }}!</p>
      <button @click="logout">Logout</button>
      <h1>Learned {{ learnedCount }} out of {{ totalCharacters }}</h1>
      <RecycleScroller class="scroller" :items="characters" :item-size="50" :gridItems="20" key-field="id"
        v-slot="{ item }">
        <CharacterCard :character="item" />
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
import store from './store'; 

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

    const characters = ref(0);
    const learnedCount = ref(0);
    const totalCharacters = charactersData.filter(c => c.serial.includes('A')).length;

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
      // updateLearned,
      handleLogin,
      logout,
      store
    };
  }
}
</script>

<style>
.vue-recycle-scroller.direction-vertical.scroller {
  height: 80vh;
  width: 100%;
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
