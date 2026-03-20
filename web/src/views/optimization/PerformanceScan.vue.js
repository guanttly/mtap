/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { ref, onMounted } from 'vue';
import { optimizationApi } from '@/api/optimization';
import { usePagination } from '@/composables/usePagination';
const { loading, items, pagination, fetchData, onTableChange } = usePagination(params => optimizationApi.listScans(params));
onMounted(() => fetchData());
const columns = [
    { title: '扫描周', dataIndex: 'scan_week', key: 'scan_week' },
    { title: '扫描时间', dataIndex: 'scanned_at', key: 'scanned_at' },
    { title: '发现机会数', key: 'count', customRender: ({ record }) => record.opportunities?.length ?? 0 },
    { title: '操作', key: 'actions' },
];
const expandedId = ref(null);
function handleExpand(_, record) {
    expandedId.value = expandedId.value === record.id ? null : record.id;
}
function handleToggle(record) {
    expandedId.value = expandedId.value === record.id ? null : record.id;
}
const detailColumns = [
    { title: '指标编码', dataIndex: 'metric_code', key: 'code' },
    { title: '指标名', dataIndex: 'metric_name', key: 'name' },
    { title: '当前值', dataIndex: 'current_value', key: 'cur', customRender: ({ text }) => text?.toFixed?.(2) ?? text },
    { title: '正常值', dataIndex: 'normal_value', key: 'norm', customRender: ({ text }) => text?.toFixed?.(2) ?? text },
    { title: '偏差%', dataIndex: 'deviation_pct', key: 'dev', customRender: ({ text }) => `${(text * 100).toFixed(1)}%` },
    { title: '建议类别', dataIndex: 'suggested_category', key: 'cat' },
];
debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
const __VLS_ctx = {};
let __VLS_components;
let __VLS_directives;
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({});
const __VLS_0 = {}.ATable;
/** @type {[typeof __VLS_components.ATable, typeof __VLS_components.aTable, typeof __VLS_components.ATable, typeof __VLS_components.aTable, ]} */ ;
// @ts-ignore
const __VLS_1 = __VLS_asFunctionalComponent(__VLS_0, new __VLS_0({
    ...{ 'onChange': {} },
    columns: (__VLS_ctx.columns),
    dataSource: (__VLS_ctx.items),
    loading: (__VLS_ctx.loading),
    pagination: (__VLS_ctx.pagination),
    rowKey: "id",
    size: "middle",
    expandable: ({ expandedRowKeys: __VLS_ctx.expandedId ? [__VLS_ctx.expandedId] : [], onExpand: __VLS_ctx.handleExpand }),
}));
const __VLS_2 = __VLS_1({
    ...{ 'onChange': {} },
    columns: (__VLS_ctx.columns),
    dataSource: (__VLS_ctx.items),
    loading: (__VLS_ctx.loading),
    pagination: (__VLS_ctx.pagination),
    rowKey: "id",
    size: "middle",
    expandable: ({ expandedRowKeys: __VLS_ctx.expandedId ? [__VLS_ctx.expandedId] : [], onExpand: __VLS_ctx.handleExpand }),
}, ...__VLS_functionalComponentArgsRest(__VLS_1));
let __VLS_4;
let __VLS_5;
let __VLS_6;
const __VLS_7 = {
    onChange: (__VLS_ctx.onTableChange)
};
__VLS_3.slots.default;
{
    const { expandedRowRender: __VLS_thisSlot } = __VLS_3.slots;
    const [{ record }] = __VLS_getSlotParams(__VLS_thisSlot);
    const __VLS_8 = {}.ATable;
    /** @type {[typeof __VLS_components.ATable, typeof __VLS_components.aTable, ]} */ ;
    // @ts-ignore
    const __VLS_9 = __VLS_asFunctionalComponent(__VLS_8, new __VLS_8({
        dataSource: (record.opportunities),
        columns: (__VLS_ctx.detailColumns),
        pagination: (false),
        size: "small",
        rowKey: "metric_code",
    }));
    const __VLS_10 = __VLS_9({
        dataSource: (record.opportunities),
        columns: (__VLS_ctx.detailColumns),
        pagination: (false),
        size: "small",
        rowKey: "metric_code",
    }, ...__VLS_functionalComponentArgsRest(__VLS_9));
}
{
    const { bodyCell: __VLS_thisSlot } = __VLS_3.slots;
    const [{ column, record }] = __VLS_getSlotParams(__VLS_thisSlot);
    if (column.key === 'actions') {
        const __VLS_12 = {}.AButton;
        /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
        // @ts-ignore
        const __VLS_13 = __VLS_asFunctionalComponent(__VLS_12, new __VLS_12({
            ...{ 'onClick': {} },
            type: "link",
            size: "small",
        }));
        const __VLS_14 = __VLS_13({
            ...{ 'onClick': {} },
            type: "link",
            size: "small",
        }, ...__VLS_functionalComponentArgsRest(__VLS_13));
        let __VLS_16;
        let __VLS_17;
        let __VLS_18;
        const __VLS_19 = {
            onClick: (...[$event]) => {
                if (!(column.key === 'actions'))
                    return;
                __VLS_ctx.handleToggle(record);
            }
        };
        __VLS_15.slots.default;
        var __VLS_15;
    }
}
var __VLS_3;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            loading: loading,
            items: items,
            pagination: pagination,
            onTableChange: onTableChange,
            columns: columns,
            expandedId: expandedId,
            handleExpand: handleExpand,
            handleToggle: handleToggle,
            detailColumns: detailColumns,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
