<!-- 核心目的：顶部栏 -->
<!-- 模块功能：用户信息、通知图标、退出登录 -->
<script setup lang="ts">
import { useAuth } from '@/composables/useAuth'

defineProps<{ title?: string }>()

const { currentUser, logout } = useAuth()
</script>

<template>
  <a-layout-header style="background: #fff; padding: 0 24px; height: 56px; line-height: 56px; display: flex; align-items: center; justify-content: space-between; box-shadow: 0 1px 4px rgba(0,21,41,.08); z-index: 10; position: relative;">
    <span style="font-size: 15px; font-weight: 600; color: #262626;">{{ title }}</span>
    <div style="display: flex; align-items: center; gap: 16px;">
      <a-dropdown>
        <a style="display: flex; align-items: center; gap: 8px; color: #434343; cursor: pointer; padding: 4px 8px; border-radius: 6px; transition: background .2s;" class="hover:bg-gray-50">
          <a-avatar size="small" style="background: linear-gradient(135deg, #1677ff, #0958d9); font-size: 13px;">{{ currentUser?.real_name?.[0] ?? 'U' }}</a-avatar>
          <div style="display: flex; flex-direction: column; align-items: flex-start; line-height: 1.3;">
            <span style="font-size: 13px; font-weight: 500;">{{ currentUser?.real_name ?? '用户' }}</span>
            <span style="color: #8c8c8c; font-size: 11px;">{{ currentUser?.role_name }}</span>
          </div>
          <span class="i-ant-design:down-outlined" style="font-size: 11px; color: #8c8c8c;" />
        </a>
        <template #overlay>
          <a-menu>
            <a-menu-item key="logout" @click="logout">
              <span class="i-ant-design:logout-outlined mr-1" />退出登录
            </a-menu-item>
          </a-menu>
        </template>
      </a-dropdown>
    </div>
  </a-layout-header>
</template>
