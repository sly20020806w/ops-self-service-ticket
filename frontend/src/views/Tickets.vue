<template>
  <div class="page">
    <NavBar />
    <div style="display:flex;justify-content:space-between;gap:12px;flex-wrap:wrap;align-items:end">
      <div>
        <div class="brand">工单中心</div>
        <div class="sub">按状态过滤；创建时绑定服务树叶子节点；提交后走审批→网关→执行。</div>
      </div>
      <div style="display:flex;gap:8px;flex-wrap:wrap">
        <el-select v-model="status" clearable placeholder="状态" style="width:180px" @change="load">
          <el-option v-for="s in statuses" :key="s" :label="s" :value="s" />
        </el-select>
        <el-button @click="load">刷新</el-button>
        <el-button type="primary" color="#2dd4bf" @click="showCreate = true">新建工单</el-button>
      </div>
    </div>

    <div class="panel" style="margin-top:16px">
      <el-table :data="list" @row-click="(r) => $router.push(`/tickets/${r.id}`)">
        <el-table-column prop="ticketNo" label="单号" width="170" />
        <el-table-column prop="title" label="标题" min-width="160" />
        <el-table-column prop="status" label="状态" width="140" />
        <el-table-column label="优先级" width="90">
          <template #default="{ row }"><span :class="'tag-' + (row.priority || '').toLowerCase()">{{ row.priority }}</span></template>
        </el-table-column>
        <el-table-column label="服务树" min-width="180">
          <template #default="{ row }">{{ row.treeNode?.path || '-' }}</template>
        </el-table-column>
        <el-table-column label="SLA" width="90">
          <template #default="{ row }">
            <span v-if="row.slaBreached" class="tag-sla">超时</span>
            <span v-else style="color:var(--ok)">正常</span>
          </template>
        </el-table-column>
        <el-table-column label="创建人" width="120">
          <template #default="{ row }">{{ row.creator?.display || '-' }}</template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="showCreate" title="新建自主工单" width="640px">
      <el-form label-width="100px">
        <el-form-item label="流程" required>
          <el-select v-model="form.flowId" style="width:100%" @change="onFlowChange">
            <el-option v-for="f in flows" :key="f.id" :label="f.name" :value="f.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="标题" required>
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="优先级">
          <el-radio-group v-model="form.priority">
            <el-radio-button value="P0" />
            <el-radio-button value="P1" />
            <el-radio-button value="P2" />
            <el-radio-button value="P3" />
          </el-radio-group>
        </el-form-item>
        <el-form-item label="服务树叶子">
          <el-select v-model="form.treeNodeId" clearable filterable style="width:100%">
            <el-option v-for="n in leaves" :key="n.id" :label="n.path" :value="n.id" />
          </el-select>
        </el-form-item>
        <el-divider>动态表单</el-divider>
        <el-form-item v-for="f in fields" :key="f.key" :label="f.label" :required="f.required">
          <el-input v-if="f.type==='textarea'" type="textarea" v-model="form.formData[f.key]" />
          <el-input-number v-else-if="f.type==='number'" v-model="form.formData[f.key]" />
          <el-switch v-else-if="f.type==='switch'" v-model="form.formData[f.key]" />
          <el-select v-else-if="f.type==='select'" v-model="form.formData[f.key]" style="width:100%">
            <el-option v-for="o in (f.options||[])" :key="o" :label="o" :value="o" />
          </el-select>
          <el-input v-else v-model="form.formData[f.key]" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreate=false">取消</el-button>
        <el-button type="primary" color="#2dd4bf" :loading="saving" @click="create">创建草稿</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { api } from '../api'
import NavBar from '../components/NavBar.vue'

const router = useRouter()
const list = ref([])
const status = ref('')
const statuses = ['draft','pending_approval','approved','executing','done','rejected','closed']
const showCreate = ref(false)
const saving = ref(false)
const flows = ref([])
const forms = ref([])
const tree = ref([])
const fields = ref([])
const form = reactive({ title: '', flowId: null, priority: 'P2', treeNodeId: null, formData: {} })
const leaves = computed(() => (tree.value || []).filter((n) => n.leaf))

async function load() {
  const { data } = await api.tickets(status.value || undefined)
  list.value = data || []
}

async function onFlowChange() {
  const flow = flows.value.find((x) => x.id === form.flowId)
  if (!flow) return
  const f = forms.value.find((x) => x.id === flow.formId)
  fields.value = parseJSON(f?.fields || f?.fieldsJSON || '[]')
  const data = {}
  fields.value.forEach((x) => { data[x.key] = x.type === 'switch' ? false : (x.type === 'number' ? 0 : '') })
  form.formData = data
}

function parseJSON(v) {
  if (Array.isArray(v)) return v
  try { return JSON.parse(v || '[]') } catch { return [] }
}

async function create() {
  if (!form.title || !form.flowId) return ElMessage.warning('请填写标题和流程')
  saving.value = true
  try {
    const { data } = await api.createTicket({ ...form })
    ElMessage.success('已创建草稿，请提交审批')
    showCreate.value = false
    router.push(`/tickets/${data.id}`)
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  const [t, f, tr] = await Promise.all([api.tickets(), api.flows(), api.tree()])
  list.value = t.data || []
  flows.value = f.data || []
  tree.value = tr.data || []
  const fm = await api.forms()
  forms.value = fm.data || []
  if (flows.value[0]) {
    form.flowId = flows.value[0].id
    await onFlowChange()
  }
})
</script>
