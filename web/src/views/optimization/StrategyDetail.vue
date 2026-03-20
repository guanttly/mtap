<!-- 核心目的：StrategyDetail页面 -->
<!-- 模块功能：智能效能优化-相关功能页面 -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { optimizationApi } from '@/api/optimization'
import type { OptimizationStrategy } from '@/types/optimization'

const route = useRoute()
const router = useRouter()
const strategyId = route.params.id as string

const loading = ref(false)
const strategy = ref<OptimizationStrategy | null>(null)

async function fetchData() {
  loading.value = true
  try {
    strategy.value = await optimizationApi.getStrategy(strategyId)
  }
  finally { loading.value = false }
}

async function handleRollback() {
  await optimizationApi.rollbackStrategy(strategyId)
  message.success('已回滚')
  fetchData()
}

async function handlePromote() {
  await optimizationApi.promoteStrategy(strategyId)
  message.success('已推全量')
  fetchData()
}

onMounted(fetchData)
</script>

<template>
  <a-spin :spinning="loading">
    <template v-if="strategy">
      <a-page-header :title="strategy.title" @back="router.back()">
        <template #tags>
          <a-tag :color="strategy.status === 'promoted' ? 'success' : strategy.status === 'trial_running' ? 'processing' : 'default'">{{ strategy.status }}</a-tag>
        </template>
        <template #extra>
          <a-button v-if="strategy.status === 'trial_running'" @click="handleRollback">回滚</a-button>
          <a-button v-if="strategy.status === 'trial_running'" type="primary" @click="handlePromote">推全量</a-button>
        </template>
      </a-page-header>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-card title="基本信息" size="small">
            <a-descriptions :column="1" size="small">
              <a-descriptions-item label="策略类别">{{ strategy.category }}</a-descriptions-item>
              <a-descriptions-item label="当前值">{{ strategy.current_value }}</a-descriptions-item>
              <a-descriptions-item label="目标值">{{ strategy.target_value }}</a-descriptions-item>
              <a-descriptions-item label="预期收益">{{ strategy.expected_benefit }}</a-descriptions-item>
              <a-descriptions-item label="风险说明">{{ strategy.risk_note }}</a-descriptions-item>
            </a-descriptions>
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card title="试运行信息" size="small">
            <template v-if="strategy.trial_run">
              <a-descriptions :column="1" size="small">
                <a-descriptions-item label="状态">{{ strategy.trial_run.status }}</a-descriptions-item>
                <a-descriptions-item label="天数">{{ strategy.trial_run.trial_days }} 天</a-descriptions-item>
                <a-descriptions-item label="开始">{{ strategy.trial_run.started_at }}</a-descriptions-item>
                <a-descriptions-item label="结束">{{ strategy.trial_run.ends_at }}</a-descriptions-item>
              </a-descriptions>
            </template>
            <a-empty v-else description="尚未开始试运行" :image-style="{ height: '40px' }" />
          </a-card>
        </a-col>
      </a-row>
    </template>
    <a-empty v-else-if="!loading" description="策略不存在" />
  </a-spin>
</template>
