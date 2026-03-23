<!-- 核心目的：根组件 -->
<!-- 模块功能：路由出口、全局布局容器 -->
<script setup lang="ts">
import zhCN from 'ant-design-vue/es/locale/zh_CN'
</script>

<template>
  <a-config-provider :locale="zhCN">
    <router-view />
  </a-config-provider>
</template>
