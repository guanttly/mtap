/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { ref, onMounted, computed } from 'vue';
import { useRoute } from 'vue-router';
import { message } from 'ant-design-vue';
import { optimizationApi } from '@/api/optimization';
const route = useRoute();
const id = route.params.id;
const loading = ref(false);
const report = ref(null);
const submitting = ref(false);
const reviewForm = ref({ approved: true, comment: '' });
// 计算 ROI = (年收益 - 总投入) / 总投入
const roi = computed(() => {
    if (!report.value)
        return 0;
    const { total_investment, expected_annual_revenue } = report.value;
    if (!total_investment)
        return 0;
    return (expected_annual_revenue - total_investment) / total_investment;
});
async function fetchData() {
    loading.value = true;
    try {
        report.value = await optimizationApi.getROIReport(id);
    }
    finally {
        loading.value = false;
    }
}
async function submitReview() {
    submitting.value = true;
    try {
        await optimizationApi.submitROIResult(id, reviewForm.value);
        message.success('ROI 审核结果已提交');
        fetchData();
    }
    finally {
        submitting.value = false;
    }
}
onMounted(fetchData);
debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
const __VLS_ctx = {};
let __VLS_components;
let __VLS_directives;
const __VLS_0 = {}.ASpin;
/** @type {[typeof __VLS_components.ASpin, typeof __VLS_components.aSpin, typeof __VLS_components.ASpin, typeof __VLS_components.aSpin, ]} */ ;
// @ts-ignore
const __VLS_1 = __VLS_asFunctionalComponent(__VLS_0, new __VLS_0({
    spinning: (__VLS_ctx.loading),
}));
const __VLS_2 = __VLS_1({
    spinning: (__VLS_ctx.loading),
}, ...__VLS_functionalComponentArgsRest(__VLS_1));
var __VLS_4 = {};
__VLS_3.slots.default;
if (__VLS_ctx.report) {
    const __VLS_5 = {}.ACard;
    /** @type {[typeof __VLS_components.ACard, typeof __VLS_components.aCard, typeof __VLS_components.ACard, typeof __VLS_components.aCard, ]} */ ;
    // @ts-ignore
    const __VLS_6 = __VLS_asFunctionalComponent(__VLS_5, new __VLS_5({
        title: "ROI 分析报告",
        size: "small",
    }));
    const __VLS_7 = __VLS_6({
        title: "ROI 分析报告",
        size: "small",
    }, ...__VLS_functionalComponentArgsRest(__VLS_6));
    __VLS_8.slots.default;
    const __VLS_9 = {}.ADescriptions;
    /** @type {[typeof __VLS_components.ADescriptions, typeof __VLS_components.aDescriptions, typeof __VLS_components.ADescriptions, typeof __VLS_components.aDescriptions, ]} */ ;
    // @ts-ignore
    const __VLS_10 = __VLS_asFunctionalComponent(__VLS_9, new __VLS_9({
        column: (1),
        bordered: true,
        size: "small",
        ...{ style: {} },
    }));
    const __VLS_11 = __VLS_10({
        column: (1),
        bordered: true,
        size: "small",
        ...{ style: {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_10));
    __VLS_12.slots.default;
    const __VLS_13 = {}.ADescriptionsItem;
    /** @type {[typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_14 = __VLS_asFunctionalComponent(__VLS_13, new __VLS_13({
        label: "当前瓶颈",
    }));
    const __VLS_15 = __VLS_14({
        label: "当前瓶颈",
    }, ...__VLS_functionalComponentArgsRest(__VLS_14));
    __VLS_16.slots.default;
    (__VLS_ctx.report.current_bottleneck);
    var __VLS_16;
    var __VLS_12;
    const __VLS_17 = {}.ARow;
    /** @type {[typeof __VLS_components.ARow, typeof __VLS_components.aRow, typeof __VLS_components.ARow, typeof __VLS_components.aRow, ]} */ ;
    // @ts-ignore
    const __VLS_18 = __VLS_asFunctionalComponent(__VLS_17, new __VLS_17({
        gutter: (16),
        ...{ style: {} },
    }));
    const __VLS_19 = __VLS_18({
        gutter: (16),
        ...{ style: {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_18));
    __VLS_20.slots.default;
    const __VLS_21 = {}.ACol;
    /** @type {[typeof __VLS_components.ACol, typeof __VLS_components.aCol, typeof __VLS_components.ACol, typeof __VLS_components.aCol, ]} */ ;
    // @ts-ignore
    const __VLS_22 = __VLS_asFunctionalComponent(__VLS_21, new __VLS_21({
        span: (6),
    }));
    const __VLS_23 = __VLS_22({
        span: (6),
    }, ...__VLS_functionalComponentArgsRest(__VLS_22));
    __VLS_24.slots.default;
    const __VLS_25 = {}.AStatistic;
    /** @type {[typeof __VLS_components.AStatistic, typeof __VLS_components.aStatistic, ]} */ ;
    // @ts-ignore
    const __VLS_26 = __VLS_asFunctionalComponent(__VLS_25, new __VLS_25({
        title: "总投入(元)",
        value: (__VLS_ctx.report.total_investment),
        precision: (2),
    }));
    const __VLS_27 = __VLS_26({
        title: "总投入(元)",
        value: (__VLS_ctx.report.total_investment),
        precision: (2),
    }, ...__VLS_functionalComponentArgsRest(__VLS_26));
    var __VLS_24;
    const __VLS_29 = {}.ACol;
    /** @type {[typeof __VLS_components.ACol, typeof __VLS_components.aCol, typeof __VLS_components.ACol, typeof __VLS_components.aCol, ]} */ ;
    // @ts-ignore
    const __VLS_30 = __VLS_asFunctionalComponent(__VLS_29, new __VLS_29({
        span: (6),
    }));
    const __VLS_31 = __VLS_30({
        span: (6),
    }, ...__VLS_functionalComponentArgsRest(__VLS_30));
    __VLS_32.slots.default;
    const __VLS_33 = {}.AStatistic;
    /** @type {[typeof __VLS_components.AStatistic, typeof __VLS_components.aStatistic, ]} */ ;
    // @ts-ignore
    const __VLS_34 = __VLS_asFunctionalComponent(__VLS_33, new __VLS_33({
        title: "预期年收益(元)",
        value: (__VLS_ctx.report.expected_annual_revenue),
        precision: (2),
        valueStyle: ({ color: '#52c41a' }),
    }));
    const __VLS_35 = __VLS_34({
        title: "预期年收益(元)",
        value: (__VLS_ctx.report.expected_annual_revenue),
        precision: (2),
        valueStyle: ({ color: '#52c41a' }),
    }, ...__VLS_functionalComponentArgsRest(__VLS_34));
    var __VLS_32;
    const __VLS_37 = {}.ACol;
    /** @type {[typeof __VLS_components.ACol, typeof __VLS_components.aCol, typeof __VLS_components.ACol, typeof __VLS_components.aCol, ]} */ ;
    // @ts-ignore
    const __VLS_38 = __VLS_asFunctionalComponent(__VLS_37, new __VLS_37({
        span: (6),
    }));
    const __VLS_39 = __VLS_38({
        span: (6),
    }, ...__VLS_functionalComponentArgsRest(__VLS_38));
    __VLS_40.slots.default;
    const __VLS_41 = {}.AStatistic;
    /** @type {[typeof __VLS_components.AStatistic, typeof __VLS_components.aStatistic, ]} */ ;
    // @ts-ignore
    const __VLS_42 = __VLS_asFunctionalComponent(__VLS_41, new __VLS_41({
        title: "ROI",
        value: ((__VLS_ctx.roi * 100).toFixed(1)),
        suffix: "%",
        valueStyle: ({ color: __VLS_ctx.roi > 0 ? '#52c41a' : '#ff4d4f' }),
    }));
    const __VLS_43 = __VLS_42({
        title: "ROI",
        value: ((__VLS_ctx.roi * 100).toFixed(1)),
        suffix: "%",
        valueStyle: ({ color: __VLS_ctx.roi > 0 ? '#52c41a' : '#ff4d4f' }),
    }, ...__VLS_functionalComponentArgsRest(__VLS_42));
    var __VLS_40;
    const __VLS_45 = {}.ACol;
    /** @type {[typeof __VLS_components.ACol, typeof __VLS_components.aCol, typeof __VLS_components.ACol, typeof __VLS_components.aCol, ]} */ ;
    // @ts-ignore
    const __VLS_46 = __VLS_asFunctionalComponent(__VLS_45, new __VLS_45({
        span: (6),
    }));
    const __VLS_47 = __VLS_46({
        span: (6),
    }, ...__VLS_functionalComponentArgsRest(__VLS_46));
    __VLS_48.slots.default;
    const __VLS_49 = {}.AStatistic;
    /** @type {[typeof __VLS_components.AStatistic, typeof __VLS_components.aStatistic, ]} */ ;
    // @ts-ignore
    const __VLS_50 = __VLS_asFunctionalComponent(__VLS_49, new __VLS_49({
        title: "回本周期(月)",
        value: (__VLS_ctx.report.payback_period_months.toFixed(1)),
    }));
    const __VLS_51 = __VLS_50({
        title: "回本周期(月)",
        value: (__VLS_ctx.report.payback_period_months.toFixed(1)),
    }, ...__VLS_functionalComponentArgsRest(__VLS_50));
    var __VLS_48;
    var __VLS_20;
    const __VLS_53 = {}.ACard;
    /** @type {[typeof __VLS_components.ACard, typeof __VLS_components.aCard, typeof __VLS_components.ACard, typeof __VLS_components.aCard, ]} */ ;
    // @ts-ignore
    const __VLS_54 = __VLS_asFunctionalComponent(__VLS_53, new __VLS_53({
        title: "风险因素",
        size: "small",
        ...{ style: {} },
    }));
    const __VLS_55 = __VLS_54({
        title: "风险因素",
        size: "small",
        ...{ style: {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_54));
    __VLS_56.slots.default;
    for (const [r] of __VLS_getVForSourceType((__VLS_ctx.report.risk_factors))) {
        const __VLS_57 = {}.ATag;
        /** @type {[typeof __VLS_components.ATag, typeof __VLS_components.aTag, typeof __VLS_components.ATag, typeof __VLS_components.aTag, ]} */ ;
        // @ts-ignore
        const __VLS_58 = __VLS_asFunctionalComponent(__VLS_57, new __VLS_57({
            key: (r),
            color: "orange",
            ...{ style: {} },
        }));
        const __VLS_59 = __VLS_58({
            key: (r),
            color: "orange",
            ...{ style: {} },
        }, ...__VLS_functionalComponentArgsRest(__VLS_58));
        __VLS_60.slots.default;
        (r);
        var __VLS_60;
    }
    if (!__VLS_ctx.report.risk_factors?.length) {
        __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
            ...{ class: "text-muted" },
        });
    }
    var __VLS_56;
    if (!__VLS_ctx.report.approval_result) {
        const __VLS_61 = {}.ACard;
        /** @type {[typeof __VLS_components.ACard, typeof __VLS_components.aCard, typeof __VLS_components.ACard, typeof __VLS_components.aCard, ]} */ ;
        // @ts-ignore
        const __VLS_62 = __VLS_asFunctionalComponent(__VLS_61, new __VLS_61({
            title: "ROI 审核",
            size: "small",
        }));
        const __VLS_63 = __VLS_62({
            title: "ROI 审核",
            size: "small",
        }, ...__VLS_functionalComponentArgsRest(__VLS_62));
        __VLS_64.slots.default;
        const __VLS_65 = {}.AForm;
        /** @type {[typeof __VLS_components.AForm, typeof __VLS_components.aForm, typeof __VLS_components.AForm, typeof __VLS_components.aForm, ]} */ ;
        // @ts-ignore
        const __VLS_66 = __VLS_asFunctionalComponent(__VLS_65, new __VLS_65({
            model: (__VLS_ctx.reviewForm),
            layout: "vertical",
        }));
        const __VLS_67 = __VLS_66({
            model: (__VLS_ctx.reviewForm),
            layout: "vertical",
        }, ...__VLS_functionalComponentArgsRest(__VLS_66));
        __VLS_68.slots.default;
        const __VLS_69 = {}.AFormItem;
        /** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
        // @ts-ignore
        const __VLS_70 = __VLS_asFunctionalComponent(__VLS_69, new __VLS_69({
            label: "审核结论",
        }));
        const __VLS_71 = __VLS_70({
            label: "审核结论",
        }, ...__VLS_functionalComponentArgsRest(__VLS_70));
        __VLS_72.slots.default;
        const __VLS_73 = {}.ARadioGroup;
        /** @type {[typeof __VLS_components.ARadioGroup, typeof __VLS_components.aRadioGroup, typeof __VLS_components.ARadioGroup, typeof __VLS_components.aRadioGroup, ]} */ ;
        // @ts-ignore
        const __VLS_74 = __VLS_asFunctionalComponent(__VLS_73, new __VLS_73({
            value: (__VLS_ctx.reviewForm.approved),
        }));
        const __VLS_75 = __VLS_74({
            value: (__VLS_ctx.reviewForm.approved),
        }, ...__VLS_functionalComponentArgsRest(__VLS_74));
        __VLS_76.slots.default;
        const __VLS_77 = {}.ARadio;
        /** @type {[typeof __VLS_components.ARadio, typeof __VLS_components.aRadio, typeof __VLS_components.ARadio, typeof __VLS_components.aRadio, ]} */ ;
        // @ts-ignore
        const __VLS_78 = __VLS_asFunctionalComponent(__VLS_77, new __VLS_77({
            value: (true),
        }));
        const __VLS_79 = __VLS_78({
            value: (true),
        }, ...__VLS_functionalComponentArgsRest(__VLS_78));
        __VLS_80.slots.default;
        var __VLS_80;
        const __VLS_81 = {}.ARadio;
        /** @type {[typeof __VLS_components.ARadio, typeof __VLS_components.aRadio, typeof __VLS_components.ARadio, typeof __VLS_components.aRadio, ]} */ ;
        // @ts-ignore
        const __VLS_82 = __VLS_asFunctionalComponent(__VLS_81, new __VLS_81({
            value: (false),
        }));
        const __VLS_83 = __VLS_82({
            value: (false),
        }, ...__VLS_functionalComponentArgsRest(__VLS_82));
        __VLS_84.slots.default;
        var __VLS_84;
        var __VLS_76;
        var __VLS_72;
        const __VLS_85 = {}.AFormItem;
        /** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
        // @ts-ignore
        const __VLS_86 = __VLS_asFunctionalComponent(__VLS_85, new __VLS_85({
            label: "审核意见",
        }));
        const __VLS_87 = __VLS_86({
            label: "审核意见",
        }, ...__VLS_functionalComponentArgsRest(__VLS_86));
        __VLS_88.slots.default;
        const __VLS_89 = {}.ATextarea;
        /** @type {[typeof __VLS_components.ATextarea, typeof __VLS_components.aTextarea, ]} */ ;
        // @ts-ignore
        const __VLS_90 = __VLS_asFunctionalComponent(__VLS_89, new __VLS_89({
            value: (__VLS_ctx.reviewForm.comment),
            rows: (4),
        }));
        const __VLS_91 = __VLS_90({
            value: (__VLS_ctx.reviewForm.comment),
            rows: (4),
        }, ...__VLS_functionalComponentArgsRest(__VLS_90));
        var __VLS_88;
        const __VLS_93 = {}.AButton;
        /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
        // @ts-ignore
        const __VLS_94 = __VLS_asFunctionalComponent(__VLS_93, new __VLS_93({
            ...{ 'onClick': {} },
            type: "primary",
            loading: (__VLS_ctx.submitting),
        }));
        const __VLS_95 = __VLS_94({
            ...{ 'onClick': {} },
            type: "primary",
            loading: (__VLS_ctx.submitting),
        }, ...__VLS_functionalComponentArgsRest(__VLS_94));
        let __VLS_97;
        let __VLS_98;
        let __VLS_99;
        const __VLS_100 = {
            onClick: (__VLS_ctx.submitReview)
        };
        __VLS_96.slots.default;
        var __VLS_96;
        var __VLS_68;
        var __VLS_64;
    }
    else {
        const __VLS_101 = {}.AAlert;
        /** @type {[typeof __VLS_components.AAlert, typeof __VLS_components.aAlert, ]} */ ;
        // @ts-ignore
        const __VLS_102 = __VLS_asFunctionalComponent(__VLS_101, new __VLS_101({
            message: (`审核结果: ${__VLS_ctx.report.approval_result}`),
            type: "success",
        }));
        const __VLS_103 = __VLS_102({
            message: (`审核结果: ${__VLS_ctx.report.approval_result}`),
            type: "success",
        }, ...__VLS_functionalComponentArgsRest(__VLS_102));
    }
    var __VLS_8;
}
else if (!__VLS_ctx.loading) {
    const __VLS_105 = {}.AEmpty;
    /** @type {[typeof __VLS_components.AEmpty, typeof __VLS_components.aEmpty, ]} */ ;
    // @ts-ignore
    const __VLS_106 = __VLS_asFunctionalComponent(__VLS_105, new __VLS_105({}));
    const __VLS_107 = __VLS_106({}, ...__VLS_functionalComponentArgsRest(__VLS_106));
}
var __VLS_3;
/** @type {__VLS_StyleScopedClasses['text-muted']} */ ;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            loading: loading,
            report: report,
            submitting: submitting,
            reviewForm: reviewForm,
            roi: roi,
            submitReview: submitReview,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
