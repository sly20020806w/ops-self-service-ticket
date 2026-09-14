<template>
  <div v-if="route.path === '/login'">
    <router-view />
  </div>
  <div v-else>
    <NavBar />
    <router-view />
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'
import NavBar from './components/NavBar.vue'
const route = useRoute()
</script>
