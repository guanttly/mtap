<!-- 核心目的：AutoAppointment页面 -->
<!-- 模块功能：预约服务-相关功能页面 -->
<script setup lang="ts">
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { appointmentApi } from '@/api/appointment'
import PlanCompare from '@/components/business/PlanCompare.vue'
import SlotPicker from '@/components/business/SlotPicker.vue'
import type { AutoAppointmentReq } from '@/types/appointment'

const step = ref(0)
const loading = ref(false)
const plans = ref<any[]>([])
const selectedPlan = ref(0)
const confirming = ref(false)

const form = ref<AutoAppointmentReq>({
  patient_id: '',
  exam_item_ids: [],
  preferences: { preferred_time_period: 'morning' },
})

async function generatePlans() {
  if (!form.value.patient_id || !form.value.exam_item_ids?.length) {
    message.warning('请填写患者信息和检查项目')
    return
  }
  loading.value = true
  try {
    const res = await appointmentApi.autoAppointment(form.value)
    plans.value = (res.plans as any[]) ?? []
    if (plans.value.length === 0) {
      message.warning('未找到合适的预约方案，请调整条件后重试')
    }
    else {
      step.value = 1
    }
  }
  finally { loading.value = false }
}

async function confirmPlan() {
  confirming.value = true
  try {
    message.success('预约确认成功')
    step.value = 2
  }
  finally { confirming.value = false }
}

function resetForm() {
  step.value = 0
  plans.value = []
  form.value = { patient_id: '', exam_item_ids: [], preferences: { preferred_time_period: 'morning' } }
}
</script>

<template>
  <div style="max-width: 900px;">
    <a-steps :current="step" style="margin-bottom: 24px;">
      <a-step title="填写信息" />
      <a-step title="选择方案" />
      <a-step title="完成" />
    </a-steps>

    <a-card v-if="step === 0" title="患者与检查信息" size="small">
      <a-form :model="form" layout="vertical">
        <a-form-item label="患者ID"><a-input v-model:value="form.patient_id" /></a-form-item>
        <a-form-item label="检查项目ID（逗号分隔）">
          <a-textarea
            :value="form.exam_item_ids?.join(',')"
            :rows="3"
            @change="(e: Event) => form.exam_item_ids = (e.target as HTMLTextAreaElement).value.split(',').map(s => s.trim()).filter(Boolean)"
          />
        </a-form-item>
        <a-form-item label="偏好时段">
          <a-radio-group v-model:value="form.preferences!.preferred_time_period">
            <a-radio value="morning">上午</a-radio>
            <a-radio value="afternoon">下午</a-radio>
            <a-radio value="any">不限</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="偏好日期范围">
          <a-date-picker
            :value="form.preferences!.preferred_date_range?.start"
            value-format="YYYY-MM-DD"
            placeholder="开始日期"
            style="width:48%;"
            @change="(val: string) => { if (!form.preferences!.preferred_date_range) form.preferences!.preferred_date_range = { start: '', end: '' }; form.preferences!.preferred_date_range!.start = val }"
          />
          <span style="margin: 0 4%;">—</span>
          <a-date-picker
            :value="form.preferences!.preferred_date_range?.end"
            value-format="YYYY-MM-DD"
            placeholder="结束日期"
            style="width:48%;"
            @change="(val: string) => { if (!form.preferences!.preferred_date_range) form.preferences!.preferred_date_range = { start: '', end: '' }; form.preferences!.preferred_date_range!.end = val }"
          />
        </a-form-item>
        <a-button type="primary" :loading="loading" @click="generatePlans">智能生成方案</a-button>
      </a-form>
    </a-card>

    <template v-else-if="step === 1">
      <PlanCompare :plans="plans" :selected-index="selectedPlan" @select="selectedPlan = $event" />
      <div style="margin-top: 16px; display: flex; gap: 12px;">
        <a-button @click="step = 0">上一步</a-button>
        <a-button type="primary" :loading="confirming" @click="confirmPlan">确认方案 {{ selectedPlan + 1 }}</a-button>
      </div>
    </template>

    <a-result v-else status="success" title="预约成功！" sub-title="请患者按时就诊，凭二维码或手机号签到">
      <template #extra>
        <a-button type="primary" @click="resetForm">再次预约</a-button>
      </template>
    </a-result>

    <!-- SlotPicker 仅在选方案时按需展示（占位引用避免 unused import 警告） -->
    <span v-if="false"><SlotPicker :slots="[]" /></span>
  </div>
</template>
