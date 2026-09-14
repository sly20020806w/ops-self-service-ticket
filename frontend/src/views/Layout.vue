<template>
  <div class="shell">
    <div class="topbar fade-up">
      <div class="brand">自主工单 <span>Self-Service</span></div>
      <div style="display:flex;gap:10px;align-items:center;flex-wrap:wrap">
        <el-tag type="info" effect="plain">{{ user.displayName }} · {{ user.role }}</el-tag>
        <el-button text @click="$router.push('/dashboard')">看板</el-button>
        <el-button text @click="$router.push('/tickets')">工单</el-button>
        <el-button text @click="$router.push('/tickets/create')">新建</el-button>
        <el-button text @click="$router.push('/workflows')">流程</el-button>
        <el-button text @click="$router.push('/forms')">表单</el-button>
        <el-button text @click="$router.push('/audit')">审计</el-button>
        <el-button type="danger" plain @click="logout">退出</el-button>
      </div>
    </div>
    <router-view />
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
