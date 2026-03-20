import { defineStore } from 'pinia'
import { ref } from 'vue'
import { appointmentApi } from '@/api/appointment'
import type { Appointment, BlacklistRecord } from '@/types/appointment'

export const useAppointmentStore = defineStore('appointment', () => {
  const appointments = ref<Appointment[]>([])
  const blacklist = ref<BlacklistRecord[]>([])
  const total = ref(0)
  const loading = ref(false)

  async function fetchAppointments(params: Record<string, unknown> = {}) {
    loading.value = true
    try {
      const res = await appointmentApi.listAppointments(params)
      appointments.value = res.items
      total.value = res.total
      return res
    }
    finally { loading.value = false }
  }

  async function fetchBlacklist(page = 1, pageSize = 20) {
    loading.value = true
    try {
      const res = await appointmentApi.listBlacklist({ page, page_size: pageSize })
      blacklist.value = res.items
      return res
    }
    finally { loading.value = false }
  }

  return { appointments, blacklist, total, loading, fetchAppointments, fetchBlacklist }
})
