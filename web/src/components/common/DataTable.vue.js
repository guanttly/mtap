/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
export default ((__VLS_props, __VLS_ctx, __VLS_expose, __VLS_setup = (async () => {
    const __VLS_props = defineProps();
    const emit = defineEmits();
    function handleChange(pag) {
        emit('change', pag);
    }
    debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
    const __VLS_fnComponent = (await import('vue')).defineComponent({
        __typeEmits: {},
    });
    const __VLS_ctx = {};
    let __VLS_components;
    let __VLS_directives;
    const __VLS_0 = {}.ATable;
    /** @type {[typeof __VLS_components.ATable, typeof __VLS_components.aTable, typeof __VLS_components.ATable, typeof __VLS_components.aTable, ]} */ ;
    // @ts-ignore
    const __VLS_1 = __VLS_asFunctionalComponent(__VLS_0, new __VLS_0({
        ...{ 'onChange': {} },
        columns: (__VLS_ctx.columns),
        dataSource: (__VLS_ctx.dataSource),
        loading: (__VLS_ctx.loading),
        rowKey: (__VLS_ctx.rowKey ?? 'id'),
        pagination: (__VLS_ctx.pagination !== undefined ? __VLS_ctx.pagination : { showSizeChanger: true, showTotal: (t) => `共 ${t} 条` }),
        size: "middle",
    }));
    const __VLS_2 = __VLS_1({
        ...{ 'onChange': {} },
        columns: (__VLS_ctx.columns),
        dataSource: (__VLS_ctx.dataSource),
        loading: (__VLS_ctx.loading),
        rowKey: (__VLS_ctx.rowKey ?? 'id'),
        pagination: (__VLS_ctx.pagination !== undefined ? __VLS_ctx.pagination : { showSizeChanger: true, showTotal: (t) => `共 ${t} 条` }),
        size: "middle",
    }, ...__VLS_functionalComponentArgsRest(__VLS_1));
    let __VLS_4;
    let __VLS_5;
    let __VLS_6;
    const __VLS_7 = {
        onChange: (__VLS_ctx.handleChange)
    };
    var __VLS_8 = {};
    __VLS_3.slots.default;
    for (const [_, name] of __VLS_getVForSourceType((__VLS_ctx.$slots))) {
        {
            const { [__VLS_tryAsConstant(name)]: __VLS_thisSlot } = __VLS_3.slots;
            const [slotData] = __VLS_getSlotParams(__VLS_thisSlot);
            var __VLS_9 = {
                ...(slotData ?? {}),
            };
            var __VLS_10 = __VLS_tryAsConstant(name);
        }
    }
    var __VLS_3;
    // @ts-ignore
    var __VLS_11 = __VLS_10, __VLS_12 = __VLS_9;
    var __VLS_dollars;
    const __VLS_self = (await import('vue')).defineComponent({
        setup() {
            return {
                handleChange: handleChange,
            };
        },
        __typeEmits: {},
        __typeProps: {},
    });
    return {};
})()) => ({})); /* PartiallyEnd: #4569/main.vue */
