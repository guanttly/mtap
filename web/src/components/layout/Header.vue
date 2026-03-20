<!-- 核心目的：顶部栏 -->
<!-- 模块功能：用户信息、通知图标、退出登录 -->
<script setup lang="ts">
import { useAuth } from '@/composables/useAuth'

defineProps<{ title?: string }>()

const { currentUser, logout } = useAuth()
</script>

<template>
  <a-layout-header style="background: #fff; padding: 0 16px; display: flex; align-items: center; justify-content: space-between; box-shadow: 0 1px 4px rgba(0,21,41,.08);">
    <span style="font-size: 16px; font-weight: 600; color: #1f2328;">{{ title }}</span>
    <div style="display: flex; align-items: center; gap: 16px;">
      <a-dropdown>
        <a style="display: flex; align-items: center; gap: 6px; color: #1f2328; cursor: pointer;">
          <a-avatar size="small" style="background: #1890ff;">{{ currentUser?.real_name?.[0] ?? 'U' }}</a-avatar>
          <span>{{ currentUser?.real_name ?? '用户' }}</span>
          <span style="color: #8c8c8c; font-size: 12px;">{{ currentUser?.role_name }}</span>
        </a>
        <template #overlay>
          <a-menu>
            <a-menu-item key="logout" @click="logout">退出登录</a-menu-item>
          </a-menu>
        </template>
      </a-dropdown>
    </div>
  </a-layout-header>
</template>
