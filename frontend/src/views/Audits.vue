<template>
  <div class="page">
    <div class="brand">审计日志 & IM/Webhook 通知</div>
    <div class="sub">增强项：全链路留痕；课程通知能力在此可对接企业微信/钉钉 Webhook</div>

    <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-top:16px">
      <div class="panel">
        <h3 style="margin-top:0">Audit</h3>
        <el-table :data="audits" size="small" max-height="520">
          <el-table-column prop="createdAt" label="时间" width="170" />
          <el-table-column prop="action" label="动作" width="140" />
          <el-table-column prop="resource" label="资源" />
        </el-table>
      </div>
      <div class="panel">
        <h3 style="margin-top:0">Notify</h3>
        <el-table :data="notifies" size="small" max-height="520">
          <el-table-column prop="createdAt" label="时间" width="170" />
          <el-table-column prop="channel" label="通道" width="90" />
          <el-table-column prop="target" label="目标" width="100" />
          <el-table-column prop="content" label="内容" />
          <el-table-column prop="status" label="状态" width="80" />
        </el-table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { api } from '../api'

const audits = ref([])
const notifies = ref([])

onMounted(async () => {
  const [a, b] = await Promise.all([api.audits(), api.notifies()])
  audits.value = a.data
  notifies.value = b.data
})
</script>

<style scoped>
@media (max-width: 900px) {
  div[style*='grid-template-columns'] { grid-template-columns: 1fr !important; }
}
</style>
