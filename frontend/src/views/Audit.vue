<template>
  <div class="fade-up" style="display:grid;grid-template-columns:1fr 1fr;gap:16px">
    <div class="glass" style="padding:18px">
      <h2 style="margin-top:0">操作审计</h2>
      <el-table :data="audit" size="small" max-height="520">
        <el-table-column prop="createdAt" label="时间" width="170" />
        <el-table-column prop="username" label="用户" width="100" />
        <el-table-column prop="action" label="动作" width="90" />
        <el-table-column prop="resource" label="资源" width="90" />
        <el-table-column prop="resourceId" label="ID" />
      </el-table>
    </div>
    <div class="glass" style="padding:18px">
      <h2 style="margin-top:0">IM 通知日志</h2>
      <el-table :data="notify" size="small" max-height="520">
        <el-table-column prop="createdAt" label="时间" width="170" />
        <el-table-column prop="channel" label="通道" width="90" />
        <el-table-column prop="target" label="对象" width="100" />
        <el-table-column prop="success" label="成功" width="70">
          <template #default="{ row }"><el-tag :type="row.success?'success':'danger'" size="small">{{ row.success }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="content" label="内容" show-overflow-tooltip />
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { api } from '../api'
const audit = ref([])
const notify = ref([])
onMounted(async () => {
  audit.value = (await api.audit()).list || []
  notify.value = (await api.notifyLogs()).list || []
})
</script>
