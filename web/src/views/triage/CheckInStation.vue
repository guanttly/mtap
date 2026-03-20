<!-- 核心目的：CheckInStation页面 -->
<!-- 模块功能：分诊管理-相关功能页面 -->
<script setup lang="ts">
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { triageApi } from '@/api/triage'
import type { CheckInResult } from '@/types/triage'

const qrCode = ref('')
const loading = ref(false)
const result = ref<CheckInResult | null>(null)

// 护士手动签到
const nurseForm = ref({ appointment_id: '', patient_id: '', remark: '' })
const nurseLoading = ref(false)
const activeTab = ref('kiosk')

async function kioskCheckIn() {
  if (!qrCode.value) return
  loading.value = true
  result.value = null
  try {
    result.value = await triageApi.kioskCheckIn(qrCode.value)
    message.success('签到成功')
    qrCode.value = ''
  }
  catch (e: any) {
    message.error(e?.message ?? '签到失败')
  }
  finally { loading.value = false }
}

async function nurseCheckIn() {
  if (!nurseForm.value.appointment_id) return
  nurseLoading.value = true
  result.value = null
  try {
    result.value = await triageApi.nurseCheckIn(nurseForm.value)
    message.success('签到成功')
  }
  catch (e: any) {
    message.error(e?.message ?? '签到失败')
  }
  finally { nurseLoading.value = false }
}
</script>

<template>
  <div style="max-width: 600px;">
    <a-tabs v-model:activeKey="activeTab">
      <a-tab-pane key="kiosk" tab="自助签到">
        <a-card size="small">
          <div style="text-align: center; padding: 24px 0;">
            <div style="font-size: 48px; margin-bottom: 24px;">📱</div>
            <p style="color: #8c8c8c; margin-bottom: 24px;">请扫描预约二维码或输入二维码内容</p>
            <a-input-search
              v-model:value="qrCode"
              size="large"
              placeholder="扫描或粘贴二维码内容"
              enter-button="签 到"
              :loading="loading"
              style="max-width: 420px;"
              @search="kioskCheckIn"
            />
          </div>
        </a-card>
      </a-tab-pane>
      <a-tab-pane key="nurse" tab="护士辅助签到">
        <a-card size="small">
          <a-form :model="nurseForm" layout="vertical">
            <a-form-item label="预约ID"><a-input v-model:value="nurseForm.appointment_id" /></a-form-item>
            <a-form-item label="患者ID"><a-input v-model:value="nurseForm.patient_id" /></a-form-item>
            <a-form-item label="备注"><a-input v-model:value="nurseForm.remark" /></a-form-item>
            <a-button type="primary" :loading="nurseLoading" @click="nurseCheckIn">确认签到</a-button>
          </a-form>
        </a-card>
      </a-tab-pane>
    </a-tabs>

    <a-card v-if="result" title="签到结果" size="small" style="margin-top: 16px;">
      <a-result
        status="success"
        title="签到成功"
        :sub-title="`队列编号：${result.queue_number}${result.is_late ? '（迟到）' : ''}`"
      >
        <template #extra>
          <a-descriptions :column="1" size="small" bordered>
            <a-descriptions-item label="签到ID">{{ result.check_in_id }}</a-descriptions-item>
            <a-descriptions-item label="队列编号">{{ result.queue_number }}</a-descriptions-item>
            <a-descriptions-item label="预计等待">{{ result.estimated_wait }} 分钟</a-descriptions-item>
            <a-descriptions-item label="就诊地点">{{ result.room_location }}</a-descriptions-item>
          </a-descriptions>
        </template>
      </a-result>
    </a-card>
  </div>
</template>
