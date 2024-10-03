<template>
  <div id="app">
    <HeaderView />
    <router-view></router-view>
    <FooterView />
  </div>
</template>

<script setup lang="ts">
import { provide } from 'vue'
import { useCharacterStore } from './store'
import HeaderView from './components/HeaderView.vue'
import FooterView from './components/FooterView.vue'

const store = useCharacterStore()

// Provide the store to all child components
provide('characterStore', store)
</script>

<style>
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

.character-card {
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid black;
  padding: 5px;
  margin: 1px;
  cursor: pointer;
  background-color: #f0f0f0;
}

.character-card.learned {
  background-color: lightgreen;
}

.big-table {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(40px, 1fr));
  gap: 5px;
  border: 1px solid #ccc;
  border-radius: 5px;
  margin: 10px auto;
  padding: 10px;
}

.sorting-controls {
  margin-bottom: 20px;
}

.sorting-controls select {
  margin-left: 10px;
  padding: 5px;
  font-size: 16px;
}
</style>