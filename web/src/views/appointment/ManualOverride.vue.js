/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { ref } from 'vue';
import { message } from 'ant-design-vue';
import { appointmentApi } from '@/api/appointment';
const form = ref({
    appointment_id: '',
    action: 'reschedule',
    reason: '',
    new_slot_id: '',
});
const loading = ref(false);
const appointment = ref(null);
const searching = ref(false);
async function lookupAppointment() {
    if (!form.value.appointment_id)
        return;
    searching.value = true;
    try {
        appointment.value = await appointmentApi.getAppointment(form.value.appointment_id);
    }
    catch {
        message.error('未找到该预约');
    }
    finally {
        searching.value = false;
    }
}
async function handleSubmit() {
    if (!form.value.appointment_id)
        return;
    loading.value = true;
    try {
        const { appointment_id, action, reason, new_slot_id } = form.value;
        if (action === 'cancel') {
            await appointmentApi.cancel(appointment_id, reason);
        }
        else if (action === 'confirm') {
            await appointmentApi.confirm(appointment_id);
        }
        else if (action === 'mark_paid') {
            await appointmentApi.markPaid(appointment_id);
        }
        else if (action === 'reschedule') {
            await appointmentApi.reschedule(appointment_id, { new_slot_id, reason });
        }
        message.success('操作成功');
    }
    finally {
        loading.value = false;
    }
}
debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
const __VLS_ctx = {};
let __VLS_components;
let __VLS_directives;
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ style: {} },
});
const __VLS_0 = {}.ACard;
/** @type {[typeof __VLS_components.ACard, typeof __VLS_components.aCard, typeof __VLS_components.ACard, typeof __VLS_components.aCard, ]} */ ;
// @ts-ignore
const __VLS_1 = __VLS_asFunctionalComponent(__VLS_0, new __VLS_0({
    title: "人工干预",
    size: "small",
}));
const __VLS_2 = __VLS_1({
    title: "人工干预",
    size: "small",
}, ...__VLS_functionalComponentArgsRest(__VLS_1));
__VLS_3.slots.default;
const __VLS_4 = {}.AForm;
/** @type {[typeof __VLS_components.AForm, typeof __VLS_components.aForm, typeof __VLS_components.AForm, typeof __VLS_components.aForm, ]} */ ;
// @ts-ignore
const __VLS_5 = __VLS_asFunctionalComponent(__VLS_4, new __VLS_4({
    model: (__VLS_ctx.form),
    layout: "vertical",
}));
const __VLS_6 = __VLS_5({
    model: (__VLS_ctx.form),
    layout: "vertical",
}, ...__VLS_functionalComponentArgsRest(__VLS_5));
__VLS_7.slots.default;
const __VLS_8 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_9 = __VLS_asFunctionalComponent(__VLS_8, new __VLS_8({
    label: "预约号",
}));
const __VLS_10 = __VLS_9({
    label: "预约号",
}, ...__VLS_functionalComponentArgsRest(__VLS_9));
__VLS_11.slots.default;
const __VLS_12 = {}.AInputSearch;
/** @type {[typeof __VLS_components.AInputSearch, typeof __VLS_components.aInputSearch, ]} */ ;
// @ts-ignore
const __VLS_13 = __VLS_asFunctionalComponent(__VLS_12, new __VLS_12({
    ...{ 'onSearch': {} },
    value: (__VLS_ctx.form.appointment_id),
    placeholder: "输入预约ID",
    enterButton: "查询",
    loading: (__VLS_ctx.searching),
}));
const __VLS_14 = __VLS_13({
    ...{ 'onSearch': {} },
    value: (__VLS_ctx.form.appointment_id),
    placeholder: "输入预约ID",
    enterButton: "查询",
    loading: (__VLS_ctx.searching),
}, ...__VLS_functionalComponentArgsRest(__VLS_13));
let __VLS_16;
let __VLS_17;
let __VLS_18;
const __VLS_19 = {
    onSearch: (__VLS_ctx.lookupAppointment)
};
var __VLS_15;
var __VLS_11;
if (__VLS_ctx.appointment) {
    const __VLS_20 = {}.ADescriptions;
    /** @type {[typeof __VLS_components.ADescriptions, typeof __VLS_components.aDescriptions, typeof __VLS_components.ADescriptions, typeof __VLS_components.aDescriptions, ]} */ ;
    // @ts-ignore
    const __VLS_21 = __VLS_asFunctionalComponent(__VLS_20, new __VLS_20({
        column: (2),
        size: "small",
        bordered: true,
        ...{ style: {} },
    }));
    const __VLS_22 = __VLS_21({
        column: (2),
        size: "small",
        bordered: true,
        ...{ style: {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_21));
    __VLS_23.slots.default;
    const __VLS_24 = {}.ADescriptionsItem;
    /** @type {[typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_25 = __VLS_asFunctionalComponent(__VLS_24, new __VLS_24({
        label: "患者",
    }));
    const __VLS_26 = __VLS_25({
        label: "患者",
    }, ...__VLS_functionalComponentArgsRest(__VLS_25));
    __VLS_27.slots.default;
    (__VLS_ctx.appointment.patient_name);
    var __VLS_27;
    const __VLS_28 = {}.ADescriptionsItem;
    /** @type {[typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_29 = __VLS_asFunctionalComponent(__VLS_28, new __VLS_28({
        label: "状态",
    }));
    const __VLS_30 = __VLS_29({
        label: "状态",
    }, ...__VLS_functionalComponentArgsRest(__VLS_29));
    __VLS_31.slots.default;
    (__VLS_ctx.appointment.status);
    var __VLS_31;
    const __VLS_32 = {}.ADescriptionsItem;
    /** @type {[typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_33 = __VLS_asFunctionalComponent(__VLS_32, new __VLS_32({
        label: "预约号",
    }));
    const __VLS_34 = __VLS_33({
        label: "预约号",
    }, ...__VLS_functionalComponentArgsRest(__VLS_33));
    __VLS_35.slots.default;
    (__VLS_ctx.appointment.appointment_no);
    var __VLS_35;
    const __VLS_36 = {}.ADescriptionsItem;
    /** @type {[typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_37 = __VLS_asFunctionalComponent(__VLS_36, new __VLS_36({
        label: "来源",
    }));
    const __VLS_38 = __VLS_37({
        label: "来源",
    }, ...__VLS_functionalComponentArgsRest(__VLS_37));
    __VLS_39.slots.default;
    (__VLS_ctx.appointment.source);
    var __VLS_39;
    var __VLS_23;
    const __VLS_40 = {}.AFormItem;
    /** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
    // @ts-ignore
    const __VLS_41 = __VLS_asFunctionalComponent(__VLS_40, new __VLS_40({
        label: "操作类型",
    }));
    const __VLS_42 = __VLS_41({
        label: "操作类型",
    }, ...__VLS_functionalComponentArgsRest(__VLS_41));
    __VLS_43.slots.default;
    const __VLS_44 = {}.ARadioGroup;
    /** @type {[typeof __VLS_components.ARadioGroup, typeof __VLS_components.aRadioGroup, typeof __VLS_components.ARadioGroup, typeof __VLS_components.aRadioGroup, ]} */ ;
    // @ts-ignore
    const __VLS_45 = __VLS_asFunctionalComponent(__VLS_44, new __VLS_44({
        value: (__VLS_ctx.form.action),
        buttonStyle: "solid",
    }));
    const __VLS_46 = __VLS_45({
        value: (__VLS_ctx.form.action),
        buttonStyle: "solid",
    }, ...__VLS_functionalComponentArgsRest(__VLS_45));
    __VLS_47.slots.default;
    const __VLS_48 = {}.ARadioButton;
    /** @type {[typeof __VLS_components.ARadioButton, typeof __VLS_components.aRadioButton, typeof __VLS_components.ARadioButton, typeof __VLS_components.aRadioButton, ]} */ ;
    // @ts-ignore
    const __VLS_49 = __VLS_asFunctionalComponent(__VLS_48, new __VLS_48({
        value: "confirm",
    }));
    const __VLS_50 = __VLS_49({
        value: "confirm",
    }, ...__VLS_functionalComponentArgsRest(__VLS_49));
    __VLS_51.slots.default;
    var __VLS_51;
    const __VLS_52 = {}.ARadioButton;
    /** @type {[typeof __VLS_components.ARadioButton, typeof __VLS_components.aRadioButton, typeof __VLS_components.ARadioButton, typeof __VLS_components.aRadioButton, ]} */ ;
    // @ts-ignore
    const __VLS_53 = __VLS_asFunctionalComponent(__VLS_52, new __VLS_52({
        value: "mark_paid",
    }));
    const __VLS_54 = __VLS_53({
        value: "mark_paid",
    }, ...__VLS_functionalComponentArgsRest(__VLS_53));
    __VLS_55.slots.default;
    var __VLS_55;
    const __VLS_56 = {}.ARadioButton;
    /** @type {[typeof __VLS_components.ARadioButton, typeof __VLS_components.aRadioButton, typeof __VLS_components.ARadioButton, typeof __VLS_components.aRadioButton, ]} */ ;
    // @ts-ignore
    const __VLS_57 = __VLS_asFunctionalComponent(__VLS_56, new __VLS_56({
        value: "reschedule",
    }));
    const __VLS_58 = __VLS_57({
        value: "reschedule",
    }, ...__VLS_functionalComponentArgsRest(__VLS_57));
    __VLS_59.slots.default;
    var __VLS_59;
    const __VLS_60 = {}.ARadioButton;
    /** @type {[typeof __VLS_components.ARadioButton, typeof __VLS_components.aRadioButton, typeof __VLS_components.ARadioButton, typeof __VLS_components.aRadioButton, ]} */ ;
    // @ts-ignore
    const __VLS_61 = __VLS_asFunctionalComponent(__VLS_60, new __VLS_60({
        value: "cancel",
    }));
    const __VLS_62 = __VLS_61({
        value: "cancel",
    }, ...__VLS_functionalComponentArgsRest(__VLS_61));
    __VLS_63.slots.default;
    var __VLS_63;
    var __VLS_47;
    var __VLS_43;
    if (__VLS_ctx.form.action === 'reschedule') {
        const __VLS_64 = {}.AFormItem;
        /** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
        // @ts-ignore
        const __VLS_65 = __VLS_asFunctionalComponent(__VLS_64, new __VLS_64({
            label: "新号源ID",
        }));
        const __VLS_66 = __VLS_65({
            label: "新号源ID",
        }, ...__VLS_functionalComponentArgsRest(__VLS_65));
        __VLS_67.slots.default;
        const __VLS_68 = {}.AInput;
        /** @type {[typeof __VLS_components.AInput, typeof __VLS_components.aInput, ]} */ ;
        // @ts-ignore
        const __VLS_69 = __VLS_asFunctionalComponent(__VLS_68, new __VLS_68({
            value: (__VLS_ctx.form.new_slot_id),
            placeholder: "输入新的号源ID",
        }));
        const __VLS_70 = __VLS_69({
            value: (__VLS_ctx.form.new_slot_id),
            placeholder: "输入新的号源ID",
        }, ...__VLS_functionalComponentArgsRest(__VLS_69));
        var __VLS_67;
    }
    if (['cancel', 'reschedule'].includes(__VLS_ctx.form.action)) {
        const __VLS_72 = {}.AFormItem;
        /** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
        // @ts-ignore
        const __VLS_73 = __VLS_asFunctionalComponent(__VLS_72, new __VLS_72({
            label: "操作原因",
        }));
        const __VLS_74 = __VLS_73({
            label: "操作原因",
        }, ...__VLS_functionalComponentArgsRest(__VLS_73));
        __VLS_75.slots.default;
        const __VLS_76 = {}.ATextarea;
        /** @type {[typeof __VLS_components.ATextarea, typeof __VLS_components.aTextarea, ]} */ ;
        // @ts-ignore
        const __VLS_77 = __VLS_asFunctionalComponent(__VLS_76, new __VLS_76({
            value: (__VLS_ctx.form.reason),
            rows: (3),
        }));
        const __VLS_78 = __VLS_77({
            value: (__VLS_ctx.form.reason),
            rows: (3),
        }, ...__VLS_functionalComponentArgsRest(__VLS_77));
        var __VLS_75;
    }
    const __VLS_80 = {}.AButton;
    /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
    // @ts-ignore
    const __VLS_81 = __VLS_asFunctionalComponent(__VLS_80, new __VLS_80({
        ...{ 'onClick': {} },
        type: "primary",
        loading: (__VLS_ctx.loading),
    }));
    const __VLS_82 = __VLS_81({
        ...{ 'onClick': {} },
        type: "primary",
        loading: (__VLS_ctx.loading),
    }, ...__VLS_functionalComponentArgsRest(__VLS_81));
    let __VLS_84;
    let __VLS_85;
    let __VLS_86;
    const __VLS_87 = {
        onClick: (__VLS_ctx.handleSubmit)
    };
    __VLS_83.slots.default;
    var __VLS_83;
}
else {
    const __VLS_88 = {}.AEmpty;
    /** @type {[typeof __VLS_components.AEmpty, typeof __VLS_components.aEmpty, ]} */ ;
    // @ts-ignore
    const __VLS_89 = __VLS_asFunctionalComponent(__VLS_88, new __VLS_88({
        description: "请先查询预约",
        imageStyle: ({ height: '60px' }),
    }));
    const __VLS_90 = __VLS_89({
        description: "请先查询预约",
        imageStyle: ({ height: '60px' }),
    }, ...__VLS_functionalComponentArgsRest(__VLS_89));
}
var __VLS_7;
var __VLS_3;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            form: form,
            loading: loading,
            appointment: appointment,
            searching: searching,
            lookupAppointment: lookupAppointment,
            handleSubmit: handleSubmit,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
