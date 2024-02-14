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
    async loadCharacters({ commit }, username) {
        console.log("loadCharacters", username);
        const response = await fetch(
            `http://localhost:8081/fetch-characters?username=${username}`
        );
        const characters = await response.json();
        for (const character of characters) {
            const { id, learned } = character;
            commit("setCharacterLearned", {
                id: id,
                learned: learned,
            });
            if (learned) {
                console.log("learned", id);
            }
        }
    },
    async updateCharacterLearned({ commit, state }, { id, learned }) {
        console.log("Character has learned value of: ", learned)
        let response = await fetch(`http://localhost:8081/learn-character`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                username: state.username,
                character: state.characters[id].char,
                learned: learned,
                characterId: id,
            }),
        });

        // Parse the JSON response and commit the mutation
        let character = await response.json();
        commit("setCharacterLearned", {
            id: id,
            learned: character.learned,
        });

    }
    },
});





















// Note to self!! The database and the local state are not in sync.  Need to fix this when I get home.
















