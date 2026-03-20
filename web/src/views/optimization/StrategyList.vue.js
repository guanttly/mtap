/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { ref, onMounted } from 'vue';
import { message, Modal } from 'ant-design-vue';
import { useRouter } from 'vue-router';
import { optimizationApi } from '@/api/optimization';
import { usePagination } from '@/composables/usePagination';
const router = useRouter();
const { loading, items, pagination, fetchData, onTableChange } = usePagination(params => optimizationApi.listStrategies(params));
onMounted(() => fetchData());
const approveModal = ref(false);
const approveTarget = ref(null);
const approveForm = ref({ trial_days: 7, gray_scope: null });
const approving = ref(false);
const STATUS_MAP = {
    pending_review: { color: 'default', label: '待审批' },
    trial_running: { color: 'processing', label: '试运行中' },
    trial_running_b: { color: 'processing', label: '试运行B' },
    submitted_approval: { color: 'gold', label: '已提交审批' },
    pending_eval: { color: 'blue', label: '待评估' },
    promoted: { color: 'success', label: '已推全量' },
    normalized: { color: 'success', label: '已常态化' },
    rolled_back: { color: 'warning', label: '已回滚' },
    rejected: { color: 'error', label: '已拒绝' },
    archived: { color: 'default', label: '已归档' },
};
const columns = [
    { title: '策略标题', dataIndex: 'title', key: 'title' },
    { title: '类别', dataIndex: 'category', key: 'category' },
    { title: '预期收益', dataIndex: 'expected_benefit', key: 'benefit' },
    { title: '状态', dataIndex: 'status', key: 'status' },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at' },
    { title: '操作', key: 'actions' },
];
function openApprove(record) {
    approveTarget.value = record;
    approveForm.value = { trial_days: 7, gray_scope: null };
    approveModal.value = true;
}
async function handleApprove() {
    if (!approveTarget.value)
        return;
    approving.value = true;
    try {
        await optimizationApi.approveStrategy(approveTarget.value.id, approveForm.value);
        message.success('已批准并进入试运行');
        approveModal.value = false;
        fetchData();
    }
    finally {
        approving.value = false;
    }
}
function handleReject(record) {
    Modal.confirm({
        title: '拒绝策略',
        content: '确定要拒绝该优化策略吗？',
        okType: 'danger',
        onOk: async () => {
            await optimizationApi.rejectStrategy(record.id, '人工审核拒绝');
            message.success('已拒绝');
            fetchData();
        },
    });
}
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
            color: (__VLS_ctx.STATUS_MAP[record.status]?.color),
        }));
        const __VLS_10 = __VLS_9({
            color: (__VLS_ctx.STATUS_MAP[record.status]?.color),
        }, ...__VLS_functionalComponentArgsRest(__VLS_9));
        __VLS_11.slots.default;
        (__VLS_ctx.STATUS_MAP[record.status]?.label ?? record.status);
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
                __VLS_ctx.router.push(`/optimization/strategies/${record.id}`);
            }
        };
        __VLS_19.slots.default;
        var __VLS_19;
        if (record.status === 'pending_review') {
            const __VLS_24 = {}.AButton;
            /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
            // @ts-ignore
            const __VLS_25 = __VLS_asFunctionalComponent(__VLS_24, new __VLS_24({
                ...{ 'onClick': {} },
                type: "link",
                size: "small",
            }));
            const __VLS_26 = __VLS_25({
                ...{ 'onClick': {} },
                type: "link",
                size: "small",
            }, ...__VLS_functionalComponentArgsRest(__VLS_25));
            let __VLS_28;
            let __VLS_29;
            let __VLS_30;
            const __VLS_31 = {
                onClick: (...[$event]) => {
                    if (!(column.key === 'actions'))
                        return;
                    if (!(record.status === 'pending_review'))
                        return;
                    __VLS_ctx.openApprove(record);
                }
            };
            __VLS_27.slots.default;
            var __VLS_27;
        }
        if (record.status === 'pending_review') {
            const __VLS_32 = {}.AButton;
            /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
            // @ts-ignore
            const __VLS_33 = __VLS_asFunctionalComponent(__VLS_32, new __VLS_32({
                ...{ 'onClick': {} },
                type: "link",
                danger: true,
                size: "small",
            }));
            const __VLS_34 = __VLS_33({
                ...{ 'onClick': {} },
                type: "link",
                danger: true,
                size: "small",
            }, ...__VLS_functionalComponentArgsRest(__VLS_33));
            let __VLS_36;
            let __VLS_37;
            let __VLS_38;
            const __VLS_39 = {
                onClick: (...[$event]) => {
                    if (!(column.key === 'actions'))
                        return;
                    if (!(record.status === 'pending_review'))
                        return;
                    __VLS_ctx.handleReject(record);
                }
            };
            __VLS_35.slots.default;
            var __VLS_35;
        }
        if (record.status === 'trial_running') {
            const __VLS_40 = {}.AButton;
            /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
            // @ts-ignore
            const __VLS_41 = __VLS_asFunctionalComponent(__VLS_40, new __VLS_40({
                ...{ 'onClick': {} },
                type: "link",
                size: "small",
            }));
            const __VLS_42 = __VLS_41({
                ...{ 'onClick': {} },
                type: "link",
                size: "small",
            }, ...__VLS_functionalComponentArgsRest(__VLS_41));
            let __VLS_44;
            let __VLS_45;
            let __VLS_46;
            const __VLS_47 = {
                onClick: (...[$event]) => {
                    if (!(column.key === 'actions'))
                        return;
                    if (!(record.status === 'trial_running'))
                        return;
                    __VLS_ctx.router.push(`/optimization/trials/${record.trial_run?.id}`);
                }
            };
            __VLS_43.slots.default;
            var __VLS_43;
        }
        var __VLS_15;
    }
}
var __VLS_3;
const __VLS_48 = {}.AModal;
/** @type {[typeof __VLS_components.AModal, typeof __VLS_components.aModal, typeof __VLS_components.AModal, typeof __VLS_components.aModal, ]} */ ;
// @ts-ignore
const __VLS_49 = __VLS_asFunctionalComponent(__VLS_48, new __VLS_48({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.approveModal),
    title: "批准策略",
    confirmLoading: (__VLS_ctx.approving),
}));
const __VLS_50 = __VLS_49({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.approveModal),
    title: "批准策略",
    confirmLoading: (__VLS_ctx.approving),
}, ...__VLS_functionalComponentArgsRest(__VLS_49));
let __VLS_52;
let __VLS_53;
let __VLS_54;
const __VLS_55 = {
    onOk: (__VLS_ctx.handleApprove)
};
__VLS_51.slots.default;
const __VLS_56 = {}.AForm;
/** @type {[typeof __VLS_components.AForm, typeof __VLS_components.aForm, typeof __VLS_components.AForm, typeof __VLS_components.aForm, ]} */ ;
// @ts-ignore
const __VLS_57 = __VLS_asFunctionalComponent(__VLS_56, new __VLS_56({
    model: (__VLS_ctx.approveForm),
    layout: "vertical",
}));
const __VLS_58 = __VLS_57({
    model: (__VLS_ctx.approveForm),
    layout: "vertical",
}, ...__VLS_functionalComponentArgsRest(__VLS_57));
__VLS_59.slots.default;
const __VLS_60 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_61 = __VLS_asFunctionalComponent(__VLS_60, new __VLS_60({
    label: "试运行天数",
}));
const __VLS_62 = __VLS_61({
    label: "试运行天数",
}, ...__VLS_functionalComponentArgsRest(__VLS_61));
__VLS_63.slots.default;
const __VLS_64 = {}.AInputNumber;
/** @type {[typeof __VLS_components.AInputNumber, typeof __VLS_components.aInputNumber, ]} */ ;
// @ts-ignore
const __VLS_65 = __VLS_asFunctionalComponent(__VLS_64, new __VLS_64({
    value: (__VLS_ctx.approveForm.trial_days),
    min: (1),
    max: (30),
    ...{ style: {} },
}));
const __VLS_66 = __VLS_65({
    value: (__VLS_ctx.approveForm.trial_days),
    min: (1),
    max: (30),
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_65));
var __VLS_63;
var __VLS_59;
var __VLS_51;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            router: router,
            loading: loading,
            items: items,
            pagination: pagination,
            onTableChange: onTableChange,
            approveModal: approveModal,
            approveForm: approveForm,
            approving: approving,
            STATUS_MAP: STATUS_MAP,
            columns: columns,
            openApprove: openApprove,
            handleApprove: handleApprove,
            handleReject: handleReject,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
