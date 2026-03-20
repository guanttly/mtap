/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { ref, onMounted, onUnmounted } from 'vue';
import { analyticsApi } from '@/api/analytics';
import { useWebSocket } from '@/composables/useWebSocket';
import * as echarts from 'echarts/core';
import { LineChart, BarChart } from 'echarts/charts';
import { GridComponent, TooltipComponent, LegendComponent, TitleComponent } from 'echarts/components';
import { CanvasRenderer } from 'echarts/renderers';
echarts.use([LineChart, BarChart, GridComponent, TooltipComponent, LegendComponent, TitleComponent, CanvasRenderer]);
const loading = ref(false);
const snapshot = ref(null);
const trendChartEl = ref();
let trendChart = null;
// WebSocket 实时推送
const { status: wsStatus, connect, disconnect, on } = useWebSocket('/ws/dashboard');
on('dashboard_update', (payload) => {
    snapshot.value = payload;
    setTimeout(renderTrendChart, 50);
});
async function fetchData() {
    loading.value = true;
    try {
        snapshot.value = await analyticsApi.getDashboard();
        setTimeout(renderTrendChart, 100);
    }
    catch { }
    finally {
        loading.value = false;
    }
}
function renderTrendChart() {
    if (!trendChartEl.value || !snapshot.value?.wait_trend)
        return;
    if (!trendChart)
        trendChart = echarts.init(trendChartEl.value);
    const trend = snapshot.value.wait_trend;
    trendChart.setOption({
        tooltip: { trigger: 'axis' },
        xAxis: { type: 'category', data: trend.map((t) => t.time) },
        yAxis: { type: 'value', name: '分钟' },
        series: [{ name: '平均等待', type: 'line', data: trend.map((t) => t.avg_wait_min), smooth: true, lineStyle: { color: '#1890ff' }, areaStyle: { color: 'rgba(24,144,255,.1)' } }],
    });
}
onMounted(() => {
    fetchData();
    connect();
});
onUnmounted(() => {
    disconnect();
    trendChart?.dispose();
});
debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
const __VLS_ctx = {};
let __VLS_components;
let __VLS_directives;
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ style: {} },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ style: {} },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
    ...{ style: {} },
});
const __VLS_0 = {}.ABadge;
/** @type {[typeof __VLS_components.ABadge, typeof __VLS_components.aBadge, ]} */ ;
// @ts-ignore
const __VLS_1 = __VLS_asFunctionalComponent(__VLS_0, new __VLS_0({
    status: (__VLS_ctx.wsStatus === 'connected' ? 'success' : __VLS_ctx.wsStatus === 'connecting' ? 'processing' : 'error'),
    text: (__VLS_ctx.wsStatus === 'connected' ? '实时推送' : __VLS_ctx.wsStatus === 'connecting' ? '连接中' : '断开'),
    ...{ style: {} },
}));
const __VLS_2 = __VLS_1({
    status: (__VLS_ctx.wsStatus === 'connected' ? 'success' : __VLS_ctx.wsStatus === 'connecting' ? 'processing' : 'error'),
    text: (__VLS_ctx.wsStatus === 'connected' ? '实时推送' : __VLS_ctx.wsStatus === 'connecting' ? '连接中' : '断开'),
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_1));
const __VLS_4 = {}.AButton;
/** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
// @ts-ignore
const __VLS_5 = __VLS_asFunctionalComponent(__VLS_4, new __VLS_4({
    ...{ 'onClick': {} },
    loading: (__VLS_ctx.loading),
}));
const __VLS_6 = __VLS_5({
    ...{ 'onClick': {} },
    loading: (__VLS_ctx.loading),
}, ...__VLS_functionalComponentArgsRest(__VLS_5));
let __VLS_8;
let __VLS_9;
let __VLS_10;
const __VLS_11 = {
    onClick: (__VLS_ctx.fetchData)
};
__VLS_7.slots.default;
var __VLS_7;
const __VLS_12 = {}.ASpin;
/** @type {[typeof __VLS_components.ASpin, typeof __VLS_components.aSpin, typeof __VLS_components.ASpin, typeof __VLS_components.aSpin, ]} */ ;
// @ts-ignore
const __VLS_13 = __VLS_asFunctionalComponent(__VLS_12, new __VLS_12({
    spinning: (__VLS_ctx.loading),
}));
const __VLS_14 = __VLS_13({
    spinning: (__VLS_ctx.loading),
}, ...__VLS_functionalComponentArgsRest(__VLS_13));
__VLS_15.slots.default;
if (__VLS_ctx.snapshot) {
    const __VLS_16 = {}.ARow;
    /** @type {[typeof __VLS_components.ARow, typeof __VLS_components.aRow, typeof __VLS_components.ARow, typeof __VLS_components.aRow, ]} */ ;
    // @ts-ignore
    const __VLS_17 = __VLS_asFunctionalComponent(__VLS_16, new __VLS_16({
        gutter: (16),
        ...{ style: {} },
    }));
    const __VLS_18 = __VLS_17({
        gutter: (16),
        ...{ style: {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_17));
    __VLS_19.slots.default;
    const __VLS_20 = {}.ACol;
    /** @type {[typeof __VLS_components.ACol, typeof __VLS_components.aCol, typeof __VLS_components.ACol, typeof __VLS_components.aCol, ]} */ ;
    // @ts-ignore
    const __VLS_21 = __VLS_asFunctionalComponent(__VLS_20, new __VLS_20({
        span: (6),
    }));
    const __VLS_22 = __VLS_21({
        span: (6),
    }, ...__VLS_functionalComponentArgsRest(__VLS_21));
    __VLS_23.slots.default;
    const __VLS_24 = {}.ACard;
    /** @type {[typeof __VLS_components.ACard, typeof __VLS_components.aCard, typeof __VLS_components.ACard, typeof __VLS_components.aCard, ]} */ ;
    // @ts-ignore
    const __VLS_25 = __VLS_asFunctionalComponent(__VLS_24, new __VLS_24({
        size: "small",
    }));
    const __VLS_26 = __VLS_25({
        size: "small",
    }, ...__VLS_functionalComponentArgsRest(__VLS_25));
    __VLS_27.slots.default;
    const __VLS_28 = {}.AStatistic;
    /** @type {[typeof __VLS_components.AStatistic, typeof __VLS_components.aStatistic, ]} */ ;
    // @ts-ignore
    const __VLS_29 = __VLS_asFunctionalComponent(__VLS_28, new __VLS_28({
        title: "今日总号源",
        value: (__VLS_ctx.snapshot.slot_usage.total_slots),
    }));
    const __VLS_30 = __VLS_29({
        title: "今日总号源",
        value: (__VLS_ctx.snapshot.slot_usage.total_slots),
    }, ...__VLS_functionalComponentArgsRest(__VLS_29));
    var __VLS_27;
    var __VLS_23;
    const __VLS_32 = {}.ACol;
    /** @type {[typeof __VLS_components.ACol, typeof __VLS_components.aCol, typeof __VLS_components.ACol, typeof __VLS_components.aCol, ]} */ ;
    // @ts-ignore
    const __VLS_33 = __VLS_asFunctionalComponent(__VLS_32, new __VLS_32({
        span: (6),
    }));
    const __VLS_34 = __VLS_33({
        span: (6),
    }, ...__VLS_functionalComponentArgsRest(__VLS_33));
    __VLS_35.slots.default;
    const __VLS_36 = {}.ACard;
    /** @type {[typeof __VLS_components.ACard, typeof __VLS_components.aCard, typeof __VLS_components.ACard, typeof __VLS_components.aCard, ]} */ ;
    // @ts-ignore
    const __VLS_37 = __VLS_asFunctionalComponent(__VLS_36, new __VLS_36({
        size: "small",
    }));
    const __VLS_38 = __VLS_37({
        size: "small",
    }, ...__VLS_functionalComponentArgsRest(__VLS_37));
    __VLS_39.slots.default;
    const __VLS_40 = {}.AStatistic;
    /** @type {[typeof __VLS_components.AStatistic, typeof __VLS_components.aStatistic, ]} */ ;
    // @ts-ignore
    const __VLS_41 = __VLS_asFunctionalComponent(__VLS_40, new __VLS_40({
        title: "已使用",
        value: (__VLS_ctx.snapshot.slot_usage.used_slots),
        valueStyle: ({ color: '#1890ff' }),
    }));
    const __VLS_42 = __VLS_41({
        title: "已使用",
        value: (__VLS_ctx.snapshot.slot_usage.used_slots),
        valueStyle: ({ color: '#1890ff' }),
    }, ...__VLS_functionalComponentArgsRest(__VLS_41));
    var __VLS_39;
    var __VLS_35;
    const __VLS_44 = {}.ACol;
    /** @type {[typeof __VLS_components.ACol, typeof __VLS_components.aCol, typeof __VLS_components.ACol, typeof __VLS_components.aCol, ]} */ ;
    // @ts-ignore
    const __VLS_45 = __VLS_asFunctionalComponent(__VLS_44, new __VLS_44({
        span: (6),
    }));
    const __VLS_46 = __VLS_45({
        span: (6),
    }, ...__VLS_functionalComponentArgsRest(__VLS_45));
    __VLS_47.slots.default;
    const __VLS_48 = {}.ACard;
    /** @type {[typeof __VLS_components.ACard, typeof __VLS_components.aCard, typeof __VLS_components.ACard, typeof __VLS_components.aCard, ]} */ ;
    // @ts-ignore
    const __VLS_49 = __VLS_asFunctionalComponent(__VLS_48, new __VLS_48({
        size: "small",
    }));
    const __VLS_50 = __VLS_49({
        size: "small",
    }, ...__VLS_functionalComponentArgsRest(__VLS_49));
    __VLS_51.slots.default;
    const __VLS_52 = {}.AStatistic;
    /** @type {[typeof __VLS_components.AStatistic, typeof __VLS_components.aStatistic, ]} */ ;
    // @ts-ignore
    const __VLS_53 = __VLS_asFunctionalComponent(__VLS_52, new __VLS_52({
        title: "剩余可约",
        value: (__VLS_ctx.snapshot.slot_usage.available_slots),
        valueStyle: ({ color: '#52c41a' }),
    }));
    const __VLS_54 = __VLS_53({
        title: "剩余可约",
        value: (__VLS_ctx.snapshot.slot_usage.available_slots),
        valueStyle: ({ color: '#52c41a' }),
    }, ...__VLS_functionalComponentArgsRest(__VLS_53));
    var __VLS_51;
    var __VLS_47;
    const __VLS_56 = {}.ACol;
    /** @type {[typeof __VLS_components.ACol, typeof __VLS_components.aCol, typeof __VLS_components.ACol, typeof __VLS_components.aCol, ]} */ ;
    // @ts-ignore
    const __VLS_57 = __VLS_asFunctionalComponent(__VLS_56, new __VLS_56({
        span: (6),
    }));
    const __VLS_58 = __VLS_57({
        span: (6),
    }, ...__VLS_functionalComponentArgsRest(__VLS_57));
    __VLS_59.slots.default;
    const __VLS_60 = {}.ACard;
    /** @type {[typeof __VLS_components.ACard, typeof __VLS_components.aCard, typeof __VLS_components.ACard, typeof __VLS_components.aCard, ]} */ ;
    // @ts-ignore
    const __VLS_61 = __VLS_asFunctionalComponent(__VLS_60, new __VLS_60({
        size: "small",
    }));
    const __VLS_62 = __VLS_61({
        size: "small",
    }, ...__VLS_functionalComponentArgsRest(__VLS_61));
    __VLS_63.slots.default;
    const __VLS_64 = {}.AStatistic;
    /** @type {[typeof __VLS_components.AStatistic, typeof __VLS_components.aStatistic, ]} */ ;
    // @ts-ignore
    const __VLS_65 = __VLS_asFunctionalComponent(__VLS_64, new __VLS_64({
        title: "使用率",
        value: ((__VLS_ctx.snapshot.slot_usage.usage_rate * 100).toFixed(1)),
        suffix: "%",
        valueStyle: ({ color: __VLS_ctx.snapshot.slot_usage.usage_rate > 0.9 ? '#ff4d4f' : '#faad14' }),
    }));
    const __VLS_66 = __VLS_65({
        title: "使用率",
        value: ((__VLS_ctx.snapshot.slot_usage.usage_rate * 100).toFixed(1)),
        suffix: "%",
        valueStyle: ({ color: __VLS_ctx.snapshot.slot_usage.usage_rate > 0.9 ? '#ff4d4f' : '#faad14' }),
    }, ...__VLS_functionalComponentArgsRest(__VLS_65));
    var __VLS_63;
    var __VLS_59;
    var __VLS_19;
    if (__VLS_ctx.snapshot.alerts?.length) {
        for (const [alert] of __VLS_getVForSourceType((__VLS_ctx.snapshot.alerts))) {
            const __VLS_68 = {}.AAlert;
            /** @type {[typeof __VLS_components.AAlert, typeof __VLS_components.aAlert, ]} */ ;
            // @ts-ignore
            const __VLS_69 = __VLS_asFunctionalComponent(__VLS_68, new __VLS_68({
                key: (alert.message),
                message: (alert.message),
                type: (alert.type === 'critical' ? 'error' : alert.type === 'warning' ? 'warning' : 'info'),
                showIcon: true,
                ...{ style: {} },
            }));
            const __VLS_70 = __VLS_69({
                key: (alert.message),
                message: (alert.message),
                type: (alert.type === 'critical' ? 'error' : alert.type === 'warning' ? 'warning' : 'info'),
                showIcon: true,
                ...{ style: {} },
            }, ...__VLS_functionalComponentArgsRest(__VLS_69));
        }
    }
    const __VLS_72 = {}.ACard;
    /** @type {[typeof __VLS_components.ACard, typeof __VLS_components.aCard, typeof __VLS_components.ACard, typeof __VLS_components.aCard, ]} */ ;
    // @ts-ignore
    const __VLS_73 = __VLS_asFunctionalComponent(__VLS_72, new __VLS_72({
        title: "今日等待时长趋势",
        size: "small",
        ...{ style: {} },
    }));
    const __VLS_74 = __VLS_73({
        title: "今日等待时长趋势",
        size: "small",
        ...{ style: {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_73));
    __VLS_75.slots.default;
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div)({
        ref: "trendChartEl",
        ...{ style: {} },
    });
    /** @type {typeof __VLS_ctx.trendChartEl} */ ;
    var __VLS_75;
    const __VLS_76 = {}.ACard;
    /** @type {[typeof __VLS_components.ACard, typeof __VLS_components.aCard, typeof __VLS_components.ACard, typeof __VLS_components.aCard, ]} */ ;
    // @ts-ignore
    const __VLS_77 = __VLS_asFunctionalComponent(__VLS_76, new __VLS_76({
        title: "设备运行状态",
        size: "small",
    }));
    const __VLS_78 = __VLS_77({
        title: "设备运行状态",
        size: "small",
    }, ...__VLS_functionalComponentArgsRest(__VLS_77));
    __VLS_79.slots.default;
    const __VLS_80 = {}.ATable;
    /** @type {[typeof __VLS_components.ATable, typeof __VLS_components.aTable, typeof __VLS_components.ATable, typeof __VLS_components.aTable, ]} */ ;
    // @ts-ignore
    const __VLS_81 = __VLS_asFunctionalComponent(__VLS_80, new __VLS_80({
        dataSource: (__VLS_ctx.snapshot.device_status),
        columns: ([
            { title: '设备', dataIndex: 'device_name', key: 'device_name' },
            { title: '状态', dataIndex: 'status', key: 'status' },
            { title: '当前队列', dataIndex: 'queue_count', key: 'queue_count' },
        ]),
        pagination: (false),
        rowKey: "device_id",
        size: "small",
    }));
    const __VLS_82 = __VLS_81({
        dataSource: (__VLS_ctx.snapshot.device_status),
        columns: ([
            { title: '设备', dataIndex: 'device_name', key: 'device_name' },
            { title: '状态', dataIndex: 'status', key: 'status' },
            { title: '当前队列', dataIndex: 'queue_count', key: 'queue_count' },
        ]),
        pagination: (false),
        rowKey: "device_id",
        size: "small",
    }, ...__VLS_functionalComponentArgsRest(__VLS_81));
    __VLS_83.slots.default;
    {
        const { bodyCell: __VLS_thisSlot } = __VLS_83.slots;
        const [{ column, record }] = __VLS_getSlotParams(__VLS_thisSlot);
        if (column.key === 'status') {
            const __VLS_84 = {}.ABadge;
            /** @type {[typeof __VLS_components.ABadge, typeof __VLS_components.aBadge, ]} */ ;
            // @ts-ignore
            const __VLS_85 = __VLS_asFunctionalComponent(__VLS_84, new __VLS_84({
                status: (record.status === 'active' ? 'success' : record.status === 'fault' ? 'error' : 'warning'),
                text: (record.status),
            }));
            const __VLS_86 = __VLS_85({
                status: (record.status === 'active' ? 'success' : record.status === 'fault' ? 'error' : 'warning'),
                text: (record.status),
            }, ...__VLS_functionalComponentArgsRest(__VLS_85));
        }
    }
    var __VLS_83;
    var __VLS_79;
}
else if (!__VLS_ctx.loading) {
    const __VLS_88 = {}.AEmpty;
    /** @type {[typeof __VLS_components.AEmpty, typeof __VLS_components.aEmpty, ]} */ ;
    // @ts-ignore
    const __VLS_89 = __VLS_asFunctionalComponent(__VLS_88, new __VLS_88({
        description: "点击刷新获取数据",
    }));
    const __VLS_90 = __VLS_89({
        description: "点击刷新获取数据",
    }, ...__VLS_functionalComponentArgsRest(__VLS_89));
}
var __VLS_15;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            loading: loading,
            snapshot: snapshot,
            trendChartEl: trendChartEl,
            wsStatus: wsStatus,
            fetchData: fetchData,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
