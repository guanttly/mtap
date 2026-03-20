/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { ref, onMounted, onUnmounted } from 'vue';
import { triageApi } from '@/api/triage';
const roomId = ref('room-001');
const status = ref(null);
const loading = ref(false);
let timer = null;
const STATUS_COLOR = {
    waiting: 'default', called: 'processing', in_exam: 'warning', completed: 'success', missed: 'error',
};
const STATUS_LABEL = {
    waiting: '等待中', called: '已叫号', in_exam: '检查中', completed: '已完成', missed: '已过号',
};
async function fetchStatus() {
    loading.value = true;
    try {
        status.value = await triageApi.getQueueStatus(roomId.value);
    }
    finally {
        loading.value = false;
    }
}
onMounted(() => {
    fetchStatus();
    timer = setInterval(fetchStatus, 15000);
});
onUnmounted(() => { if (timer)
    clearInterval(timer); });
debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
const __VLS_ctx = {};
let __VLS_components;
let __VLS_directives;
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ style: {} },
});
const __VLS_0 = {}.AInput;
/** @type {[typeof __VLS_components.AInput, typeof __VLS_components.aInput, ]} */ ;
// @ts-ignore
const __VLS_1 = __VLS_asFunctionalComponent(__VLS_0, new __VLS_0({
    value: (__VLS_ctx.roomId),
    placeholder: "检查室ID",
    ...{ style: {} },
}));
const __VLS_2 = __VLS_1({
    value: (__VLS_ctx.roomId),
    placeholder: "检查室ID",
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_1));
const __VLS_4 = {}.AButton;
/** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
// @ts-ignore
const __VLS_5 = __VLS_asFunctionalComponent(__VLS_4, new __VLS_4({
    ...{ 'onClick': {} },
    type: "primary",
    loading: (__VLS_ctx.loading),
}));
const __VLS_6 = __VLS_5({
    ...{ 'onClick': {} },
    type: "primary",
    loading: (__VLS_ctx.loading),
}, ...__VLS_functionalComponentArgsRest(__VLS_5));
let __VLS_8;
let __VLS_9;
let __VLS_10;
const __VLS_11 = {
    onClick: (__VLS_ctx.fetchStatus)
};
__VLS_7.slots.default;
var __VLS_7;
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
    ...{ class: "text-muted" },
    ...{ style: {} },
});
if (__VLS_ctx.status) {
    const __VLS_12 = {}.ARow;
    /** @type {[typeof __VLS_components.ARow, typeof __VLS_components.aRow, typeof __VLS_components.ARow, typeof __VLS_components.aRow, ]} */ ;
    // @ts-ignore
    const __VLS_13 = __VLS_asFunctionalComponent(__VLS_12, new __VLS_12({
        gutter: (16),
        ...{ style: {} },
    }));
    const __VLS_14 = __VLS_13({
        gutter: (16),
        ...{ style: {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_13));
    __VLS_15.slots.default;
    const __VLS_16 = {}.ACol;
    /** @type {[typeof __VLS_components.ACol, typeof __VLS_components.aCol, typeof __VLS_components.ACol, typeof __VLS_components.aCol, ]} */ ;
    // @ts-ignore
    const __VLS_17 = __VLS_asFunctionalComponent(__VLS_16, new __VLS_16({
        span: (8),
    }));
    const __VLS_18 = __VLS_17({
        span: (8),
    }, ...__VLS_functionalComponentArgsRest(__VLS_17));
    __VLS_19.slots.default;
    const __VLS_20 = {}.AStatistic;
    /** @type {[typeof __VLS_components.AStatistic, typeof __VLS_components.aStatistic, ]} */ ;
    // @ts-ignore
    const __VLS_21 = __VLS_asFunctionalComponent(__VLS_20, new __VLS_20({
        title: "等待人数",
        value: (__VLS_ctx.status.waiting_count),
    }));
    const __VLS_22 = __VLS_21({
        title: "等待人数",
        value: (__VLS_ctx.status.waiting_count),
    }, ...__VLS_functionalComponentArgsRest(__VLS_21));
    var __VLS_19;
    const __VLS_24 = {}.ACol;
    /** @type {[typeof __VLS_components.ACol, typeof __VLS_components.aCol, typeof __VLS_components.ACol, typeof __VLS_components.aCol, ]} */ ;
    // @ts-ignore
    const __VLS_25 = __VLS_asFunctionalComponent(__VLS_24, new __VLS_24({
        span: (8),
    }));
    const __VLS_26 = __VLS_25({
        span: (8),
    }, ...__VLS_functionalComponentArgsRest(__VLS_25));
    __VLS_27.slots.default;
    const __VLS_28 = {}.AStatistic;
    /** @type {[typeof __VLS_components.AStatistic, typeof __VLS_components.aStatistic, ]} */ ;
    // @ts-ignore
    const __VLS_29 = __VLS_asFunctionalComponent(__VLS_28, new __VLS_28({
        title: "平均等待(分)",
        value: (__VLS_ctx.status.average_wait),
        precision: (1),
    }));
    const __VLS_30 = __VLS_29({
        title: "平均等待(分)",
        value: (__VLS_ctx.status.average_wait),
        precision: (1),
    }, ...__VLS_functionalComponentArgsRest(__VLS_29));
    var __VLS_27;
    const __VLS_32 = {}.ACol;
    /** @type {[typeof __VLS_components.ACol, typeof __VLS_components.aCol, typeof __VLS_components.ACol, typeof __VLS_components.aCol, ]} */ ;
    // @ts-ignore
    const __VLS_33 = __VLS_asFunctionalComponent(__VLS_32, new __VLS_32({
        span: (8),
    }));
    const __VLS_34 = __VLS_33({
        span: (8),
    }, ...__VLS_functionalComponentArgsRest(__VLS_33));
    __VLS_35.slots.default;
    const __VLS_36 = {}.AStatistic;
    /** @type {[typeof __VLS_components.AStatistic, typeof __VLS_components.aStatistic, ]} */ ;
    // @ts-ignore
    const __VLS_37 = __VLS_asFunctionalComponent(__VLS_36, new __VLS_36({
        title: "已完成",
        value: (__VLS_ctx.status.entries?.filter(e => e.status === 'completed').length ?? 0),
    }));
    const __VLS_38 = __VLS_37({
        title: "已完成",
        value: (__VLS_ctx.status.entries?.filter(e => e.status === 'completed').length ?? 0),
    }, ...__VLS_functionalComponentArgsRest(__VLS_37));
    var __VLS_35;
    var __VLS_15;
    const __VLS_40 = {}.ACard;
    /** @type {[typeof __VLS_components.ACard, typeof __VLS_components.aCard, typeof __VLS_components.ACard, typeof __VLS_components.aCard, ]} */ ;
    // @ts-ignore
    const __VLS_41 = __VLS_asFunctionalComponent(__VLS_40, new __VLS_40({
        title: "当前队列",
        size: "small",
    }));
    const __VLS_42 = __VLS_41({
        title: "当前队列",
        size: "small",
    }, ...__VLS_functionalComponentArgsRest(__VLS_41));
    __VLS_43.slots.default;
    const __VLS_44 = {}.AList;
    /** @type {[typeof __VLS_components.AList, typeof __VLS_components.aList, typeof __VLS_components.AList, typeof __VLS_components.aList, ]} */ ;
    // @ts-ignore
    const __VLS_45 = __VLS_asFunctionalComponent(__VLS_44, new __VLS_44({
        dataSource: (__VLS_ctx.status.entries),
        size: "small",
    }));
    const __VLS_46 = __VLS_45({
        dataSource: (__VLS_ctx.status.entries),
        size: "small",
    }, ...__VLS_functionalComponentArgsRest(__VLS_45));
    __VLS_47.slots.default;
    {
        const { renderItem: __VLS_thisSlot } = __VLS_47.slots;
        const [{ item, index }] = __VLS_getSlotParams(__VLS_thisSlot);
        const __VLS_48 = {}.AListItem;
        /** @type {[typeof __VLS_components.AListItem, typeof __VLS_components.aListItem, typeof __VLS_components.AListItem, typeof __VLS_components.aListItem, ]} */ ;
        // @ts-ignore
        const __VLS_49 = __VLS_asFunctionalComponent(__VLS_48, new __VLS_48({}));
        const __VLS_50 = __VLS_49({}, ...__VLS_functionalComponentArgsRest(__VLS_49));
        __VLS_51.slots.default;
        const __VLS_52 = {}.AListItemMeta;
        /** @type {[typeof __VLS_components.AListItemMeta, typeof __VLS_components.aListItemMeta, typeof __VLS_components.AListItemMeta, typeof __VLS_components.aListItemMeta, ]} */ ;
        // @ts-ignore
        const __VLS_53 = __VLS_asFunctionalComponent(__VLS_52, new __VLS_52({}));
        const __VLS_54 = __VLS_53({}, ...__VLS_functionalComponentArgsRest(__VLS_53));
        __VLS_55.slots.default;
        {
            const { title: __VLS_thisSlot } = __VLS_55.slots;
            __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
                ...{ style: {} },
            });
            (index + 1);
            (item.patient_name_masked);
        }
        {
            const { description: __VLS_thisSlot } = __VLS_55.slots;
            (item.queue_number);
        }
        var __VLS_55;
        const __VLS_56 = {}.ABadge;
        /** @type {[typeof __VLS_components.ABadge, typeof __VLS_components.aBadge, ]} */ ;
        // @ts-ignore
        const __VLS_57 = __VLS_asFunctionalComponent(__VLS_56, new __VLS_56({
            status: __VLS_ctx.STATUS_COLOR[item.status],
            text: (__VLS_ctx.STATUS_LABEL[item.status]),
        }));
        const __VLS_58 = __VLS_57({
            status: __VLS_ctx.STATUS_COLOR[item.status],
            text: (__VLS_ctx.STATUS_LABEL[item.status]),
        }, ...__VLS_functionalComponentArgsRest(__VLS_57));
        var __VLS_51;
    }
    {
        const { empty: __VLS_thisSlot } = __VLS_47.slots;
        const __VLS_60 = {}.AEmpty;
        /** @type {[typeof __VLS_components.AEmpty, typeof __VLS_components.aEmpty, ]} */ ;
        // @ts-ignore
        const __VLS_61 = __VLS_asFunctionalComponent(__VLS_60, new __VLS_60({
            description: "队列为空",
        }));
        const __VLS_62 = __VLS_61({
            description: "队列为空",
        }, ...__VLS_functionalComponentArgsRest(__VLS_61));
    }
    var __VLS_47;
    var __VLS_43;
}
else if (!__VLS_ctx.loading) {
    const __VLS_64 = {}.AEmpty;
    /** @type {[typeof __VLS_components.AEmpty, typeof __VLS_components.aEmpty, ]} */ ;
    // @ts-ignore
    const __VLS_65 = __VLS_asFunctionalComponent(__VLS_64, new __VLS_64({
        description: "暂无数据，请选择检查室",
    }));
    const __VLS_66 = __VLS_65({
        description: "暂无数据，请选择检查室",
    }, ...__VLS_functionalComponentArgsRest(__VLS_65));
}
/** @type {__VLS_StyleScopedClasses['text-muted']} */ ;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            roomId: roomId,
            status: status,
            loading: loading,
            STATUS_COLOR: STATUS_COLOR,
            STATUS_LABEL: STATUS_LABEL,
            fetchStatus: fetchStatus,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
