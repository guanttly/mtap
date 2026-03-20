import { createRouter, createWebHistory } from 'vue-router';
import { useUserStore } from '@/stores/user';
const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/login',
            name: 'Login',
            component: () => import('@/views/login/LoginPage.vue'),
            meta: { public: true },
        },
        {
            path: '/',
            component: () => import('@/components/layout/AppLayout.vue'),
            redirect: '/rule/conflicts',
            children: [
                // 规则引擎
                { path: 'rule/conflicts', name: 'ConflictRuleList', component: () => import('@/views/rule/ConflictRuleList.vue'), meta: { title: '冲突规则' } },
                { path: 'rule/conflict-packages', name: 'ConflictPackageList', component: () => import('@/views/rule/ConflictPackageList.vue'), meta: { title: '冲突包' } },
                { path: 'rule/dependencies', name: 'DependencyRuleList', component: () => import('@/views/rule/DependencyRuleList.vue'), meta: { title: '依赖规则' } },
                { path: 'rule/priority-tags', name: 'PriorityTagList', component: () => import('@/views/rule/PriorityTagList.vue'), meta: { title: '优先级标签' } },
                { path: 'rule/sorting-strategy', name: 'SortingStrategyForm', component: () => import('@/views/rule/SortingStrategyForm.vue'), meta: { title: '排序策略' } },
                { path: 'rule/patient-adapt', name: 'PatientAdaptRuleList', component: () => import('@/views/rule/PatientAdaptRuleList.vue'), meta: { title: '患者适配规则' } },
                { path: 'rule/source-controls', name: 'SourceControlList', component: () => import('@/views/rule/SourceControlList.vue'), meta: { title: '开单来源控制' } },
                // 资源管理
                { path: 'resource/devices', name: 'DeviceList', component: () => import('@/views/resource/DeviceList.vue'), meta: { title: '设备管理' } },
                { path: 'resource/exam-items', name: 'ExamItemList', component: () => import('@/views/resource/ExamItemList.vue'), meta: { title: '检查项目' } },
                { path: 'resource/schedules', name: 'ScheduleCalendar', component: () => import('@/views/resource/ScheduleCalendar.vue'), meta: { title: '排班管理' } },
                { path: 'resource/slot-pools', name: 'SlotPoolView', component: () => import('@/views/resource/SlotPoolView.vue'), meta: { title: '号源池' } },
                { path: 'resource/item-aliases', name: 'ItemAliasManager', component: () => import('@/views/resource/ItemAliasManager.vue'), meta: { title: '项目别名' } },
                // 预约服务
                { path: 'appointment/list', name: 'AppointmentList', component: () => import('@/views/appointment/AppointmentList.vue'), meta: { title: '预约列表' } },
                { path: 'appointment/auto', name: 'AutoAppointment', component: () => import('@/views/appointment/AutoAppointment.vue'), meta: { title: '自动预约' } },
                { path: 'appointment/combo', name: 'ComboAppointment', component: () => import('@/views/appointment/ComboAppointment.vue'), meta: { title: '组合预约' } },
                { path: 'appointment/manual', name: 'ManualOverride', component: () => import('@/views/appointment/ManualOverride.vue'), meta: { title: '人工干预' } },
                { path: 'appointment/blacklist', name: 'BlacklistManager', component: () => import('@/views/appointment/BlacklistManager.vue'), meta: { title: '黑名单' } },
                // 分诊执行
                { path: 'triage/checkin', name: 'CheckInStation', component: () => import('@/views/triage/CheckInStation.vue'), meta: { title: '签到站' } },
                { path: 'triage/queue', name: 'WaitingQueueView', component: () => import('@/views/triage/WaitingQueueView.vue'), meta: { title: '候诊队列' } },
                { path: 'triage/call', name: 'NurseCallPanel', component: () => import('@/views/triage/NurseCallPanel.vue'), meta: { title: '护士呼叫' } },
                { path: 'triage/screen', name: 'TriageScreen', component: () => import('@/views/triage/TriageScreen.vue'), meta: { title: '分诊大屏' } },
                // 统计分析
                { path: 'analytics/dashboard', name: 'Dashboard', component: () => import('@/views/analytics/Dashboard.vue'), meta: { title: '实时监控' } },
                { path: 'analytics/report', name: 'ReportExport', component: () => import('@/views/analytics/ReportExport.vue'), meta: { title: '报表导出' } },
                // 效能优化
                { path: 'optimization/metrics', name: 'MetricsDashboard', component: () => import('@/views/optimization/MetricsDashboard.vue'), meta: { title: '效率指标' } },
                { path: 'optimization/alerts', name: 'BottleneckAlerts', component: () => import('@/views/optimization/BottleneckAlerts.vue'), meta: { title: '瓶颈告警' } },
                { path: 'optimization/strategies', name: 'StrategyList', component: () => import('@/views/optimization/StrategyList.vue'), meta: { title: '优化策略' } },
                { path: 'optimization/strategies/:id', name: 'StrategyDetail', component: () => import('@/views/optimization/StrategyDetail.vue'), meta: { title: '策略详情' } },
                { path: 'optimization/trials/:id', name: 'TrialMonitor', component: () => import('@/views/optimization/TrialMonitor.vue'), meta: { title: '试运行监控' } },
                { path: 'optimization/evaluations/:id', name: 'EvaluationReport', component: () => import('@/views/optimization/EvaluationReport.vue'), meta: { title: '评估报告' } },
                { path: 'optimization/roi/:id', name: 'ROIReport', component: () => import('@/views/optimization/ROIReport.vue'), meta: { title: 'ROI报告' } },
                { path: 'optimization/scans', name: 'PerformanceScan', component: () => import('@/views/optimization/PerformanceScan.vue'), meta: { title: '效能扫描' } },
                // 系统管理
                { path: 'admin/users', name: 'UserManagement', component: () => import('@/views/admin/UserList.vue'), meta: { title: '用户管理' } },
                { path: 'admin/roles', name: 'RoleManagement', component: () => import('@/views/admin/RoleList.vue'), meta: { title: '角色管理' } },
            ],
        },
        { path: '/:pathMatch(.*)*', redirect: '/' },
    ],
});
// 路由守卫
router.beforeEach((to) => {
    const userStore = useUserStore();
    if (!to.meta.public && !userStore.isLoggedIn) {
        return { name: 'Login', query: { redirect: to.fullPath } };
    }
    if (to.name === 'Login' && userStore.isLoggedIn) {
        return { path: '/' };
    }
});
export default router;
