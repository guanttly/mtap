<!-- 核心目的：设备管理页面 -->
<!-- 模块功能：设备CRUD、状态管理、支持检查类型配置 -->
<script setup lang="ts">
import type { Device } from '@/types/resource'
import { message, Modal } from 'ant-design-vue'
import { onMounted, ref } from 'vue'
import { resourceApi } from '@/api/resource'
import { usePagination } from '@/composables/usePagination'

const { loading, items, pagination, fetchData, onTableChange } = usePagination<Device>(
  params => resourceApi.listDevices(params),
)
onMounted(() => fetchData())

const STATUS_COLOR: Record<string, string> = { online: 'success', offline: 'default', maintenance: 'warning' }
const STATUS_LABEL: Record<string, string> = { online: '在线', offline: '离线', maintenance: '维护中' }

const columns = [
  { title: '设备名称', dataIndex: 'name', key: 'name' },
  { title: '型号', dataIndex: 'model', key: 'model' },
  { title: '厂商', dataIndex: 'manufacturer', key: 'manufacturer' },
  { title: '支持检查类型', key: 'exam_types' },
  { title: '每日最大号位', dataIndex: 'max_daily_slots', key: 'max_daily_slots' },
  { title: '状态', dataIndex: 'status', key: 'status' },
  { title: '操作', key: 'actions' },
]

const showModal = ref(false)
const saving = ref(false)
const editForm = ref<Partial<Device>>({ status: 'online', supported_exam_types: [] })
const examTypesText = ref('')

function openCreate() {
  editForm.value = { status: 'online', supported_exam_types: [], max_daily_slots: 20 }
  examTypesText.value = ''
  showModal.value = true
}

function openEdit(record: Device) {
  editForm.value = { ...record }
  examTypesText.value = (record.supported_exam_types ?? []).join(',')
  showModal.value = true
}

async function handleSave() {
  if (!editForm.value.name?.trim()) {
    message.warning('请填写设备名称')
    return
  }
  editForm.value.supported_exam_types = examTypesText.value
    .split(',')
    .map(s => s.trim())
    .filter(Boolean)

  saving.value = true
  try {
    if (editForm.value.id) {
      await resourceApi.updateDevice(editForm.value.id, editForm.value)
      message.success('更新成功')
    }
    else {
      await resourceApi.createDevice(editForm.value)
      message.success('创建成功')
    }
    showModal.value = false
    fetchData()
  }
  finally { saving.value = false }
}

function handleDelete(record: Device) {
  Modal.confirm({
    title: '确认删除',
    content: `确定删除设备「${record.name}」吗？该操作将同时影响关联排班。`,
    okType: 'danger',
    onOk: async () => {
      await resourceApi.deleteDevice(record.id)
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
        新增设备
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
        <template v-if="column.key === 'exam_types'">
          <a-space wrap>
            <a-tag v-for="t in (record as Device).supported_exam_types" :key="t" color="blue">
              {{ t }}
            </a-tag>
            <span v-if="!(record as Device).supported_exam_types?.length" style="color: #bfbfbf;">—</span>
          </a-space>
        </template>
        <template v-else-if="column.key === 'status'">
          <a-badge :status="STATUS_COLOR[(record as Device).status] as any" :text="STATUS_LABEL[(record as Device).status]" />
        </template>
        <template v-else-if="column.key === 'actions'">
          <a-space>
            <a-button type="link" size="small" @click="openEdit(record as Device)">编辑</a-button>
            <a-button type="link" danger size="small" @click="handleDelete(record as Device)">删除</a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="showModal"
      :title="editForm.id ? '编辑设备' : '新增设备'"
      :confirm-loading="saving"
      @ok="handleSave"
    >
      <a-form :model="editForm" layout="vertical">
        <a-form-item label="设备名称" required>
          <a-input v-model:value="editForm.name" placeholder="如：MRI-西门子Vida 3T" />
        </a-form-item>
        <a-form-item label="型号">
          <a-input v-model:value="editForm.model" />
        </a-form-item>
        <a-form-item label="厂商">
          <a-input v-model:value="editForm.manufacturer" />
        </a-form-item>
        <a-form-item label="每日最大号位数">
          <a-input-number v-model:value="editForm.max_daily_slots" :min="1" :max="500" style="width: 100%;" />
        </a-form-item>
        <a-form-item label="支持检查类型（逗号分隔）">
          <a-textarea v-model:value="examTypesText" :rows="2" placeholder="如：MRI平扫,MRI增强,MRI功能成像" />
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="editForm.status">
            <a-select-option v-for="(label, val) in STATUS_LABEL" :key="val" :value="val">
              {{ label }}
            </a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>
