<!-- 核心目的：DependencyRuleList页面 -->
<!-- 模块功能：规则引擎-DependencyRule列表管理CRUD与配置 -->
<script setup lang="ts">
import { h, onMounted, ref } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { ruleApi } from '@/api/rule'
import { usePagination } from '@/composables/usePagination'
import type { DependencyRule } from '@/types/rule'

const { loading, items, pagination, fetchData, onTableChange } = usePagination<DependencyRule>(
  params => ruleApi.listDependencyRules(params),
)
onMounted(() => fetchData())

const showModal = ref(false)
const saving = ref(false)
const isEdit = ref(false)
const editId = ref('')
const editForm = ref<Partial<DependencyRule>>({ type: 'mandatory', validity_hours: 72 })

const TYPE_LABEL: Record<string, string> = { mandatory: '强制依赖', recommended: '建议依赖' }
const TYPE_COLOR: Record<string, string> = { mandatory: 'red', recommended: 'blue' }

const columns = [
  { title: '前置项目', dataIndex: 'pre_item_name', key: 'pre_item_name' },
  { title: '后续项目', dataIndex: 'post_item_name', key: 'post_item_name' },
  { title: '依赖类型', dataIndex: 'type', key: 'type', customRender: ({ text }: { text: string }) => h('a-tag', { color: TYPE_COLOR[text] ?? 'default' }, TYPE_LABEL[text] ?? text) },
  { title: '有效期(小时)', dataIndex: 'validity_hours', key: 'validity_hours' },
  { title: '状态', dataIndex: 'status', key: 'status', customRender: ({ text }: { text: string }) => h('a-badge', { status: text === 'active' ? 'success' : 'default', text: text === 'active' ? '启用' : '停用' }) },
  { title: '操作', key: 'actions', width: 140 },
]

function openCreate() {
  isEdit.value = false
  editId.value = ''
  editForm.value = { type: 'mandatory', validity_hours: 72 }
  showModal.value = true
}

function openEdit(record: DependencyRule) {
  isEdit.value = true
  editId.value = record.id
  editForm.value = { type: record.type, validity_hours: record.validity_hours, status: record.status }
  showModal.value = true
}

async function handleSave() {
  saving.value = true
  try {
    if (isEdit.value) {
      await ruleApi.updateDependencyRule(editId.value, {
        type: editForm.value.type,
        validity_hours: editForm.value.validity_hours,
        status: editForm.value.status,
      })
      message.success('更新成功')
    }
    else {
      if (!editForm.value.pre_item_id || !editForm.value.post_item_id) {
        message.warning('请填写前置项目ID和后续项目ID')
        return
      }
      await ruleApi.createDependencyRule(editForm.value)
      message.success('创建成功')
    }
    showModal.value = false
    fetchData()
  }
  finally { saving.value = false }
}

function handleDelete(record: DependencyRule) {
  Modal.confirm({
    title: '确认删除',
    content: `确定删除「${record.pre_item_name ?? record.pre_item_id} → ${record.post_item_name ?? record.post_item_id}」依赖规则吗？`,
    okType: 'danger',
    onOk: async () => {
      await ruleApi.deleteDependencyRule(record.id)
      message.success('删除成功')
      fetchData()
    },
  })
}
</script>

<template>
  <a-card class="list-card" :bordered="false">
    <template #title>依赖规则</template>
    <template #extra>
      <a-space>
        <a-button @click="fetchData">刷新</a-button>
        <a-button type="primary" @click="openCreate">新增依赖规则</a-button>
      </a-space>
    </template>
    <a-table :columns="columns" :data-source="items" :loading="loading" :pagination="pagination" row-key="id" size="middle" @change="onTableChange">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button type="link" size="small" @click="openEdit(record as DependencyRule)">编辑</a-button>
            <a-button type="link" danger size="small" @click="handleDelete(record as DependencyRule)">删除</a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="showModal"
      :title="isEdit ? '编辑依赖规则' : '新增依赖规则'"
      :confirm-loading="saving"
      @ok="handleSave"
    >
      <a-form :model="editForm" layout="vertical">
        <template v-if="!isEdit">
          <a-form-item label="前置项目ID" required>
            <a-input v-model:value="editForm.pre_item_id" placeholder="必须先完成该项目" />
          </a-form-item>
          <a-form-item label="后续项目ID" required>
            <a-input v-model:value="editForm.post_item_id" placeholder="依赖前置项目的项目" />
          </a-form-item>
        </template>
        <a-form-item label="依赖类型">
          <a-select v-model:value="editForm.type">
            <a-select-option value="mandatory">强制依赖（必须完成前置）</a-select-option>
            <a-select-option value="recommended">建议依赖（推荐先做前置）</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="有效期（小时）">
          <a-input-number v-model:value="editForm.validity_hours" :min="1" :max="8760" style="width:100%;" />
          <div style="color:#8c8c8c;font-size:12px;margin-top:4px;">前置检查结果在多少小时内有效（0表示永久有效）</div>
        </a-form-item>
        <a-form-item v-if="isEdit" label="状态">
          <a-radio-group v-model:value="editForm.status">
            <a-radio value="active">启用</a-radio>
            <a-radio value="inactive">停用</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>
