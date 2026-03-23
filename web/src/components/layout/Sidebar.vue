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
  <a-layout-sider v-model:collapsed="collapsed" collapsible :width="220" style="background: #001529; height: 100vh; position: sticky; top: 0; flex-shrink: 0;">
    <div style="height: 48px; display: flex; align-items: center; justify-content: center; color: #fff; font-weight: bold; font-size: 16px; overflow: hidden; padding: 0 8px;">
      <span v-if="!collapsed">MTAP</span>
      <span v-else>M</span>
    </div>
    <!-- calc: 100vh - 48px logo - 48px collapse trigger -->
    <div style="height: calc(100vh - 96px); overflow-y: auto; overflow-x: hidden;">
      <a-menu
        theme="dark"
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
