<!-- 核心目的：SortingStrategyForm页面 -->
<!-- 模块功能：规则引擎-SortingStrategy配置表单CRUD与配置 -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { ruleApi } from '@/api/rule'
import type { SortingStrategy } from '@/types/rule'

const loading = ref(false)
const saving = ref(false)
const form = ref<Partial<SortingStrategy>>({
  type: 'shortest_wait',
  scope_campuses: [],
  scope_depts: [],
  scope_devices: [],
})

// 用逗号分隔字符串辅助输入
const campusesText = ref('')
const deptsText = ref('')
const devicesText = ref('')

async function fetchData() {
  loading.value = true
  try {
    const strategy = await ruleApi.getSortingStrategy()
    form.value = strategy
    campusesText.value = strategy.scope_campuses?.join(',') ?? ''
    deptsText.value = strategy.scope_depts?.join(',') ?? ''
    devicesText.value = strategy.scope_devices?.join(',') ?? ''
  }
  catch {}
  finally { loading.value = false }
}
onMounted(fetchData)

async function handleSave() {
  form.value.scope_campuses = campusesText.value.split(',').map(s => s.trim()).filter(Boolean)
  form.value.scope_depts = deptsText.value.split(',').map(s => s.trim()).filter(Boolean)
  form.value.scope_devices = devicesText.value.split(',').map(s => s.trim()).filter(Boolean)
  saving.value = true
  try {
    await ruleApi.saveSortingStrategy(form.value)
    message.success('策略保存成功')
  }
  finally { saving.value = false }
}
</script>

<template>
  <div style="max-width: 600px;">
    <a-spin :spinning="loading">
      <a-form :model="form" layout="vertical">
        <a-card title="排序策略" size="small" style="margin-bottom: 16px;">
          <a-form-item label="策略类型">
            <a-radio-group v-model:value="form.type">
              <a-radio value="shortest_wait">最短等待</a-radio>
              <a-radio value="nearest">最近就近</a-radio>
              <a-radio value="priority">优先级优先</a-radio>
            </a-radio-group>
          </a-form-item>
          <a-form-item label="生效日期范围">
            <a-date-picker v-model:value="form.start_date" value-format="YYYY-MM-DD" placeholder="开始日期" style="width:48%;" />
            <span style="margin: 0 4%;">—</span>
            <a-date-picker v-model:value="form.end_date" value-format="YYYY-MM-DD" placeholder="结束日期" style="width:48%;" />
          </a-form-item>
        </a-card>
        <a-card title="生效范围（留空表示全局）" size="small" style="margin-bottom: 16px;">
          <a-form-item label="院区（逗号分隔ID）">
            <a-input v-model:value="campusesText" placeholder="campus_id1,campus_id2" />
          </a-form-item>
          <a-form-item label="科室（逗号分隔ID）">
            <a-input v-model:value="deptsText" placeholder="dept_id1,dept_id2" />
          </a-form-item>
          <a-form-item label="设备（逗号分隔ID）">
            <a-input v-model:value="devicesText" placeholder="device_id1,device_id2" />
          </a-form-item>
        </a-card>
        <a-button type="primary" :loading="saving" @click="handleSave">保存策略</a-button>
      </a-form>
    </a-spin>
  </div>
</template>
