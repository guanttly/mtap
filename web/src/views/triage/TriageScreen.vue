<!-- 核心目的：TriageScreen页面 -->
<!-- 模块功能：分诊管理-相关功能页面 -->
<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { triageApi } from '@/api/triage'
import { useWebSocket } from '@/composables/useWebSocket'
import type { QueueStatus, CallResult } from '@/types/triage'

const rooms = ['room-001', 'room-002', 'room-003']
const statusMap = ref<Record<string, QueueStatus>>({})
const currentCall = ref<CallResult | null>(null)
const now = ref(new Date())

const { connect, on } = useWebSocket<CallResult>('/ws/triage')

async function fetchAll() {
  for (const roomId of rooms) {
    try {
      statusMap.value[roomId] = await triageApi.getQueueStatus(roomId)
    }
    catch {}
  }
}

let clockTimer: ReturnType<typeof setInterval>
let refreshTimer: ReturnType<typeof setInterval>

onMounted(() => {
  fetchAll()
  clockTimer = setInterval(() => { now.value = new Date() }, 1000)
  refreshTimer = setInterval(fetchAll, 20000)
  connect()
  on('call', (payload) => { currentCall.value = payload })
})
onUnmounted(() => {
  clearInterval(clockTimer)
  clearInterval(refreshTimer)
})

function timeStr() {
  return now.value.toLocaleTimeString('zh-CN', { hour12: false })
}
function dateStr() {
  return now.value.toLocaleDateString('zh-CN', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })
}
</script>

<template>
  <div style="background: #001529; min-height: 100vh; color: #fff; padding: 0; font-family: 'Microsoft YaHei', sans-serif;">
    <!-- 顶部 -->
    <div style="background: #003a70; padding: 16px 32px; display: flex; align-items: center; justify-content: space-between;">
      <div style="font-size: 28px; font-weight: 800; letter-spacing: 4px;">医 技 分 诊 系 统</div>
      <div style="text-align: right;">
        <div style="font-size: 32px; font-weight: 700; font-variant-numeric: tabular-nums;">{{ timeStr() }}</div>
        <div style="font-size: 14px; opacity: .7;">{{ dateStr() }}</div>
      </div>
    </div>

    <!-- 当前呼叫 -->
    <div v-if="currentCall" style="background: #1890ff; padding: 24px 32px; text-align: center; animation: pulse 2s infinite;">
      <div style="font-size: 16px; margin-bottom: 4px; opacity: .85;">请 前 往 就 诊</div>
      <div style="font-size: 56px; font-weight: 900; letter-spacing: 6px;">{{ currentCall.patient_name_masked }}</div>
      <div style="font-size: 20px; margin-top: 4px;">{{ currentCall.room_name }} · 第 {{ currentCall.queue_number }} 号</div>
    </div>

    <!-- 各检查室队列 -->
    <div style="display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; padding: 24px 32px;">
      <div
        v-for="roomId in rooms"
        :key="roomId"
        style="background: rgba(255,255,255,.06); border-radius: 12px; padding: 20px;"
      >
        <div style="font-size: 18px; font-weight: 700; border-bottom: 1px solid rgba(255,255,255,.15); padding-bottom: 10px; margin-bottom: 12px;">
          {{ statusMap[roomId]?.room_name ?? roomId }}
        </div>
        <div v-if="statusMap[roomId]">
          <div style="display: flex; gap: 16px; margin-bottom: 16px;">
            <div style="text-align: center;">
              <div style="font-size: 28px; font-weight: 700; color: #40a9ff;">{{ statusMap[roomId].waiting_count }}</div>
              <div style="font-size: 12px; opacity: .6;">等待</div>
            </div>
            <div style="text-align: center;">
              <div style="font-size: 28px; font-weight: 700; color: #52c41a;">{{ statusMap[roomId].entries?.filter(e => e.status === 'completed').length ?? 0 }}</div>
              <div style="font-size: 12px; opacity: .6;">已完成</div>
            </div>
          </div>
          <div v-for="(entry, idx) in statusMap[roomId].entries?.slice(0, 5)" :key="entry.id" style="display: flex; align-items: center; gap: 10px; padding: 8px 0; border-bottom: 1px solid rgba(255,255,255,.06);">
            <span :style="{ fontWeight: '700', fontSize: '20px', color: idx === 0 ? '#40a9ff' : 'rgba(255,255,255,.5)', width: '24px' }">{{ idx + 1 }}</span>
            <span :style="{ flex: 1, opacity: idx === 0 ? 1 : 0.6 }">{{ entry.patient_name_masked }}</span>
          </div>
        </div>
        <a-spin v-else style="margin: 20px auto; display: block;" />
      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: .92; }
}
</style>
