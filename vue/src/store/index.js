import { createStore } from "vuex";

export default createStore({
  state: {
    characters: {},
    username: "",
    loggedIn: false
  },
  mutations: {
    setCharacterLearned(state, { id, learned }) {
        state.characters[id] = learned;    
    },
    setUsername(state, username) {
        state.username = username;
    },
    setLoggedIn(state, loggedIn) {
        state.loggedIn = loggedIn;
    }
  },

  actions: {
    async loadCharacters({ commit, dispatch }, username) {
        console.log("fetch-characters?username=", username);
        const response = await fetch(
            `http://localhost:8081/fetch-characters?username=${username}`
        );
        const characterCards = await response.json();
        for (const card of characterCards) {
            commit("setCharacterLearned", {
                character: card.character,
                characterId: card.characterId,
                learned: card.learned,
            });

            if (card.learned) {
                localStorage.setItem(card.character, card.learned);
            } else {
                localStorage.removeItem(card.character);
            }

            dispatch('updateCharacterLearned', {
                id: card.characterId,
                learned: card.learned
            });
        }
    },
    async updateCharacterLearned({ commit, state }, { id, learned }) {
      console.log("Character has learned value of: ", learned);
      let response = await fetch(`http://localhost:8081/learn-character`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          username: state.username,
          character: id,
          learned: learned,
          characterId: id,
        }),
      });

      // Parse the JSON response and commit the mutation
      let character = await response.json();
      commit("setCharacterLearned", {
        characterId: id,
        learned: character.learned,
      });
    }
    },
});

