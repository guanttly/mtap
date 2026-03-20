/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { ref, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import { appointmentApi } from '@/api/appointment';
import { usePagination } from '@/composables/usePagination';
const { loading, items, pagination, fetchData, onTableChange, search } = usePagination(params => appointmentApi.listAppointments(params));
onMounted(() => fetchData());
const STATUS_COLOR = {
    pending: 'default', confirmed: 'processing', paid: 'blue', checked_in: 'cyan',
    completed: 'success', cancelled: 'error', no_show: 'warning',
};
const STATUS_LABEL = {
    pending: '待确认', confirmed: '已确认', paid: '已缴费', checked_in: '已签到',
    completed: '已完成', cancelled: '已取消', no_show: '爽约',
};
const columns = [
    { title: '预约号', dataIndex: 'appointment_no', key: 'no', width: 160 },
    { title: '患者', dataIndex: 'patient_name', key: 'patient_name' },
    { title: '状态', dataIndex: 'status', key: 'status' },
    { title: '来源', dataIndex: 'source', key: 'source' },
    { title: '操作', key: 'actions' },
];
const searchKeyword = ref('');
function handleSearch() {
    search({ keyword: searchKeyword.value });
}
async function handleCancel(record) {
    await appointmentApi.cancel(record.id, '人工取消');
    message.success('已取消');
    fetchData();
}
debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
const __VLS_ctx = {};
let __VLS_components;
let __VLS_directives;
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "search-bar" },
    ...{ style: {} },
});
const __VLS_0 = {}.AInputSearch;
/** @type {[typeof __VLS_components.AInputSearch, typeof __VLS_components.aInputSearch, ]} */ ;
// @ts-ignore
const __VLS_1 = __VLS_asFunctionalComponent(__VLS_0, new __VLS_0({
    ...{ 'onSearch': {} },
    value: (__VLS_ctx.searchKeyword),
    placeholder: "患者姓名/就诊号",
    ...{ style: {} },
}));
const __VLS_2 = __VLS_1({
    ...{ 'onSearch': {} },
    value: (__VLS_ctx.searchKeyword),
    placeholder: "患者姓名/就诊号",
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_1));
let __VLS_4;
let __VLS_5;
let __VLS_6;
const __VLS_7 = {
    onSearch: (__VLS_ctx.handleSearch)
};
var __VLS_3;
const __VLS_8 = {}.ASelect;
/** @type {[typeof __VLS_components.ASelect, typeof __VLS_components.aSelect, typeof __VLS_components.ASelect, typeof __VLS_components.aSelect, ]} */ ;
// @ts-ignore
const __VLS_9 = __VLS_asFunctionalComponent(__VLS_8, new __VLS_8({
    ...{ 'onChange': {} },
    ...{ style: {} },
    placeholder: "状态筛选",
    allowClear: true,
}));
const __VLS_10 = __VLS_9({
    ...{ 'onChange': {} },
    ...{ style: {} },
    placeholder: "状态筛选",
    allowClear: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_9));
let __VLS_12;
let __VLS_13;
let __VLS_14;
const __VLS_15 = {
    onChange: ((v) => __VLS_ctx.search({ status: v }))
};
__VLS_11.slots.default;
for (const [label, val] of __VLS_getVForSourceType((__VLS_ctx.STATUS_LABEL))) {
    const __VLS_16 = {}.ASelectOption;
    /** @type {[typeof __VLS_components.ASelectOption, typeof __VLS_components.aSelectOption, typeof __VLS_components.ASelectOption, typeof __VLS_components.aSelectOption, ]} */ ;
    // @ts-ignore
    const __VLS_17 = __VLS_asFunctionalComponent(__VLS_16, new __VLS_16({
        key: (val),
        value: (val),
    }));
    const __VLS_18 = __VLS_17({
        key: (val),
        value: (val),
    }, ...__VLS_functionalComponentArgsRest(__VLS_17));
    __VLS_19.slots.default;
    (label);
    var __VLS_19;
}
var __VLS_11;
const __VLS_20 = {}.ATable;
/** @type {[typeof __VLS_components.ATable, typeof __VLS_components.aTable, typeof __VLS_components.ATable, typeof __VLS_components.aTable, ]} */ ;
// @ts-ignore
const __VLS_21 = __VLS_asFunctionalComponent(__VLS_20, new __VLS_20({
    ...{ 'onChange': {} },
    columns: (__VLS_ctx.columns),
    dataSource: (__VLS_ctx.items),
    loading: (__VLS_ctx.loading),
    pagination: (__VLS_ctx.pagination),
    rowKey: "id",
    size: "middle",
}));
const __VLS_22 = __VLS_21({
    ...{ 'onChange': {} },
    columns: (__VLS_ctx.columns),
    dataSource: (__VLS_ctx.items),
    loading: (__VLS_ctx.loading),
    pagination: (__VLS_ctx.pagination),
    rowKey: "id",
    size: "middle",
}, ...__VLS_functionalComponentArgsRest(__VLS_21));
let __VLS_24;
let __VLS_25;
let __VLS_26;
const __VLS_27 = {
    onChange: (__VLS_ctx.onTableChange)
};
__VLS_23.slots.default;
{
    const { bodyCell: __VLS_thisSlot } = __VLS_23.slots;
    const [{ column, record }] = __VLS_getSlotParams(__VLS_thisSlot);
    if (column.key === 'status') {
        const __VLS_28 = {}.ATag;
        /** @type {[typeof __VLS_components.ATag, typeof __VLS_components.aTag, typeof __VLS_components.ATag, typeof __VLS_components.aTag, ]} */ ;
        // @ts-ignore
        const __VLS_29 = __VLS_asFunctionalComponent(__VLS_28, new __VLS_28({
            color: (__VLS_ctx.STATUS_COLOR[record.status]),
        }));
        const __VLS_30 = __VLS_29({
            color: (__VLS_ctx.STATUS_COLOR[record.status]),
        }, ...__VLS_functionalComponentArgsRest(__VLS_29));
        __VLS_31.slots.default;
        (__VLS_ctx.STATUS_LABEL[record.status]);
        var __VLS_31;
    }
    if (column.key === 'actions') {
        const __VLS_32 = {}.ASpace;
        /** @type {[typeof __VLS_components.ASpace, typeof __VLS_components.aSpace, typeof __VLS_components.ASpace, typeof __VLS_components.aSpace, ]} */ ;
        // @ts-ignore
        const __VLS_33 = __VLS_asFunctionalComponent(__VLS_32, new __VLS_32({}));
        const __VLS_34 = __VLS_33({}, ...__VLS_functionalComponentArgsRest(__VLS_33));
        __VLS_35.slots.default;
        const __VLS_36 = {}.AButton;
        /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
        // @ts-ignore
        const __VLS_37 = __VLS_asFunctionalComponent(__VLS_36, new __VLS_36({
            type: "link",
            size: "small",
        }));
        const __VLS_38 = __VLS_37({
            type: "link",
            size: "small",
        }, ...__VLS_functionalComponentArgsRest(__VLS_37));
        __VLS_39.slots.default;
        var __VLS_39;
        if (['pending', 'confirmed', 'paid'].includes(record.status)) {
            const __VLS_40 = {}.AButton;
            /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
            // @ts-ignore
            const __VLS_41 = __VLS_asFunctionalComponent(__VLS_40, new __VLS_40({
                ...{ 'onClick': {} },
                type: "link",
                danger: true,
                size: "small",
            }));
            const __VLS_42 = __VLS_41({
                ...{ 'onClick': {} },
                type: "link",
                danger: true,
                size: "small",
            }, ...__VLS_functionalComponentArgsRest(__VLS_41));
            let __VLS_44;
            let __VLS_45;
            let __VLS_46;
            const __VLS_47 = {
                onClick: (...[$event]) => {
                    if (!(column.key === 'actions'))
                        return;
                    if (!(['pending', 'confirmed', 'paid'].includes(record.status)))
                        return;
                    __VLS_ctx.handleCancel(record);
                }
            };
            __VLS_43.slots.default;
            var __VLS_43;
        }
        var __VLS_35;
    }
}
var __VLS_23;
/** @type {__VLS_StyleScopedClasses['search-bar']} */ ;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            loading: loading,
            items: items,
            pagination: pagination,
            onTableChange: onTableChange,
            search: search,
            STATUS_COLOR: STATUS_COLOR,
            STATUS_LABEL: STATUS_LABEL,
            columns: columns,
            searchKeyword: searchKeyword,
            handleSearch: handleSearch,
            handleCancel: handleCancel,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
