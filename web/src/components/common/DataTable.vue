<!-- 核心目的：通用数据表格 -->
<!-- 模块功能：分页、排序、筛选、行操作的通用表格组件 -->
<script setup lang="ts" generic="T extends Record<string, unknown>">
import type { TableColumnType } from 'ant-design-vue'

defineProps<{
  columns: TableColumnType[]
  dataSource: T[]
  loading?: boolean
  rowKey?: string
  pagination?: Record<string, unknown> | false
}>()

const emit = defineEmits<{
  change: [pagination: Record<string, unknown>]
  'row-action': [action: string, record: T]
}>()

function handleChange(pag: Record<string, unknown>) {
  emit('change', pag)
}
</script>

<template>
  <a-table
    :columns="columns"
    :data-source="dataSource"
    :loading="loading"
    :row-key="rowKey ?? 'id'"
    :pagination="pagination !== undefined ? pagination : { showSizeChanger: true, showTotal: (t: number) => `共 ${t} 条` }"
    size="middle"
    @change="handleChange"
  >
    <template v-for="(_, name) in $slots" #[name]="slotData">
      <slot :name="name" v-bind="slotData ?? {}" />
    </template>
  </a-table>
</template>
