<!-- 核心目的：BottleneckAlerts页面 -->
<!-- 模块功能：智能效能优化-相关功能页面 -->
<!-- 核心目的：BottleneckAlerts页面 -->
<!-- 模块功能：智能效能优化-相关功能页面 -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { optimizationApi } from '@/api/optimization'
import { usePagination } from '@/composables/usePagination'
import type { BottleneckAlert } from '@/types/optimization'

const { loading, items, pagination, fetchData, onTableChange } = usePagination<BottleneckAlert>(
  params => optimizationApi.listAlerts(params),
)
onMounted(() => fetchData())

const dismissModal = ref(false)
const dismissTarget = ref<BottleneckAlert | null>(null)
const dismissReason = ref('')
const dismissing = ref(false)

function openDismiss(record: BottleneckAlert) {
  dismissTarget.value = record
  dismissReason.value = ''
  dismissModal.value = true
}

async function handleDismiss() {
  if (!dismissTarget.value) return
  dismissing.value = true
  try {
    await optimizationApi.dismissAlert(dismissTarget.value.id, dismissReason.value)
    message.success('告警已确认')
    dismissModal.value = false
    fetchData()
  }
  finally { dismissing.value = false }
}

const SEVERITY_COLOR: Record<string, string> = { critical: 'red', high: 'orange', medium: 'gold', low: 'blue' }

const columns = [
  { title: '告警类型', dataIndex: 'alert_type', key: 'type' },
  { title: '严重程度', dataIndex: 'severity', key: 'severity' },
  { title: '摘要', dataIndex: 'summary', key: 'summary' },
  { title: '设备', dataIndex: 'device_name', key: 'device' },
  { title: '当前值', dataIndex: 'current_value', key: 'value' },
  { title: '阈值', dataIndex: 'threshold_value', key: 'threshold' },
  { title: '状态', dataIndex: 'status', key: 'status' },
  { title: '操作', key: 'actions' },
]
</script>

<template>
  <div>
    <a-table :columns="columns" :data-source="items" :loading="loading" :pagination="pagination" row-key="id" size="middle" @change="onTableChange">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'severity'">
          <a-tag :color="SEVERITY_COLOR[(record as BottleneckAlert).severity]">{{ (record as BottleneckAlert).severity }}</a-tag>
        </template>
        <template v-if="column.key === 'status'">
          <a-badge :status="(record as BottleneckAlert).status === 'open' ? 'error' : 'default'" :text="(record as BottleneckAlert).status === 'open' ? '待处理' : '已处理'" />
        </template>
        <template v-if="column.key === 'actions'">
          <a-button v-if="(record as BottleneckAlert).status === 'open'" type="link" size="small" @click="openDismiss(record as BottleneckAlert)">确认</a-button>
        </template>
      </template>
    </a-table>

    <a-modal v-model:open="dismissModal" title="确认告警" :confirm-loading="dismissing" @ok="handleDismiss">
      <a-form layout="vertical">
        <a-form-item label="处理说明">
          <a-textarea v-model:value="dismissReason" :rows="4" placeholder="请说明处理措施..." />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>
