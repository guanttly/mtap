/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { Modal } from 'ant-design-vue';
function confirm(opts = {}) {
    return new Promise((resolve, reject) => {
        Modal.confirm({
            title: opts.title ?? '确认操作',
            content: opts.content ?? '确定要执行此操作吗？',
            okText: opts.okText ?? '确定',
            cancelText: '取消',
            okType: opts.danger ? 'danger' : 'primary',
            onOk: () => resolve(),
            onCancel: () => reject(new Error('cancelled')),
        });
    });
}
const __VLS_exposed = { confirm };
defineExpose(__VLS_exposed);
debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
const __VLS_ctx = {};
let __VLS_components;
let __VLS_directives;
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({});
var __VLS_0 = {};
// @ts-ignore
var __VLS_1 = __VLS_0;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
const __VLS_component = (await import('vue')).defineComponent({
    setup() {
        return {
            ...__VLS_exposed,
        };
    },
});
export default {};
; /* PartiallyEnd: #4569/main.vue */
