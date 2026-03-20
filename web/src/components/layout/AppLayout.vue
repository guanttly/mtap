<!-- 核心目的：应用主布局 -->
<!-- 模块功能：侧边栏 + 顶栏 + 内容区域的整体布局框架 -->
<script setup lang="ts">
import { computed, ref, onErrorCaptured } from 'vue'
import { useRoute } from 'vue-router'
import Sidebar from './Sidebar.vue'
import AppHeader from './Header.vue'

const route = useRoute()
const pageTitle = computed(() => String(route.meta.title ?? ''))

// 错误边界：子页面组件崩溃时，通过切换 key 强制重建 RouterView
// 防止 vnode 损坏导致路由表永久失灵
const routerKey = ref(0)
onErrorCaptured((err: unknown) => {
  const msg = err instanceof Error ? err.message : String(err)
  // 只拦截 vnode 相关的 TypeError，其他错误继续向上传播
  if (err instanceof TypeError && (msg.includes('vnode') || msg.includes('parentNode'))) {
    routerKey.value++
    return false
  }
})
</script>

<template>
  <a-layout style="height: 100vh; overflow: hidden">
    <Sidebar />
    <a-layout style="overflow: hidden">
      <AppHeader :title="pageTitle" />
      <a-layout-content style="margin: 0; overflow-y: auto; min-height: 0">
        <div class="page-container">
          <router-view :key="routerKey" />
        </div>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>
