<template>
  <div class="glass fade-up" style="padding:20px;max-width:860px">
    <h2 style="margin-top:0">新建自主工单</h2>
    <p style="color:#5b6b7c;margin-top:-6px">支持模板预填、动态表单、服务树叶子绑定、直接提交或存草稿。</p>

    <el-form label-width="110px">
      <el-form-item label="快捷模板">
        <el-select v-model="templateId" clearable placeholder="可选" style="width:100%" @change="applyTemplate">
          <el-option v-for="t in templates" :key="t.id" :label="t.name" :value="t.id" />
        </el-select>
      </el-form-item>
      <el-form-item label="流程" required>
        <el-select v-model="form.workflowId" style="width:100%" @change="onWorkflowChange">
          <el-option v-for="w in workflows" :key="w.id" :label="`${w.name} (SLA ${w.slaHours}h)`" :value="w.id" />
        </el-select>
      </el-form-item>
      <el-form-item label="标题" required>
        <el-input v-model="form.title" />
      </el-form-item>
      <el-form-item label="优先级">
        <el-radio-group v-model="form.priority">
          <el-radio-button value="low">low</el-radio-button>
          <el-radio-button value="medium">medium</el-radio-button>
          <el-radio-button value="high">high</el-radio-button>
          <el-radio-button value="urgent">urgent</el-radio-button>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="服务树叶子">
        <el-select v-model="form.treeNodeId" clearable filterable style="width:100%" placeholder="课程要求：绑定叶子节点">
          <el-option v-for="n in leafNodes" :key="n.id" :label="n.path" :value="n.id" />
        </el-select>
      </el-form-item>

      <el-divider>动态表单字段</el-divider>
      <el-form-item v-for="f in fields" :key="f.key" :label="f.label" :required="f.required">
        <el-input v-if="f.type==='input'" v-model="form.formData[f.key]" :placeholder="f.placeholder" />
        <el-input v-else-if="f.type==='textarea'" type="textarea" v-model="form.formData[f.key]" :rows="3" />
        <el-input-number v-else-if="f.type==='number'" v-model="form.formData[f.key]" />
        <el-switch v-else-if="f.type==='switch'" v-model="form.formData[f.key]" />
        <el-select v-else-if="f.type==='select'" v-model="form.formData[f.key]" style="width:100%">
          <el-option v-for="o in (f.options||[])" :key="o" :label="o" :value="o" />
        </el-select>
        <el-input v-else v-model="form.formData[f.key]" />
      </el-form-item>

      <el-form-item>
        <el-button @click="save(false)" :loading="loading">存草稿</el-button>
        <el-button type="primary" color="#0e7c7b" @click="save(true)" :loading="loading">提交审批</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { api } from '../api'

const router = useRouter()
const loading = ref(false)
const workflows = ref([])
const templates = ref([])
const tree = ref([])
const fields = ref([])
const templateId = ref()
const form = reactive({
  workflowId: null, title: '', priority: 'medium', treeNodeId: null, formData: {},
})

const leafNodes = computed(() => (tree.value || []).filter((n) => n.isLeaf))

onMounted(async () => {
  const [w, t, tr] = await Promise.all([api.workflows(), api.templates(), api.tree()])
  workflows.value = w.list || []
  templates.value = t.list || []
  tree.value = tr.list || []
  if (workflows.value[0]) {
    form.workflowId = workflows.value[0].id
    await onWorkflowChange()
  }
})

async function onWorkflowChange() {
  const wf = workflows.value.find((x) => x.id === form.workflowId)
  if (!wf) return
  const f = await api.form(wf.formId)
  fields.value = JSON.parse(f.schemaJson || '[]')
  const data = {}
  fields.value.forEach((x) => { data[x.key] = x.default ?? (x.type === 'switch' ? false : '') })
  form.formData = data
}

async function applyTemplate(id) {
  const t = templates.value.find((x) => x.id === id)
  if (!t) return
  form.workflowId = t.workflowId
  form.title = t.titlePattern
  form.priority = t.priority
  await onWorkflowChange()
  form.formData = { ...form.formData, ...JSON.parse(t.formDataJson || '{}') }
}

async function save(submit) {
  if (!form.title || !form.workflowId) return ElMessage.warning('请填写标题和流程')
  loading.value = true
  try {
    const t = await api.createTicket({ ...form, submit, templateId: templateId.value || null })
    ElMessage.success(submit ? '已提交' : '草稿已保存')
    router.push(`/tickets/${t.id}`)
  } finally {
    loading.value = false
  }
}
</script>
