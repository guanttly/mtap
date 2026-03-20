<!-- 核心目的：方案对比组件 -->
<!-- 模块功能：多套预约方案的对比视图（总耗时、往返次数、冲突提示） -->
<script setup lang="ts">
defineProps<{
  plans: Array<{
    index: number
    total_duration_min: number
    visit_count: number
    items: Array<{ exam_item_name: string, device_name: string, date: string, start_time: string, end_time: string }>
    conflicts?: string[]
    warnings?: string[]
  }>
  selectedIndex?: number
}>()

const emit = defineEmits<{
  select: [index: number]
}>()
</script>

<template>
  <div style="display: flex; gap: 12px; overflow-x: auto; padding-bottom: 4px;">
    <a-card
      v-for="plan in plans"
      :key="plan.index"
      size="small"
      :style="{
        minWidth: '240px',
        cursor: 'pointer',
        border: plan.index === selectedIndex ? '2px solid #1890ff' : '1px solid #d9d9d9',
        flexShrink: 0,
      }"
      @click="emit('select', plan.index)"
    >
      <template #title>
        <div style="display: flex; align-items: center; justify-content: space-between;">
          <span>方案 {{ plan.index + 1 }}</span>
          <a-tag v-if="plan.index === selectedIndex" color="blue" style="margin: 0;">已选</a-tag>
        </div>
      </template>
      <div style="display: flex; gap: 16px; margin-bottom: 10px;">
        <div style="text-align: center;">
          <div style="font-size: 20px; font-weight: 700; color: #1890ff;">{{ Math.round(plan.total_duration_min / 60 * 10) / 10 }}</div>
          <div style="font-size: 11px; color: #8c8c8c;">小时</div>
        </div>
        <div style="text-align: center;">
          <div style="font-size: 20px; font-weight: 700; color: #52c41a;">{{ plan.visit_count }}</div>
          <div style="font-size: 11px; color: #8c8c8c;">次就诊</div>
        </div>
      </div>
      <a-timeline :style="{ fontSize: '12px' }">
        <a-timeline-item v-for="item in plan.items" :key="item.exam_item_name">
          <div>{{ item.exam_item_name }}</div>
          <div style="color: #8c8c8c;">{{ item.date }} {{ item.start_time }} · {{ item.device_name }}</div>
        </a-timeline-item>
      </a-timeline>
      <template v-if="plan.conflicts?.length">
        <a-alert v-for="msg in plan.conflicts" :key="msg" :message="msg" type="error" banner style="font-size: 12px; margin-top: 4px;" />
      </template>
      <template v-if="plan.warnings?.length">
        <a-alert v-for="msg in plan.warnings" :key="msg" :message="msg" type="warning" banner style="font-size: 12px; margin-top: 4px;" />
      </template>
    </a-card>
    <a-empty v-if="!plans.length" description="暂无方案" />
  </div>
</template>
