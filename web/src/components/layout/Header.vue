<!-- 核心目的：顶部导航栏 -->
<!-- 模块功能：平台品牌展示、通知铃铛、用户信息与退出 -->
<script setup lang="ts">
import { useAuth } from '@/composables/useAuth'
import SvgIcon from '@/components/common/SvgIcon.vue'

defineProps<{ title?: string }>()

const { currentUser, logout } = useAuth()
</script>

<template>
  <a-layout-header class="app-header">
    <!-- 装饰气泡 -->
    <span class="deco deco-a" />
    <span class="deco deco-b" />
    <span class="deco deco-c" />
    <span class="deco deco-d" />

    <!-- 左：品牌 -->
    <div class="header-left">
      <SvgIcon name="medicine-box-outlined" :size="26" color="rgba(255,255,255,0.95)" />
      <span class="header-brand">MTAP 医疗预约管理平台</span>
    </div>

    <!-- 右：操作区 -->
    <div class="header-right">
      <!-- 通知铃铛 -->
      <a-badge :count="0" :offset="[-2, 4]">
        <span class="icon-btn">
          <SvgIcon name="bell-outlined" :size="18" color="rgba(255,255,255,0.85)" />
        </span>
      </a-badge>

      <!-- 用户下拉 -->
      <a-dropdown placement="bottomRight">
        <a class="user-trigger" @click.prevent>
          <a-avatar class="user-avatar" :size="30">
            {{ currentUser?.real_name?.[0] ?? 'U' }}
          </a-avatar>
          <div class="user-info">
            <span class="user-name">{{ currentUser?.real_name ?? '用户' }}</span>
            <span class="user-role">{{ currentUser?.role_name }}</span>
          </div>
          <SvgIcon name="down-outlined" :size="11" color="rgba(255,255,255,0.55)" />
        </a>
        <template #overlay>
          <a-menu>
            <a-menu-item key="logout" @click="logout">
              <SvgIcon name="logout-outlined" :size="14" style="margin-right:6px;vertical-align:-2px" />
              退出登录
            </a-menu-item>
          </a-menu>
        </template>
      </a-dropdown>
    </div>
  </a-layout-header>
</template>

<style scoped>
/* ===== 顶栏容器 ===== */
.app-header {
  position: relative;
  height: 56px !important;
  line-height: 56px !important;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  overflow: hidden;
  flex-shrink: 0;
  background: linear-gradient(90deg, #0b3d91 0%, #1260c8 55%, #1a6ad4 100%) !important;
  box-shadow: 0 2px 10px rgba(11, 61, 145, 0.45);
  z-index: 100;
}

/* ===== 装饰气泡 ===== */
.deco {
  position: absolute;
  border-radius: 50%;
  pointer-events: none;
}
.deco-a { width: 160px; height: 160px; top: -70px; right: 260px; background: rgba(255,255,255,0.06); }
.deco-b { width: 100px; height: 100px; top: -10px; right: 440px; background: rgba(100,210,255,0.09); }
.deco-c { width: 70px;  height: 70px;  bottom: -30px; right: 360px; background: rgba(100,210,255,0.08); }
.deco-d { width: 50px;  height: 50px;  top: 5px; right: 580px; background: rgba(255,255,255,0.05); }

/* ===== 左侧品牌 ===== */
.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
  z-index: 1;
}
.header-brand {
  font-size: 17px;
  font-weight: 700;
  color: #fff;
  letter-spacing: 0.04em;
  white-space: nowrap;
  text-shadow: 0 1px 4px rgba(0,0,0,0.2);
}

/* ===== 右侧操作区 ===== */
.header-right {
  display: flex;
  align-items: center;
  gap: 10px;
  z-index: 1;
}

.icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  cursor: pointer;
  transition: background 0.2s;
}
.icon-btn:hover { background: rgba(255,255,255,0.14); }

.user-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 10px;
  border-radius: 8px;
  cursor: pointer;
  text-decoration: none !important;
  transition: background 0.2s;
}
.user-trigger:hover { background: rgba(255,255,255,0.12); }

.user-avatar {
  background: rgba(255,255,255,0.22) !important;
  color: #fff !important;
  font-size: 13px !important;
  border: 2px solid rgba(255,255,255,0.38) !important;
  flex-shrink: 0;
}

.user-info {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  line-height: 1.25;
}
.user-name {
  font-size: 13px;
  font-weight: 500;
  color: #fff;
}
.user-role {
  font-size: 11px;
  color: rgba(255,255,255,0.6);
}
</style>
