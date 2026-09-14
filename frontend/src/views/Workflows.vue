<template>
  <div class="glass fade-up" style="padding:18px">
    <h2 style="margin-top:0">流程管理</h2>
    <p style="color:#5b6b7c">课程对齐：工作流节点、审批/执行角色、排他网关分流。</p>
    <el-collapse>
      <el-collapse-item v-for="w in list" :key="w.id" :title="`${w.name} · ${w.code} · SLA ${w.slaHours}h`">
        <p>{{ w.description }}</p>
        <el-table :data="w.nodes || []" size="small">
          <el-table-column prop="sort" label="#" width="60" />
          <el-table-column prop="key" label="Key" width="120" />
          <el-table-column prop="name" label="名称" />
          <el-table-column prop="type" label="类型" width="100" />
          <el-table-column prop="assigneeRole" label="角色" width="110" />
          <el-table-column prop="gatewayExpr" label="网关表达式" />
          <el-table-column prop="defaultNext" label="默认下一跳" width="120" />
        </el-table>
      </el-collapse-item>
    </el-collapse>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { api } from '../api'
const list = ref([])
onMounted(async () => { list.value = (await api.workflows()).list || [] })
</script>
