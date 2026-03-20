/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { computed } from 'vue';
const props = defineProps();
const emit = defineEmits();
const availableSlots = computed(() => props.slots.filter(s => s.status === 'available'));
const otherSlots = computed(() => props.slots.filter(s => s.status !== 'available'));
function selectSlot(slot) {
    emit('update:selectedSlotId', slot.id);
    emit('select', slot);
}
function slotBorderColor(slot) {
    if (slot.id === props.selectedSlotId)
        return '#1890ff';
    return '#d9d9d9';
}
function slotBg(slot) {
    if (slot.id === props.selectedSlotId)
        return '#e6f4ff';
    if (slot.status === 'available')
        return '#f6ffed';
    return '#f5f5f5';
}
debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
const __VLS_ctx = {};
let __VLS_components;
let __VLS_directives;
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({});
if (__VLS_ctx.slots.length === 0) {
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({});
    const __VLS_0 = {}.AEmpty;
    /** @type {[typeof __VLS_components.AEmpty, typeof __VLS_components.aEmpty, ]} */ ;
    // @ts-ignore
    const __VLS_1 = __VLS_asFunctionalComponent(__VLS_0, new __VLS_0({
        description: "无可用号源",
        imageStyle: ({ height: '40px' }),
    }));
    const __VLS_2 = __VLS_1({
        description: "无可用号源",
        imageStyle: ({ height: '40px' }),
    }, ...__VLS_functionalComponentArgsRest(__VLS_1));
}
else {
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ style: {} },
    });
    for (const [slot] of __VLS_getVForSourceType(([...__VLS_ctx.availableSlots, ...__VLS_ctx.otherSlots]))) {
        __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
            ...{ onClick: (...[$event]) => {
                    if (!!(__VLS_ctx.slots.length === 0))
                        return;
                    slot.status === 'available' && __VLS_ctx.selectSlot(slot);
                } },
            key: (slot.id),
            ...{ style: ({
                    background: __VLS_ctx.slotBg(slot),
                    border: `1px solid ${__VLS_ctx.slotBorderColor(slot)}`,
                    borderRadius: '6px',
                    padding: '8px 12px',
                    cursor: slot.status === 'available' ? 'pointer' : 'not-allowed',
                    opacity: slot.status === 'available' ? 1 : 0.5,
                    minWidth: '120px',
                    textAlign: 'center',
                }) },
        });
        __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
            ...{ style: {} },
        });
        (slot.start_time);
        (slot.end_time);
        __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
            ...{ style: {} },
        });
        (slot.status === 'available' ? '可预约' : slot.status);
        __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
            ...{ style: {} },
        });
        (slot.standard_duration);
    }
}
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            availableSlots: availableSlots,
            otherSlots: otherSlots,
            selectSlot: selectSlot,
            slotBorderColor: slotBorderColor,
            slotBg: slotBg,
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
