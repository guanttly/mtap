import { defineStore } from 'pinia'
import { ref } from 'vue'
import { triageApi } from '@/api/triage'
import type { QueueStatus, CheckInResult } from '@/types/triage'

export const useTriageStore = defineStore('triage', () => {
  const queueStatus = ref<QueueStatus | null>(null)
  const lastCheckIn = ref<CheckInResult | null>(null)
  const loading = ref(false)

  async function fetchQueueStatus(roomId: string) {
    loading.value = true
    try {
      const res = await triageApi.getQueueStatus(roomId)
      queueStatus.value = res
      return res
    }
    finally { loading.value = false }
  }

  async function checkIn(qrCodeData: string) {
    loading.value = true
    try {
      const res = await triageApi.kioskCheckIn(qrCodeData)
      lastCheckIn.value = res
      return res
    }
    finally { loading.value = false }
  }

  return { queueStatus, lastCheckIn, loading, fetchQueueStatus, checkIn }
})
