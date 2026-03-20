/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { ref, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import { appointmentApi } from '@/api/appointment';
import { usePagination } from '@/composables/usePagination';
const { loading, items, pagination, fetchData, onTableChange } = usePagination(params => appointmentApi.listBlacklist(params));
onMounted(() => fetchData());
const appealModal = ref(false);
const appealTarget = ref(null);
const appealReason = ref('');
const appealing = ref(false);
function openAppeal(record) {
    appealTarget.value = record;
    appealReason.value = '';
    appealModal.value = true;
}
async function submitAppeal() {
    if (!appealTarget.value || !appealReason.value.trim())
        return;
    appealing.value = true;
    try {
        await appointmentApi.submitAppeal(appealTarget.value.id, appealReason.value);
        message.success('申诉已提交');
        appealModal.value = false;
        fetchData();
    }
    finally {
        appealing.value = false;
    }
}
async function handleRemove(record) {
    await appointmentApi.removeFromBlacklist(record.id);
    message.success('已移出黑名单');
    fetchData();
}
const columns = [
    { title: '患者', dataIndex: 'patient_name', key: 'patient_name' },
    { title: '加入原因', dataIndex: 'reason', key: 'reason' },
    { title: '爽约次数', dataIndex: 'no_show_count', key: 'no_show_count' },
    { title: '加入时间', dataIndex: 'added_at', key: 'added_at' },
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
    if (column.key === 'status') {
        const __VLS_8 = {}.ATag;
        /** @type {[typeof __VLS_components.ATag, typeof __VLS_components.aTag, typeof __VLS_components.ATag, typeof __VLS_components.aTag, ]} */ ;
        // @ts-ignore
        const __VLS_9 = __VLS_asFunctionalComponent(__VLS_8, new __VLS_8({
            color: (record.status === 'active' ? 'red' : 'default'),
        }));
        const __VLS_10 = __VLS_9({
            color: (record.status === 'active' ? 'red' : 'default'),
        }, ...__VLS_functionalComponentArgsRest(__VLS_9));
        __VLS_11.slots.default;
        (record.status === 'active' ? '黑名单中' : '已解除');
        var __VLS_11;
    }
    if (column.key === 'actions') {
        const __VLS_12 = {}.ASpace;
        /** @type {[typeof __VLS_components.ASpace, typeof __VLS_components.aSpace, typeof __VLS_components.ASpace, typeof __VLS_components.aSpace, ]} */ ;
        // @ts-ignore
        const __VLS_13 = __VLS_asFunctionalComponent(__VLS_12, new __VLS_12({}));
        const __VLS_14 = __VLS_13({}, ...__VLS_functionalComponentArgsRest(__VLS_13));
        __VLS_15.slots.default;
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
                __VLS_ctx.openAppeal(record);
            }
        };
        __VLS_19.slots.default;
        var __VLS_19;
        const __VLS_24 = {}.AButton;
        /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
        // @ts-ignore
        const __VLS_25 = __VLS_asFunctionalComponent(__VLS_24, new __VLS_24({
            ...{ 'onClick': {} },
            type: "link",
            danger: true,
            size: "small",
        }));
        const __VLS_26 = __VLS_25({
            ...{ 'onClick': {} },
            type: "link",
            danger: true,
            size: "small",
        }, ...__VLS_functionalComponentArgsRest(__VLS_25));
        let __VLS_28;
        let __VLS_29;
        let __VLS_30;
        const __VLS_31 = {
            onClick: (...[$event]) => {
                if (!(column.key === 'actions'))
                    return;
                __VLS_ctx.handleRemove(record);
            }
        };
        __VLS_27.slots.default;
        var __VLS_27;
        var __VLS_15;
    }
}
var __VLS_3;
const __VLS_32 = {}.AModal;
/** @type {[typeof __VLS_components.AModal, typeof __VLS_components.aModal, typeof __VLS_components.AModal, typeof __VLS_components.aModal, ]} */ ;
// @ts-ignore
const __VLS_33 = __VLS_asFunctionalComponent(__VLS_32, new __VLS_32({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.appealModal),
    title: "提交申诉",
    confirmLoading: (__VLS_ctx.appealing),
}));
const __VLS_34 = __VLS_33({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.appealModal),
    title: "提交申诉",
    confirmLoading: (__VLS_ctx.appealing),
}, ...__VLS_functionalComponentArgsRest(__VLS_33));
let __VLS_36;
let __VLS_37;
let __VLS_38;
const __VLS_39 = {
    onOk: (__VLS_ctx.submitAppeal)
};
__VLS_35.slots.default;
const __VLS_40 = {}.AForm;
/** @type {[typeof __VLS_components.AForm, typeof __VLS_components.aForm, typeof __VLS_components.AForm, typeof __VLS_components.aForm, ]} */ ;
// @ts-ignore
const __VLS_41 = __VLS_asFunctionalComponent(__VLS_40, new __VLS_40({
    layout: "vertical",
}));
const __VLS_42 = __VLS_41({
    layout: "vertical",
}, ...__VLS_functionalComponentArgsRest(__VLS_41));
__VLS_43.slots.default;
const __VLS_44 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_45 = __VLS_asFunctionalComponent(__VLS_44, new __VLS_44({
    label: "申诉理由",
}));
const __VLS_46 = __VLS_45({
    label: "申诉理由",
}, ...__VLS_functionalComponentArgsRest(__VLS_45));
__VLS_47.slots.default;
const __VLS_48 = {}.ATextarea;
/** @type {[typeof __VLS_components.ATextarea, typeof __VLS_components.aTextarea, ]} */ ;
// @ts-ignore
const __VLS_49 = __VLS_asFunctionalComponent(__VLS_48, new __VLS_48({
    value: (__VLS_ctx.appealReason),
    rows: (4),
    placeholder: "请详细说明申诉理由...",
}));
const __VLS_50 = __VLS_49({
    value: (__VLS_ctx.appealReason),
    rows: (4),
    placeholder: "请详细说明申诉理由...",
}, ...__VLS_functionalComponentArgsRest(__VLS_49));
var __VLS_47;
var __VLS_43;
var __VLS_35;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            loading: loading,
            items: items,
            pagination: pagination,
            onTableChange: onTableChange,
            appealModal: appealModal,
            appealReason: appealReason,
            appealing: appealing,
            openAppeal: openAppeal,
            submitAppeal: submitAppeal,
            handleRemove: handleRemove,
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
