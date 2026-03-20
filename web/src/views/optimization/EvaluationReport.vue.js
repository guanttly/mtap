/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { optimizationApi } from '@/api/optimization';
const route = useRoute();
const id = route.params.id;
const loading = ref(false);
const report = ref(null);
async function fetchData() {
    loading.value = true;
    try {
        report.value = await optimizationApi.getEvaluation(id);
    }
    finally {
        loading.value = false;
    }
}
onMounted(fetchData);
const baselineColumns = [
    { title: '指标', dataIndex: 'key' },
    {
        title: '数值',
        dataIndex: 'value',
        customRender: ({ text }) => text?.toFixed?.(2) ?? text,
    },
];
const trialColumns = [
    { title: '指标', dataIndex: 'key' },
    {
        title: '数值',
        dataIndex: 'value',
        customRender: ({ text }) => text?.toFixed?.(2) ?? text,
    },
    {
        title: '变化%',
        dataIndex: 'change',
        customRender: ({ text }) => `${text >= 0 ? '+' : ''}${(text * 100).toFixed(1)}%`,
    },
];
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
        title: "评估报告",
        size: "small",
    }));
    const __VLS_7 = __VLS_6({
        title: "评估报告",
        size: "small",
    }, ...__VLS_functionalComponentArgsRest(__VLS_6));
    __VLS_8.slots.default;
    const __VLS_9 = {}.ADescriptions;
    /** @type {[typeof __VLS_components.ADescriptions, typeof __VLS_components.aDescriptions, typeof __VLS_components.ADescriptions, typeof __VLS_components.aDescriptions, ]} */ ;
    // @ts-ignore
    const __VLS_10 = __VLS_asFunctionalComponent(__VLS_9, new __VLS_9({
        column: (2),
        bordered: true,
        size: "small",
        ...{ style: {} },
    }));
    const __VLS_11 = __VLS_10({
        column: (2),
        bordered: true,
        size: "small",
        ...{ style: {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_10));
    __VLS_12.slots.default;
    const __VLS_13 = {}.ADescriptionsItem;
    /** @type {[typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_14 = __VLS_asFunctionalComponent(__VLS_13, new __VLS_13({
        label: "策略ID",
    }));
    const __VLS_15 = __VLS_14({
        label: "策略ID",
    }, ...__VLS_functionalComponentArgsRest(__VLS_14));
    __VLS_16.slots.default;
    (__VLS_ctx.report.strategy_id);
    var __VLS_16;
    const __VLS_17 = {}.ADescriptionsItem;
    /** @type {[typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_18 = __VLS_asFunctionalComponent(__VLS_17, new __VLS_17({
        label: "试运行ID",
    }));
    const __VLS_19 = __VLS_18({
        label: "试运行ID",
    }, ...__VLS_functionalComponentArgsRest(__VLS_18));
    __VLS_20.slots.default;
    (__VLS_ctx.report.trial_run_id);
    var __VLS_20;
    const __VLS_21 = {}.ADescriptionsItem;
    /** @type {[typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_22 = __VLS_asFunctionalComponent(__VLS_21, new __VLS_21({
        label: "生成时间",
        span: (2),
    }));
    const __VLS_23 = __VLS_22({
        label: "生成时间",
        span: (2),
    }, ...__VLS_functionalComponentArgsRest(__VLS_22));
    __VLS_24.slots.default;
    (__VLS_ctx.report.generated_at);
    var __VLS_24;
    const __VLS_25 = {}.ADescriptionsItem;
    /** @type {[typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_26 = __VLS_asFunctionalComponent(__VLS_25, new __VLS_25({
        label: "是否合格",
    }));
    const __VLS_27 = __VLS_26({
        label: "是否合格",
    }, ...__VLS_functionalComponentArgsRest(__VLS_26));
    __VLS_28.slots.default;
    const __VLS_29 = {}.ATag;
    /** @type {[typeof __VLS_components.ATag, typeof __VLS_components.aTag, typeof __VLS_components.ATag, typeof __VLS_components.aTag, ]} */ ;
    // @ts-ignore
    const __VLS_30 = __VLS_asFunctionalComponent(__VLS_29, new __VLS_29({
        color: (__VLS_ctx.report.is_qualified ? 'success' : 'error'),
    }));
    const __VLS_31 = __VLS_30({
        color: (__VLS_ctx.report.is_qualified ? 'success' : 'error'),
    }, ...__VLS_functionalComponentArgsRest(__VLS_30));
    __VLS_32.slots.default;
    (__VLS_ctx.report.is_qualified ? '合格' : '不合格');
    var __VLS_32;
    var __VLS_28;
    const __VLS_33 = {}.ADescriptionsItem;
    /** @type {[typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_34 = __VLS_asFunctionalComponent(__VLS_33, new __VLS_33({
        label: "合格阈值",
    }));
    const __VLS_35 = __VLS_34({
        label: "合格阈值",
    }, ...__VLS_functionalComponentArgsRest(__VLS_34));
    __VLS_36.slots.default;
    ((__VLS_ctx.report.qualify_threshold * 100).toFixed(1));
    var __VLS_36;
    const __VLS_37 = {}.ADescriptionsItem;
    /** @type {[typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_38 = __VLS_asFunctionalComponent(__VLS_37, new __VLS_37({
        label: "推荐操作",
        span: (2),
    }));
    const __VLS_39 = __VLS_38({
        label: "推荐操作",
        span: (2),
    }, ...__VLS_functionalComponentArgsRest(__VLS_38));
    __VLS_40.slots.default;
    const __VLS_41 = {}.ATag;
    /** @type {[typeof __VLS_components.ATag, typeof __VLS_components.aTag, typeof __VLS_components.ATag, typeof __VLS_components.aTag, ]} */ ;
    // @ts-ignore
    const __VLS_42 = __VLS_asFunctionalComponent(__VLS_41, new __VLS_41({
        color: (__VLS_ctx.report.recommendation === 'promote' ? 'success' : __VLS_ctx.report.recommendation === 'rollback' ? 'error' : 'warning'),
    }));
    const __VLS_43 = __VLS_42({
        color: (__VLS_ctx.report.recommendation === 'promote' ? 'success' : __VLS_ctx.report.recommendation === 'rollback' ? 'error' : 'warning'),
    }, ...__VLS_functionalComponentArgsRest(__VLS_42));
    __VLS_44.slots.default;
    (__VLS_ctx.report.recommendation);
    var __VLS_44;
    var __VLS_40;
    var __VLS_12;
    const __VLS_45 = {}.ARow;
    /** @type {[typeof __VLS_components.ARow, typeof __VLS_components.aRow, typeof __VLS_components.ARow, typeof __VLS_components.aRow, ]} */ ;
    // @ts-ignore
    const __VLS_46 = __VLS_asFunctionalComponent(__VLS_45, new __VLS_45({
        gutter: (16),
    }));
    const __VLS_47 = __VLS_46({
        gutter: (16),
    }, ...__VLS_functionalComponentArgsRest(__VLS_46));
    __VLS_48.slots.default;
    const __VLS_49 = {}.ACol;
    /** @type {[typeof __VLS_components.ACol, typeof __VLS_components.aCol, typeof __VLS_components.ACol, typeof __VLS_components.aCol, ]} */ ;
    // @ts-ignore
    const __VLS_50 = __VLS_asFunctionalComponent(__VLS_49, new __VLS_49({
        span: (12),
    }));
    const __VLS_51 = __VLS_50({
        span: (12),
    }, ...__VLS_functionalComponentArgsRest(__VLS_50));
    __VLS_52.slots.default;
    const __VLS_53 = {}.ACard;
    /** @type {[typeof __VLS_components.ACard, typeof __VLS_components.aCard, typeof __VLS_components.ACard, typeof __VLS_components.aCard, ]} */ ;
    // @ts-ignore
    const __VLS_54 = __VLS_asFunctionalComponent(__VLS_53, new __VLS_53({
        title: "基准指标",
        size: "small",
    }));
    const __VLS_55 = __VLS_54({
        title: "基准指标",
        size: "small",
    }, ...__VLS_functionalComponentArgsRest(__VLS_54));
    __VLS_56.slots.default;
    const __VLS_57 = {}.ATable;
    /** @type {[typeof __VLS_components.ATable, typeof __VLS_components.aTable, ]} */ ;
    // @ts-ignore
    const __VLS_58 = __VLS_asFunctionalComponent(__VLS_57, new __VLS_57({
        dataSource: (Object.entries(__VLS_ctx.report.baseline_metrics).map(([k, v]) => ({ key: k, value: v }))),
        columns: (__VLS_ctx.baselineColumns),
        pagination: (false),
        size: "small",
    }));
    const __VLS_59 = __VLS_58({
        dataSource: (Object.entries(__VLS_ctx.report.baseline_metrics).map(([k, v]) => ({ key: k, value: v }))),
        columns: (__VLS_ctx.baselineColumns),
        pagination: (false),
        size: "small",
    }, ...__VLS_functionalComponentArgsRest(__VLS_58));
    var __VLS_56;
    var __VLS_52;
    const __VLS_61 = {}.ACol;
    /** @type {[typeof __VLS_components.ACol, typeof __VLS_components.aCol, typeof __VLS_components.ACol, typeof __VLS_components.aCol, ]} */ ;
    // @ts-ignore
    const __VLS_62 = __VLS_asFunctionalComponent(__VLS_61, new __VLS_61({
        span: (12),
    }));
    const __VLS_63 = __VLS_62({
        span: (12),
    }, ...__VLS_functionalComponentArgsRest(__VLS_62));
    __VLS_64.slots.default;
    const __VLS_65 = {}.ACard;
    /** @type {[typeof __VLS_components.ACard, typeof __VLS_components.aCard, typeof __VLS_components.ACard, typeof __VLS_components.aCard, ]} */ ;
    // @ts-ignore
    const __VLS_66 = __VLS_asFunctionalComponent(__VLS_65, new __VLS_65({
        title: "试运行指标",
        size: "small",
    }));
    const __VLS_67 = __VLS_66({
        title: "试运行指标",
        size: "small",
    }, ...__VLS_functionalComponentArgsRest(__VLS_66));
    __VLS_68.slots.default;
    const __VLS_69 = {}.ATable;
    /** @type {[typeof __VLS_components.ATable, typeof __VLS_components.aTable, ]} */ ;
    // @ts-ignore
    const __VLS_70 = __VLS_asFunctionalComponent(__VLS_69, new __VLS_69({
        dataSource: (Object.entries(__VLS_ctx.report.trial_metrics).map(([k, v]) => ({ key: k, value: v, change: __VLS_ctx.report.change_pct[k] }))),
        columns: (__VLS_ctx.trialColumns),
        pagination: (false),
        size: "small",
    }));
    const __VLS_71 = __VLS_70({
        dataSource: (Object.entries(__VLS_ctx.report.trial_metrics).map(([k, v]) => ({ key: k, value: v, change: __VLS_ctx.report.change_pct[k] }))),
        columns: (__VLS_ctx.trialColumns),
        pagination: (false),
        size: "small",
    }, ...__VLS_functionalComponentArgsRest(__VLS_70));
    var __VLS_68;
    var __VLS_64;
    var __VLS_48;
    var __VLS_8;
}
else if (!__VLS_ctx.loading) {
    const __VLS_73 = {}.AEmpty;
    /** @type {[typeof __VLS_components.AEmpty, typeof __VLS_components.aEmpty, ]} */ ;
    // @ts-ignore
    const __VLS_74 = __VLS_asFunctionalComponent(__VLS_73, new __VLS_73({}));
    const __VLS_75 = __VLS_74({}, ...__VLS_functionalComponentArgsRest(__VLS_74));
}
var __VLS_3;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            loading: loading,
            report: report,
            baselineColumns: baselineColumns,
            trialColumns: trialColumns,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
