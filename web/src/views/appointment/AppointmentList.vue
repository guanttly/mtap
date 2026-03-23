<!-- 核心目的：AppointmentList页面 -->
<!-- 模块功能：预约服务-相关功能页面 -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { appointmentApi } from '@/api/appointment'
import { usePagination } from '@/composables/usePagination'
import type { Appointment } from '@/types/appointment'

const { loading, items, pagination, fetchData, onTableChange, search } = usePagination<Appointment>(
  params => appointmentApi.listAppointments(params),
)
onMounted(() => fetchData())

const STATUS_COLOR: Record<string, string> = {
  pending: 'default', confirmed: 'processing', paid: 'blue', checked_in: 'cyan',
  completed: 'success', cancelled: 'error', no_show: 'warning',
}
const STATUS_LABEL: Record<string, string> = {
  pending: '待确认', confirmed: '已确认', paid: '已缴费', checked_in: '已签到',
  completed: '已完成', cancelled: '已取消', no_show: '爽约',
}

const columns = [
  { title: '预约号', dataIndex: 'appointment_no', key: 'no', width: 160 },
  { title: '患者', dataIndex: 'patient_name', key: 'patient_name' },
  { title: '状态', dataIndex: 'status', key: 'status' },
  { title: '来源', dataIndex: 'source', key: 'source' },
  { title: '操作', key: 'actions' },
]

const searchKeyword = ref('')

function handleSearch() {
  search({ keyword: searchKeyword.value })
}

async function handleCancel(record: Appointment) {
  await appointmentApi.cancel(record.id, '人工取消')
  message.success('已取消')
  fetchData()
}
</script>

<template>
  <a-card class="list-card" :bordered="false">
    <template #title>预约列表</template>

    <!-- 搜索工具栏 -->
    <div class="list-toolbar">
      <a-input-search
        v-model:value="searchKeyword"
        placeholder="患者姓名 / 就诊号"
        style="width: 260px;"
        allow-clear
        @search="handleSearch"
      />
      <a-select
        style="width: 140px;"
        placeholder="状态筛选"
        allow-clear
        @change="(v: string) => search({ status: v })"
      >
        <a-select-option v-for="(label, val) in STATUS_LABEL" :key="val" :value="val">
          {{ label }}
        </a-select-option>
      </a-select>
    </div>

    <a-table
      :columns="columns"
      :data-source="items"
      :loading="loading"
      :pagination="pagination"
      row-key="id"
      size="middle"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag :color="STATUS_COLOR[(record as Appointment).status]">
            {{ STATUS_LABEL[(record as Appointment).status] }}
          </a-tag>
        </template>
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button type="link" size="small">详情</a-button>
            <a-button
              v-if="['pending','confirmed','paid'].includes((record as Appointment).status)"
              type="link"
              danger
              size="small"
              @click="handleCancel(record as Appointment)"
            >取消</a-button>
          </a-space>
        </template>
      </template>
    </a-table>
  </a-card>
</template>
