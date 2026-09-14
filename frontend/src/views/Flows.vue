<template>
  <div class="page">
    <NavBar />
    <div class="brand">表单设计器 & 流程定义</div>
    <div class="sub">动态表单字段、流程节点、排他网关、服务树叶子绑定。</div>

    <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-top:16px">
      <div class="panel">
        <h3 style="margin-top:0">表单模板</h3>
        <el-table :data="forms" size="small">
          <el-table-column prop="id" label="ID" width="60" />
          <el-table-column prop="name" label="名称" />
          <el-table-column prop="description" label="说明" />
        </el-table>
        <pre class="code">{{ pretty(forms[0]?.fields || forms[0]?.fieldsJSON) }}</pre>
      </div>
      <div class="panel">
        <h3 style="margin-top:0">流程定义</h3>
        <el-table :data="flows" size="small">
          <el-table-column prop="id" label="ID" width="60" />
          <el-table-column prop="name" label="名称" />
          <el-table-column prop="formId" label="表单" width="80" />
        </el-table>
        <h4>Nodes</h4>
        <pre class="code">{{ pretty(flows[0]?.nodes || flows[0]?.nodesJSON) }}</pre>
        <h4>Edges（含 exclusive gateway）</h4>
        <pre class="code">{{ pretty(flows[0]?.edges || flows[0]?.edgesJSON) }}</pre>
      </div>
    </div>

    <div class="panel" style="margin-top:16px">
      <h3 style="margin-top:0">服务树（可绑定叶子节点）</h3>
      <el-table :data="tree" size="small">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="path" label="路径" />
        <el-table-column prop="leaf" label="叶子" width="80">
          <template #default="{ row }">{{ row.leaf ? '是' : '否' }}</template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { api } from '../api'
import NavBar from '../components/NavBar.vue'

const forms = ref([])
const flows = ref([])
const tree = ref([])

function pretty(v) {
  if (v == null) return ''
  try {
    const obj = typeof v === 'string' ? JSON.parse(v) : v
    return JSON.stringify(obj, null, 2)
  } catch {
    return String(v)
  }
}

onMounted(async () => {
  const [a, b, c] = await Promise.all([api.forms(), api.flows(), api.tree()])
  forms.value = a.data
  flows.value = b.data
  tree.value = c.data
})
</script>

<style scoped>
.code {
  background: #0a101c;
  border: 1px solid var(--line);
  border-radius: 8px;
  padding: 10px;
  overflow: auto;
  max-height: 260px;
  color: #9de7dc;
  font-size: 12px;
}
@media (max-width: 900px) {
  div[style*='grid-template-columns'] { grid-template-columns: 1fr !important; }
}
</style>
