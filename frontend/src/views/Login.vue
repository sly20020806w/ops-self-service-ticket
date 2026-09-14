<template>
  <div class="login-page">
    <div class="panel login-box">
      <div class="brand">自主工单 <span>Ticket</span></div>
      <div class="sub">自主工单平台：表单、流程、排他网关、服务树、状态动作、IM 通知；含 SLA、审计与运营看板。</div>
      <el-form style="margin-top:18px" @submit.prevent="onLogin">
        <el-form-item label="账号">
          <el-select v-model="username" style="width:100%">
            <el-option label="alice · 申请人" value="alice" />
            <el-option label="bob · 审批人" value="bob" />
            <el-option label="carol · 执行人" value="carol" />
            <el-option label="admin · 管理员" value="admin" />
          </el-select>
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="password" type="password" show-password />
        </el-form-item>
        <el-button type="primary" native-type="submit" style="width:100%" :loading="loading" color="#2dd4bf">进入平台</el-button>
      </el-form>
      <div class="sub" style="margin-top:12px;font-size:12px">演示密码：admin123</div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { api } from '../api'

const router = useRouter()
const username = ref('alice')
const password = ref('admin123')
const loading = ref(false)

async function onLogin() {
  loading.value = true
  try {
    const { data } = await api.login(username.value, password.value)
    localStorage.setItem('token', data.token)
    localStorage.setItem('user', JSON.stringify(data.user))
    ElMessage.success(`欢迎 ${data.user.display}`)
    router.push('/')
  } finally {
    loading.value = false
  }
}
</script>
