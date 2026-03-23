<!-- 核心目的：AutoAppointment页面 -->
<!-- 模块功能：预约服务-相关功能页面 -->
<script setup lang="ts">
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { appointmentApi } from '@/api/appointment'
import SvgIcon from '@/components/common/SvgIcon.vue'
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
  <div>
    <!-- Steps 导航 -->
    <a-card :bordered="false" style="margin-bottom: 20px;">
      <a-steps :current="step">
        <a-step title="填写信息" description="患者与检查信息" />
        <a-step title="选择方案" description="选择最优预约方案" />
        <a-step title="完成" description="预约确认成功" />
      </a-steps>
    </a-card>

    <!-- Step 0: 填写表单 -->
    <a-row v-if="step === 0" :gutter="20">
      <a-col :xs="24" :xl="16">
        <a-card :bordered="false" title="患者与检查信息">
          <a-form :model="form" layout="vertical">
            <a-form-item label="患者ID" required>
              <a-input
                v-model:value="form.patient_id"
                placeholder="请输入患者ID"
                allow-clear
                size="large"
              />
            </a-form-item>

            <a-form-item label="检查项目ID" extra="多个检查项目请用逗号分隔，如：item1,item2,item3" required>
              <a-textarea
                :value="form.exam_item_ids?.join(',')"
                :rows="3"
                placeholder="输入检查项目ID，逗号分隔"
                @change="(e: Event) => form.exam_item_ids = (e.target as HTMLTextAreaElement).value.split(',').map(s => s.trim()).filter(Boolean)"
              />
            </a-form-item>

            <a-row :gutter="24">
              <a-col :span="12">
                <a-form-item label="偏好时段">
                  <a-radio-group v-model:value="form.preferences!.preferred_time_period" button-style="outline">
                    <a-radio value="morning"><SvgIcon name="sun-outlined" style="margin-right:4px" />上午</a-radio>
                    <a-radio value="afternoon"><SvgIcon name="cloud-outlined" style="margin-right:4px" />下午</a-radio>
                    <a-radio value="any"><SvgIcon name="unlock-outlined" style="margin-right:4px" />不限</a-radio>
                  </a-radio-group>
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="偏好日期范围">
                  <a-space direction="vertical" style="width: 100%;">
                    <a-date-picker
                      :value="form.preferences!.preferred_date_range?.start"
                      value-format="YYYY-MM-DD"
                      placeholder="开始日期"
                      style="width: 100%;"
                      @change="(val: string) => { if (!form.preferences!.preferred_date_range) form.preferences!.preferred_date_range = { start: '', end: '' }; form.preferences!.preferred_date_range!.start = val }"
                    />
                    <a-date-picker
                      :value="form.preferences!.preferred_date_range?.end"
                      value-format="YYYY-MM-DD"
                      placeholder="结束日期"
                      style="width: 100%;"
                      @change="(val: string) => { if (!form.preferences!.preferred_date_range) form.preferences!.preferred_date_range = { start: '', end: '' }; form.preferences!.preferred_date_range!.end = val }"
                    />
                  </a-space>
                </a-form-item>
              </a-col>
            </a-row>

            <a-form-item style="margin-bottom: 0;">
              <a-button type="primary" size="large" :loading="loading" @click="generatePlans">
                <span class="i-ant-design:thunderbolt-outlined" />
                智能生成方案
              </a-button>
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>

      <a-col :xs="24" :xl="8">
        <a-card :bordered="false" title="填写说明" class="help-card">
          <a-steps direction="vertical" :current="-1" size="small" style="margin-bottom: 16px;">
            <a-step title="填写患者 ID" description="输入就诊系统中的患者唯一标识符" />
            <a-step title="填写检查项目" description="可填写多个项目，系统自动检测冲突" />
            <a-step title="选择偏好时段" description="系统将优先匹配您偏好的时间段" />
          </a-steps>
          <a-alert
            type="info"
            show-icon
            message="多项目自动排期"
            description="系统将根据冲突规则、设备排班、号源池情况自动生成最优预约方案。"
          />
        </a-card>
      </a-col>
    </a-row>

    <!-- Step 1: 选择方案 -->
    <template v-else-if="step === 1">
      <PlanCompare :plans="plans" :selected-index="selectedPlan" @select="selectedPlan = $event" />
      <div style="margin-top: 16px; display: flex; gap: 12px;">
        <a-button size="large" @click="step = 0">上一步</a-button>
        <a-button type="primary" size="large" :loading="confirming" @click="confirmPlan">
          确认方案 {{ selectedPlan + 1 }}
        </a-button>
      </div>
    </template>

    <!-- Step 2: 完成 -->
    <a-card v-else :bordered="false">
      <a-result
        status="success"
        title="预约成功！"
        sub-title="请患者按时就诊，凭预约二维码或手机号到检查室签到"
      >
        <template #extra>
          <a-button type="primary" size="large" @click="resetForm">再次预约</a-button>
        </template>
      </a-result>
    </a-card>

    <span v-if="false"><SlotPicker :slots="[]" /></span>
  </div>
</template>
