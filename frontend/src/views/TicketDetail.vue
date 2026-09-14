<template>
  <div class="page" v-if="ticket">
    <NavBar />
    <div class="panel">
      <div style="display:flex;justify-content:space-between;gap:12px;flex-wrap:wrap">
        <div>
          <div class="sub">{{ ticket.ticketNo }}</div>
          <div class="brand" style="font-size:24px;margin-top:4px">{{ ticket.title }}</div>
          <div style="margin-top:8px;display:flex;gap:8px;flex-wrap:wrap">
            <el-tag>{{ ticket.status }}</el-tag>
            <el-tag type="warning">{{ ticket.priority }}</el-tag>
            <el-tag :type="ticket.slaBreached ? 'danger' : 'success'">SLA {{ ticket.slaBreached ? '超时' : '正常' }}</el-tag>
            <el-tag type="info">节点 {{ ticket.currentNode }}</el-tag>
          </div>
        </div>
        <div style="display:flex;gap:8px;flex-wrap:wrap;align-items:start">
          <el-button v-for="a in available" :key="a" :type="btnType(a)" @click="doAction(a)">{{ a }}</el-button>
        </div>
      </div>

      <el-descriptions :column="2" border style="margin-top:16px">
        <el-descriptions-item label="创建人">{{ ticket.creator?.display }}</el-descriptions-item>
        <el-descriptions-item label="服务树">{{ ticket.treeNode?.path || '-' }}</el-descriptions-item>
        <el-descriptions-item label="SLA 截止">{{ ticket.slaDeadline || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ ticket.createdAt }}</el-descriptions-item>
      </el-descriptions>

      <h3>表单数据</h3>
      <el-table :data="formRows" size="small">
        <el-table-column prop="k" label="字段" width="160" />
        <el-table-column prop="v" label="值" />
      </el-table>

      <h3>备注后操作</h3>
      <el-input v-model="comment" type="textarea" :rows="2" placeholder="审批意见 / 执行说明 / 评论" />
    </div>

    <div class="panel" style="margin-top:16px">
      <h3 style="margin-top:0">流转记录（课程：动作时间线）</h3>
      <el-timeline>
        <el-timeline-item v-for="a in actions" :key="a.id" :timestamp="a.createdAt" placement="top">
          <b>{{ a.action }}</b> · {{ a.actor?.display || a.actorId }}
          <div class="sub">{{ a.fromNode }} → {{ a.toNode }} {{ a.comment ? '· ' + a.comment : '' }}</div>
        </el-timeline-item>
      </el-timeline>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { api } from '../api'
import NavBar from '../components/NavBar.vue'

const route = useRoute()
const ticket = ref(null)
const available = ref([])
const actions = ref([])
const comment = ref('')

const formRows = computed(() => {
  try {
    const raw = ticket.value?.formData || ticket.value?.formDataJSON || {}
    const obj = typeof raw === 'string' ? JSON.parse(raw) : raw
    return Object.entries(obj || {}).map(([k, v]) => ({ k, v: String(v) }))
  } catch { return [] }
})

function btnType(a) {
  if (a === 'approve' || a === 'execute') return 'primary'
  if (a === 'reject') return 'danger'
  return 'default'
}

async function load() {
  const { data } = await api.ticket(route.params.id)
  ticket.value = data.ticket
  available.value = data.availableActions || []
  actions.value = data.actions || []
}

async function doAction(action) {
  if (action !== 'comment') {
    try {
      await ElMessageBox.confirm(`确认执行动作：${action}？`, '二次确认')
    } catch { return }
  } else if (!comment.value.trim()) {
    return ElMessage.warning('请填写评论')
  }
  await api.action(route.params.id, action, comment.value)
  ElMessage.success('操作成功')
  comment.value = ''
  await load()
}

onMounted(load)
</script>
