<!-- 核心目的：NurseCallPanel页面 -->
<!-- 模块功能：分诊管理-相关功能页面 -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { triageApi } from '@/api/triage'
import type { QueueStatus, CallResult } from '@/types/triage'

const roomId = ref('room-001')
const status = ref<QueueStatus | null>(null)
const callResult = ref<CallResult | null>(null)
const loading = ref(false)
const callLoading = ref(false)

async function fetchStatus() {
  loading.value = true
  try {
    const res = await triageApi.getQueueStatus(roomId.value)
    // 确保 entries 始终是数组，防止 a-list 收到 null/undefined 崩溃
    status.value = { ...res, entries: res?.entries ?? [] }
  }
  catch {
    // 接口异常时保持现有数据，不覆盖 status
  }
  finally { loading.value = false }
}

async function callNext() {
  callLoading.value = true
  try {
    callResult.value = await triageApi.callNext(roomId.value)
    message.success(`呼叫: ${callResult.value.patient_name_masked}`)
    fetchStatus()
  }
  catch { /* 错误已由拦截器弹出 */ }
  finally { callLoading.value = false }
}

async function recall() {
  callLoading.value = true
  try {
    callResult.value = await triageApi.recall(roomId.value)
    message.info('已重新呼叫')
  }
  catch { /* 错误已由拦截器弹出 */ }
  finally { callLoading.value = false }
}

async function miss() {
  try {
    await triageApi.missAndRequeue(roomId.value)
    message.warning('标记为过号，已重新排队')
    fetchStatus()
  }
  catch { /* 错误已由拦截器弹出 */ }
}

onMounted(fetchStatus)
</script>

<template>
  <div style="max-width: 800px;">
    <div style="display: flex; gap: 12px; margin-bottom: 16px; align-items: center;">
      <a-input v-model:value="roomId" style="width: 180px;" />
      <a-button :loading="loading" @click="fetchStatus">刷新队列</a-button>
    </div>

    <a-row :gutter="16">
      <a-col :span="14">
        <a-card title="叫号控制台" size="small">
          <div v-if="callResult" style="text-align: center; padding: 24px; background: #e6f4ff; border-radius: 8px; margin-bottom: 16px;">
            <div style="font-size: 14px; color: #8c8c8c; margin-bottom: 4px;">当前呼叫</div>
            <div style="font-size: 36px; font-weight: 800; color: #1890ff;">{{ callResult.patient_name_masked }}</div>
            <div style="color: #52c41a;">{{ callResult.room_name }} · 第 {{ callResult.queue_number }} 号</div>
          </div>
          <div v-else style="text-align: center; padding: 24px; color: #8c8c8c;">尚未开始叫号</div>
          <div style="display: flex; gap: 12px; justify-content: center; margin-top: 8px;">
            <a-button type="primary" size="large" :loading="callLoading" @click="callNext">呼叫下一位</a-button>
            <a-button size="large" @click="recall">重新呼叫</a-button>
            <a-button size="large" danger @click="miss">标记过号</a-button>
          </div>
        </a-card>
      </a-col>
      <a-col :span="10">
        <a-card title="等候队列" size="small" :loading="loading">
          <a-list v-if="status" :data-source="status.entries ?? []" size="small">
            <template #renderItem="{ item, index }">
              <a-list-item v-if="item">
                <span style="font-size: 18px; font-weight: 700; color: #1890ff; width: 28px; display: inline-block;">{{ index + 1 }}</span>
                <span style="flex: 1;">{{ item.patient_name_masked }}</span>
                <span class="text-muted" style="font-size: 12px;">{{ item.queue_number }}</span>
              </a-list-item>
            </template>
            <template #empty><a-empty description="队列为空" :image-style="{ height: '40px' }" /></template>
          </a-list>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>
