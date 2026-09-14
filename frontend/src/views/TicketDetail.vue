<template>
  <div v-if="ticket" class="fade-up" style="display:grid;grid-template-columns:1.4fr 1fr;gap:16px">
    <div class="glass" style="padding:20px">
      <div style="display:flex;justify-content:space-between;gap:12px;flex-wrap:wrap">
        <div>
          <div style="color:#5b6b7c;font-size:13px">{{ ticket.ticketNo }}</div>
          <h2 style="margin:4px 0 8px">{{ ticket.title }}</h2>
          <el-space wrap>
            <el-tag>{{ ticket.status }}</el-tag>
            <el-tag :type="ticket.slaBreached ? 'danger' : 'success'">SLA {{ ticket.slaBreached ? '超时' : '正常' }}</el-tag>
            <el-tag type="warning">{{ ticket.priority }}</el-tag>
            <el-tag type="info">节点 {{ ticket.currentNode || '-' }}</el-tag>
          </el-space>
        </div>
        <el-space wrap>
          <el-button v-if="can('submit')" @click="doAct('submit')">提交</el-button>
          <el-button v-if="can('approve')" type="primary" color="#0e7c7b" @click="ask('approve')">通过</el-button>
          <el-button v-if="can('reject')" type="danger" @click="ask('reject')">驳回</el-button>
          <el-button v-if="can('execute')" type="warning" @click="ask('execute')">执行完成</el-button>
          <el-button v-if="can('close')" @click="ask('close')">关闭</el-button>
          <el-button v-if="can('cancel')" @click="ask('cancel')">撤销</el-button>
        </el-space>
      </div>

      <el-divider />
      <el-descriptions :column="2" border>
        <el-descriptions-item label="创建人">{{ ticket.creator?.displayName }}</el-descriptions-item>
        <el-descriptions-item label="处理人">{{ ticket.assignee?.displayName || '-' }}</el-descriptions-item>
        <el-descriptions-item label="流程">{{ ticket.workflow?.name }}</el-descriptions-item>
        <el-descriptions-item label="服务树">{{ ticket.treeNode?.path || '-' }}</el-descriptions-item>
        <el-descriptions-item label="SLA 截止">{{ ticket.slaDeadline || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ ticket.createdAt }}</el-descriptions-item>
      </el-descriptions>

      <h3>表单数据</h3>
      <el-table :data="formRows" size="small">
        <el-table-column prop="k" label="字段" width="160" />
        <el-table-column prop="v" label="值" />
      </el-table>

      <h3>评论</h3>
      <el-input v-model="comment" type="textarea" :rows="2" placeholder="补充说明" />
      <el-button style="margin-top:8px" @click="doComment" :disabled="!can('comment') && available.length">发表评论</el-button>
    </div>

    <div class="glass" style="padding:20px">
      <h3 style="margin-top:0">流转时间线</h3>
      <el-timeline>
        <el-timeline-item v-for="a in actions" :key="a.id" :timestamp="a.createdAt" placement="top">
          <b>{{ a.action }}</b> · {{ a.actor?.displayName || a.actorId }}
          <div style="color:#5b6b7c;font-size:13px">{{ a.fromStatus }} → {{ a.toStatus }} {{ a.comment ? '· ' + a.comment : '' }}</div>
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

const route = useRoute()
const ticket = ref(null)
const available = ref([])
const actions = ref([])
const comment = ref('')

const formRows = computed(() => {
  try {
    return Object.entries(JSON.parse(ticket.value?.formDataJson || '{}')).map(([k, v]) => ({ k, v: String(v) }))
  } catch { return [] }
})

function can(a) { return available.value.includes(a) }

async function load() {
  const id = route.params.id
  const [detail, acts] = await Promise.all([api.ticket(id), api.ticketActions(id)])
  ticket.value = detail.ticket
  available.value = detail.availableActions || []
  actions.value = acts.list || []
}

async function ask(kind) {
  const { value } = await ElMessageBox.prompt('请填写备注（可选）', kind, { inputValue: comment.value, confirmButtonText: '确认' }).catch(() => ({ value: null }))
  if (value === null) return
  comment.value = value
  await doAct(kind)
}

async function doAct(kind) {
  const id = route.params.id
  const payload = { comment: comment.value }
  const map = { submit: () => api.submit(id), approve: () => api.approve(id, payload), reject: () => api.reject(id, payload), execute: () => api.execute(id, payload), close: () => api.close(id, payload), cancel: () => api.cancel(id, payload) }
  await map[kind]()
  ElMessage.success('操作成功')
  comment.value = ''
  await load()
}

async function doComment() {
  await api.comment(route.params.id, { comment: comment.value })
  ElMessage.success('已评论')
  comment.value = ''
  await load()
}

onMounted(load)
</script>
