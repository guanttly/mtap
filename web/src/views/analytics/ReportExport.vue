<!-- 核心目的：ReportExport页面 -->
<!-- 模块功能：统计分析-相关功能页面 -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { analyticsApi } from '@/api/analytics'
import type { Report, ReportType, ReportFormat } from '@/types/analytics'
import { usePagination } from '@/composables/usePagination'

const { loading, items, pagination, fetchData, onTableChange } = usePagination<Report>(
  params => analyticsApi.listReports(params),
)
onMounted(() => fetchData())

const showModal = ref(false)
const generating = ref(false)
const genForm = ref<{ report_type: ReportType, date_start: string, date_end: string, format: ReportFormat }>({
  report_type: 'daily_summary',
  date_start: '',
  date_end: '',
  format: 'xlsx',
})

const STATUS_COLOR: Record<string, string> = { generating: 'processing', ready: 'success', failed: 'error' }
const STATUS_LABEL: Record<string, string> = { generating: '生成中', ready: '已就绪', failed: '失败' }

async function generateReport() {
  generating.value = true
  try {
    await analyticsApi.generateReport(genForm.value)
    message.success('报表生成任务已提交')
    showModal.value = false
    fetchData()
  }
  finally { generating.value = false }
}

async function downloadReport(record: Report) {
  try {
    const blob = await analyticsApi.exportReport(record.id, record.format)
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `report_${record.id}.${record.format}`
    a.click()
    URL.revokeObjectURL(url)
  }
  catch { message.error('下载失败') }
}

const columns = [
  { title: '报表类型', dataIndex: 'report_type', key: 'type' },
  { title: '开始日期', dataIndex: 'date_start', key: 'start' },
  { title: '结束日期', dataIndex: 'date_end', key: 'end' },
  { title: '格式', dataIndex: 'format', key: 'format' },
  { title: '状态', dataIndex: 'status', key: 'status' },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at' },
  { title: '操作', key: 'actions' },
]
</script>

<template>
  <a-card class="list-card" :bordered="false">
    <template #title>报表导出</template>
    <template #extra>
      <a-button type="primary" @click="showModal = true">生成新报表</a-button>
    </template>
    <a-table :columns="columns" :data-source="items" :loading="loading" :pagination="pagination" row-key="id" size="middle" @change="onTableChange">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag :color="STATUS_COLOR[(record as Report).status]">{{ STATUS_LABEL[(record as Report).status] }}</a-tag>
        </template>
        <template v-if="column.key === 'actions'">
          <a-button
            v-if="(record as Report).status === 'ready'"
            type="link" size="small"
            @click="downloadReport(record as Report)"
          >下载</a-button>
        </template>
      </template>
    </a-table>

    <a-modal v-model:open="showModal" title="生成报表" :confirm-loading="generating" @ok="generateReport">
      <a-form :model="genForm" layout="vertical">
        <a-form-item label="报表类型">
          <a-select v-model:value="genForm.report_type">
            <a-select-option value="daily">日报</a-select-option>
            <a-select-option value="weekly">周报</a-select-option>
            <a-select-option value="monthly">月报</a-select-option>
            <a-select-option value="custom">自定义</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="开始日期"><a-date-picker v-model:value="genForm.date_start" value-format="YYYY-MM-DD" style="width:100%;" /></a-form-item>
        <a-form-item label="结束日期"><a-date-picker v-model:value="genForm.date_end" value-format="YYYY-MM-DD" style="width:100%;" /></a-form-item>
        <a-form-item label="导出格式">
          <a-radio-group v-model:value="genForm.format">
            <a-radio value="excel">Excel</a-radio>
            <a-radio value="pdf">PDF</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>
