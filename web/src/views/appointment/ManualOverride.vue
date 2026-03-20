<!-- 核心目的：ManualOverride页面 -->
<!-- 模块功能：预约服务-相关功能页面 -->
<script setup lang="ts">
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { appointmentApi } from '@/api/appointment'
import SlotPicker from '@/components/business/SlotPicker.vue'

const form = ref({
  appointment_id: '',
  action: 'reschedule' as 'reschedule' | 'cancel' | 'confirm' | 'mark_paid',
  reason: '',
  new_slot_id: '',
})
const loading = ref(false)
const appointment = ref<any>(null)
const searching = ref(false)

async function lookupAppointment() {
  if (!form.value.appointment_id) return
  searching.value = true
  try {
    appointment.value = await appointmentApi.getAppointment(form.value.appointment_id)
  }
  catch { message.error('未找到该预约') }
  finally { searching.value = false }
}

async function handleSubmit() {
  if (!form.value.appointment_id) return
  loading.value = true
  try {
    const { appointment_id, action, reason, new_slot_id } = form.value
    if (action === 'cancel') {
      await appointmentApi.cancel(appointment_id, reason)
    }
    else if (action === 'confirm') {
      await appointmentApi.confirm(appointment_id)
    }
    else if (action === 'mark_paid') {
      await appointmentApi.markPaid(appointment_id)
    }
    else if (action === 'reschedule') {
      await appointmentApi.reschedule(appointment_id, { new_slot_id, reason })
    }
    message.success('操作成功')
  }
  finally { loading.value = false }
}
</script>

<template>
  <div style="max-width: 700px;">
    <a-card title="人工干预" size="small">
      <a-form :model="form" layout="vertical">
        <a-form-item label="预约号">
          <a-input-search v-model:value="form.appointment_id" placeholder="输入预约ID" enter-button="查询" :loading="searching" @search="lookupAppointment" />
        </a-form-item>

        <template v-if="appointment">
          <a-descriptions :column="2" size="small" bordered style="margin-bottom: 16px;">
            <a-descriptions-item label="患者">{{ appointment.patient_name }}</a-descriptions-item>
            <a-descriptions-item label="状态">{{ appointment.status }}</a-descriptions-item>
            <a-descriptions-item label="预约号">{{ appointment.appointment_no }}</a-descriptions-item>
            <a-descriptions-item label="来源">{{ appointment.source }}</a-descriptions-item>
          </a-descriptions>

          <a-form-item label="操作类型">
            <a-radio-group v-model:value="form.action" button-style="solid">
              <a-radio-button value="confirm">确认预约</a-radio-button>
              <a-radio-button value="mark_paid">标记缴费</a-radio-button>
              <a-radio-button value="reschedule">改期</a-radio-button>
              <a-radio-button value="cancel">取消</a-radio-button>
            </a-radio-group>
          </a-form-item>

          <a-form-item v-if="form.action === 'reschedule'" label="新号源ID">
            <a-input v-model:value="form.new_slot_id" placeholder="输入新的号源ID" />
          </a-form-item>

          <a-form-item v-if="['cancel', 'reschedule'].includes(form.action)" label="操作原因">
            <a-textarea v-model:value="form.reason" :rows="3" />
          </a-form-item>

          <a-button type="primary" :loading="loading" @click="handleSubmit">执行操作</a-button>
        </template>
        <a-empty v-else description="请先查询预约" :image-style="{ height: '60px' }" />
      </a-form>
    </a-card>
  </div>
</template>
