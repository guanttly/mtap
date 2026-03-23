<!-- 核心目的：本地 SVG 图标渲染组件 -->
<!-- 模块功能：从 assets/icons 注册表加载图标，支持尺寸和颜色自定义 -->
<script setup lang="ts">
import { computed } from 'vue'
import { iconBodies, SVG_VIEWBOX, type IconName } from '@/assets/icons/index'

const props = withDefaults(defineProps<{
  /** 图标名称，对应 @/assets/icons/index.ts 中的 IconName */
  name: IconName
  /** 图标尺寸，支持数字 (px) 或带单位的字符串，默认 1em */
  size?: number | string
  /** 图标颜色，默认继承父级 currentColor */
  color?: string
  /** 自定义 class */
  class?: string
}>(), {
  size: '1em',
  color: 'currentColor',
})

const sizeValue = computed(() =>
  typeof props.size === 'number' ? `${props.size}px` : props.size,
)

const body = computed(() => iconBodies[props.name] ?? '')
</script>

<template>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    :viewBox="SVG_VIEWBOX"
    :width="sizeValue"
    :height="sizeValue"
    :fill="color"
    :class="['svg-icon', props.class]"
    aria-hidden="true"
    v-html="body"
  />
</template>

<style scoped>
.svg-icon {
  display: inline-block;
  vertical-align: -0.15em;
  flex-shrink: 0;
  line-height: 1;
}
</style>
