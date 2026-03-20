<!-- 核心目的：患者属性适配规则管理页面 -->
<!-- 模块功能：规则引擎-根据患者年龄/性别/孕产妇状态过滤可用设备/时段 -->
<script setup lang="ts">
import type { PatientAdaptRule } from '@/types/rule'
import { message, Modal } from 'ant-design-vue'
import { h, onMounted, ref } from 'vue'
import { ruleApi } from '@/api/rule'

const loading = ref(false)
const items = ref<PatientAdaptRule[]>([])
const showModal = ref(false)
const saving = ref(false)
const isEdit = ref(false)
const editId = ref('')
const editForm = ref<Partial<PatientAdaptRule>>({})

const COND_TYPE_MAP: Record<string, string> = {
  age: '年龄范围',
  gender: '性别',
  pregnancy: '孕产妇',
}
const ACTION_MAP: Record<string, string> = {
  filter_device: '过滤设备',
  filter_slot: '过滤时段',
  filter_doctor: '过滤医生',
}
const STATUS_COLOR: Record<string, string> = { active: 'success', inactive: 'default' }
const STATUS_LABEL: Record<string, string> = { active: '启用', inactive: '停用' }

const columns = [
  {
    title: '条件类型',
    dataIndex: 'condition_type',
    key: 'condition_type',
    customRender: ({ text }: { text: string }) => COND_TYPE_MAP[text] ?? text,
  },
  { title: '条件值', dataIndex: 'condition_value', key: 'condition_value' },
  {
    title: '适配动作',
    dataIndex: 'action',
    key: 'action',
    customRender: ({ text }: { text: string }) => h('a-tag', {}, ACTION_MAP[text] ?? text),
  },
  {
    title: '动作参数',
    key: 'action_params',
    customRender: ({ record }: { record: PatientAdaptRule }) =>
      Object.entries(record.action_params ?? {}).map(([k, v]) => `${k}=${v}`).join('; ') || '-',
  },
  { title: '优先级', dataIndex: 'priority', key: 'priority' },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    customRender: ({ text }: { text: string }) =>
      h('a-badge', { status: STATUS_COLOR[text] ?? 'default', text: STATUS_LABEL[text] ?? text }),
  },
  { title: '操作', key: 'actions', width: 160 },
]

async function fetchData() {
  loading.value = true
  try {
    const res = await ruleApi.listPatientAdaptRules()
    items.value = res.items ?? []
  }
  finally { loading.value = false }
}

onMounted(() => fetchData())

function openCreate() {
  isEdit.value = false
  editId.value = ''
  editForm.value = {
    condition_type: 'age',
    condition_value: '',
    action: 'filter_device',
    action_params: {},
    priority: 0,
  }
  showModal.value = true
}

function openEdit(record: PatientAdaptRule) {
  isEdit.value = true
  editId.value = record.id
  editForm.value = {
    condition_type: record.condition_type,
    condition_value: record.condition_value,
    action: record.action,
    action_params: { ...record.action_params },
    priority: record.priority,
    status: record.status,
  }
  showModal.value = true
}

// action_params 以 "key=value;key2=value2" 格式编辑
const actionParamsStr = ref('')
function syncParamsFromStr() {
  const result: Record<string, string> = {}
  actionParamsStr.value.split(';').forEach((pair) => {
    const [k, ...rest] = pair.split('=')
    if (k?.trim()) {
      result[k.trim()] = rest.join('=').trim()
    }
  })
  editForm.value.action_params = result
}
function openCreateOrEdit(record?: PatientAdaptRule) {
  if (record) {
    openEdit(record)
    actionParamsStr.value = Object.entries(record.action_params ?? {}).map(([k, v]) => `${k}=${v}`).join(';')
  }
  else {
    openCreate()
    actionParamsStr.value = ''
  }
}

async function handleSave() {
  if (!editForm.value.condition_value?.trim()) {
    message.warning('请填写条件值')
    return
  }
  syncParamsFromStr()
  saving.value = true
  try {
    // 后端以整个列表保存，此处演示单条追加/更新（实际可全量提交）
    const payload = [editForm.value]
    await ruleApi.savePatientAdaptRules(payload)
    message.success(isEdit.value ? '更新成功' : '创建成功')
    showModal.value = false
    fetchData()
  }
  finally { saving.value = false }
}

function handleToggleStatus(record: PatientAdaptRule) {
  const next = record.status === 'active' ? 'inactive' : 'active'
  Modal.confirm({
    title: `确认${next === 'inactive' ? '停用' : '启用'}`,
    content: `确定要${next === 'inactive' ? '停用' : '启用'}该适配规则吗？`,
    onOk: async () => {
      await ruleApi.savePatientAdaptRules([{ ...record, status: next }])
      message.success('操作成功')
      fetchData()
    },
  })
}
</script>

<template>
  <div>
    <div class="action-bar mb-4 flex items-center gap-2">
      <a-button type="primary" @click="openCreateOrEdit()">
        新增适配规则
      </a-button>
      <a-button @click="fetchData">
        刷新
      </a-button>
    </div>

    <a-table
      :columns="columns"
      :data-source="items"
      :loading="loading"
      row-key="id"
      size="middle"
      :pagination="false"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button size="small" @click="openCreateOrEdit(record as PatientAdaptRule)">
              编辑
            </a-button>
            <a-button
              size="small"
              :type="record.status === 'active' ? 'default' : 'primary'"
              @click="handleToggleStatus(record as PatientAdaptRule)"
            >
              {{ record.status === 'active' ? '停用' : '启用' }}
            </a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <!-- 新增/编辑弹窗 -->
    <a-modal
      v-model:open="showModal"
      :title="isEdit ? '编辑适配规则' : '新增适配规则'"
      :confirm-loading="saving"
      @ok="handleSave"
    >
      <a-form :model="editForm" layout="vertical" class="pt-2">
        <a-form-item label="条件类型" required>
          <a-select v-model:value="editForm.condition_type">
            <a-select-option value="age">
              年龄范围
            </a-select-option>
            <a-select-option value="gender">
              性别
            </a-select-option>
            <a-select-option value="pregnancy">
              孕产妇
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="条件值" required>
          <a-input
            v-model:value="editForm.condition_value"
            :placeholder="editForm.condition_type === 'age' ? '如: 0-14 (儿童) 或 70+ (老年)' : editForm.condition_type === 'gender' ? 'male / female' : 'true'"
          />
        </a-form-item>
        <a-form-item label="适配动作" required>
          <a-select v-model:value="editForm.action">
            <a-select-option value="filter_device">
              过滤设备
            </a-select-option>
            <a-select-option value="filter_slot">
              过滤时段
            </a-select-option>
            <a-select-option value="filter_doctor">
              过滤医生
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="动作参数（key=value;key2=value2）">
          <a-textarea
            v-model:value="actionParamsStr"
            :rows="2"
            placeholder="如: device_ids=DEV001,DEV002;pool_type=department"
          />
        </a-form-item>
        <a-form-item label="优先级（数值越大越优先）">
          <a-input-number v-model:value="editForm.priority" :min="0" :max="100" style="width: 100%;" />
        </a-form-item>
        <a-form-item v-if="isEdit" label="状态">
          <a-select v-model:value="editForm.status">
            <a-select-option value="active">
              启用
            </a-select-option>
            <a-select-option value="inactive">
              停用
            </a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>
