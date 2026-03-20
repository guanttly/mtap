<!-- 核心目的：排班日历管理页面 -->
<!-- 模块功能：排班创建、批量生成、替班、暂停 -->
<script setup lang="ts">
import type { Device, Schedule } from '@/types/resource'
import { message } from 'ant-design-vue'
import { onMounted, ref } from 'vue'
import { resourceApi } from '@/api/resource'
import { usePagination } from '@/composables/usePagination'

// ---- 排班列表 ----
const filterDate = ref<string>('')
const { loading, items, pagination, fetchData, onTableChange } = usePagination<Schedule>(
  params => resourceApi.listSchedules({ ...params, date: filterDate.value || undefined }),
)
onMounted(fetchData)

const STATUS_COLOR: Record<string, string> = { normal: 'success', suspended: 'warning', substitute: 'processing' }
const STATUS_LABEL: Record<string, string> = { normal: '正常', suspended: '已暂停', substitute: '替班' }

const columns = [
  { title: '设备', dataIndex: 'device_name', key: 'device_name' },
  { title: '工作日期', dataIndex: 'work_date', key: 'work_date' },
  { title: '时间段', key: 'time' },
  { title: '时隙时长(分)', dataIndex: 'slot_minutes', key: 'slot_minutes' },
  { title: '状态', dataIndex: 'status', key: 'status' },
  { title: '操作', key: 'actions' },
]

async function suspend(record: Schedule) {
  await resourceApi.suspendSchedule(record.device_id, record.work_date, record.start_time, record.end_time, '手动暂停')
  message.success('已暂停该排班')
  fetchData()
}

// ---- 设备列表（供下拉选择） ----
const devices = ref<Device[]>([])
async function loadDevices() {
  const res = await resourceApi.listDevices({ page: 1, page_size: 200 })
  devices.value = res.items ?? []
}
onMounted(loadDevices)

// ---- 新增单条排班 ----
const showAddModal = ref(false)
const addSaving = ref(false)
const addForm = ref<Partial<Schedule>>({ status: 'normal' })

function openAdd() {
  addForm.value = { status: 'normal' }
  showAddModal.value = true
}
async function handleAdd() {
  if (!addForm.value.device_id || !addForm.value.work_date) {
    message.warning('请选择设备和工作日期')
    return
  }
  addSaving.value = true
  try {
    await resourceApi.createSchedule(addForm.value)
    message.success('排班创建成功')
    showAddModal.value = false
    fetchData()
  }
  finally { addSaving.value = false }
}

// ---- 批量生成排班 ----
const showGenModal = ref(false)
const genSaving = ref(false)
const genForm = ref({
  device_id: '',
  start_date: '',
  end_date: '',
  start_time: '08:00',
  end_time: '17:00',
  slot_minutes: 30,
  skip_weekends: false,
})

function openGenerate() {
  genForm.value = { device_id: '', start_date: '', end_date: '', start_time: '08:00', end_time: '17:00', slot_minutes: 30, skip_weekends: false }
  showGenModal.value = true
}
async function handleGenerate() {
  if (!genForm.value.device_id || !genForm.value.start_date || !genForm.value.end_date) {
    message.warning('请填写设备、开始日期和结束日期')
    return
  }
  genSaving.value = true
  try {
    await resourceApi.generateSchedule(genForm.value)
    message.success('批量生成成功')
    showGenModal.value = false
    fetchData()
  }
  finally { genSaving.value = false }
}

// ---- 替班 ----
const showSubModal = ref(false)
const subSaving = ref(false)
const subForm = ref({ source_device_id: '', target_device_id: '', date: '' })

function openSubstitute() {
  subForm.value = { source_device_id: '', target_device_id: '', date: '' }
  showSubModal.value = true
}
async function handleSubstitute() {
  if (!subForm.value.source_device_id || !subForm.value.target_device_id || !subForm.value.date) {
    message.warning('请填写所有替班信息')
    return
  }
  subSaving.value = true
  try {
    await resourceApi.substituteSchedule(subForm.value.source_device_id, subForm.value.target_device_id, subForm.value.date)
    message.success('替班操作成功')
    showSubModal.value = false
    fetchData()
  }
  finally { subSaving.value = false }
}
</script>

<template>
  <div>
    <!-- 操作栏 -->
    <div class="mb-4 flex flex-wrap items-center gap-2">
      <a-button type="primary" @click="openAdd">新增排班</a-button>
      <a-button @click="openGenerate">批量生成</a-button>
      <a-button @click="openSubstitute">替班管理</a-button>
      <a-date-picker
        v-model:value="filterDate"
        placeholder="按日期筛选"
        value-format="YYYY-MM-DD"
        style="width: 160px;"
        allow-clear
        @change="fetchData"
      />
      <a-button @click="fetchData">刷新</a-button>
    </div>

    <!-- 排班列表 -->
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
        <template v-if="column.key === 'time'">
          {{ (record as Schedule).start_time }} — {{ (record as Schedule).end_time }}
        </template>
        <template v-else-if="column.key === 'status'">
          <a-badge :status="STATUS_COLOR[(record as Schedule).status] as any" :text="STATUS_LABEL[(record as Schedule).status]" />
        </template>
        <template v-else-if="column.key === 'actions'">
          <a-button
            v-if="(record as Schedule).status === 'normal'"
            type="link"
            danger
            size="small"
            @click="suspend(record as Schedule)"
          >
            暂停
          </a-button>
        </template>
      </template>
    </a-table>

    <!-- 新增排班弹窗 -->
    <a-modal v-model:open="showAddModal" title="新增单条排班" :confirm-loading="addSaving" @ok="handleAdd">
      <a-form :model="addForm" layout="vertical">
        <a-form-item label="设备" required>
          <a-select v-model:value="addForm.device_id" placeholder="请选择设备" show-search option-filter-prop="label">
            <a-select-option v-for="d in devices" :key="d.id" :value="d.id" :label="d.name">{{ d.name }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="工作日期" required>
          <a-date-picker v-model:value="addForm.work_date" value-format="YYYY-MM-DD" style="width: 100%;" />
        </a-form-item>
        <a-form-item label="工作时间">
          <a-time-picker v-model:value="addForm.start_time" format="HH:mm" value-format="HH:mm" placeholder="开始" style="width: 48%;" />
          <span style="margin: 0 4%;">—</span>
          <a-time-picker v-model:value="addForm.end_time" format="HH:mm" value-format="HH:mm" placeholder="结束" style="width: 48%;" />
        </a-form-item>
        <a-form-item label="时隙时长（分钟）">
          <a-input-number v-model:value="addForm.slot_minutes" :min="5" :max="120" style="width: 100%;" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 批量生成排班弹窗 -->
    <a-modal v-model:open="showGenModal" title="批量生成排班" :confirm-loading="genSaving" ok-text="立即生成" @ok="handleGenerate">
      <a-form :model="genForm" layout="vertical">
        <a-form-item label="设备" required>
          <a-select v-model:value="genForm.device_id" placeholder="请选择设备" show-search option-filter-prop="label">
            <a-select-option v-for="d in devices" :key="d.id" :value="d.id" :label="d.name">{{ d.name }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="日期范围" required>
          <a-date-picker v-model:value="genForm.start_date" value-format="YYYY-MM-DD" placeholder="开始日期" style="width: 48%;" />
          <span style="margin: 0 4%;">—</span>
          <a-date-picker v-model:value="genForm.end_date" value-format="YYYY-MM-DD" placeholder="结束日期" style="width: 48%;" />
        </a-form-item>
        <a-form-item label="每日工作时间">
          <a-time-picker v-model:value="genForm.start_time" format="HH:mm" value-format="HH:mm" placeholder="开始" style="width: 48%;" />
          <span style="margin: 0 4%;">—</span>
          <a-time-picker v-model:value="genForm.end_time" format="HH:mm" value-format="HH:mm" placeholder="结束" style="width: 48%;" />
        </a-form-item>
        <a-form-item label="时隙时长（分钟）">
          <a-input-number v-model:value="genForm.slot_minutes" :min="5" :max="120" style="width: 100%;" />
        </a-form-item>
        <a-form-item label="跳过周末">
          <a-switch v-model:checked="genForm.skip_weekends" checked-children="是" un-checked-children="否" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 替班弹窗 -->
    <a-modal v-model:open="showSubModal" title="替班管理" :confirm-loading="subSaving" ok-text="确认替班" @ok="handleSubstitute">
      <a-alert message="替班将把指定日期原设备的所有预约转移给替班设备" type="info" show-icon class="mb-4" />
      <a-form :model="subForm" layout="vertical">
        <a-form-item label="原设备" required>
          <a-select v-model:value="subForm.source_device_id" placeholder="选择原设备" show-search option-filter-prop="label">
            <a-select-option v-for="d in devices" :key="d.id" :value="d.id" :label="d.name">{{ d.name }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="替班设备" required>
          <a-select v-model:value="subForm.target_device_id" placeholder="选择替班设备" show-search option-filter-prop="label">
            <a-select-option v-for="d in devices" :key="d.id" :value="d.id" :label="d.name">{{ d.name }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="替班日期" required>
          <a-date-picker v-model:value="subForm.date" value-format="YYYY-MM-DD" style="width: 100%;" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>
