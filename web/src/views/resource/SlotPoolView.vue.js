/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { ref, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import { resourceApi } from '@/api/resource';
const pools = ref([]);
const slots = ref([]);
const loading = ref(false);
const slotsLoading = ref(false);
const selectedPool = ref(null);
const STATUS_COLOR = { available: 'success', booked: 'processing', locked: 'warning', released: 'default' };
async function fetchPools() {
    loading.value = true;
    try {
        const res = await resourceApi.listSlotPools();
        pools.value = res.items;
    }
    finally {
        loading.value = false;
    }
}
async function fetchSlots(pool) {
    selectedPool.value = pool;
    slotsLoading.value = true;
    try {
        const res = await resourceApi.listTimeSlots({ slot_pool_id: pool.id });
        slots.value = res.items;
    }
    finally {
        slotsLoading.value = false;
    }
}
async function releaseSlot(slot) {
    await resourceApi.releaseSlot(slot.id);
    message.success('已强制释放');
    if (selectedPool.value)
        fetchSlots(selectedPool.value);
}
onMounted(fetchPools);
const slotColumns = [
    { title: '日期', dataIndex: 'date', key: 'date' },
    { title: '开始', dataIndex: 'start_time', key: 'start_time' },
    { title: '结束', dataIndex: 'end_time', key: 'end_time' },
    { title: '总数', dataIndex: 'total_count', key: 'total_count' },
    { title: '剩余', dataIndex: 'remaining_count', key: 'remaining_count' },
    { title: '状态', dataIndex: 'status', key: 'status' },
    { title: '操作', key: 'actions' },
];
debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
const __VLS_ctx = {};
let __VLS_components;
let __VLS_directives;
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ style: {} },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ style: {} },
});
const __VLS_0 = {}.ACard;
/** @type {[typeof __VLS_components.ACard, typeof __VLS_components.aCard, typeof __VLS_components.ACard, typeof __VLS_components.aCard, ]} */ ;
// @ts-ignore
const __VLS_1 = __VLS_asFunctionalComponent(__VLS_0, new __VLS_0({
    title: "号源池列表",
    size: "small",
    loading: (__VLS_ctx.loading),
}));
const __VLS_2 = __VLS_1({
    title: "号源池列表",
    size: "small",
    loading: (__VLS_ctx.loading),
}, ...__VLS_functionalComponentArgsRest(__VLS_1));
__VLS_3.slots.default;
const __VLS_4 = {}.AList;
/** @type {[typeof __VLS_components.AList, typeof __VLS_components.aList, typeof __VLS_components.AList, typeof __VLS_components.aList, ]} */ ;
// @ts-ignore
const __VLS_5 = __VLS_asFunctionalComponent(__VLS_4, new __VLS_4({
    dataSource: (__VLS_ctx.pools),
    size: "small",
}));
const __VLS_6 = __VLS_5({
    dataSource: (__VLS_ctx.pools),
    size: "small",
}, ...__VLS_functionalComponentArgsRest(__VLS_5));
__VLS_7.slots.default;
{
    const { renderItem: __VLS_thisSlot } = __VLS_7.slots;
    const [{ item }] = __VLS_getSlotParams(__VLS_thisSlot);
    const __VLS_8 = {}.AListItem;
    /** @type {[typeof __VLS_components.AListItem, typeof __VLS_components.aListItem, typeof __VLS_components.AListItem, typeof __VLS_components.aListItem, ]} */ ;
    // @ts-ignore
    const __VLS_9 = __VLS_asFunctionalComponent(__VLS_8, new __VLS_8({
        ...{ 'onClick': {} },
        ...{ style: ({ cursor: 'pointer', background: __VLS_ctx.selectedPool?.id === item.id ? '#e6f4ff' : 'transparent', borderRadius: '6px', padding: '8px 12px' }) },
    }));
    const __VLS_10 = __VLS_9({
        ...{ 'onClick': {} },
        ...{ style: ({ cursor: 'pointer', background: __VLS_ctx.selectedPool?.id === item.id ? '#e6f4ff' : 'transparent', borderRadius: '6px', padding: '8px 12px' }) },
    }, ...__VLS_functionalComponentArgsRest(__VLS_9));
    let __VLS_12;
    let __VLS_13;
    let __VLS_14;
    const __VLS_15 = {
        onClick: (...[$event]) => {
            __VLS_ctx.fetchSlots(item);
        }
    };
    __VLS_11.slots.default;
    const __VLS_16 = {}.AListItemMeta;
    /** @type {[typeof __VLS_components.AListItemMeta, typeof __VLS_components.aListItemMeta, ]} */ ;
    // @ts-ignore
    const __VLS_17 = __VLS_asFunctionalComponent(__VLS_16, new __VLS_16({
        title: (item.name),
        description: (`${item.device_name} · ${item.exam_item_name}`),
    }));
    const __VLS_18 = __VLS_17({
        title: (item.name),
        description: (`${item.device_name} · ${item.exam_item_name}`),
    }, ...__VLS_functionalComponentArgsRest(__VLS_17));
    var __VLS_11;
}
var __VLS_7;
var __VLS_3;
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ style: {} },
});
const __VLS_20 = {}.ACard;
/** @type {[typeof __VLS_components.ACard, typeof __VLS_components.aCard, typeof __VLS_components.ACard, typeof __VLS_components.aCard, ]} */ ;
// @ts-ignore
const __VLS_21 = __VLS_asFunctionalComponent(__VLS_20, new __VLS_20({
    title: (__VLS_ctx.selectedPool ? `「${__VLS_ctx.selectedPool.name}」号源明细` : '请选择号源池'),
    size: "small",
    loading: (__VLS_ctx.slotsLoading),
}));
const __VLS_22 = __VLS_21({
    title: (__VLS_ctx.selectedPool ? `「${__VLS_ctx.selectedPool.name}」号源明细` : '请选择号源池'),
    size: "small",
    loading: (__VLS_ctx.slotsLoading),
}, ...__VLS_functionalComponentArgsRest(__VLS_21));
__VLS_23.slots.default;
if (__VLS_ctx.selectedPool) {
    const __VLS_24 = {}.ATable;
    /** @type {[typeof __VLS_components.ATable, typeof __VLS_components.aTable, typeof __VLS_components.ATable, typeof __VLS_components.aTable, ]} */ ;
    // @ts-ignore
    const __VLS_25 = __VLS_asFunctionalComponent(__VLS_24, new __VLS_24({
        columns: (__VLS_ctx.slotColumns),
        dataSource: (__VLS_ctx.slots),
        rowKey: "id",
        size: "small",
        pagination: ({ pageSize: 20 }),
    }));
    const __VLS_26 = __VLS_25({
        columns: (__VLS_ctx.slotColumns),
        dataSource: (__VLS_ctx.slots),
        rowKey: "id",
        size: "small",
        pagination: ({ pageSize: 20 }),
    }, ...__VLS_functionalComponentArgsRest(__VLS_25));
    __VLS_27.slots.default;
    {
        const { bodyCell: __VLS_thisSlot } = __VLS_27.slots;
        const [{ column, record }] = __VLS_getSlotParams(__VLS_thisSlot);
        if (column.key === 'status') {
            const __VLS_28 = {}.ABadge;
            /** @type {[typeof __VLS_components.ABadge, typeof __VLS_components.aBadge, ]} */ ;
            // @ts-ignore
            const __VLS_29 = __VLS_asFunctionalComponent(__VLS_28, new __VLS_28({
                status: __VLS_ctx.STATUS_COLOR[record.status],
                text: (record.status),
            }));
            const __VLS_30 = __VLS_29({
                status: __VLS_ctx.STATUS_COLOR[record.status],
                text: (record.status),
            }, ...__VLS_functionalComponentArgsRest(__VLS_29));
        }
        if (column.key === 'actions') {
            if (record.status === 'booked') {
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
                        if (!(__VLS_ctx.selectedPool))
                            return;
                        if (!(column.key === 'actions'))
                            return;
                        if (!(record.status === 'booked'))
                            return;
                        __VLS_ctx.releaseSlot(record);
                    }
                };
                __VLS_35.slots.default;
                var __VLS_35;
            }
        }
    }
    var __VLS_27;
}
else {
    const __VLS_40 = {}.AEmpty;
    /** @type {[typeof __VLS_components.AEmpty, typeof __VLS_components.aEmpty, ]} */ ;
    // @ts-ignore
    const __VLS_41 = __VLS_asFunctionalComponent(__VLS_40, new __VLS_40({
        description: "选择左侧号源池查看明细",
    }));
    const __VLS_42 = __VLS_41({
        description: "选择左侧号源池查看明细",
    }, ...__VLS_functionalComponentArgsRest(__VLS_41));
}
var __VLS_23;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            pools: pools,
            slots: slots,
            loading: loading,
            slotsLoading: slotsLoading,
            selectedPool: selectedPool,
            STATUS_COLOR: STATUS_COLOR,
            fetchSlots: fetchSlots,
            releaseSlot: releaseSlot,
            slotColumns: slotColumns,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
