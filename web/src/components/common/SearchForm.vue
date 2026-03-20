<!-- 核心目的：通用搜索表单 -->
<!-- 模块功能：动态筛选条件、搜索/重置按钮 -->
<script setup lang="ts">
import { reactive } from 'vue'

const props = withDefaults(defineProps<{
  loading?: boolean
}>(), { loading: false })

const emit = defineEmits<{
  search: [params: Record<string, unknown>]
  reset: []
}>()

const model = reactive<Record<string, unknown>>({})

function handleSearch() {
  const params: Record<string, unknown> = {}
  for (const [k, v] of Object.entries(model)) {
    if (v !== undefined && v !== null && v !== '') params[k] = v
  }
  emit('search', params)
}

function handleReset() {
  for (const k of Object.keys(model)) model[k] = undefined
  emit('reset')
  emit('search', {})
}

defineExpose({ model })
</script>

<template>
  <div class="search-bar">
    <a-form layout="inline" :model="model">
      <slot :model="model" />
      <a-form-item>
        <a-space>
          <a-button type="primary" :loading="loading" @click="handleSearch">查询</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-space>
      </a-form-item>
    </a-form>
  </div>
</template>
