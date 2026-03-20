<!-- 核心目的：患者信息卡片 -->
<!-- 模块功能：展示患者基本信息（脱敏姓名、年龄、性别、来源标签） -->
<script setup lang="ts">
import { computed } from 'vue'
import { maskName } from '@/utils/desensitize'

const props = defineProps<{
  patientId: string
  name: string
  age?: number
  gender?: 'M' | 'F' | string
  source?: string
  tags?: string[]
}>()

const maskedName = computed(() => maskName(props.name))
const genderText = computed(() => props.gender === 'M' ? '男' : props.gender === 'F' ? '女' : props.gender ?? '-')
const genderColor = computed(() => props.gender === 'M' ? 'blue' : props.gender === 'F' ? 'pink' : 'default')
</script>

<template>
  <a-card size="small" :bodyStyle="{ padding: '12px' }">
    <div style="display: flex; align-items: center; gap: 12px;">
      <a-avatar style="background: #1890ff; flex-shrink: 0;">{{ maskedName[0] }}</a-avatar>
      <div style="flex: 1; min-width: 0;">
        <div style="font-weight: 600; font-size: 15px;">{{ maskedName }}</div>
        <div style="display: flex; gap: 8px; align-items: center; margin-top: 4px; flex-wrap: wrap;">
          <a-tag :color="genderColor" style="margin: 0;">{{ genderText }}</a-tag>
          <span v-if="age" class="text-muted">{{ age }}岁</span>
          <a-tag v-if="source" color="processing" style="margin: 0;">{{ source }}</a-tag>
          <a-tag v-for="tag in tags" :key="tag" color="gold" style="margin: 0;">{{ tag }}</a-tag>
        </div>
      </div>
      <span class="text-muted" style="font-size: 12px; flex-shrink: 0;">{{ patientId }}</span>
    </div>
  </a-card>
</template>
