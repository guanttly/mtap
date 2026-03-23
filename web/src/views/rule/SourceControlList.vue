<!-- 核心目的：开单来源控制管理页面 -->
<!-- 模块功能：规则引擎-控制门诊/住院/转诊来源的号源池分配比例与溢出策略 -->
<script setup lang="ts">
import type { SourceControl } from '@/types/rule'
import { message, Modal } from 'ant-design-vue'
import { h, onMounted, ref } from 'vue'
import { ruleApi } from '@/api/rule'

const loading = ref(false)
const items = ref<SourceControl[]>([])
const showModal = ref(false)
const saving = ref(false)
const isEdit = ref(false)
const editForm = ref<Partial<SourceControl>>({})

const SOURCE_TYPE_MAP: Record<string, string> = {
  outpatient: '门诊',
  inpatient: '住院',
  referral: '转诊',
}

const STATUS_COLOR: Record<string, string> = { active: 'success', inactive: 'default' }
const STATUS_LABEL: Record<string, string> = { active: '启用', inactive: '停用' }

const columns = [
  {
    title: '开单来源',
    dataIndex: 'source_type',
    key: 'source_type',
    customRender: ({ text }: { text: string }) =>
      h('a-tag', { color: 'blue' }, SOURCE_TYPE_MAP[text] ?? text),
  },
  { title: '绑定号源池 ID', dataIndex: 'slot_pool_id', key: 'slot_pool_id' },
  {
    title: '分配比例',
    dataIndex: 'allocation_ratio',
    key: 'allocation_ratio',
    customRender: ({ text }: { text: number }) => `${(text * 100).toFixed(0)}%`,
  },
  {
    title: '溢出启用',
    dataIndex: 'overflow_enabled',
    key: 'overflow_enabled',
    customRender: ({ text }: { text: boolean }) =>
      h('a-tag', { color: text ? 'orange' : 'default' }, text ? '是' : '否'),
  },
  {
    title: '溢出目标池 ID',
    dataIndex: 'overflow_target_pool_id',
    key: 'overflow_target_pool_id',
    customRender: ({ text }: { text: string }) => text || '-',
  },
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
    const res = await ruleApi.listSourceControls()
    items.value = (res.items ?? []) as SourceControl[]
  }
  finally { loading.value = false }
}

onMounted(() => fetchData())

function openCreate() {
  isEdit.value = false
  editForm.value = {
    source_type: 'outpatient',
    slot_pool_id: '',
    allocation_ratio: 1,
    overflow_enabled: false,
    overflow_target_pool_id: '',
  }
  showModal.value = true
}

function openEdit(record: SourceControl) {
  isEdit.value = true
  editForm.value = { ...record }
  showModal.value = true
}

async function handleSave() {
  if (!editForm.value.slot_pool_id?.trim()) {
    message.warning('请填写号源池 ID')
    return
  }
  const ratio = editForm.value.allocation_ratio ?? 0
  if (ratio < 0 || ratio > 1) {
    message.warning('分配比例须在 0~1 之间')
    return
  }
  saving.value = true
  try {
    await ruleApi.saveSourceControls([editForm.value])
    message.success(isEdit.value ? '更新成功' : '创建成功')
    showModal.value = false
    fetchData()
  }
  finally { saving.value = false }
}

function handleToggleStatus(record: SourceControl) {
  const next = record.status === 'active' ? 'inactive' : 'active'
  Modal.confirm({
    title: `确认${next === 'inactive' ? '停用' : '启用'}`,
    content: `确定要${next === 'inactive' ? '停用' : '启用'}「${SOURCE_TYPE_MAP[record.source_type]}」来源控制规则吗？`,
    onOk: async () => {
      await ruleApi.saveSourceControls([{ ...record, status: next }])
      message.success('操作成功')
      fetchData()
    },
  })
}
</script>

<template>
  <a-card class="list-card" :bordered="false">
    <template #title>开单来源控制</template>
    <template #extra>
      <a-space>
        <a-button @click="fetchData">刷新</a-button>
        <a-button type="primary" @click="openCreate">新增来源控制</a-button>
      </a-space>
    </template>

    <div class="list-toolbar">
      <a-alert
        type="info"
        show-icon
        message="将不同开单来源（门诊/住院/转诊）的号源预约流量，按比例路由到指定号源池，并可设置溢出目标池。"
        style="flex: 1;"
      />
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
            <a-button size="small" @click="openEdit(record as SourceControl)">
              编辑
            </a-button>
            <a-button
              size="small"
              :type="record.status === 'active' ? 'default' : 'primary'"
              @click="handleToggleStatus(record as SourceControl)"
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
      :title="isEdit ? '编辑来源控制' : '新增来源控制'"
      :confirm-loading="saving"
      @ok="handleSave"
    >
      <a-form :model="editForm" layout="vertical" class="pt-2">
        <a-form-item label="开单来源" required>
          <a-select v-model:value="editForm.source_type">
            <a-select-option value="outpatient">
              门诊
            </a-select-option>
            <a-select-option value="inpatient">
              住院
            </a-select-option>
            <a-select-option value="referral">
              转诊
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="绑定号源池 ID" required>
          <a-input v-model:value="editForm.slot_pool_id" placeholder="填写号源池 ID" />
        </a-form-item>
        <a-form-item label="分配比例（0~1，如 0.6 代表 60%）" required>
          <a-input-number
            v-model:value="editForm.allocation_ratio"
            :min="0"
            :max="1"
            :step="0.05"
            :precision="2"
            style="width: 100%;"
          />
        </a-form-item>
        <a-form-item label="允许溢出">
          <a-switch v-model:checked="editForm.overflow_enabled" />
          <span class="ml-2 text-gray-500 text-sm">开启后，超出配额的预约将自动路由到溢出目标池</span>
        </a-form-item>
        <a-form-item v-if="editForm.overflow_enabled" label="溢出目标号源池 ID">
          <a-input v-model:value="editForm.overflow_target_pool_id" placeholder="填写溢出目标号源池 ID" />
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
  </a-card>
</template>
