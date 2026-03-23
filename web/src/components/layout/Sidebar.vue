<!-- 核心目的：侧边导航栏 -->
<!-- 模块功能：菜单树渲染、路由导航、折叠/展开 -->
<script setup lang="ts">
import { h, ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import SvgIcon from '@/components/common/SvgIcon.vue'
import type { IconName } from '@/assets/icons/index'

/** 统一生成菜单图标渲染函数，使用本地 SVG 资源 */
function icon(name: IconName) {
  return () => h(SvgIcon, { name, size: '1em' })
}

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const collapsed = ref(false)

const menuItems = computed(() => [
  {
    key: '/rule',
    icon: icon('apartment-outlined'),
    label: '规则引擎',
    children: [
      { key: '/rule/conflicts', icon: icon('interaction-outlined'), label: '冲突规则' },
      { key: '/rule/conflict-packages', icon: icon('inbox-outlined'), label: '冲突包' },
      { key: '/rule/dependencies', icon: icon('link-outlined'), label: '依赖规则' },
      { key: '/rule/priority-tags', icon: icon('tag-outlined'), label: '优先级标签' },
      { key: '/rule/sorting-strategy', icon: icon('sort-ascending-outlined'), label: '排序策略' },
      { key: '/rule/patient-adapt', icon: icon('solution-outlined'), label: '患者适配规则' },
      { key: '/rule/source-controls', icon: icon('control-outlined'), label: '来源控制' },
    ],
  },
  {
    key: '/resource',
    icon: icon('database-outlined'),
    label: '资源管理',
    children: [
      { key: '/resource/devices', icon: icon('desktop-outlined'), label: '设备管理' },
      { key: '/resource/exam-items', icon: icon('file-search-outlined'), label: '检查项目' },
      { key: '/resource/slot-pools', icon: icon('number-outlined'), label: '号源池' },
      { key: '/resource/schedules', icon: icon('schedule-outlined'), label: '排班日历' },
      { key: '/resource/item-aliases', icon: icon('edit-outlined'), label: '项目别名' },
    ],
  },
  {
    key: '/appointment',
    icon: icon('calendar-outlined'),
    label: '预约服务',
    children: [
      { key: '/appointment/list', icon: icon('unordered-list-outlined'), label: '预约列表' },
      { key: '/appointment/auto', icon: icon('robot-outlined'), label: '智能预约' },
      { key: '/appointment/combo', icon: icon('shopping-outlined'), label: '套餐预约' },
      { key: '/appointment/manual', icon: icon('tool-outlined'), label: '人工干预' },
      { key: '/appointment/blacklist', icon: icon('stop-outlined'), label: '黑名单' },
    ],
  },
  {
    key: '/triage',
    icon: icon('team-outlined'),
    label: '分诊叫号',
    children: [
      { key: '/triage/checkin', icon: icon('check-circle-outlined'), label: '签到台' },
      { key: '/triage/queue', icon: icon('ordered-list-outlined'), label: '等候队列' },
      { key: '/triage/call', icon: icon('notification-outlined'), label: '叫号台' },
      { key: '/triage/screen', icon: icon('fund-projection-screen-outlined'), label: '大屏显示' },
    ],
  },
  {
    key: '/analytics',
    icon: icon('bar-chart-outlined'),
    label: '统计分析',
    children: [
      { key: '/analytics/dashboard', icon: icon('dashboard-outlined'), label: '数据看板' },
      { key: '/analytics/report', icon: icon('file-excel-outlined'), label: '报表导出' },
    ],
  },
  {
    key: '/optimization',
    icon: icon('thunderbolt-outlined'),
    label: '效能优化',
    children: [
      { key: '/optimization/metrics', icon: icon('rise-outlined'), label: '效率指标' },
      { key: '/optimization/alerts', icon: icon('alert-outlined'), label: '瓶颈告警' },
      { key: '/optimization/strategies', icon: icon('bulb-outlined'), label: '优化策略' },
      { key: '/optimization/scans', icon: icon('scan-outlined'), label: '周期扫描' },
    ],
  },
  {
    key: '/admin',
    icon: icon('setting-outlined'),
    label: '系统管理',
    children: [
      { key: '/admin/users', icon: icon('user-outlined'), label: '用户管理' },
      { key: '/admin/roles', icon: icon('idcard-outlined'), label: '角色管理' },
    ],
  },
])

const selectedKeys = computed(() => [route.path])

// 根据当前路由自动打开对应父级菜单，并在路由变化时保持同步
function getParentKey(path: string) {
  const seg = '/' + path.split('/')[1]
  return seg
}
const openKeys = ref<string[]>([getParentKey(route.path)])
watch(() => route.path, (path) => {
  openKeys.value = [getParentKey(path)]
})

function onOpenChange(keys: string[]) {
  // 手风琴：找到本次新展开的那个 key，只保留它；若是收起操作则置空
  const latest = keys.find(k => !openKeys.value.includes(k))
  openKeys.value = latest ? [latest] : []
}

function onMenuClick({ key }: { key: string }) {
  if (key === '/triage/screen') {
    window.open('/triage/screen', '_blank')
    return
  }
  router.push(key)
}
</script>

<template>
  <a-layout-sider
    v-model:collapsed="collapsed"
    :trigger="null"
    :width="220"
    :collapsed-width="64"
    class="app-sider"
  >
    <!-- 浮动折叠按钮：贴在侧边栏右边缘 -->
    <button
      class="sider-collapse-toggle"
      @click="collapsed = !collapsed"
      :title="collapsed ? '展开侧边栏' : '收起侧边栏'"
    >
      <svg width="11" height="11" viewBox="0 0 1024 1024" class="collapse-icon" :class="{ 'is-collapsed': collapsed }">
        <path fill="currentColor" d="M724 218.3V141c0-6.7-7.7-10.4-12.9-6.3L260.3 486.8a31.86 31.86 0 0 0 0 50.3l450.8 352.1c5.3 4.1 12.9.4 12.9-6.3v-77.3c0-4.9-2.3-9.6-6.1-12.6l-360-281 360-281.1c3.8-3 6.1-7.7 6.1-12.6z"/>
      </svg>
    </button>
    <!-- 菜单区域 -->
    <div class="sider-scroll">
      <a-menu
        theme="light"
        mode="inline"
        :selected-keys="selectedKeys"
        :open-keys="openKeys"
        :items="menuItems"
        @click="onMenuClick"
        @openChange="onOpenChange"
      />
    </div>
  </a-layout-sider>
</template>

<style scoped>
/* 浮动折叠按钮：半圆形贴在右边缘 */
.sider-collapse-toggle {
  position: absolute;
  top: 14px;
  right: -13px;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  border: 1px solid #e8ecf0;
  border-radius: 50%;
  background: #fff;
  color: #b0bdd0;
  cursor: pointer;
  box-shadow: 2px 0 6px rgba(0, 0, 0, 0.06);
  transition: background 0.15s, color 0.15s, box-shadow 0.15s;
  padding: 0;
}
.sider-collapse-toggle:hover {
  background: #f0f5ff;
  color: #1260c8;
  box-shadow: 2px 0 8px rgba(18, 96, 200, 0.15);
}
.collapse-icon {
  flex-shrink: 0;
  transition: transform 0.22s ease;
}
.collapse-icon.is-collapsed {
  transform: rotate(180deg);
}
/* 菜单滚动区 */
.sider-scroll {
  height: calc(100vh - 56px);
  overflow-y: auto;
  overflow-x: hidden;
}
</style>
