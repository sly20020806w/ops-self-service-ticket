<template>
  <div class="nav">
    <router-link to="/">看板</router-link>
    <router-link to="/tickets">工单</router-link>
    <router-link to="/flows">表单&流程</router-link>
    <router-link to="/audits">审计&通知</router-link>
    <span style="flex:1"></span>
    <span class="sub" style="margin:0">{{ user.display }} · {{ user.role }}</span>
    <button @click="logout">退出</button>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
const router = useRouter()
const user = computed(() => {
  try { return JSON.parse(localStorage.getItem('user') || '{}') } catch { return {} }
})
function logout() {
  localStorage.clear()
  router.push('/login')
}
</script>
