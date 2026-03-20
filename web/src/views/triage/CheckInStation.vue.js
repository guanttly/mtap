/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { ref } from 'vue';
import { message } from 'ant-design-vue';
import { triageApi } from '@/api/triage';
const qrCode = ref('');
const loading = ref(false);
const result = ref(null);
// 护士手动签到
const nurseForm = ref({ appointment_id: '', patient_id: '', remark: '' });
const nurseLoading = ref(false);
const activeTab = ref('kiosk');
async function kioskCheckIn() {
    if (!qrCode.value)
        return;
    loading.value = true;
    result.value = null;
    try {
        result.value = await triageApi.kioskCheckIn(qrCode.value);
        message.success('签到成功');
        qrCode.value = '';
    }
    catch (e) {
        message.error(e?.message ?? '签到失败');
    }
    finally {
        loading.value = false;
    }
}
async function nurseCheckIn() {
    if (!nurseForm.value.appointment_id)
        return;
    nurseLoading.value = true;
    result.value = null;
    try {
        result.value = await triageApi.nurseCheckIn(nurseForm.value);
        message.success('签到成功');
    }
    catch (e) {
        message.error(e?.message ?? '签到失败');
    }
    finally {
        nurseLoading.value = false;
    }
}
debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
const __VLS_ctx = {};
let __VLS_components;
let __VLS_directives;
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ style: {} },
});
const __VLS_0 = {}.ATabs;
/** @type {[typeof __VLS_components.ATabs, typeof __VLS_components.aTabs, typeof __VLS_components.ATabs, typeof __VLS_components.aTabs, ]} */ ;
// @ts-ignore
const __VLS_1 = __VLS_asFunctionalComponent(__VLS_0, new __VLS_0({
    activeKey: (__VLS_ctx.activeTab),
}));
const __VLS_2 = __VLS_1({
    activeKey: (__VLS_ctx.activeTab),
}, ...__VLS_functionalComponentArgsRest(__VLS_1));
__VLS_3.slots.default;
const __VLS_4 = {}.ATabPane;
/** @type {[typeof __VLS_components.ATabPane, typeof __VLS_components.aTabPane, typeof __VLS_components.ATabPane, typeof __VLS_components.aTabPane, ]} */ ;
// @ts-ignore
const __VLS_5 = __VLS_asFunctionalComponent(__VLS_4, new __VLS_4({
    key: "kiosk",
    tab: "自助签到",
}));
const __VLS_6 = __VLS_5({
    key: "kiosk",
    tab: "自助签到",
}, ...__VLS_functionalComponentArgsRest(__VLS_5));
__VLS_7.slots.default;
const __VLS_8 = {}.ACard;
/** @type {[typeof __VLS_components.ACard, typeof __VLS_components.aCard, typeof __VLS_components.ACard, typeof __VLS_components.aCard, ]} */ ;
// @ts-ignore
const __VLS_9 = __VLS_asFunctionalComponent(__VLS_8, new __VLS_8({
    size: "small",
}));
const __VLS_10 = __VLS_9({
    size: "small",
}, ...__VLS_functionalComponentArgsRest(__VLS_9));
__VLS_11.slots.default;
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ style: {} },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ style: {} },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.p, __VLS_intrinsicElements.p)({
    ...{ style: {} },
});
const __VLS_12 = {}.AInputSearch;
/** @type {[typeof __VLS_components.AInputSearch, typeof __VLS_components.aInputSearch, ]} */ ;
// @ts-ignore
const __VLS_13 = __VLS_asFunctionalComponent(__VLS_12, new __VLS_12({
    ...{ 'onSearch': {} },
    value: (__VLS_ctx.qrCode),
    size: "large",
    placeholder: "扫描或粘贴二维码内容",
    enterButton: "签 到",
    loading: (__VLS_ctx.loading),
    ...{ style: {} },
}));
const __VLS_14 = __VLS_13({
    ...{ 'onSearch': {} },
    value: (__VLS_ctx.qrCode),
    size: "large",
    placeholder: "扫描或粘贴二维码内容",
    enterButton: "签 到",
    loading: (__VLS_ctx.loading),
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_13));
let __VLS_16;
let __VLS_17;
let __VLS_18;
const __VLS_19 = {
    onSearch: (__VLS_ctx.kioskCheckIn)
};
var __VLS_15;
var __VLS_11;
var __VLS_7;
const __VLS_20 = {}.ATabPane;
/** @type {[typeof __VLS_components.ATabPane, typeof __VLS_components.aTabPane, typeof __VLS_components.ATabPane, typeof __VLS_components.aTabPane, ]} */ ;
// @ts-ignore
const __VLS_21 = __VLS_asFunctionalComponent(__VLS_20, new __VLS_20({
    key: "nurse",
    tab: "护士辅助签到",
}));
const __VLS_22 = __VLS_21({
    key: "nurse",
    tab: "护士辅助签到",
}, ...__VLS_functionalComponentArgsRest(__VLS_21));
__VLS_23.slots.default;
const __VLS_24 = {}.ACard;
/** @type {[typeof __VLS_components.ACard, typeof __VLS_components.aCard, typeof __VLS_components.ACard, typeof __VLS_components.aCard, ]} */ ;
// @ts-ignore
const __VLS_25 = __VLS_asFunctionalComponent(__VLS_24, new __VLS_24({
    size: "small",
}));
const __VLS_26 = __VLS_25({
    size: "small",
}, ...__VLS_functionalComponentArgsRest(__VLS_25));
__VLS_27.slots.default;
const __VLS_28 = {}.AForm;
/** @type {[typeof __VLS_components.AForm, typeof __VLS_components.aForm, typeof __VLS_components.AForm, typeof __VLS_components.aForm, ]} */ ;
// @ts-ignore
const __VLS_29 = __VLS_asFunctionalComponent(__VLS_28, new __VLS_28({
    model: (__VLS_ctx.nurseForm),
    layout: "vertical",
}));
const __VLS_30 = __VLS_29({
    model: (__VLS_ctx.nurseForm),
    layout: "vertical",
}, ...__VLS_functionalComponentArgsRest(__VLS_29));
__VLS_31.slots.default;
const __VLS_32 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_33 = __VLS_asFunctionalComponent(__VLS_32, new __VLS_32({
    label: "预约ID",
}));
const __VLS_34 = __VLS_33({
    label: "预约ID",
}, ...__VLS_functionalComponentArgsRest(__VLS_33));
__VLS_35.slots.default;
const __VLS_36 = {}.AInput;
/** @type {[typeof __VLS_components.AInput, typeof __VLS_components.aInput, ]} */ ;
// @ts-ignore
const __VLS_37 = __VLS_asFunctionalComponent(__VLS_36, new __VLS_36({
    value: (__VLS_ctx.nurseForm.appointment_id),
}));
const __VLS_38 = __VLS_37({
    value: (__VLS_ctx.nurseForm.appointment_id),
}, ...__VLS_functionalComponentArgsRest(__VLS_37));
var __VLS_35;
const __VLS_40 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_41 = __VLS_asFunctionalComponent(__VLS_40, new __VLS_40({
    label: "患者ID",
}));
const __VLS_42 = __VLS_41({
    label: "患者ID",
}, ...__VLS_functionalComponentArgsRest(__VLS_41));
__VLS_43.slots.default;
const __VLS_44 = {}.AInput;
/** @type {[typeof __VLS_components.AInput, typeof __VLS_components.aInput, ]} */ ;
// @ts-ignore
const __VLS_45 = __VLS_asFunctionalComponent(__VLS_44, new __VLS_44({
    value: (__VLS_ctx.nurseForm.patient_id),
}));
const __VLS_46 = __VLS_45({
    value: (__VLS_ctx.nurseForm.patient_id),
}, ...__VLS_functionalComponentArgsRest(__VLS_45));
var __VLS_43;
const __VLS_48 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_49 = __VLS_asFunctionalComponent(__VLS_48, new __VLS_48({
    label: "备注",
}));
const __VLS_50 = __VLS_49({
    label: "备注",
}, ...__VLS_functionalComponentArgsRest(__VLS_49));
__VLS_51.slots.default;
const __VLS_52 = {}.AInput;
/** @type {[typeof __VLS_components.AInput, typeof __VLS_components.aInput, ]} */ ;
// @ts-ignore
const __VLS_53 = __VLS_asFunctionalComponent(__VLS_52, new __VLS_52({
    value: (__VLS_ctx.nurseForm.remark),
}));
const __VLS_54 = __VLS_53({
    value: (__VLS_ctx.nurseForm.remark),
}, ...__VLS_functionalComponentArgsRest(__VLS_53));
var __VLS_51;
const __VLS_56 = {}.AButton;
/** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
// @ts-ignore
const __VLS_57 = __VLS_asFunctionalComponent(__VLS_56, new __VLS_56({
    ...{ 'onClick': {} },
    type: "primary",
    loading: (__VLS_ctx.nurseLoading),
}));
const __VLS_58 = __VLS_57({
    ...{ 'onClick': {} },
    type: "primary",
    loading: (__VLS_ctx.nurseLoading),
}, ...__VLS_functionalComponentArgsRest(__VLS_57));
let __VLS_60;
let __VLS_61;
let __VLS_62;
const __VLS_63 = {
    onClick: (__VLS_ctx.nurseCheckIn)
};
__VLS_59.slots.default;
var __VLS_59;
var __VLS_31;
var __VLS_27;
var __VLS_23;
var __VLS_3;
if (__VLS_ctx.result) {
    const __VLS_64 = {}.ACard;
    /** @type {[typeof __VLS_components.ACard, typeof __VLS_components.aCard, typeof __VLS_components.ACard, typeof __VLS_components.aCard, ]} */ ;
    // @ts-ignore
    const __VLS_65 = __VLS_asFunctionalComponent(__VLS_64, new __VLS_64({
        title: "签到结果",
        size: "small",
        ...{ style: {} },
    }));
    const __VLS_66 = __VLS_65({
        title: "签到结果",
        size: "small",
        ...{ style: {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_65));
    __VLS_67.slots.default;
    const __VLS_68 = {}.AResult;
    /** @type {[typeof __VLS_components.AResult, typeof __VLS_components.aResult, typeof __VLS_components.AResult, typeof __VLS_components.aResult, ]} */ ;
    // @ts-ignore
    const __VLS_69 = __VLS_asFunctionalComponent(__VLS_68, new __VLS_68({
        status: "success",
        title: "签到成功",
        subTitle: (`队列编号：${__VLS_ctx.result.queue_number}${__VLS_ctx.result.is_late ? '（迟到）' : ''}`),
    }));
    const __VLS_70 = __VLS_69({
        status: "success",
        title: "签到成功",
        subTitle: (`队列编号：${__VLS_ctx.result.queue_number}${__VLS_ctx.result.is_late ? '（迟到）' : ''}`),
    }, ...__VLS_functionalComponentArgsRest(__VLS_69));
    __VLS_71.slots.default;
    {
        const { extra: __VLS_thisSlot } = __VLS_71.slots;
        const __VLS_72 = {}.ADescriptions;
        /** @type {[typeof __VLS_components.ADescriptions, typeof __VLS_components.aDescriptions, typeof __VLS_components.ADescriptions, typeof __VLS_components.aDescriptions, ]} */ ;
        // @ts-ignore
        const __VLS_73 = __VLS_asFunctionalComponent(__VLS_72, new __VLS_72({
            column: (1),
            size: "small",
            bordered: true,
        }));
        const __VLS_74 = __VLS_73({
            column: (1),
            size: "small",
            bordered: true,
        }, ...__VLS_functionalComponentArgsRest(__VLS_73));
        __VLS_75.slots.default;
        const __VLS_76 = {}.ADescriptionsItem;
        /** @type {[typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, ]} */ ;
        // @ts-ignore
        const __VLS_77 = __VLS_asFunctionalComponent(__VLS_76, new __VLS_76({
            label: "签到ID",
        }));
        const __VLS_78 = __VLS_77({
            label: "签到ID",
        }, ...__VLS_functionalComponentArgsRest(__VLS_77));
        __VLS_79.slots.default;
        (__VLS_ctx.result.check_in_id);
        var __VLS_79;
        const __VLS_80 = {}.ADescriptionsItem;
        /** @type {[typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, ]} */ ;
        // @ts-ignore
        const __VLS_81 = __VLS_asFunctionalComponent(__VLS_80, new __VLS_80({
            label: "队列编号",
        }));
        const __VLS_82 = __VLS_81({
            label: "队列编号",
        }, ...__VLS_functionalComponentArgsRest(__VLS_81));
        __VLS_83.slots.default;
        (__VLS_ctx.result.queue_number);
        var __VLS_83;
        const __VLS_84 = {}.ADescriptionsItem;
        /** @type {[typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, ]} */ ;
        // @ts-ignore
        const __VLS_85 = __VLS_asFunctionalComponent(__VLS_84, new __VLS_84({
            label: "预计等待",
        }));
        const __VLS_86 = __VLS_85({
            label: "预计等待",
        }, ...__VLS_functionalComponentArgsRest(__VLS_85));
        __VLS_87.slots.default;
        (__VLS_ctx.result.estimated_wait);
        var __VLS_87;
        const __VLS_88 = {}.ADescriptionsItem;
        /** @type {[typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, ]} */ ;
        // @ts-ignore
        const __VLS_89 = __VLS_asFunctionalComponent(__VLS_88, new __VLS_88({
            label: "就诊地点",
        }));
        const __VLS_90 = __VLS_89({
            label: "就诊地点",
        }, ...__VLS_functionalComponentArgsRest(__VLS_89));
        __VLS_91.slots.default;
        (__VLS_ctx.result.room_location);
        var __VLS_91;
        var __VLS_75;
    }
    var __VLS_71;
    var __VLS_67;
}
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            qrCode: qrCode,
            loading: loading,
            result: result,
            nurseForm: nurseForm,
            nurseLoading: nurseLoading,
            activeTab: activeTab,
            kioskCheckIn: kioskCheckIn,
            nurseCheckIn: nurseCheckIn,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
