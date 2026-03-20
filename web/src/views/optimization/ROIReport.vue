<!-- 核心目的：ROIReport页面 -->
<!-- 模块功能：智能效能优化-相关功能页面 -->
<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import { optimizationApi } from '@/api/optimization'
import type { ROIReport } from '@/types/optimization'

const route = useRoute()
const id = route.params.id as string
const loading = ref(false)
const report = ref<ROIReport | null>(null)
const submitting = ref(false)
const reviewForm = ref({ approved: true, comment: '' })

// 计算 ROI = (年收益 - 总投入) / 总投入
const roi = computed(() => {
  if (!report.value) return 0
  const { total_investment, expected_annual_revenue } = report.value
  if (!total_investment) return 0
  return (expected_annual_revenue - total_investment) / total_investment
})

async function fetchData() {
  loading.value = true
  try {
    report.value = await optimizationApi.getROIReport(id)
  }
  finally { loading.value = false }
}

async function submitReview() {
  submitting.value = true
  try {
    await optimizationApi.submitROIResult(id, reviewForm.value)
    message.success('ROI 审核结果已提交')
    fetchData()
  }
  finally { submitting.value = false }
}

onMounted(fetchData)
</script>

<template>
  <a-spin :spinning="loading">
    <template v-if="report">
      <a-card title="ROI 分析报告" size="small">
        <a-descriptions :column="1" bordered size="small" style="margin-bottom: 16px;">
          <a-descriptions-item label="当前瓶颈">{{ report.current_bottleneck }}</a-descriptions-item>
        </a-descriptions>
        <a-row :gutter="16" style="margin-bottom: 16px;">
          <a-col :span="6">
            <a-statistic title="总投入(元)" :value="report.total_investment" :precision="2" />
          </a-col>
          <a-col :span="6">
            <a-statistic title="预期年收益(元)" :value="report.expected_annual_revenue" :precision="2" :value-style="{ color: '#52c41a' }" />
          </a-col>
          <a-col :span="6">
            <a-statistic title="ROI" :value="(roi * 100).toFixed(1)" suffix="%" :value-style="{ color: roi > 0 ? '#52c41a' : '#ff4d4f' }" />
          </a-col>
          <a-col :span="6">
            <a-statistic title="回本周期(月)" :value="report.payback_period_months.toFixed(1)" />
          </a-col>
        </a-row>

        <a-card title="风险因素" size="small" style="margin-bottom: 16px;">
          <a-tag v-for="r in report.risk_factors" :key="r" color="orange" style="margin: 4px;">{{ r }}</a-tag>
          <span v-if="!report.risk_factors?.length" class="text-muted">无</span>
        </a-card>

        <a-card v-if="!report.approval_result" title="ROI 审核" size="small">
          <a-form :model="reviewForm" layout="vertical">
            <a-form-item label="审核结论">
              <a-radio-group v-model:value="reviewForm.approved">
                <a-radio :value="true">批准</a-radio>
                <a-radio :value="false">拒绝</a-radio>
              </a-radio-group>
            </a-form-item>
            <a-form-item label="审核意见">
              <a-textarea v-model:value="reviewForm.comment" :rows="4" />
            </a-form-item>
            <a-button type="primary" :loading="submitting" @click="submitReview">提交审核</a-button>
          </a-form>
        </a-card>
        <a-alert v-else :message="`审核结果: ${report.approval_result}`" type="success" />
      </a-card>
    </template>
    <a-empty v-else-if="!loading" />
  </a-spin>
</template>
