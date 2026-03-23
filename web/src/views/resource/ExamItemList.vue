<!-- 核心目的：检查项目管理页面 -->
<!-- 模块功能：检查项目CRUD、空腹标记管理 -->
<script setup lang="ts">
import type { ExamItem } from '@/types/resource'
import { message, Modal } from 'ant-design-vue'
import { h, onMounted, ref } from 'vue'
import { resourceApi } from '@/api/resource'
import { usePagination } from '@/composables/usePagination'

const { loading, items, pagination, fetchData, onTableChange } = usePagination<ExamItem>(
  params => resourceApi.listExamItems(params),
)

onMounted(() => fetchData())

const showModal = ref(false)
const saving = ref(false)
const editForm = ref<Partial<ExamItem & { id?: string }>>({})
const isEdit = ref(false)

const columns = [
  { title: '项目名称', dataIndex: 'name', key: 'name' },
  {
    title: '标准时长(分钟)',
    dataIndex: 'duration_min',
    key: 'duration_min',
  },
  {
    title: '空腹要求',
    dataIndex: 'is_fasting',
    key: 'is_fasting',
    customRender: ({ text }: { text: boolean }) =>
      h('a-tag', { color: text ? 'orange' : 'default' }, text ? '需空腹' : '无要求'),
  },
  { title: '空腹说明', dataIndex: 'fasting_desc', key: 'fasting_desc' },
  { title: '操作', key: 'actions' },
]

function openCreate() {
  editForm.value = { is_fasting: false, duration_min: 30, fasting_desc: '' }
  isEdit.value = false
  showModal.value = true
}

function openEdit(record: ExamItem) {
  editForm.value = { ...record }
  isEdit.value = true
  showModal.value = true
}

async function handleSave() {
  if (!editForm.value.name?.trim()) {
    message.warning('请填写项目名称')
    return
  }
  saving.value = true
  try {
    if (isEdit.value && (editForm.value as any).id) {
      await resourceApi.updateExamItem((editForm.value as any).id, editForm.value)
      message.success('更新成功')
    }
    else {
      await resourceApi.createExamItem(editForm.value)
      message.success('创建成功')
    }
    showModal.value = false
    fetchData()
  }
  finally {
    saving.value = false
  }
}

function handleDelete(record: ExamItem) {
  Modal.confirm({
    title: '确认删除',
    content: `确定删除检查项目「${record.name}」吗？此操作不可恢复。`,
    okType: 'danger',
    onOk: async () => {
      await resourceApi.deleteExamItem(record.id)
      message.success('删除成功')
      fetchData()
    },
  })
}
</script>

<template>
  <a-card class="list-card" :bordered="false">
    <template #title>检查项目</template>
    <template #extra>
      <a-space>
        <a-button @click="fetchData">刷新</a-button>
        <a-button type="primary" @click="openCreate">新增检查项目</a-button>
      </a-space>
    </template>

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
            <a-button type="link" size="small" @click="openEdit(record as ExamItem)">
              编辑
            </a-button>
            <a-button type="link" danger size="small" @click="handleDelete(record as ExamItem)">
              删除
            </a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="showModal"
      :title="isEdit ? '编辑检查项目' : '新增检查项目'"
      :confirm-loading="saving"
      @ok="handleSave"
    >
      <a-form :model="editForm" layout="vertical">
        <a-form-item label="项目名称" required>
          <a-input v-model:value="editForm.name" placeholder="请输入检查项目名称" />
        </a-form-item>
        <a-form-item label="标准时长（分钟）">
          <a-input-number
            v-model:value="editForm.duration_min"
            :min="1"
            :max="480"
            style="width: 100%"
          />
        </a-form-item>
        <a-form-item label="空腹要求">
          <a-switch
            v-model:checked="editForm.is_fasting"
            checked-children="需空腹"
            un-checked-children="无要求"
          />
        </a-form-item>
        <a-form-item v-if="editForm.is_fasting" label="空腹说明">
          <a-textarea
            v-model:value="editForm.fasting_desc"
            :rows="3"
            placeholder="请描述空腹要求（如：检查前8小时禁食禁水）"
            :maxlength="200"
            show-count
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>
