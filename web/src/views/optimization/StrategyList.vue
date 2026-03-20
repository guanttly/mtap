<!-- 核心目的：StrategyList页面 -->
<!-- 模块功能：智能效能优化-相关功能页面 -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { useRouter } from 'vue-router'
import { optimizationApi } from '@/api/optimization'
import { usePagination } from '@/composables/usePagination'
import type { OptimizationStrategy } from '@/types/optimization'

const router = useRouter()
const { loading, items, pagination, fetchData, onTableChange } = usePagination<OptimizationStrategy>(
  params => optimizationApi.listStrategies(params),
)
onMounted(() => fetchData())

const approveModal = ref(false)
const approveTarget = ref<OptimizationStrategy | null>(null)
const approveForm = ref({ trial_days: 7, gray_scope: null })
const approving = ref(false)

const STATUS_MAP: Partial<Record<string, { color: string, label: string }>> = {
  pending_review: { color: 'default', label: '待审批' },
  trial_running: { color: 'processing', label: '试运行中' },
  trial_running_b: { color: 'processing', label: '试运行B' },
  submitted_approval: { color: 'gold', label: '已提交审批' },
  pending_eval: { color: 'blue', label: '待评估' },
  promoted: { color: 'success', label: '已推全量' },
  normalized: { color: 'success', label: '已常态化' },
  rolled_back: { color: 'warning', label: '已回滚' },
  rejected: { color: 'error', label: '已拒绝' },
  archived: { color: 'default', label: '已归档' },
}

const columns = [
  { title: '策略标题', dataIndex: 'title', key: 'title' },
  { title: '类别', dataIndex: 'category', key: 'category' },
  { title: '预期收益', dataIndex: 'expected_benefit', key: 'benefit' },
  { title: '状态', dataIndex: 'status', key: 'status' },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at' },
  { title: '操作', key: 'actions' },
]

function openApprove(record: OptimizationStrategy) {
  approveTarget.value = record
  approveForm.value = { trial_days: 7, gray_scope: null }
  approveModal.value = true
}

async function handleApprove() {
  if (!approveTarget.value) return
  approving.value = true
  try {
    await optimizationApi.approveStrategy(approveTarget.value.id, approveForm.value)
    message.success('已批准并进入试运行')
    approveModal.value = false
    fetchData()
  }
  finally { approving.value = false }
}

function handleReject(record: OptimizationStrategy) {
  Modal.confirm({
    title: '拒绝策略',
    content: '确定要拒绝该优化策略吗？',
    okType: 'danger',
    onOk: async () => {
      await optimizationApi.rejectStrategy(record.id, '人工审核拒绝')
      message.success('已拒绝')
      fetchData()
    },
  })
}
</script>

<template>
  <div>
    <a-table :columns="columns" :data-source="items" :loading="loading" :pagination="pagination" row-key="id" size="middle" @change="onTableChange">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag :color="STATUS_MAP[(record as OptimizationStrategy).status]?.color">
            {{ STATUS_MAP[(record as OptimizationStrategy).status]?.label ?? (record as OptimizationStrategy).status }}
          </a-tag>
        </template>
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button type="link" size="small" @click="router.push(`/optimization/strategies/${(record as OptimizationStrategy).id}`)">详情</a-button>
            <a-button v-if="(record as OptimizationStrategy).status === 'pending_review'" type="link" size="small" @click="openApprove(record as OptimizationStrategy)">批准</a-button>
            <a-button v-if="(record as OptimizationStrategy).status === 'pending_review'" type="link" danger size="small" @click="handleReject(record as OptimizationStrategy)">拒绝</a-button>
            <a-button v-if="(record as OptimizationStrategy).status === 'trial_running'" type="link" size="small" @click="router.push(`/optimization/trials/${(record as OptimizationStrategy).trial_run?.id}`)">监控</a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <a-modal v-model:open="approveModal" title="批准策略" :confirm-loading="approving" @ok="handleApprove">
      <a-form :model="approveForm" layout="vertical">
        <a-form-item label="试运行天数"><a-input-number v-model:value="approveForm.trial_days" :min="1" :max="30" style="width:100%;" /></a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>
