<script setup>
import { onMounted } from 'vue'
import Navigation from './components/Navigation.vue'
</script>

<template>
  <div id="app">
    <div class="cursor"></div>
    <div class="cursor-glow"></div>
    <Navigation />
    <router-view />
  </div>
</template>

<style>
@import './styles/main.css';
@import url('https://fonts.googleapis.com/css2?family=Cormorant+Garamond:wght@300;400;600&family=Outfit:wght@200;300;400&display=swap');

#app {
  min-height: 100vh;
}
</style>
