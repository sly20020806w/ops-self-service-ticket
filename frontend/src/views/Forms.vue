<template>
  <div class="glass fade-up" style="padding:18px">
    <h2 style="margin-top:0">表单设计器产物</h2>
    <p style="color:#5b6b7c">课程对齐：动态字段（input/textarea/select/number/switch），工单创建时按 Schema 渲染。</p>
    <el-table :data="list">
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="code" label="编码" width="140" />
      <el-table-column prop="version" label="版本" width="80" />
      <el-table-column label="字段数" width="90">
        <template #default="{ row }">{{ parse(row.schemaJson).length }}</template>
      </el-table-column>
      <el-table-column label="字段预览">
        <template #default="{ row }">
          <el-tag v-for="f in parse(row.schemaJson)" :key="f.key" style="margin:2px" size="small">{{ f.label }}({{ f.type }})</el-tag>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { api } from '../api'
const list = ref([])
function parse(s) { try { return JSON.parse(s || '[]') } catch { return [] } }
onMounted(async () => { list.value = (await api.forms()).list || [] })
</script>
