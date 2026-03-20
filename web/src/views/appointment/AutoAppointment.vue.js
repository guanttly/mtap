/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { ref } from 'vue';
import { message } from 'ant-design-vue';
import { appointmentApi } from '@/api/appointment';
import PlanCompare from '@/components/business/PlanCompare.vue';
import SlotPicker from '@/components/business/SlotPicker.vue';
const step = ref(0);
const loading = ref(false);
const plans = ref([]);
const selectedPlan = ref(0);
const confirming = ref(false);
const form = ref({
    patient_id: '',
    exam_item_ids: [],
    preferences: { preferred_time_period: 'morning' },
});
async function generatePlans() {
    if (!form.value.patient_id || !form.value.exam_item_ids?.length) {
        message.warning('请填写患者信息和检查项目');
        return;
    }
    loading.value = true;
    try {
        const res = await appointmentApi.autoAppointment(form.value);
        plans.value = res.plans ?? [];
        if (plans.value.length === 0) {
            message.warning('未找到合适的预约方案，请调整条件后重试');
        }
        else {
            step.value = 1;
        }
    }
    finally {
        loading.value = false;
    }
}
async function confirmPlan() {
    confirming.value = true;
    try {
        message.success('预约确认成功');
        step.value = 2;
    }
    finally {
        confirming.value = false;
    }
}
function resetForm() {
    step.value = 0;
    plans.value = [];
    form.value = { patient_id: '', exam_item_ids: [], preferences: { preferred_time_period: 'morning' } };
}
debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
const __VLS_ctx = {};
let __VLS_components;
let __VLS_directives;
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ style: {} },
});
const __VLS_0 = {}.ASteps;
/** @type {[typeof __VLS_components.ASteps, typeof __VLS_components.aSteps, typeof __VLS_components.ASteps, typeof __VLS_components.aSteps, ]} */ ;
// @ts-ignore
const __VLS_1 = __VLS_asFunctionalComponent(__VLS_0, new __VLS_0({
    current: (__VLS_ctx.step),
    ...{ style: {} },
}));
const __VLS_2 = __VLS_1({
    current: (__VLS_ctx.step),
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_1));
__VLS_3.slots.default;
const __VLS_4 = {}.AStep;
/** @type {[typeof __VLS_components.AStep, typeof __VLS_components.aStep, ]} */ ;
// @ts-ignore
const __VLS_5 = __VLS_asFunctionalComponent(__VLS_4, new __VLS_4({
    title: "填写信息",
}));
const __VLS_6 = __VLS_5({
    title: "填写信息",
}, ...__VLS_functionalComponentArgsRest(__VLS_5));
const __VLS_8 = {}.AStep;
/** @type {[typeof __VLS_components.AStep, typeof __VLS_components.aStep, ]} */ ;
// @ts-ignore
const __VLS_9 = __VLS_asFunctionalComponent(__VLS_8, new __VLS_8({
    title: "选择方案",
}));
const __VLS_10 = __VLS_9({
    title: "选择方案",
}, ...__VLS_functionalComponentArgsRest(__VLS_9));
const __VLS_12 = {}.AStep;
/** @type {[typeof __VLS_components.AStep, typeof __VLS_components.aStep, ]} */ ;
// @ts-ignore
const __VLS_13 = __VLS_asFunctionalComponent(__VLS_12, new __VLS_12({
    title: "完成",
}));
const __VLS_14 = __VLS_13({
    title: "完成",
}, ...__VLS_functionalComponentArgsRest(__VLS_13));
var __VLS_3;
if (__VLS_ctx.step === 0) {
    const __VLS_16 = {}.ACard;
    /** @type {[typeof __VLS_components.ACard, typeof __VLS_components.aCard, typeof __VLS_components.ACard, typeof __VLS_components.aCard, ]} */ ;
    // @ts-ignore
    const __VLS_17 = __VLS_asFunctionalComponent(__VLS_16, new __VLS_16({
        title: "患者与检查信息",
        size: "small",
    }));
    const __VLS_18 = __VLS_17({
        title: "患者与检查信息",
        size: "small",
    }, ...__VLS_functionalComponentArgsRest(__VLS_17));
    __VLS_19.slots.default;
    const __VLS_20 = {}.AForm;
    /** @type {[typeof __VLS_components.AForm, typeof __VLS_components.aForm, typeof __VLS_components.AForm, typeof __VLS_components.aForm, ]} */ ;
    // @ts-ignore
    const __VLS_21 = __VLS_asFunctionalComponent(__VLS_20, new __VLS_20({
        model: (__VLS_ctx.form),
        layout: "vertical",
    }));
    const __VLS_22 = __VLS_21({
        model: (__VLS_ctx.form),
        layout: "vertical",
    }, ...__VLS_functionalComponentArgsRest(__VLS_21));
    __VLS_23.slots.default;
    const __VLS_24 = {}.AFormItem;
    /** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
    // @ts-ignore
    const __VLS_25 = __VLS_asFunctionalComponent(__VLS_24, new __VLS_24({
        label: "患者ID",
    }));
    const __VLS_26 = __VLS_25({
        label: "患者ID",
    }, ...__VLS_functionalComponentArgsRest(__VLS_25));
    __VLS_27.slots.default;
    const __VLS_28 = {}.AInput;
    /** @type {[typeof __VLS_components.AInput, typeof __VLS_components.aInput, ]} */ ;
    // @ts-ignore
    const __VLS_29 = __VLS_asFunctionalComponent(__VLS_28, new __VLS_28({
        value: (__VLS_ctx.form.patient_id),
    }));
    const __VLS_30 = __VLS_29({
        value: (__VLS_ctx.form.patient_id),
    }, ...__VLS_functionalComponentArgsRest(__VLS_29));
    var __VLS_27;
    const __VLS_32 = {}.AFormItem;
    /** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
    // @ts-ignore
    const __VLS_33 = __VLS_asFunctionalComponent(__VLS_32, new __VLS_32({
        label: "检查项目ID（逗号分隔）",
    }));
    const __VLS_34 = __VLS_33({
        label: "检查项目ID（逗号分隔）",
    }, ...__VLS_functionalComponentArgsRest(__VLS_33));
    __VLS_35.slots.default;
    const __VLS_36 = {}.ATextarea;
    /** @type {[typeof __VLS_components.ATextarea, typeof __VLS_components.aTextarea, ]} */ ;
    // @ts-ignore
    const __VLS_37 = __VLS_asFunctionalComponent(__VLS_36, new __VLS_36({
        ...{ 'onChange': {} },
        value: (__VLS_ctx.form.exam_item_ids?.join(',')),
        rows: (3),
    }));
    const __VLS_38 = __VLS_37({
        ...{ 'onChange': {} },
        value: (__VLS_ctx.form.exam_item_ids?.join(',')),
        rows: (3),
    }, ...__VLS_functionalComponentArgsRest(__VLS_37));
    let __VLS_40;
    let __VLS_41;
    let __VLS_42;
    const __VLS_43 = {
        onChange: ((e) => __VLS_ctx.form.exam_item_ids = e.target.value.split(',').map(s => s.trim()).filter(Boolean))
    };
    var __VLS_39;
    var __VLS_35;
    const __VLS_44 = {}.AFormItem;
    /** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
    // @ts-ignore
    const __VLS_45 = __VLS_asFunctionalComponent(__VLS_44, new __VLS_44({
        label: "偏好时段",
    }));
    const __VLS_46 = __VLS_45({
        label: "偏好时段",
    }, ...__VLS_functionalComponentArgsRest(__VLS_45));
    __VLS_47.slots.default;
    const __VLS_48 = {}.ARadioGroup;
    /** @type {[typeof __VLS_components.ARadioGroup, typeof __VLS_components.aRadioGroup, typeof __VLS_components.ARadioGroup, typeof __VLS_components.aRadioGroup, ]} */ ;
    // @ts-ignore
    const __VLS_49 = __VLS_asFunctionalComponent(__VLS_48, new __VLS_48({
        value: (__VLS_ctx.form.preferences.preferred_time_period),
    }));
    const __VLS_50 = __VLS_49({
        value: (__VLS_ctx.form.preferences.preferred_time_period),
    }, ...__VLS_functionalComponentArgsRest(__VLS_49));
    __VLS_51.slots.default;
    const __VLS_52 = {}.ARadio;
    /** @type {[typeof __VLS_components.ARadio, typeof __VLS_components.aRadio, typeof __VLS_components.ARadio, typeof __VLS_components.aRadio, ]} */ ;
    // @ts-ignore
    const __VLS_53 = __VLS_asFunctionalComponent(__VLS_52, new __VLS_52({
        value: "morning",
    }));
    const __VLS_54 = __VLS_53({
        value: "morning",
    }, ...__VLS_functionalComponentArgsRest(__VLS_53));
    __VLS_55.slots.default;
    var __VLS_55;
    const __VLS_56 = {}.ARadio;
    /** @type {[typeof __VLS_components.ARadio, typeof __VLS_components.aRadio, typeof __VLS_components.ARadio, typeof __VLS_components.aRadio, ]} */ ;
    // @ts-ignore
    const __VLS_57 = __VLS_asFunctionalComponent(__VLS_56, new __VLS_56({
        value: "afternoon",
    }));
    const __VLS_58 = __VLS_57({
        value: "afternoon",
    }, ...__VLS_functionalComponentArgsRest(__VLS_57));
    __VLS_59.slots.default;
    var __VLS_59;
    const __VLS_60 = {}.ARadio;
    /** @type {[typeof __VLS_components.ARadio, typeof __VLS_components.aRadio, typeof __VLS_components.ARadio, typeof __VLS_components.aRadio, ]} */ ;
    // @ts-ignore
    const __VLS_61 = __VLS_asFunctionalComponent(__VLS_60, new __VLS_60({
        value: "any",
    }));
    const __VLS_62 = __VLS_61({
        value: "any",
    }, ...__VLS_functionalComponentArgsRest(__VLS_61));
    __VLS_63.slots.default;
    var __VLS_63;
    var __VLS_51;
    var __VLS_47;
    const __VLS_64 = {}.AFormItem;
    /** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
    // @ts-ignore
    const __VLS_65 = __VLS_asFunctionalComponent(__VLS_64, new __VLS_64({
        label: "偏好日期范围",
    }));
    const __VLS_66 = __VLS_65({
        label: "偏好日期范围",
    }, ...__VLS_functionalComponentArgsRest(__VLS_65));
    __VLS_67.slots.default;
    const __VLS_68 = {}.ADatePicker;
    /** @type {[typeof __VLS_components.ADatePicker, typeof __VLS_components.aDatePicker, ]} */ ;
    // @ts-ignore
    const __VLS_69 = __VLS_asFunctionalComponent(__VLS_68, new __VLS_68({
        ...{ 'onChange': {} },
        value: (__VLS_ctx.form.preferences.preferred_date_range?.start),
        valueFormat: "YYYY-MM-DD",
        placeholder: "开始日期",
        ...{ style: {} },
    }));
    const __VLS_70 = __VLS_69({
        ...{ 'onChange': {} },
        value: (__VLS_ctx.form.preferences.preferred_date_range?.start),
        valueFormat: "YYYY-MM-DD",
        placeholder: "开始日期",
        ...{ style: {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_69));
    let __VLS_72;
    let __VLS_73;
    let __VLS_74;
    const __VLS_75 = {
        onChange: ((val) => { if (!__VLS_ctx.form.preferences.preferred_date_range)
            __VLS_ctx.form.preferences.preferred_date_range = { start: '', end: '' }; __VLS_ctx.form.preferences.preferred_date_range.start = val; })
    };
    var __VLS_71;
    __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
        ...{ style: {} },
    });
    const __VLS_76 = {}.ADatePicker;
    /** @type {[typeof __VLS_components.ADatePicker, typeof __VLS_components.aDatePicker, ]} */ ;
    // @ts-ignore
    const __VLS_77 = __VLS_asFunctionalComponent(__VLS_76, new __VLS_76({
        ...{ 'onChange': {} },
        value: (__VLS_ctx.form.preferences.preferred_date_range?.end),
        valueFormat: "YYYY-MM-DD",
        placeholder: "结束日期",
        ...{ style: {} },
    }));
    const __VLS_78 = __VLS_77({
        ...{ 'onChange': {} },
        value: (__VLS_ctx.form.preferences.preferred_date_range?.end),
        valueFormat: "YYYY-MM-DD",
        placeholder: "结束日期",
        ...{ style: {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_77));
    let __VLS_80;
    let __VLS_81;
    let __VLS_82;
    const __VLS_83 = {
        onChange: ((val) => { if (!__VLS_ctx.form.preferences.preferred_date_range)
            __VLS_ctx.form.preferences.preferred_date_range = { start: '', end: '' }; __VLS_ctx.form.preferences.preferred_date_range.end = val; })
    };
    var __VLS_79;
    var __VLS_67;
    const __VLS_84 = {}.AButton;
    /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
    // @ts-ignore
    const __VLS_85 = __VLS_asFunctionalComponent(__VLS_84, new __VLS_84({
        ...{ 'onClick': {} },
        type: "primary",
        loading: (__VLS_ctx.loading),
    }));
    const __VLS_86 = __VLS_85({
        ...{ 'onClick': {} },
        type: "primary",
        loading: (__VLS_ctx.loading),
    }, ...__VLS_functionalComponentArgsRest(__VLS_85));
    let __VLS_88;
    let __VLS_89;
    let __VLS_90;
    const __VLS_91 = {
        onClick: (__VLS_ctx.generatePlans)
    };
    __VLS_87.slots.default;
    var __VLS_87;
    var __VLS_23;
    var __VLS_19;
}
else if (__VLS_ctx.step === 1) {
    /** @type {[typeof PlanCompare, ]} */ ;
    // @ts-ignore
    const __VLS_92 = __VLS_asFunctionalComponent(PlanCompare, new PlanCompare({
        ...{ 'onSelect': {} },
        plans: (__VLS_ctx.plans),
        selectedIndex: (__VLS_ctx.selectedPlan),
    }));
    const __VLS_93 = __VLS_92({
        ...{ 'onSelect': {} },
        plans: (__VLS_ctx.plans),
        selectedIndex: (__VLS_ctx.selectedPlan),
    }, ...__VLS_functionalComponentArgsRest(__VLS_92));
    let __VLS_95;
    let __VLS_96;
    let __VLS_97;
    const __VLS_98 = {
        onSelect: (...[$event]) => {
            if (!!(__VLS_ctx.step === 0))
                return;
            if (!(__VLS_ctx.step === 1))
                return;
            __VLS_ctx.selectedPlan = $event;
        }
    };
    var __VLS_94;
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ style: {} },
    });
    const __VLS_99 = {}.AButton;
    /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
    // @ts-ignore
    const __VLS_100 = __VLS_asFunctionalComponent(__VLS_99, new __VLS_99({
        ...{ 'onClick': {} },
    }));
    const __VLS_101 = __VLS_100({
        ...{ 'onClick': {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_100));
    let __VLS_103;
    let __VLS_104;
    let __VLS_105;
    const __VLS_106 = {
        onClick: (...[$event]) => {
            if (!!(__VLS_ctx.step === 0))
                return;
            if (!(__VLS_ctx.step === 1))
                return;
            __VLS_ctx.step = 0;
        }
    };
    __VLS_102.slots.default;
    var __VLS_102;
    const __VLS_107 = {}.AButton;
    /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
    // @ts-ignore
    const __VLS_108 = __VLS_asFunctionalComponent(__VLS_107, new __VLS_107({
        ...{ 'onClick': {} },
        type: "primary",
        loading: (__VLS_ctx.confirming),
    }));
    const __VLS_109 = __VLS_108({
        ...{ 'onClick': {} },
        type: "primary",
        loading: (__VLS_ctx.confirming),
    }, ...__VLS_functionalComponentArgsRest(__VLS_108));
    let __VLS_111;
    let __VLS_112;
    let __VLS_113;
    const __VLS_114 = {
        onClick: (__VLS_ctx.confirmPlan)
    };
    __VLS_110.slots.default;
    (__VLS_ctx.selectedPlan + 1);
    var __VLS_110;
}
else {
    const __VLS_115 = {}.AResult;
    /** @type {[typeof __VLS_components.AResult, typeof __VLS_components.aResult, typeof __VLS_components.AResult, typeof __VLS_components.aResult, ]} */ ;
    // @ts-ignore
    const __VLS_116 = __VLS_asFunctionalComponent(__VLS_115, new __VLS_115({
        status: "success",
        title: "预约成功！",
        subTitle: "请患者按时就诊，凭二维码或手机号签到",
    }));
    const __VLS_117 = __VLS_116({
        status: "success",
        title: "预约成功！",
        subTitle: "请患者按时就诊，凭二维码或手机号签到",
    }, ...__VLS_functionalComponentArgsRest(__VLS_116));
    __VLS_118.slots.default;
    {
        const { extra: __VLS_thisSlot } = __VLS_118.slots;
        const __VLS_119 = {}.AButton;
        /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
        // @ts-ignore
        const __VLS_120 = __VLS_asFunctionalComponent(__VLS_119, new __VLS_119({
            ...{ 'onClick': {} },
            type: "primary",
        }));
        const __VLS_121 = __VLS_120({
            ...{ 'onClick': {} },
            type: "primary",
        }, ...__VLS_functionalComponentArgsRest(__VLS_120));
        let __VLS_123;
        let __VLS_124;
        let __VLS_125;
        const __VLS_126 = {
            onClick: (__VLS_ctx.resetForm)
        };
        __VLS_122.slots.default;
        var __VLS_122;
    }
    var __VLS_118;
}
if (false) {
    __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({});
    /** @type {[typeof SlotPicker, ]} */ ;
    // @ts-ignore
    const __VLS_127 = __VLS_asFunctionalComponent(SlotPicker, new SlotPicker({
        slots: ([]),
    }));
    const __VLS_128 = __VLS_127({
        slots: ([]),
    }, ...__VLS_functionalComponentArgsRest(__VLS_127));
}
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            PlanCompare: PlanCompare,
            SlotPicker: SlotPicker,
            step: step,
            loading: loading,
            plans: plans,
            selectedPlan: selectedPlan,
            confirming: confirming,
            form: form,
            generatePlans: generatePlans,
            confirmPlan: confirmPlan,
            resetForm: resetForm,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
