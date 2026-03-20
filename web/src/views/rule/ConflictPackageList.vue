<!-- 核心目的：ConflictPackageList页面 -->
<!-- 模块功能：规则引擎-ConflictPackage列表管理CRUD与配置 -->
<script setup lang="ts">
import { h, onMounted, ref } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { ruleApi } from '@/api/rule'
import { usePagination } from '@/composables/usePagination'
import type { ConflictPackage } from '@/types/rule'

const { loading, items, pagination, fetchData, onTableChange } = usePagination<ConflictPackage>(
  params => ruleApi.listConflictPackages(params),
)
onMounted(() => fetchData())

const showModal = ref(false)
const saving = ref(false)
const isEdit = ref(false)
const editId = ref('')
const editForm = ref<Partial<ConflictPackage>>({ level: 'warning', interval_unit: 'hour', items: [] })
const itemIdsText = ref('')

const LEVEL_MAP: Record<string, string> = { forbid: '禁止', warning: '警告' }
const LEVEL_COLOR: Record<string, string> = { forbid: 'red', warning: 'orange' }

const columns = [
  { title: '包名称', dataIndex: 'name', key: 'name' },
  { title: '最短间隔', key: 'interval', customRender: ({ record }: { record: ConflictPackage }) => `${record.min_interval} ${record.interval_unit === 'hour' ? '小时' : '天'}` },
  { title: '冲突级别', dataIndex: 'level', key: 'level', customRender: ({ text }: { text: string }) => h('a-tag', { color: LEVEL_COLOR[text] ?? 'default' }, LEVEL_MAP[text] ?? text) },
  { title: '包含项目数', key: 'count', customRender: ({ record }: { record: ConflictPackage }) => record.items?.length ?? 0 },
  { title: '状态', dataIndex: 'status', key: 'status', customRender: ({ text }: { text: string }) => h('a-badge', { status: text === 'active' ? 'success' : 'default', text: text === 'active' ? '启用' : '停用' }) },
  { title: '操作', key: 'actions', width: 140 },
]

function openCreate() {
  isEdit.value = false
  editId.value = ''
  editForm.value = { level: 'warning', interval_unit: 'hour', items: [] }
  itemIdsText.value = ''
  showModal.value = true
}

function openEdit(record: ConflictPackage) {
  isEdit.value = true
  editId.value = record.id
  editForm.value = { name: record.name, min_interval: record.min_interval, interval_unit: record.interval_unit, level: record.level }
  itemIdsText.value = (record.items ?? []).map(it => it.exam_item_id).join(',')
  showModal.value = true
}

async function handleSave() {
  if (!editForm.value.name?.trim()) {
    message.warning('请填写包名称')
    return
  }
  const itemIds = itemIdsText.value.split(',').map(s => s.trim()).filter(Boolean)
  if (itemIds.length < 2) {
    message.warning('冲突包至少需要2个检查项目')
    return
  }
  saving.value = true
  try {
    if (isEdit.value) {
      await ruleApi.updateConflictPackage(editId.value, {
        name: editForm.value.name,
        item_ids: itemIds,
        min_interval: editForm.value.min_interval,
        level: editForm.value.level,
      })
      message.success('更新成功')
    }
    else {
      await ruleApi.createConflictPackage({
        name: editForm.value.name!,
        item_ids: itemIds,
        min_interval: editForm.value.min_interval,
        interval_unit: editForm.value.interval_unit,
        level: editForm.value.level!,
      })
      message.success('创建成功')
    }
    showModal.value = false
    fetchData()
  }
  finally { saving.value = false }
}

function handleDelete(record: ConflictPackage) {
  Modal.confirm({
    title: '确认删除',
    content: `确定删除冲突包「${record.name}」吗？`,
    okType: 'danger',
    onOk: async () => {
      await ruleApi.deleteConflictPackage(record.id)
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
        新增冲突包
      </a-button>
      <a-button @click="fetchData">
        刷新
      </a-button>
    </div>
    <a-table :columns="columns" :data-source="items" :loading="loading" :pagination="pagination" row-key="id" size="middle" @change="onTableChange">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button type="link" size="small" @click="openEdit(record as ConflictPackage)">编辑</a-button>
            <a-button type="link" danger size="small" @click="handleDelete(record as ConflictPackage)">删除</a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="showModal"
      :title="isEdit ? '编辑冲突包' : '新增冲突包'"
      :confirm-loading="saving"
      @ok="handleSave"
    >
      <a-form :model="editForm" layout="vertical">
        <a-form-item label="包名称" required>
          <a-input v-model:value="editForm.name" placeholder="请输入冲突包名称" :maxlength="30" show-count />
        </a-form-item>
        <a-form-item label="最短间隔">
          <a-input-number v-model:value="editForm.min_interval" :min="0" :max="720" style="width:60%;" />
          <a-select v-model:value="editForm.interval_unit" style="width:38%; margin-left:2%;">
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
        <a-form-item label="包含检查项ID（逗号分隔，至少2个）" required>
          <a-textarea
            v-model:value="itemIdsText"
            :rows="3"
            placeholder="exam_item_id1,exam_item_id2,exam_item_id3..."
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>
