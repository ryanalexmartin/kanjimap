<template>
  <header>
    <h1>KanjiMap</h1>
    <h2>漢字地圖</h2>
    <h2>by Ryan Alex Martin</h2>
    <h2>馬丁瑞安創造</h2>
    <a href="https://ryanalexmartin.com">ryanalexmartin.com</a>
    <br />
    <a href="https://github.com/ryanalexmartin/kanjimap"
      >github.com/ryanalexmartin/kanjimap</a
    >
    <br />
    <br />
    <div>
      Click on a character to mark it as learned. Click again to mark it as
      unlearned. I recommend something like
      <a href="https://github.com/gkovacs/LiuChanFirefox/tree/firefox"
        >LiuChan</a
      >
      to help you learn the characters.
      <br />
      點擊一個字符標記為已學。再次點擊標記為未學。
    </div>
    <div class="learn-from-passage">
      Paste Chinese text here to learn characters automatically:
      <br />
      將中文文本粘貼到此處以自動學習字符：
      <textarea v-model="state.passage"></textarea>
      <button @click="learnFromPassage">Learn from Passage 從文本學習</button>
    </div>
  </header>
  <div id="app">
    <div v-if="state.isLoggedIn">
      <p>歡迎 Welcome, {{ state.username }}!</p>
      <button @click="logout">Logout 登出</button>
      <h1>Learned {{ learnedCount }} out of {{ totalCharacters }}</h1>
      <h1>已學 {{ learnedCount }} / {{ totalCharacters }}</h1>

      <!-- Add sorting dropdown -->
      <div class="sorting-controls">
        <label for="sort-select">Sort by: </label>
        <select
          id="sort-select"
          v-model="state.sortBy"
          @change="sortCharacters"
        >
          <option value="default">Default</option>
          <option value="frequency">Frequency</option>
        </select>
      </div>

      <div class="big-table">
        <div v-for="item in sortedCharacters" :key="item.id">
          <div
            class="character-card"
            :class="{ learned: item.learned }"
            @click="updateCharacterLearned(item, !item.learned)"
          >
            {{ item.char }}
            <span class="frequency" v-if="state.sortBy === 'frequency'">
              <!-- ({{ item.frequency }}) -->
            </span>
          </div>
        </div>
      </div>
    </div>
    <div v-else>
      <LoginView
        v-if="!state.isRegistering"
        @login="handleLogin"
        @register="showRegisterView"
      />
      <RegisterView v-else @registered="showLoginView" />
    </div>
  </div>
</template>

<script>
import { ref, onMounted, reactive, computed } from "vue";
import charactersData from "./data/variant-WordData.json";
import LoginView from "./components/LoginView.vue";
import RegisterView from "./components/RegisterView.vue";

export default {
  components: {
    LoginView,
    RegisterView,
  },
  setup() {
    const state = reactive({
      isRegistering: false,
      isLoggedIn: localStorage.getItem("token") !== null,
      username: localStorage.getItem("username"),
      passage: "",
      sortBy: "default",
    });

    // Initialize characters as an empty array
    const characters = ref([]);
    const learnedCount = ref(0);
    const totalCharacters = charactersData.filter(c => c.serial.includes('A')).length;

    // learn all characters in the 'passage' text field
    // get 'passage' via v-model
    const learnFromPassage = async () => {
      // parse for chinese characters
      const characters = state.passage.match(/[\u4e00-\u9fa5]/g);
      if (characters) {
        for (const char of characters) {
          const character = charactersData.find((c) => c.word === char);
          if (character) {
            updateCharacterLearned(
              {
                id: character.serial,
                char: character.word,
                learned: true,
              },
              true
            );
          }
        }
      }
      // refresh the characters
      fetchCharacters(state.username);
    };

    const fetchCharacters = async (username) => {
      try {
        const response = await fetch(
          `${process.env.VUE_APP_API_URL}:${process.env.VUE_APP_API_PORT}/fetch-characters?username=${username}`,
          {
            method: 'GET',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${localStorage.getItem('token')}`
            },
          }
        );
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const characterCards = await response.json();
        if (characterCards && Array.isArray(characterCards)) {
          characters.value = characterCards.map(card => ({
            id: card.characterId,
            char: card.character,
            learned: card.learned,
            frequency: card.frequency || 0, // Default to 0 if frequency is not provided
          }));
        } else {
          console.error('Received invalid data from server:', characterCards);
          characters.value = [];
        }
      } catch (error) {
        console.error('Error fetching characters:', error);
        characters.value = [];
      }
      learnedCount.value = characters.value.filter(c => c.learned).length;
    };

    const sortCharacters = () => {
      // This function will be called when the sorting option changes
      // The actual sorting is done in the computed property 'sortedCharacters'
    };

    const sortedCharacters = computed(() => {
      if (!characters.value || characters.value.length === 0) {
        return [];
      }

      // Filter characters starting with 'A'
      const filteredCharacters = characters.value.filter(char => char.id.startsWith('A'));

      if (state.sortBy === 'frequency') {
        return filteredCharacters.sort((a, b) => (b.frequency || 0) - (a.frequency || 0));
      } else {
        // Default sorting (by id)
        return filteredCharacters.sort((a, b) => a.id.localeCompare(b.id));
      }
    });

    if (state.isLoggedIn) {
      fetchCharacters(state.username);
    }

    const handleLogin = async (username) => {
      // Add the missing 'commit' parameter
      state.isLoggedIn = true;
      state.username = username;
      await fetchCharacters(username);
    };
    const logout = () => {
      state.isLoggedIn = false;
      state.username = "";
      localStorage.clear();
    };
    const showRegisterView = () => {
      state.isRegistering = true;
    };
    const showLoginView = () => {
      state.isRegistering = false;
    };


    const updateCharacterLearned = (character, learned) => {
      character.learned = learned;
      localStorage.setItem(character.char, learned);
      learnedCount.value = characters.value.filter((c) => c.learned).length;
      // Send a request to the server to update the learned status (JSON)
      let response = fetch(
        `${process.env.VUE_APP_API_URL}:${process.env.VUE_APP_API_PORT}/learn-character`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            username: state.username,
            chinese_character: character.char,
            characterId: character.id,
            learned: learned,
          }),
        }
      );

      console.log(response);
    };

    onMounted(() => {
      characters.value = charactersData
        .filter((char) => char.serial.startsWith("A"))
        .map((char) => {
          return {
            id: char.serial,
            char: char.word,
            learned: localStorage.getItem(char.word) === "true",
          };
        });
      learnedCount.value = characters.value.filter((c) => c.learned).length;
    });

    return {
      state,
      showRegisterView,
      showLoginView,
      characters,
      sortedCharacters,
      learnedCount,
      totalCharacters,
      updateCharacterLearned,
      handleLogin,
      logout,
      learnFromPassage,
      sortCharacters,
    };
  },
};
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

.learn-from-passage {
  margin-top: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.learn-from-passage textarea {
  width: 300px;
  height: 200px;
  font-size: 12px;
  padding: 5px;
  margin-right: 10px;
  border-radius: 5px;
  resize: none;
}

.sorting-controls {
  margin-bottom: 20px;
}

.sorting-controls select {
  margin-left: 10px;
  padding: 5px;
  font-size: 16px;
}

.character-card {
  position: relative;
}

.character-card .frequency {
  position: absolute;
  bottom: 2px;
  right: 2px;
  font-size: 10px;
  color: #666;
}
</style>
