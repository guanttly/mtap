<!-- 核心目的：Dashboard页面 -->
<!-- 模块功能：统计分析-实时监控看板，支持 WebSocket 推送 + 手动刷新 -->
<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { analyticsApi } from '@/api/analytics'
import type { DashboardSnapshot, WaitTrendPoint } from '@/types/analytics'
import { useWebSocket } from '@/composables/useWebSocket'
import * as echarts from 'echarts/core'
import { LineChart, BarChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent, TitleComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

echarts.use([LineChart, BarChart, GridComponent, TooltipComponent, LegendComponent, TitleComponent, CanvasRenderer])

const loading = ref(false)
const snapshot = ref<DashboardSnapshot | null>(null)
const trendChartEl = ref<HTMLElement>()
let trendChart: echarts.ECharts | null = null

// WebSocket 实时推送
const { status: wsStatus, connect, disconnect, on } = useWebSocket<DashboardSnapshot>('/ws/dashboard')
on('dashboard_update', (payload) => {
  snapshot.value = payload
  setTimeout(renderTrendChart, 50)
})

async function fetchData() {
  loading.value = true
  try {
    snapshot.value = await analyticsApi.getDashboard()
    setTimeout(renderTrendChart, 100)
  }
  catch {}
  finally { loading.value = false }
}

function renderTrendChart() {
  if (!trendChartEl.value || !snapshot.value?.wait_trend) return
  if (!trendChart) trendChart = echarts.init(trendChartEl.value)
  const trend: WaitTrendPoint[] = snapshot.value.wait_trend
  trendChart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: trend.map((t: WaitTrendPoint) => t.time) },
    yAxis: { type: 'value', name: '分钟' },
    series: [{ name: '平均等待', type: 'line', data: trend.map((t: WaitTrendPoint) => t.avg_wait_min), smooth: true, lineStyle: { color: '#1890ff' }, areaStyle: { color: 'rgba(24,144,255,.1)' } }],
  })
}

onMounted(() => {
  fetchData()
  connect()
})
onUnmounted(() => {
  disconnect()
  trendChart?.dispose()
})
</script>

<template>
  <div>
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
      <div style="display: flex; align-items: center; gap: 12px;">
        <span style="font-size: 16px; font-weight: 600;">实时数据看板</span>
        <a-badge
          :status="wsStatus === 'connected' ? 'success' : wsStatus === 'connecting' ? 'processing' : 'error'"
          :text="wsStatus === 'connected' ? '实时推送' : wsStatus === 'connecting' ? '连接中' : '断开'"
          style="font-size: 12px;"
        />
      </div>
      <a-button :loading="loading" @click="fetchData">手动刷新</a-button>
    </div>

    <a-spin :spinning="loading">
      <template v-if="snapshot">
        <!-- 概览卡片 -->
        <a-row :gutter="16" style="margin-bottom: 16px;">
          <a-col :span="6">
            <a-card size="small">
              <a-statistic title="今日总号源" :value="snapshot.slot_usage.total_slots" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card size="small">
              <a-statistic title="已使用" :value="snapshot.slot_usage.used_slots" :value-style="{ color: '#1890ff' }" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card size="small">
              <a-statistic title="剩余可约" :value="snapshot.slot_usage.available_slots" :value-style="{ color: '#52c41a' }" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card size="small">
              <a-statistic title="使用率" :value="(snapshot.slot_usage.usage_rate * 100).toFixed(1)" suffix="%" :value-style="{ color: snapshot.slot_usage.usage_rate > 0.9 ? '#ff4d4f' : '#faad14' }" />
            </a-card>
          </a-col>
        </a-row>

        <!-- 告警 -->
        <template v-if="snapshot.alerts?.length">
          <a-alert v-for="alert in snapshot.alerts" :key="alert.message" :message="alert.message" :type="alert.type === 'critical' ? 'error' : alert.type === 'warning' ? 'warning' : 'info'" show-icon style="margin-bottom: 8px;" />
        </template>

        <!-- 等待趋势 -->
        <a-card title="今日等待时长趋势" size="small" style="margin-bottom: 16px;">
          <div ref="trendChartEl" style="height: 280px;" />
        </a-card>

        <!-- 设备状态 -->
        <a-card title="设备运行状态" size="small">
          <a-table
            :data-source="snapshot.device_status ?? []"
            :columns="[
              { title: '设备', dataIndex: 'device_name', key: 'device_name' },
              { title: '状态', dataIndex: 'status', key: 'status' },
              { title: '当前队列', dataIndex: 'queue_count', key: 'queue_count' },
            ]"
            :pagination="false"
            row-key="device_id"
            size="small"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'status'">
                <a-badge :status="record.status === 'active' ? 'success' : record.status === 'fault' ? 'error' : 'warning'" :text="record.status" />
              </template>
            </template>
          </a-table>
        </a-card>
      </template>
      <a-empty v-else-if="!loading" description="点击刷新获取数据" />
    </a-spin>
  </div>
</template>
