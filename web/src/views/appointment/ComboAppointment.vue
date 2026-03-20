<!-- 核心目的：ComboAppointment页面 -->
<!-- 模块功能：预约服务-相关功能页面 -->
<script setup lang="ts">
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { appointmentApi } from '@/api/appointment'

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
  <div style="max-width: 600px;">
    <a-card title="套餐预约" size="small">
      <a-form :model="form" layout="vertical">
        <a-form-item label="患者ID"><a-input v-model:value="form.patient_id" placeholder="输入患者ID" /></a-form-item>
        <a-form-item label="套餐ID"><a-input v-model:value="form.package_id" placeholder="输入检查套餐ID" /></a-form-item>
        <a-form-item label="来源">
          <a-select v-model:value="form.source">
            <a-select-option value="outpatient">门诊</a-select-option>
            <a-select-option value="inpatient">住院</a-select-option>
            <a-select-option value="physical">体检</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="备注"><a-textarea v-model:value="form.remark" :rows="3" /></a-form-item>
        <a-button type="primary" :loading="loading" @click="handleSubmit">提交套餐预约</a-button>
      </a-form>
    </a-card>

    <a-card v-if="result" title="预约结果" size="small" style="margin-top: 16px;">
      <a-descriptions :column="2" bordered size="small">
        <a-descriptions-item label="预约号">{{ result.appointment_no }}</a-descriptions-item>
        <a-descriptions-item label="状态">{{ result.status }}</a-descriptions-item>
        <a-descriptions-item label="患者">{{ result.patient_name }}</a-descriptions-item>
        <a-descriptions-item label="创建时间">{{ result.created_at }}</a-descriptions-item>
      </a-descriptions>
    </a-card>
  </div>
</template>
