<!-- 核心目的：应用主布局 -->
<!-- 模块功能：侧边栏 + 顶栏 + 内容区域的整体布局框架 -->
<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import Sidebar from './Sidebar.vue'
import AppHeader from './Header.vue'

const route = useRoute()
const pageTitle = computed(() => String(route.meta.title ?? ''))
</script>

<template>
  <a-layout style="min-height: 100vh">
    <Sidebar />
    <a-layout>
      <AppHeader :title="pageTitle" />
      <a-layout-content style="margin: 0; overflow: auto">
        <div class="page-container">
          <router-view />
        </div>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>
