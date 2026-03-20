// 核心目的：统计分析 Pinia Store
// 模块功能：缓存大屏快照、报表列表状态管理
import type { DashboardSnapshot, Report } from '@/types/analytics'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { analyticsApi } from '@/api/analytics'

export const useAnalyticsStore = defineStore('analytics', () => {
  const snapshot = ref<DashboardSnapshot | null>(null)
  const reports = ref<Report[]>([])
  const totalReports = ref(0)
  const loading = ref(false)

  async function fetchDashboard(campusId?: string) {
    loading.value = true
    try {
      snapshot.value = await analyticsApi.getDashboard(campusId)
      return snapshot.value
    }
    finally { loading.value = false }
  }

  async function fetchReports(params: Record<string, unknown> = {}) {
    loading.value = true
    try {
      const res = await analyticsApi.listReports(params)
      reports.value = res.items
      totalReports.value = res.total
      return res
    }
    finally { loading.value = false }
  }

  async function generateReport(data: Parameters<typeof analyticsApi.generateReport>[0]) {
    return analyticsApi.generateReport(data)
  }

  function clearSnapshot() {
    snapshot.value = null
  }

  return {
    snapshot,
    reports,
    totalReports,
    loading,
    fetchDashboard,
    fetchReports,
    generateReport,
    clearSnapshot,
  }
})
