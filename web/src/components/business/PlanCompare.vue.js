const __VLS_props = defineProps();
const emit = defineEmits();
debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
const __VLS_ctx = {};
let __VLS_components;
let __VLS_directives;
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ style: {} },
});
for (const [plan] of __VLS_getVForSourceType((__VLS_ctx.plans))) {
    const __VLS_0 = {}.ACard;
    /** @type {[typeof __VLS_components.ACard, typeof __VLS_components.aCard, typeof __VLS_components.ACard, typeof __VLS_components.aCard, ]} */ ;
    // @ts-ignore
    const __VLS_1 = __VLS_asFunctionalComponent(__VLS_0, new __VLS_0({
        ...{ 'onClick': {} },
        key: (plan.index),
        size: "small",
        ...{ style: ({
                minWidth: '240px',
                cursor: 'pointer',
                border: plan.index === __VLS_ctx.selectedIndex ? '2px solid #1890ff' : '1px solid #d9d9d9',
                flexShrink: 0,
            }) },
    }));
    const __VLS_2 = __VLS_1({
        ...{ 'onClick': {} },
        key: (plan.index),
        size: "small",
        ...{ style: ({
                minWidth: '240px',
                cursor: 'pointer',
                border: plan.index === __VLS_ctx.selectedIndex ? '2px solid #1890ff' : '1px solid #d9d9d9',
                flexShrink: 0,
            }) },
    }, ...__VLS_functionalComponentArgsRest(__VLS_1));
    let __VLS_4;
    let __VLS_5;
    let __VLS_6;
    const __VLS_7 = {
        onClick: (...[$event]) => {
            __VLS_ctx.emit('select', plan.index);
        }
    };
    __VLS_3.slots.default;
    {
        const { title: __VLS_thisSlot } = __VLS_3.slots;
        __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
            ...{ style: {} },
        });
        __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({});
        (plan.index + 1);
        if (plan.index === __VLS_ctx.selectedIndex) {
            const __VLS_8 = {}.ATag;
            /** @type {[typeof __VLS_components.ATag, typeof __VLS_components.aTag, typeof __VLS_components.ATag, typeof __VLS_components.aTag, ]} */ ;
            // @ts-ignore
            const __VLS_9 = __VLS_asFunctionalComponent(__VLS_8, new __VLS_8({
                color: "blue",
                ...{ style: {} },
            }));
            const __VLS_10 = __VLS_9({
                color: "blue",
                ...{ style: {} },
            }, ...__VLS_functionalComponentArgsRest(__VLS_9));
            __VLS_11.slots.default;
            var __VLS_11;
        }
    }
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ style: {} },
    });
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ style: {} },
    });
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ style: {} },
    });
    (Math.round(plan.total_duration_min / 60 * 10) / 10);
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ style: {} },
    });
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ style: {} },
    });
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ style: {} },
    });
    (plan.visit_count);
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ style: {} },
    });
    const __VLS_12 = {}.ATimeline;
    /** @type {[typeof __VLS_components.ATimeline, typeof __VLS_components.aTimeline, typeof __VLS_components.ATimeline, typeof __VLS_components.aTimeline, ]} */ ;
    // @ts-ignore
    const __VLS_13 = __VLS_asFunctionalComponent(__VLS_12, new __VLS_12({
        ...{ style: ({ fontSize: '12px' }) },
    }));
    const __VLS_14 = __VLS_13({
        ...{ style: ({ fontSize: '12px' }) },
    }, ...__VLS_functionalComponentArgsRest(__VLS_13));
    __VLS_15.slots.default;
    for (const [item] of __VLS_getVForSourceType((plan.items))) {
        const __VLS_16 = {}.ATimelineItem;
        /** @type {[typeof __VLS_components.ATimelineItem, typeof __VLS_components.aTimelineItem, typeof __VLS_components.ATimelineItem, typeof __VLS_components.aTimelineItem, ]} */ ;
        // @ts-ignore
        const __VLS_17 = __VLS_asFunctionalComponent(__VLS_16, new __VLS_16({
            key: (item.exam_item_name),
        }));
        const __VLS_18 = __VLS_17({
            key: (item.exam_item_name),
        }, ...__VLS_functionalComponentArgsRest(__VLS_17));
        __VLS_19.slots.default;
        __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({});
        (item.exam_item_name);
        __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
            ...{ style: {} },
        });
        (item.date);
        (item.start_time);
        (item.device_name);
        var __VLS_19;
    }
    var __VLS_15;
    if (plan.conflicts?.length) {
        for (const [msg] of __VLS_getVForSourceType((plan.conflicts))) {
            const __VLS_20 = {}.AAlert;
            /** @type {[typeof __VLS_components.AAlert, typeof __VLS_components.aAlert, ]} */ ;
            // @ts-ignore
            const __VLS_21 = __VLS_asFunctionalComponent(__VLS_20, new __VLS_20({
                key: (msg),
                message: (msg),
                type: "error",
                banner: true,
                ...{ style: {} },
            }));
            const __VLS_22 = __VLS_21({
                key: (msg),
                message: (msg),
                type: "error",
                banner: true,
                ...{ style: {} },
            }, ...__VLS_functionalComponentArgsRest(__VLS_21));
        }
    }
    if (plan.warnings?.length) {
        for (const [msg] of __VLS_getVForSourceType((plan.warnings))) {
            const __VLS_24 = {}.AAlert;
            /** @type {[typeof __VLS_components.AAlert, typeof __VLS_components.aAlert, ]} */ ;
            // @ts-ignore
            const __VLS_25 = __VLS_asFunctionalComponent(__VLS_24, new __VLS_24({
                key: (msg),
                message: (msg),
                type: "warning",
                banner: true,
                ...{ style: {} },
            }));
            const __VLS_26 = __VLS_25({
                key: (msg),
                message: (msg),
                type: "warning",
                banner: true,
                ...{ style: {} },
            }, ...__VLS_functionalComponentArgsRest(__VLS_25));
        }
    }
    var __VLS_3;
}
if (!__VLS_ctx.plans.length) {
    const __VLS_28 = {}.AEmpty;
    /** @type {[typeof __VLS_components.AEmpty, typeof __VLS_components.aEmpty, ]} */ ;
    // @ts-ignore
    const __VLS_29 = __VLS_asFunctionalComponent(__VLS_28, new __VLS_28({
        description: "暂无方案",
    }));
    const __VLS_30 = __VLS_29({
        description: "暂无方案",
    }, ...__VLS_functionalComponentArgsRest(__VLS_29));
}
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            emit: emit,
        };
    },
    __typeEmits: {},
    __typeProps: {},
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
    __typeEmits: {},
    __typeProps: {},
});
; /* PartiallyEnd: #4569/main.vue */
