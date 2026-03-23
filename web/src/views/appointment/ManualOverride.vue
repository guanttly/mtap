<!-- 核心目的：ManualOverride页面 -->
<!-- 模块功能：预约服务-相关功能页面 -->
<script setup lang="ts">
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { appointmentApi } from '@/api/appointment'
import SlotPicker from '@/components/business/SlotPicker.vue'
import SvgIcon from '@/components/common/SvgIcon.vue'

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
  <a-row :gutter="20">
    <a-col :xs="24" :xl="15">
      <a-card :bordered="false" title="人工干预">
        <a-form :model="form" layout="vertical">
          <a-form-item label="预约号查询">
            <a-input-search
              v-model:value="form.appointment_id"
              placeholder="输入预约ID进行查询"
              enter-button="查询"
              size="large"
              :loading="searching"
              @search="lookupAppointment"
            />
          </a-form-item>

          <template v-if="appointment">
            <a-divider style="margin: 16px 0;" />

            <div style="margin-bottom: 20px;">
              <a-descriptions :column="2" size="small" bordered>
                <a-descriptions-item label="患者">{{ appointment.patient_name }}</a-descriptions-item>
                <a-descriptions-item label="状态">
                  <a-tag color="blue">{{ appointment.status }}</a-tag>
                </a-descriptions-item>
                <a-descriptions-item label="预约号">{{ appointment.appointment_no }}</a-descriptions-item>
                <a-descriptions-item label="来源">{{ appointment.source }}</a-descriptions-item>
              </a-descriptions>
            </div>

            <a-form-item label="操作类型">
              <a-radio-group v-model:value="form.action" button-style="solid" size="large">
                <a-radio-button value="confirm"><SvgIcon name="check-circle-outlined" style="margin-right:4px" />确认预约</a-radio-button>
                <a-radio-button value="mark_paid"><SvgIcon name="pay-circle-outlined" style="margin-right:4px" />标记缴费</a-radio-button>
                <a-radio-button value="reschedule"><SvgIcon name="calendar-outlined" style="margin-right:4px" />改期</a-radio-button>
                <a-radio-button value="cancel"><SvgIcon name="close-outlined" style="margin-right:4px" />取消</a-radio-button>
              </a-radio-group>
            </a-form-item>

            <a-form-item v-if="form.action === 'reschedule'" label="新号源ID">
              <a-input v-model:value="form.new_slot_id" placeholder="输入新的号源ID" />
            </a-form-item>

            <a-form-item v-if="['cancel', 'reschedule'].includes(form.action)" label="操作原因">
              <a-textarea v-model:value="form.reason" :rows="3" placeholder="请填写操作原因，将记录到操作日志" />
            </a-form-item>

            <a-form-item style="margin-bottom: 0;">
              <a-button type="primary" size="large" :loading="loading" @click="handleSubmit">
                <span class="i-ant-design:tool-outlined" />
                执行操作
              </a-button>
            </a-form-item>
          </template>

          <a-empty v-else description="请先通过预约号查询预约记录" style="padding: 32px 0;" />
        </a-form>
      </a-card>
    </a-col>

    <a-col :xs="24" :xl="9">
      <a-card :bordered="false" title="操作说明" class="help-card">
        <a-timeline>
          <a-timeline-item color="green">确认预约：将预约状态从「待确认」改为「已确认」</a-timeline-item>
          <a-timeline-item color="blue">标记缴费：将预约标记为「已缴费」状态</a-timeline-item>
          <a-timeline-item color="orange">改期：重新分配号源，需填写改期原因</a-timeline-item>
          <a-timeline-item color="red">取消：取消该预约，需填写取消原因</a-timeline-item>
        </a-timeline>
        <a-alert
          type="warning"
          show-icon
          message="注意"
          description="所有操作将被记录到系统操作日志，请谨慎执行。"
          style="margin-top: 8px;"
        />
      </a-card>
    </a-col>
  </a-row>
</template>
