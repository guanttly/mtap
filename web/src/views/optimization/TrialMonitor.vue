<!-- 核心目的：TrialMonitor页面 -->
<!-- 模块功能：智能效能优化-相关功能页面 -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { optimizationApi } from '@/api/optimization'
import type { TrialRun } from '@/types/optimization'

const route = useRoute()
const trialId = route.params.id as string
const loading = ref(false)
const trial = ref<TrialRun | null>(null)

async function fetchData() {
  loading.value = true
  try {
    trial.value = await optimizationApi.getTrialMonitor(trialId)
  }
  finally { loading.value = false }
}
onMounted(fetchData)
</script>

<template>
  <a-spin :spinning="loading">
    <template v-if="trial">
      <a-row :gutter="16" style="margin-bottom: 16px;">
        <a-col :span="6">
          <a-card size="small">
            <a-statistic title="试运行状态" :value="trial.status" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card size="small">
            <a-statistic title="试运行天数" :value="trial.trial_days" suffix="天" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card size="small">
            <a-statistic title="开始时间" :value="trial.started_at" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card size="small">
            <a-statistic title="结束时间" :value="trial.ends_at" />
          </a-card>
        </a-col>
      </a-row>

      <a-card title="灰度范围" size="small" style="margin-bottom: 16px;">
        <a-descriptions :column="1" size="small">
          <a-descriptions-item label="科室">
            <a-tag v-for="d in trial.gray_scope.department_ids" :key="d">{{ d }}</a-tag>
            <span v-if="!trial.gray_scope.department_ids?.length" class="text-muted">全部</span>
          </a-descriptions-item>
          <a-descriptions-item label="设备">
            <a-tag v-for="d in trial.gray_scope.device_ids" :key="d">{{ d }}</a-tag>
            <span v-if="!trial.gray_scope.device_ids?.length" class="text-muted">全部</span>
          </a-descriptions-item>
          <a-descriptions-item label="时间段">
            <a-tag v-for="t in trial.gray_scope.time_periods" :key="t">{{ t }}</a-tag>
            <span v-if="!trial.gray_scope.time_periods?.length" class="text-muted">全天</span>
          </a-descriptions-item>
        </a-descriptions>
      </a-card>

      <a-card title="紧急回滚阈值" size="small">
        <a-statistic :value="(trial.emergency_rollback_threshold * 100).toFixed(1)" suffix="%" />
        <div style="color: #8c8c8c; font-size: 12px; margin-top: 4px;">超过此偏差将自动触发回滚</div>
      </a-card>
    </template>
    <a-empty v-else-if="!loading" description="暂无试运行数据" />
  </a-spin>
</template>
