<!-- 核心目的：EvaluationReport页面 -->
<!-- 模块功能：智能效能优化-相关功能页面 -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { optimizationApi } from '@/api/optimization'
import type { EvaluationReport } from '@/types/optimization'

const route = useRoute()
const id = route.params.id as string
const loading = ref(false)
const report = ref<EvaluationReport | null>(null)

async function fetchData() {
  loading.value = true
  try {
    report.value = await optimizationApi.getEvaluation(id)
  }
  finally { loading.value = false }
}
onMounted(fetchData)

const baselineColumns = [
  { title: '指标', dataIndex: 'key' },
  {
    title: '数值',
    dataIndex: 'value',
    customRender: ({ text }: { text: number }) => text?.toFixed?.(2) ?? text,
  },
]
const trialColumns = [
  { title: '指标', dataIndex: 'key' },
  {
    title: '数值',
    dataIndex: 'value',
    customRender: ({ text }: { text: number }) => text?.toFixed?.(2) ?? text,
  },
  {
    title: '变化%',
    dataIndex: 'change',
    customRender: ({ text }: { text: number }) =>
      `${text >= 0 ? '+' : ''}${(text * 100).toFixed(1)}%`,
  },
]
</script>

<template>
  <a-spin :spinning="loading">
    <template v-if="report">
      <a-card title="评估报告" size="small">
        <a-descriptions :column="2" bordered size="small" style="margin-bottom: 16px;">
          <a-descriptions-item label="策略ID">{{ report.strategy_id }}</a-descriptions-item>
          <a-descriptions-item label="试运行ID">{{ report.trial_run_id }}</a-descriptions-item>
          <a-descriptions-item label="生成时间" :span="2">{{ report.generated_at }}</a-descriptions-item>
          <a-descriptions-item label="是否合格">
            <a-tag :color="report.is_qualified ? 'success' : 'error'">{{ report.is_qualified ? '合格' : '不合格' }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="合格阈值">{{ (report.qualify_threshold * 100).toFixed(1) }}%</a-descriptions-item>
          <a-descriptions-item label="推荐操作" :span="2">
            <a-tag :color="report.recommendation === 'promote' ? 'success' : report.recommendation === 'rollback' ? 'error' : 'warning'">
              {{ report.recommendation }}
            </a-tag>
          </a-descriptions-item>
        </a-descriptions>

        <a-row :gutter="16">
          <a-col :span="12">
            <a-card title="基准指标" size="small">
              <a-table
                :data-source="Object.entries(report.baseline_metrics).map(([k, v]) => ({ key: k, value: v }))"
                :columns="baselineColumns"
                :pagination="false"
                size="small"
              />
            </a-card>
          </a-col>
          <a-col :span="12">
            <a-card title="试运行指标" size="small">
              <a-table
                :data-source="Object.entries(report.trial_metrics).map(([k, v]) => ({ key: k, value: v, change: report!.change_pct[k] }))"
                :columns="trialColumns"
                :pagination="false"
                size="small"
              />
            </a-card>
          </a-col>
        </a-row>
      </a-card>
    </template>
    <a-empty v-else-if="!loading" />
  </a-spin>
</template>
