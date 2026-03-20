<!-- 核心目的：排序策略管理页面 -->
<!-- 模块功能：规则引擎-排序策略列表CRUD -->
<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { ruleApi } from '@/api/rule'
import type { SortingStrategy } from '@/types/rule'

const loading = ref(false)
const saving = ref(false)
const modalVisible = ref(false)
const strategies = ref<SortingStrategy[]>([])

const defaultForm = () => ({
  type: 'shortest_wait' as 'shortest_wait' | 'nearest' | 'priority',
  campusesText: '',
  deptsText: '',
  devicesText: '',
  start_date: '',
  end_date: '',
})
const form = reactive(defaultForm())

const columns = [
  { title: '策略类型', dataIndex: 'type', key: 'type', customRender: ({ text }: { text: string }) => ({ shortest_wait: '最短等待', nearest: '最近就近', priority: '优先级优先' }[text] ?? text) },
  { title: '院区范围', dataIndex: ['scope', 'campus_ids'], key: 'campus_ids', customRender: ({ value }: { value: string[] }) => value?.join(', ') || '全局' },
  { title: '科室范围', dataIndex: ['scope', 'department_ids'], key: 'dept_ids', customRender: ({ value }: { value: string[] }) => value?.join(', ') || '全局' },
  { title: '设备范围', dataIndex: ['scope', 'device_ids'], key: 'device_ids', customRender: ({ value }: { value: string[] }) => value?.join(', ') || '全局' },
  { title: '开始日期', dataIndex: 'start_date', key: 'start_date' },
  { title: '结束日期', dataIndex: 'end_date', key: 'end_date' },
  { title: '状态', dataIndex: 'status', key: 'status' },
  { title: '操作', key: 'action' },
]

async function fetchData() {
  loading.value = true
  try {
    const res = await ruleApi.listSortingStrategies()
    strategies.value = res?.items ?? []
  }
  catch {}
  finally { loading.value = false }
}

onMounted(fetchData)

function openAdd() {
  Object.assign(form, defaultForm())
  modalVisible.value = true
}

async function handleSave() {
  saving.value = true
  try {
    await ruleApi.saveSortingStrategy({
      type: form.type,
      scope: {
        campus_ids: form.campusesText.split(',').map(s => s.trim()).filter(Boolean),
        department_ids: form.deptsText.split(',').map(s => s.trim()).filter(Boolean),
        device_ids: form.devicesText.split(',').map(s => s.trim()).filter(Boolean),
      },
      start_date: form.start_date,
      end_date: form.end_date,
    })
    message.success('策略创建成功')
    modalVisible.value = false
    await fetchData()
  }
  catch {}
  finally { saving.value = false }
}

async function handleDelete(id: string) {
  try {
    await ruleApi.deleteSortingStrategy(id)
    message.success('策略已删除')
    await fetchData()
  }
  catch {}
}
</script>

<template>
  <div>
    <div style="margin-bottom: 16px; display: flex; justify-content: flex-end;">
      <a-button type="primary" @click="openAdd">＋ 新增策略</a-button>
    </div>
    <a-table
      :columns="columns"
      :data-source="strategies"
      :loading="loading"
      row-key="id"
      :pagination="false"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <a-popconfirm title="确认删除该策略？" ok-text="删除" cancel-text="取消" @confirm="handleDelete((record as SortingStrategy).id)">
            <a-button type="link" danger size="small">删除</a-button>
          </a-popconfirm>
        </template>
        <template v-else-if="column.key === 'type'">
          {{ { shortest_wait: '最短等待', nearest: '最近就近', priority: '优先级优先' }[(record as SortingStrategy).type] }}
        </template>
        <template v-else-if="column.key === 'campus_ids'">
          {{ (record as SortingStrategy).scope?.campus_ids?.join(', ') || '全局' }}
        </template>
        <template v-else-if="column.key === 'dept_ids'">
          {{ (record as SortingStrategy).scope?.department_ids?.join(', ') || '全局' }}
        </template>
        <template v-else-if="column.key === 'device_ids'">
          {{ (record as SortingStrategy).scope?.device_ids?.join(', ') || '全局' }}
        </template>
      </template>
    </a-table>

    <a-modal v-model:open="modalVisible" title="新增排序策略" :confirm-loading="saving" @ok="handleSave" ok-text="保存" cancel-text="取消">
      <a-form :model="form" layout="vertical">
        <a-form-item label="策略类型">
          <a-radio-group v-model:value="form.type">
            <a-radio value="shortest_wait">最短等待</a-radio>
            <a-radio value="nearest">最近就近</a-radio>
            <a-radio value="priority">优先级优先</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="生效日期范围">
          <div style="display: flex; gap: 8px; align-items: center;">
            <a-date-picker v-model:value="form.start_date" value-format="YYYY-MM-DD" placeholder="开始日期" style="flex: 1;" />
            <span>—</span>
            <a-date-picker v-model:value="form.end_date" value-format="YYYY-MM-DD" placeholder="结束日期" style="flex: 1;" />
          </div>
        </a-form-item>
        <a-card title="生效范围（留空表示全局）" size="small">
          <a-form-item label="院区（逗号分隔ID）">
            <a-input v-model:value="form.campusesText" placeholder="campus_id1,campus_id2" />
          </a-form-item>
          <a-form-item label="科室（逗号分隔ID）">
            <a-input v-model:value="form.deptsText" placeholder="dept_id1,dept_id2" />
          </a-form-item>
          <a-form-item label="设备（逗号分隔ID）">
            <a-input v-model:value="form.devicesText" placeholder="device_id1,device_id2" />
          </a-form-item>
        </a-card>
      </a-form>
    </a-modal>
  </div>
</template>
