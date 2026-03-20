/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { h, ref, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useUserStore } from '@/stores/user';
const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const collapsed = ref(false);
const menuItems = computed(() => [
    {
        key: '/rule',
        icon: () => h('span', { class: 'i-ant-design:apartment-outlined' }),
        label: '规则引擎',
        children: [
            { key: '/rule/conflicts', label: '冲突规则' },
            { key: '/rule/conflict-packages', label: '冲突包' },
            { key: '/rule/dependencies', label: '依赖规则' },
            { key: '/rule/priority-tags', label: '优先级标签' },
            { key: '/rule/sorting-strategy', label: '排序策略' },
        ],
    },
    {
        key: '/resource',
        icon: () => h('span', { class: 'i-ant-design:database-outlined' }),
        label: '资源管理',
        children: [
            { key: '/resource/devices', label: '设备管理' },
            { key: '/resource/exam-items', label: '检查项目' },
            { key: '/resource/slot-pools', label: '号源池' },
            { key: '/resource/schedules', label: '排班日历' },
            { key: '/resource/item-aliases', label: '项目别名' },
        ],
    },
    {
        key: '/appointment',
        icon: () => h('span', { class: 'i-ant-design:calendar-outlined' }),
        label: '预约服务',
        children: [
            { key: '/appointment/list', label: '预约列表' },
            { key: '/appointment/auto', label: '智能预约' },
            { key: '/appointment/combo', label: '套餐预约' },
            { key: '/appointment/manual', label: '人工干预' },
            { key: '/appointment/blacklist', label: '黑名单' },
        ],
    },
    {
        key: '/triage',
        icon: () => h('span', { class: 'i-ant-design:team-outlined' }),
        label: '分诊叫号',
        children: [
            { key: '/triage/checkin', label: '签到台' },
            { key: '/triage/queue', label: '等候队列' },
            { key: '/triage/call', label: '叫号台' },
            { key: '/triage/screen', label: '大屏显示' },
        ],
    },
    {
        key: '/analytics',
        icon: () => h('span', { class: 'i-ant-design:bar-chart-outlined' }),
        label: '统计分析',
        children: [
            { key: '/analytics/dashboard', label: '数据看板' },
            { key: '/analytics/report', label: '报表导出' },
        ],
    },
    {
        key: '/optimization',
        icon: () => h('span', { class: 'i-ant-design:thunderbolt-outlined' }),
        label: '效能优化',
        children: [
            { key: '/optimization/metrics', label: '效率指标' },
            { key: '/optimization/alerts', label: '瓶颈告警' },
            { key: '/optimization/strategies', label: '优化策略' },
            { key: '/optimization/scans', label: '周期扫描' },
        ],
    },
    {
        key: '/admin',
        icon: () => h('span', { class: 'i-ant-design:setting-outlined' }),
        label: '系统管理',
        children: [
            { key: '/admin/users', label: '用户管理' },
            { key: '/admin/roles', label: '角色管理' },
        ],
    },
]);
const selectedKeys = computed(() => [route.path]);
const openKeys = ref([]);
function onMenuClick({ key }) {
    router.push(key);
}
debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
const __VLS_ctx = {};
let __VLS_components;
let __VLS_directives;
const __VLS_0 = {}.ALayoutSider;
/** @type {[typeof __VLS_components.ALayoutSider, typeof __VLS_components.aLayoutSider, typeof __VLS_components.ALayoutSider, typeof __VLS_components.aLayoutSider, ]} */ ;
// @ts-ignore
const __VLS_1 = __VLS_asFunctionalComponent(__VLS_0, new __VLS_0({
    collapsed: (__VLS_ctx.collapsed),
    collapsible: true,
    width: (220),
    ...{ style: {} },
}));
const __VLS_2 = __VLS_1({
    collapsed: (__VLS_ctx.collapsed),
    collapsible: true,
    width: (220),
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_1));
var __VLS_4 = {};
__VLS_3.slots.default;
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ style: {} },
});
if (!__VLS_ctx.collapsed) {
    __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({});
}
else {
    __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({});
}
const __VLS_5 = {}.AMenu;
/** @type {[typeof __VLS_components.AMenu, typeof __VLS_components.aMenu, ]} */ ;
// @ts-ignore
const __VLS_6 = __VLS_asFunctionalComponent(__VLS_5, new __VLS_5({
    ...{ 'onClick': {} },
    ...{ 'onOpenChange': {} },
    theme: "dark",
    mode: "inline",
    selectedKeys: (__VLS_ctx.selectedKeys),
    openKeys: (__VLS_ctx.openKeys),
    inlineCollapsed: (__VLS_ctx.collapsed),
    items: (__VLS_ctx.menuItems),
}));
const __VLS_7 = __VLS_6({
    ...{ 'onClick': {} },
    ...{ 'onOpenChange': {} },
    theme: "dark",
    mode: "inline",
    selectedKeys: (__VLS_ctx.selectedKeys),
    openKeys: (__VLS_ctx.openKeys),
    inlineCollapsed: (__VLS_ctx.collapsed),
    items: (__VLS_ctx.menuItems),
}, ...__VLS_functionalComponentArgsRest(__VLS_6));
let __VLS_9;
let __VLS_10;
let __VLS_11;
const __VLS_12 = {
    onClick: (__VLS_ctx.onMenuClick)
};
const __VLS_13 = {
    onOpenChange: ((keys) => __VLS_ctx.openKeys = keys)
};
var __VLS_8;
var __VLS_3;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            collapsed: collapsed,
            menuItems: menuItems,
            selectedKeys: selectedKeys,
            openKeys: openKeys,
            onMenuClick: onMenuClick,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
