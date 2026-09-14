<template>
  <div class="page">
    <NavBar />
    <div class="brand">老板演示看板</div>
    <div class="sub">一眼看清工单吞吐、状态分布与 SLA 风险，适合评审/述职截图。</div>

    <div class="stats" style="margin-top:16px">
      <div class="stat"><div class="k">工单总量</div><div class="v">{{ d.total || 0 }}</div></div>
      <div class="stat"><div class="k">SLA 超时</div><div class="v" style="color:var(--danger)">{{ d.slaBreached || 0 }}</div></div>
      <div class="stat"><div class="k">待审批</div><div class="v" style="color:var(--warn)">{{ statusCnt('pending_approval') }}</div></div>
      <div class="stat"><div class="k">已完成</div><div class="v" style="color:var(--ok)">{{ statusCnt('done') }}</div></div>
    </div>

    <div class="panel">
      <div style="display:flex;justify-content:space-between;gap:12px;flex-wrap:wrap;align-items:center">
        <h3 style="margin:0">状态分布</h3>
        <el-button type="primary" color="#2dd4bf" @click="$router.push('/tickets')">去创建/处理工单</el-button>
      </div>
      <el-table :data="d.byStatus || []" style="margin-top:12px" size="small">
        <el-table-column prop="Status" label="状态" />
        <el-table-column prop="Cnt" label="数量" />
      </el-table>
      <div class="sub" style="margin-top:10px">生成时间：{{ d.generatedAt || '-' }}</div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { api } from '../api'
import NavBar from '../components/NavBar.vue'

const d = ref({})
function statusCnt(s) {
  const row = (d.value.byStatus || []).find((x) => x.Status === s)
  return row ? row.Cnt : 0
}
onMounted(async () => {
  const { data } = await api.dashboard()
  d.value = data
})
</script>
