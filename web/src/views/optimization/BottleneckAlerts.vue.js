/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { ref, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import { optimizationApi } from '@/api/optimization';
import { usePagination } from '@/composables/usePagination';
const { loading, items, pagination, fetchData, onTableChange } = usePagination(params => optimizationApi.listAlerts(params));
onMounted(() => fetchData());
const dismissModal = ref(false);
const dismissTarget = ref(null);
const dismissReason = ref('');
const dismissing = ref(false);
function openDismiss(record) {
    dismissTarget.value = record;
    dismissReason.value = '';
    dismissModal.value = true;
}
async function handleDismiss() {
    if (!dismissTarget.value)
        return;
    dismissing.value = true;
    try {
        await optimizationApi.dismissAlert(dismissTarget.value.id, dismissReason.value);
        message.success('告警已确认');
        dismissModal.value = false;
        fetchData();
    }
    finally {
        dismissing.value = false;
    }
}
const SEVERITY_COLOR = { critical: 'red', high: 'orange', medium: 'gold', low: 'blue' };
const columns = [
    { title: '告警类型', dataIndex: 'alert_type', key: 'type' },
    { title: '严重程度', dataIndex: 'severity', key: 'severity' },
    { title: '摘要', dataIndex: 'summary', key: 'summary' },
    { title: '设备', dataIndex: 'device_name', key: 'device' },
    { title: '当前值', dataIndex: 'current_value', key: 'value' },
    { title: '阈值', dataIndex: 'threshold_value', key: 'threshold' },
    { title: '状态', dataIndex: 'status', key: 'status' },
    { title: '操作', key: 'actions' },
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
}));
const __VLS_2 = __VLS_1({
    ...{ 'onChange': {} },
    columns: (__VLS_ctx.columns),
    dataSource: (__VLS_ctx.items),
    loading: (__VLS_ctx.loading),
    pagination: (__VLS_ctx.pagination),
    rowKey: "id",
    size: "middle",
}, ...__VLS_functionalComponentArgsRest(__VLS_1));
let __VLS_4;
let __VLS_5;
let __VLS_6;
const __VLS_7 = {
    onChange: (__VLS_ctx.onTableChange)
};
__VLS_3.slots.default;
{
    const { bodyCell: __VLS_thisSlot } = __VLS_3.slots;
    const [{ column, record }] = __VLS_getSlotParams(__VLS_thisSlot);
    if (column.key === 'severity') {
        const __VLS_8 = {}.ATag;
        /** @type {[typeof __VLS_components.ATag, typeof __VLS_components.aTag, typeof __VLS_components.ATag, typeof __VLS_components.aTag, ]} */ ;
        // @ts-ignore
        const __VLS_9 = __VLS_asFunctionalComponent(__VLS_8, new __VLS_8({
            color: (__VLS_ctx.SEVERITY_COLOR[record.severity]),
        }));
        const __VLS_10 = __VLS_9({
            color: (__VLS_ctx.SEVERITY_COLOR[record.severity]),
        }, ...__VLS_functionalComponentArgsRest(__VLS_9));
        __VLS_11.slots.default;
        (record.severity);
        var __VLS_11;
    }
    if (column.key === 'status') {
        const __VLS_12 = {}.ABadge;
        /** @type {[typeof __VLS_components.ABadge, typeof __VLS_components.aBadge, ]} */ ;
        // @ts-ignore
        const __VLS_13 = __VLS_asFunctionalComponent(__VLS_12, new __VLS_12({
            status: (record.status === 'open' ? 'error' : 'default'),
            text: (record.status === 'open' ? '待处理' : '已处理'),
        }));
        const __VLS_14 = __VLS_13({
            status: (record.status === 'open' ? 'error' : 'default'),
            text: (record.status === 'open' ? '待处理' : '已处理'),
        }, ...__VLS_functionalComponentArgsRest(__VLS_13));
    }
    if (column.key === 'actions') {
        if (record.status === 'open') {
            const __VLS_16 = {}.AButton;
            /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
            // @ts-ignore
            const __VLS_17 = __VLS_asFunctionalComponent(__VLS_16, new __VLS_16({
                ...{ 'onClick': {} },
                type: "link",
                size: "small",
            }));
            const __VLS_18 = __VLS_17({
                ...{ 'onClick': {} },
                type: "link",
                size: "small",
            }, ...__VLS_functionalComponentArgsRest(__VLS_17));
            let __VLS_20;
            let __VLS_21;
            let __VLS_22;
            const __VLS_23 = {
                onClick: (...[$event]) => {
                    if (!(column.key === 'actions'))
                        return;
                    if (!(record.status === 'open'))
                        return;
                    __VLS_ctx.openDismiss(record);
                }
            };
            __VLS_19.slots.default;
            var __VLS_19;
        }
    }
}
var __VLS_3;
const __VLS_24 = {}.AModal;
/** @type {[typeof __VLS_components.AModal, typeof __VLS_components.aModal, typeof __VLS_components.AModal, typeof __VLS_components.aModal, ]} */ ;
// @ts-ignore
const __VLS_25 = __VLS_asFunctionalComponent(__VLS_24, new __VLS_24({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.dismissModal),
    title: "确认告警",
    confirmLoading: (__VLS_ctx.dismissing),
}));
const __VLS_26 = __VLS_25({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.dismissModal),
    title: "确认告警",
    confirmLoading: (__VLS_ctx.dismissing),
}, ...__VLS_functionalComponentArgsRest(__VLS_25));
let __VLS_28;
let __VLS_29;
let __VLS_30;
const __VLS_31 = {
    onOk: (__VLS_ctx.handleDismiss)
};
__VLS_27.slots.default;
const __VLS_32 = {}.AForm;
/** @type {[typeof __VLS_components.AForm, typeof __VLS_components.aForm, typeof __VLS_components.AForm, typeof __VLS_components.aForm, ]} */ ;
// @ts-ignore
const __VLS_33 = __VLS_asFunctionalComponent(__VLS_32, new __VLS_32({
    layout: "vertical",
}));
const __VLS_34 = __VLS_33({
    layout: "vertical",
}, ...__VLS_functionalComponentArgsRest(__VLS_33));
__VLS_35.slots.default;
const __VLS_36 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_37 = __VLS_asFunctionalComponent(__VLS_36, new __VLS_36({
    label: "处理说明",
}));
const __VLS_38 = __VLS_37({
    label: "处理说明",
}, ...__VLS_functionalComponentArgsRest(__VLS_37));
__VLS_39.slots.default;
const __VLS_40 = {}.ATextarea;
/** @type {[typeof __VLS_components.ATextarea, typeof __VLS_components.aTextarea, ]} */ ;
// @ts-ignore
const __VLS_41 = __VLS_asFunctionalComponent(__VLS_40, new __VLS_40({
    value: (__VLS_ctx.dismissReason),
    rows: (4),
    placeholder: "请说明处理措施...",
}));
const __VLS_42 = __VLS_41({
    value: (__VLS_ctx.dismissReason),
    rows: (4),
    placeholder: "请说明处理措施...",
}, ...__VLS_functionalComponentArgsRest(__VLS_41));
var __VLS_39;
var __VLS_35;
var __VLS_27;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            loading: loading,
            items: items,
            pagination: pagination,
            onTableChange: onTableChange,
            dismissModal: dismissModal,
            dismissReason: dismissReason,
            dismissing: dismissing,
            openDismiss: openDismiss,
            handleDismiss: handleDismiss,
            SEVERITY_COLOR: SEVERITY_COLOR,
            columns: columns,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
