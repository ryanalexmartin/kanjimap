<template>
  <div id="app">
    <HeaderView />
    <router-view></router-view>
    <FooterView />
  </div>
</template>

<script setup lang="ts">
import HeaderView from './components/HeaderView.vue'
import FooterView from './components/FooterView.vue'
</script>

<style>
/* Global styles */
</style>