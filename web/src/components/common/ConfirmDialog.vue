<!-- 核心目的：通用确认弹窗 -->
<!-- 模块功能：操作确认、危险操作二次确认 -->
<script setup lang="ts">
import { Modal } from 'ant-design-vue'

export interface ConfirmOptions {
  title?: string
  content?: string
  okText?: string
  danger?: boolean
}

function confirm(opts: ConfirmOptions = {}): Promise<void> {
  return new Promise((resolve, reject) => {
    Modal.confirm({
      title: opts.title ?? '确认操作',
      content: opts.content ?? '确定要执行此操作吗？',
      okText: opts.okText ?? '确定',
      cancelText: '取消',
      okType: opts.danger ? 'danger' : 'primary',
      onOk: () => resolve(),
      onCancel: () => reject(new Error('cancelled')),
    })
  })
}

defineExpose({ confirm })
</script>

<template>
  <span><slot /></span>
</template>
