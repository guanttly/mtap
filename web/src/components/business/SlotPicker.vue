<!-- 核心目的：号源选择器 -->
<!-- 模块功能：可视化号源时段选择（按时段维度） -->
<script setup lang="ts">
import { computed } from 'vue'
import type { TimeSlot } from '@/types/resource'

const props = defineProps<{
  slots: TimeSlot[]
  selectedSlotId?: string
}>()

const emit = defineEmits<{
  'update:selectedSlotId': [id: string]
  select: [slot: TimeSlot]
}>()

const availableSlots = computed(() => props.slots.filter(s => s.status === 'available'))
const otherSlots = computed(() => props.slots.filter(s => s.status !== 'available'))

function selectSlot(slot: TimeSlot) {
  emit('update:selectedSlotId', slot.id)
  emit('select', slot)
}

function slotBorderColor(slot: TimeSlot) {
  if (slot.id === props.selectedSlotId) return '#1890ff'
  return '#d9d9d9'
}

function slotBg(slot: TimeSlot) {
  if (slot.id === props.selectedSlotId) return '#e6f4ff'
  if (slot.status === 'available') return '#f6ffed'
  return '#f5f5f5'
}
</script>

<template>
  <div>
    <div v-if="slots.length === 0">
      <a-empty description="无可用号源" :image-style="{ height: '40px' }" />
    </div>
    <template v-else>
      <div style="display: flex; gap: 8px; flex-wrap: wrap;">
        <div
          v-for="slot in [...availableSlots, ...otherSlots]"
          :key="slot.id"
          :style="{
            background: slotBg(slot),
            border: `1px solid ${slotBorderColor(slot)}`,
            borderRadius: '6px',
            padding: '8px 12px',
            cursor: slot.status === 'available' ? 'pointer' : 'not-allowed',
            opacity: slot.status === 'available' ? 1 : 0.5,
            minWidth: '120px',
            textAlign: 'center',
          }"
          @click="slot.status === 'available' && selectSlot(slot)"
        >
          <div style="font-weight: 600;">{{ slot.start_time }}–{{ slot.end_time }}</div>
          <div style="font-size: 12px; color: #8c8c8c;">{{ slot.status === 'available' ? '可预约' : slot.status }}</div>
          <div style="font-size: 11px; color: #bfbfbf;">{{ slot.standard_duration }}分钟</div>
        </div>
      </div>
    </template>
  </div>
</template>
