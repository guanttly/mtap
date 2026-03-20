<!-- 核心目的：SlotPoolView页面 -->
<!-- 模块功能：资源管理-SlotPool视图 -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { resourceApi } from '@/api/resource'
import type { SlotPool, TimeSlot } from '@/types/resource'

const pools = ref<SlotPool[]>([])
const slots = ref<TimeSlot[]>([])
const loading = ref(false)
const slotsLoading = ref(false)
const selectedPool = ref<SlotPool | null>(null)

const STATUS_COLOR: Record<string, string> = { available: 'success', booked: 'processing', locked: 'warning', released: 'default' }

async function fetchPools() {
  loading.value = true
  try {
    const res = await resourceApi.listSlotPools()
    pools.value = res.items
  }
  finally { loading.value = false }
}

async function fetchSlots(pool: SlotPool) {
  selectedPool.value = pool
  slotsLoading.value = true
  try {
    const res = await resourceApi.listTimeSlots({ slot_pool_id: pool.id })
    slots.value = res.items
  }
  finally { slotsLoading.value = false }
}

async function releaseSlot(slot: TimeSlot) {
  await resourceApi.releaseSlot(slot.id)
  message.success('已强制释放')
  if (selectedPool.value) fetchSlots(selectedPool.value)
}

onMounted(fetchPools)

const slotColumns = [
  { title: '日期', dataIndex: 'date', key: 'date' },
  { title: '开始', dataIndex: 'start_time', key: 'start_time' },
  { title: '结束', dataIndex: 'end_time', key: 'end_time' },
  { title: '总数', dataIndex: 'total_count', key: 'total_count' },
  { title: '剩余', dataIndex: 'remaining_count', key: 'remaining_count' },
  { title: '状态', dataIndex: 'status', key: 'status' },
  { title: '操作', key: 'actions' },
]
</script>

<template>
  <div style="display: flex; gap: 16px; height: 100%;">
    <div style="width: 280px; flex-shrink: 0;">
      <a-card title="号源池列表" size="small" :loading="loading">
        <a-list :data-source="pools" size="small">
          <template #renderItem="{ item }">
            <a-list-item
              :style="{ cursor: 'pointer', background: selectedPool?.id === item.id ? '#e6f4ff' : 'transparent', borderRadius: '6px', padding: '8px 12px' }"
              @click="fetchSlots(item)"
            >
              <a-list-item-meta :title="item.name" :description="`${item.device_name} · ${item.exam_item_name}`" />
            </a-list-item>
          </template>
        </a-list>
      </a-card>
    </div>
    <div style="flex: 1; min-width: 0;">
      <a-card :title="selectedPool ? `「${selectedPool.name}」号源明细` : '请选择号源池'" size="small" :loading="slotsLoading">
        <a-table v-if="selectedPool" :columns="slotColumns" :data-source="slots" row-key="id" size="small" :pagination="{ pageSize: 20 }">
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'status'">
              <a-badge :status="STATUS_COLOR[(record as TimeSlot).status] as any" :text="(record as TimeSlot).status" />
            </template>
            <template v-if="column.key === 'actions'">
              <a-button
                v-if="(record as TimeSlot).status === 'booked'"
                type="link" danger size="small"
                @click="releaseSlot(record as TimeSlot)"
              >强制释放</a-button>
            </template>
          </template>
        </a-table>
        <a-empty v-else description="选择左侧号源池查看明细" />
      </a-card>
    </div>
  </div>
</template>
