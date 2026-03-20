<!-- 核心目的：WaitingQueueView页面 -->
<!-- 模块功能：分诊管理-相关功能页面 -->
<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { triageApi } from '@/api/triage'
import type { QueueStatus } from '@/types/triage'

const roomId = ref('room-001')
const status = ref<QueueStatus | null>(null)
const loading = ref(false)
let timer: ReturnType<typeof setInterval> | null = null

const STATUS_COLOR: Record<string, string> = {
  waiting: 'default', called: 'processing', in_exam: 'warning', completed: 'success', missed: 'error',
}
const STATUS_LABEL: Record<string, string> = {
  waiting: '等待中', called: '已叫号', in_exam: '检查中', completed: '已完成', missed: '已过号',
}

async function fetchStatus() {
  loading.value = true
  try {
    const res = await triageApi.getQueueStatus(roomId.value)
    status.value = { ...res, entries: res?.entries ?? [] }
  }
  catch { /* 错误已由拦截器弹出 */ }
  finally { loading.value = false }
}

onMounted(() => {
  fetchStatus()
  timer = setInterval(fetchStatus, 15000)
})
onUnmounted(() => { if (timer) clearInterval(timer) })
</script>

<template>
  <div>
    <div style="display: flex; gap: 12px; align-items: center; margin-bottom: 16px;">
      <a-input v-model:value="roomId" placeholder="检查室ID" style="width: 200px;" />
      <a-button type="primary" :loading="loading" @click="fetchStatus">刷新</a-button>
      <span class="text-muted" style="font-size: 12px;">每15秒自动刷新</span>
    </div>

    <template v-if="status">
      <a-row :gutter="16" style="margin-bottom: 16px;">
        <a-col :span="8"><a-statistic title="等待人数" :value="status.waiting_count" /></a-col>
        <a-col :span="8"><a-statistic title="平均等待(分)" :value="status.average_wait" :precision="1" /></a-col>
        <a-col :span="8"><a-statistic title="已完成" :value="status.entries?.filter(e => e.status === 'completed').length ?? 0" /></a-col>
      </a-row>

      <a-card title="当前队列" size="small">
        <a-list :data-source="status.entries ?? []" size="small">
          <template #renderItem="{ item, index }">
            <a-list-item>
              <a-list-item-meta>
                <template #title>
                  <span style="font-size: 20px; font-weight: 700; color: #1890ff; margin-right: 12px;">{{ index + 1 }}</span>
                  {{ item.patient_name_masked }}
                </template>
                <template #description>队列号 {{ item.queue_number }}</template>
              </a-list-item-meta>
              <a-badge :status="STATUS_COLOR[item.status] as any" :text="STATUS_LABEL[item.status]" />
            </a-list-item>
          </template>
          <template #empty><a-empty description="队列为空" /></template>
        </a-list>
      </a-card>
    </template>
    <a-empty v-else-if="!loading" description="暂无数据，请选择检查室" />
  </div>
</template>
