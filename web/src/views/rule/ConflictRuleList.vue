<!-- 核心目的：ConflictRuleList页面 -->
<!-- 模块功能：规则引擎-ConflictRule列表管理CRUD与配置 -->
<script setup lang="ts">
import { h, onMounted, ref } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { ruleApi } from '@/api/rule'
import { usePagination } from '@/composables/usePagination'
import type { ConflictRule } from '@/types/rule'

const { loading, items, pagination, fetchData, onTableChange } = usePagination<ConflictRule>(
  params => ruleApi.listConflictRules(params),
)

onMounted(() => fetchData())

const showModal = ref(false)
const saving = ref(false)
const isEdit = ref(false)
const editId = ref('')
const editForm = ref<Partial<ConflictRule>>({})

const LEVEL_MAP: Record<string, string> = { forbid: '禁止', warning: '警告' }
const LEVEL_COLOR: Record<string, string> = { forbid: 'red', warning: 'orange' }
const STATUS_MAP: Record<string, { color: string, label: string }> = {
  active: { color: 'success', label: '启用' },
  inactive: { color: 'default', label: '停用' },
}

const columns = [
  { title: '检查项A', dataIndex: 'item_a_name', key: 'item_a_name' },
  { title: '检查项B', dataIndex: 'item_b_name', key: 'item_b_name' },
  { title: '最短间隔', key: 'interval', customRender: ({ record }: { record: ConflictRule }) => `${record.min_interval} ${record.interval_unit === 'hour' ? '小时' : '天'}` },
  { title: '冲突级别', dataIndex: 'level', key: 'level', customRender: ({ text }: { text: string }) => h('a-tag', { color: LEVEL_COLOR[text] ?? 'default' }, LEVEL_MAP[text] ?? text) },
  { title: '状态', dataIndex: 'status', key: 'status', customRender: ({ text }: { text: string }) => h('a-badge', { status: STATUS_MAP[text]?.color ?? 'default', text: STATUS_MAP[text]?.label ?? text }) },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at', customRender: ({ text }: { text: string }) => text ? text.slice(0, 10) : '' },
  { title: '操作', key: 'actions', width: 160 },
]

function openCreate() {
  isEdit.value = false
  editId.value = ''
  editForm.value = { level: 'warning', interval_unit: 'hour', min_interval: 24 }
  showModal.value = true
}

function openEdit(record: ConflictRule) {
  isEdit.value = true
  editId.value = record.id
  editForm.value = { min_interval: record.min_interval, interval_unit: record.interval_unit, level: record.level, status: record.status }
  showModal.value = true
}

async function handleSave() {
  saving.value = true
  try {
    if (isEdit.value) {
      await ruleApi.updateConflictRule(editId.value, {
        min_interval: editForm.value.min_interval,
        level: editForm.value.level,
        status: editForm.value.status,
      })
      message.success('更新成功')
    }
    else {
      await ruleApi.createConflictRule(editForm.value)
      message.success('创建成功')
    }
    showModal.value = false
    fetchData()
  }
  finally { saving.value = false }
}

function handleToggleStatus(record: ConflictRule) {
  const next = record.status === 'active' ? 'inactive' : 'active'
  Modal.confirm({
    title: `确认${next === 'inactive' ? '停用' : '启用'}`,
    content: `确定要${next === 'inactive' ? '停用' : '启用'}该冲突规则吗？`,
    onOk: async () => {
      await ruleApi.updateConflictRule(record.id, { status: next })
      message.success('操作成功')
      fetchData()
    },
  })
}

function handleDelete(record: ConflictRule) {
  Modal.confirm({
    title: '确认删除',
    content: `确定删除「${record.item_a_name ?? record.item_a_id} ↔ ${record.item_b_name ?? record.item_b_id}」冲突规则吗？`,
    okType: 'danger',
    onOk: async () => {
      await ruleApi.deleteConflictRule(record.id)
      message.success('删除成功')
      fetchData()
    },
  })
}
</script>

<template>
  <div>
    <div class="mb-4 flex items-center gap-2">
      <a-button type="primary" @click="openCreate">
        新增冲突规则
      </a-button>
      <a-button @click="fetchData">
        刷新
      </a-button>
    </div>
    <a-table
      :columns="columns"
      :data-source="items"
      :loading="loading"
      :pagination="pagination"
      row-key="id"
      size="middle"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button type="link" size="small" @click="openEdit(record as ConflictRule)">编辑</a-button>
            <a-button type="link" size="small" @click="handleToggleStatus(record as ConflictRule)">
              {{ (record as ConflictRule).status === 'active' ? '停用' : '启用' }}
            </a-button>
            <a-button type="link" danger size="small" @click="handleDelete(record as ConflictRule)">删除</a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="showModal"
      :title="isEdit ? '编辑冲突规则' : '新增冲突规则'"
      :confirm-loading="saving"
      @ok="handleSave"
    >
      <a-form :model="editForm" layout="vertical">
        <template v-if="!isEdit">
          <a-form-item label="检查项A ID" required>
            <a-input v-model:value="editForm.item_a_id" placeholder="请输入检查项A的ID" />
          </a-form-item>
          <a-form-item label="检查项B ID" required>
            <a-input v-model:value="editForm.item_b_id" placeholder="请输入检查项B的ID" />
          </a-form-item>
        </template>
        <a-form-item label="最短间隔">
          <a-input-number v-model:value="editForm.min_interval" :min="0" :max="720" style="width: 60%;" />
          <a-select v-model:value="editForm.interval_unit" style="width: 38%; margin-left: 2%;">
            <a-select-option value="hour">小时</a-select-option>
            <a-select-option value="day">天</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="冲突级别">
          <a-radio-group v-model:value="editForm.level">
            <a-radio value="warning">警告级（可跳过）</a-radio>
            <a-radio value="forbid">禁止级（不可跳过）</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item v-if="isEdit" label="状态">
          <a-radio-group v-model:value="editForm.status">
            <a-radio value="active">启用</a-radio>
            <a-radio value="inactive">停用</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>
