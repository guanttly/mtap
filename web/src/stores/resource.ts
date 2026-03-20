import { defineStore } from 'pinia'
import { ref } from 'vue'
import { resourceApi } from '@/api/resource'
import type { Device, ExamItem, SlotPool } from '@/types/resource'

export const useResourceStore = defineStore('resource', () => {
  const devices = ref<Device[]>([])
  const examItems = ref<ExamItem[]>([])
  const slotPools = ref<SlotPool[]>([])
  const loading = ref(false)

  async function fetchDevices(params: Record<string, unknown> = {}) {
    loading.value = true
    try {
      const res = await resourceApi.listDevices(params)
      devices.value = res.items
      return res
    }
    finally { loading.value = false }
  }

  async function fetchExamItems(params: Record<string, unknown> = {}) {
    loading.value = true
    try {
      const res = await resourceApi.listExamItems(params)
      examItems.value = res.items
      return res
    }
    finally { loading.value = false }
  }

  async function fetchSlotPools() {
    loading.value = true
    try {
      const res = await resourceApi.listSlotPools()
      slotPools.value = res.items
      return res
    }
    finally { loading.value = false }
  }

  return { devices, examItems, slotPools, loading, fetchDevices, fetchExamItems, fetchSlotPools }
})
