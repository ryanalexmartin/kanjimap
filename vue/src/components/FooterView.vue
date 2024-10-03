<template>
  <footer>
    <p>&copy; {{ currentYear }} KanjiMap. All rights reserved.</p>
  </footer>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const currentYear = ref(new Date().getFullYear())
</script>

<style scoped>
footer {
  padding: 1rem;
  background-color: #f0f0f0;
  text-align: center;
}
</style>