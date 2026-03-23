<!-- 核心目的：ComboAppointment页面 -->
<!-- 模块功能：预约服务-相关功能页面 -->
<script setup lang="ts">
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { appointmentApi } from '@/api/appointment'
import SvgIcon from '@/components/common/SvgIcon.vue'

const form = ref({ patient_id: '', package_id: '', source: 'physical', remark: '' })
const loading = ref(false)
const result = ref<any>(null)

async function handleSubmit() {
  if (!form.value.patient_id || !form.value.package_id) {
    message.warning('请填写患者ID和套餐ID')
    return
  }
  loading.value = true
  try {
    result.value = await appointmentApi.comboAppointment(form.value)
    message.success('套餐预约成功')
  }
  finally { loading.value = false }
}
</script>

<template>
  <a-row :gutter="20">
    <a-col :xs="24" :xl="14">
      <a-card :bordered="false" title="套餐预约">
        <a-form :model="form" layout="vertical">
          <a-form-item label="患者ID" required>
            <a-input v-model:value="form.patient_id" placeholder="输入患者ID" allow-clear size="large" />
          </a-form-item>
          <a-form-item label="套餐ID" required>
            <a-input v-model:value="form.package_id" placeholder="输入检查套餐ID" allow-clear size="large" />
          </a-form-item>
          <a-form-item label="来源">
            <a-select v-model:value="form.source" style="width: 100%;" size="large">
              <a-select-option value="outpatient"><SvgIcon name="medicine-box-outlined" style="margin-right:4px" />门诊</a-select-option>
              <a-select-option value="inpatient"><SvgIcon name="home-outlined" style="margin-right:4px" />住院</a-select-option>
              <a-select-option value="physical"><SvgIcon name="experiment-outlined" style="margin-right:4px" />体检</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item label="备注">
            <a-textarea v-model:value="form.remark" :rows="3" placeholder="选填，补充预约说明" />
          </a-form-item>
          <a-form-item style="margin-bottom: 0;">
            <a-button type="primary" size="large" :loading="loading" @click="handleSubmit">
              <span class="i-ant-design:shopping-outlined" />
              提交套餐预约
            </a-button>
          </a-form-item>
        </a-form>
      </a-card>
    </a-col>

    <a-col :xs="24" :xl="10">
      <a-card v-if="result" :bordered="false" title="预约结果" style="margin-bottom: 16px;">
        <a-result status="success" title="套餐预约成功" style="padding: 8px 0 16px;">
          <template #extra>
            <a-descriptions :column="1" size="small" bordered>
              <a-descriptions-item label="预约号">{{ result.appointment_no }}</a-descriptions-item>
              <a-descriptions-item label="状态">{{ result.status }}</a-descriptions-item>
              <a-descriptions-item label="患者">{{ result.patient_name }}</a-descriptions-item>
              <a-descriptions-item label="创建时间">{{ result.created_at }}</a-descriptions-item>
            </a-descriptions>
          </template>
        </a-result>
      </a-card>
      <a-card v-else :bordered="false" title="套餐说明" class="help-card">
        <p style="color: #595959; font-size: 13px; margin-bottom: 14px;">套餐预约将一次性安排套餐内的所有检查项目，系统自动合理安排检查时序。</p>
        <a-alert type="info" show-icon message="来源信息用于统计分析，请如实填写" style="margin-bottom: 12px;" />
        <a-alert type="warning" show-icon message="套餐ID需在系统中预先配置，请联系管理员获取" />
      </a-card>
    </a-col>
  </a-row>
</template>
