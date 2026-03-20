<!-- 核心目的：PerformanceScan页面 -->
<!-- 模块功能：智能效能优化-相关功能页面 -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { optimizationApi } from '@/api/optimization'
import { usePagination } from '@/composables/usePagination'
import type { PerformanceScan } from '@/types/optimization'

const { loading, items, pagination, fetchData, onTableChange } = usePagination<PerformanceScan>(
  params => optimizationApi.listScans(params),
)
onMounted(() => fetchData())

const columns = [
  { title: '扫描周', dataIndex: 'scan_week', key: 'scan_week' },
  { title: '扫描时间', dataIndex: 'scanned_at', key: 'scanned_at' },
  { title: '发现机会数', key: 'count', customRender: ({ record }: { record: PerformanceScan }) => record.opportunities?.length ?? 0 },
  { title: '操作', key: 'actions' },
]

const expandedId = ref<string | null>(null)

function handleExpand(_: boolean, record: PerformanceScan) {
  expandedId.value = expandedId.value === record.id ? null : record.id
}
function handleToggle(record: PerformanceScan) {
  expandedId.value = expandedId.value === record.id ? null : record.id
}

const detailColumns = [
  { title: '指标编码', dataIndex: 'metric_code', key: 'code' },
  { title: '指标名', dataIndex: 'metric_name', key: 'name' },
  { title: '当前值', dataIndex: 'current_value', key: 'cur', customRender: ({ text }: { text: number }) => text?.toFixed?.(2) ?? text },
  { title: '正常值', dataIndex: 'normal_value', key: 'norm', customRender: ({ text }: { text: number }) => text?.toFixed?.(2) ?? text },
  { title: '偏差%', dataIndex: 'deviation_pct', key: 'dev', customRender: ({ text }: { text: number }) => `${(text * 100).toFixed(1)}%` },
  { title: '建议类别', dataIndex: 'suggested_category', key: 'cat' },
]
</script>

<template>
  <div>
    <a-table
      :columns="columns"
      :data-source="items"
      :loading="loading"
      :pagination="pagination"
      row-key="id"
      size="middle"
      :expandable="{ expandedRowKeys: expandedId ? [expandedId] : [], onExpand: handleExpand }"
      @change="onTableChange"
    >
      <template #expandedRowRender="{ record }">
        <a-table
          :data-source="(record as PerformanceScan).opportunities ?? []"
          :columns="detailColumns"
          :pagination="false"
          size="small"
          row-key="metric_code"
        />
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'actions'">
          <a-button type="link" size="small" @click="handleToggle(record as PerformanceScan)">展开详情</a-button>
        </template>
      </template>
    </a-table>
  </div>
</template>
