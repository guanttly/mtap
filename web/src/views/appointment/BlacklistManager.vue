<!-- 核心目的：BlacklistManager页面 -->
<!-- 模块功能：预约服务-相关功能页面 -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { appointmentApi } from '@/api/appointment'
import { usePagination } from '@/composables/usePagination'
import type { BlacklistRecord } from '@/types/appointment'

const { loading, items, pagination, fetchData, onTableChange } = usePagination<BlacklistRecord>(
  params => appointmentApi.listBlacklist(params),
)
onMounted(() => fetchData())

const appealModal = ref(false)
const appealTarget = ref<BlacklistRecord | null>(null)
const appealReason = ref('')
const appealing = ref(false)

function openAppeal(record: BlacklistRecord) {
  appealTarget.value = record
  appealReason.value = ''
  appealModal.value = true
}

async function submitAppeal() {
  if (!appealTarget.value || !appealReason.value.trim()) return
  appealing.value = true
  try {
    await appointmentApi.submitAppeal(appealTarget.value.id, appealReason.value)
    message.success('申诉已提交')
    appealModal.value = false
    fetchData()
  }
  finally { appealing.value = false }
}

async function handleRemove(record: BlacklistRecord) {
  await appointmentApi.removeFromBlacklist(record.id)
  message.success('已移出黑名单')
  fetchData()
}

const columns = [
  { title: '患者', dataIndex: 'patient_name', key: 'patient_name' },
  { title: '加入原因', dataIndex: 'reason', key: 'reason' },
  { title: '爽约次数', dataIndex: 'no_show_count', key: 'no_show_count' },
  { title: '加入时间', dataIndex: 'added_at', key: 'added_at' },
  { title: '状态', dataIndex: 'status', key: 'status' },
  { title: '操作', key: 'actions' },
]
</script>

<template>
  <div>
    <a-table :columns="columns" :data-source="items" :loading="loading" :pagination="pagination" row-key="id" size="middle" @change="onTableChange">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag :color="(record as BlacklistRecord).status === 'active' ? 'red' : 'default'">
            {{ (record as BlacklistRecord).status === 'active' ? '黑名单中' : '已解除' }}
          </a-tag>
        </template>
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button type="link" size="small" @click="openAppeal(record as BlacklistRecord)">提交申诉</a-button>
            <a-button type="link" danger size="small" @click="handleRemove(record as BlacklistRecord)">移除</a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <a-modal v-model:open="appealModal" title="提交申诉" :confirm-loading="appealing" @ok="submitAppeal">
      <a-form layout="vertical">
        <a-form-item label="申诉理由">
          <a-textarea v-model:value="appealReason" :rows="4" placeholder="请详细说明申诉理由..." />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>
