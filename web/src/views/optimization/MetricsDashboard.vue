<!-- 核心目的：MetricsDashboard页面 -->
<!-- 模块功能：智能效能优化-相关功能页面 -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { optimizationApi } from '@/api/optimization'
import type { EfficiencyMetric } from '@/types/optimization'

const loading = ref(false)
const metrics = ref<EfficiencyMetric[]>([])

async function fetchData() {
  loading.value = true
  try {
    const res = await optimizationApi.listMetrics()
    metrics.value = res?.items ?? []
  }
  finally { loading.value = false }
}
onMounted(fetchData)

function isNormal(m: EfficiencyMetric) {
  if (m.latest_value == null) return null
  return m.latest_value >= m.normal_min && m.latest_value <= m.normal_max
}
</script>

<template>
  <div>
    <div style="display: flex; justify-content: space-between; margin-bottom: 16px;">
      <span style="font-size: 16px; font-weight: 600;">效率指标总览</span>
      <a-button :loading="loading" @click="fetchData">刷新</a-button>
    </div>
    <a-spin :spinning="loading">
      <a-row :gutter="[16, 16]">
        <a-col v-for="metric in metrics" :key="metric.id" :span="8">
          <a-card size="small" :bodyStyle="{ padding: '16px' }">
            <div style="display: flex; justify-content: space-between; align-items: flex-start;">
              <div>
                <div style="font-size: 13px; color: #8c8c8c; margin-bottom: 6px;">{{ metric.name }} ({{ metric.code }})</div>
                <div style="font-size: 28px; font-weight: 700;">
                  {{ metric.latest_value != null ? metric.latest_value.toFixed(1) : '—' }}
                  <span style="font-size: 14px; font-weight: 400; color: #8c8c8c;">{{ metric.unit }}</span>
                </div>
                <div style="margin-top: 4px; font-size: 12px; color: #8c8c8c;">
                  正常范围: {{ metric.normal_min }} ~ {{ metric.normal_max }} {{ metric.unit }}
                </div>
              </div>
              <div v-if="metric.latest_value != null" :style="{ width: '10px', height: '10px', borderRadius: '50%', background: isNormal(metric) ? '#52c41a' : '#ff4d4f', marginTop: '4px' }" />
            </div>
            <div v-if="metric.latest_sampled_at" style="text-align: right; font-size: 11px; color: #bfbfbf; margin-top: 8px;">
              最后采样: {{ metric.latest_sampled_at }}
            </div>
          </a-card>
        </a-col>
      </a-row>
      <a-empty v-if="!loading && metrics.length === 0" description="暂无指标数据" style="margin-top: 60px;" />
    </a-spin>
  </div>
</template>
