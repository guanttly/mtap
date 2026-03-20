/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { message } from 'ant-design-vue';
import { optimizationApi } from '@/api/optimization';
const route = useRoute();
const router = useRouter();
const strategyId = route.params.id;
const loading = ref(false);
const strategy = ref(null);
async function fetchData() {
    loading.value = true;
    try {
        strategy.value = await optimizationApi.getStrategy(strategyId);
    }
    finally {
        loading.value = false;
    }
}
async function handleRollback() {
    await optimizationApi.rollbackStrategy(strategyId);
    message.success('已回滚');
    fetchData();
}
async function handlePromote() {
    await optimizationApi.promoteStrategy(strategyId);
    message.success('已推全量');
    fetchData();
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
if (__VLS_ctx.strategy) {
    const __VLS_5 = {}.APageHeader;
    /** @type {[typeof __VLS_components.APageHeader, typeof __VLS_components.aPageHeader, typeof __VLS_components.APageHeader, typeof __VLS_components.aPageHeader, ]} */ ;
    // @ts-ignore
    const __VLS_6 = __VLS_asFunctionalComponent(__VLS_5, new __VLS_5({
        ...{ 'onBack': {} },
        title: (__VLS_ctx.strategy.title),
    }));
    const __VLS_7 = __VLS_6({
        ...{ 'onBack': {} },
        title: (__VLS_ctx.strategy.title),
    }, ...__VLS_functionalComponentArgsRest(__VLS_6));
    let __VLS_9;
    let __VLS_10;
    let __VLS_11;
    const __VLS_12 = {
        onBack: (...[$event]) => {
            if (!(__VLS_ctx.strategy))
                return;
            __VLS_ctx.router.back();
        }
    };
    __VLS_8.slots.default;
    {
        const { tags: __VLS_thisSlot } = __VLS_8.slots;
        const __VLS_13 = {}.ATag;
        /** @type {[typeof __VLS_components.ATag, typeof __VLS_components.aTag, typeof __VLS_components.ATag, typeof __VLS_components.aTag, ]} */ ;
        // @ts-ignore
        const __VLS_14 = __VLS_asFunctionalComponent(__VLS_13, new __VLS_13({
            color: (__VLS_ctx.strategy.status === 'promoted' ? 'success' : __VLS_ctx.strategy.status === 'trial_running' ? 'processing' : 'default'),
        }));
        const __VLS_15 = __VLS_14({
            color: (__VLS_ctx.strategy.status === 'promoted' ? 'success' : __VLS_ctx.strategy.status === 'trial_running' ? 'processing' : 'default'),
        }, ...__VLS_functionalComponentArgsRest(__VLS_14));
        __VLS_16.slots.default;
        (__VLS_ctx.strategy.status);
        var __VLS_16;
    }
    {
        const { extra: __VLS_thisSlot } = __VLS_8.slots;
        if (__VLS_ctx.strategy.status === 'trial_running') {
            const __VLS_17 = {}.AButton;
            /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
            // @ts-ignore
            const __VLS_18 = __VLS_asFunctionalComponent(__VLS_17, new __VLS_17({
                ...{ 'onClick': {} },
            }));
            const __VLS_19 = __VLS_18({
                ...{ 'onClick': {} },
            }, ...__VLS_functionalComponentArgsRest(__VLS_18));
            let __VLS_21;
            let __VLS_22;
            let __VLS_23;
            const __VLS_24 = {
                onClick: (__VLS_ctx.handleRollback)
            };
            __VLS_20.slots.default;
            var __VLS_20;
        }
        if (__VLS_ctx.strategy.status === 'trial_running') {
            const __VLS_25 = {}.AButton;
            /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
            // @ts-ignore
            const __VLS_26 = __VLS_asFunctionalComponent(__VLS_25, new __VLS_25({
                ...{ 'onClick': {} },
                type: "primary",
            }));
            const __VLS_27 = __VLS_26({
                ...{ 'onClick': {} },
                type: "primary",
            }, ...__VLS_functionalComponentArgsRest(__VLS_26));
            let __VLS_29;
            let __VLS_30;
            let __VLS_31;
            const __VLS_32 = {
                onClick: (__VLS_ctx.handlePromote)
            };
            __VLS_28.slots.default;
            var __VLS_28;
        }
    }
    var __VLS_8;
    const __VLS_33 = {}.ARow;
    /** @type {[typeof __VLS_components.ARow, typeof __VLS_components.aRow, typeof __VLS_components.ARow, typeof __VLS_components.aRow, ]} */ ;
    // @ts-ignore
    const __VLS_34 = __VLS_asFunctionalComponent(__VLS_33, new __VLS_33({
        gutter: (16),
    }));
    const __VLS_35 = __VLS_34({
        gutter: (16),
    }, ...__VLS_functionalComponentArgsRest(__VLS_34));
    __VLS_36.slots.default;
    const __VLS_37 = {}.ACol;
    /** @type {[typeof __VLS_components.ACol, typeof __VLS_components.aCol, typeof __VLS_components.ACol, typeof __VLS_components.aCol, ]} */ ;
    // @ts-ignore
    const __VLS_38 = __VLS_asFunctionalComponent(__VLS_37, new __VLS_37({
        span: (12),
    }));
    const __VLS_39 = __VLS_38({
        span: (12),
    }, ...__VLS_functionalComponentArgsRest(__VLS_38));
    __VLS_40.slots.default;
    const __VLS_41 = {}.ACard;
    /** @type {[typeof __VLS_components.ACard, typeof __VLS_components.aCard, typeof __VLS_components.ACard, typeof __VLS_components.aCard, ]} */ ;
    // @ts-ignore
    const __VLS_42 = __VLS_asFunctionalComponent(__VLS_41, new __VLS_41({
        title: "基本信息",
        size: "small",
    }));
    const __VLS_43 = __VLS_42({
        title: "基本信息",
        size: "small",
    }, ...__VLS_functionalComponentArgsRest(__VLS_42));
    __VLS_44.slots.default;
    const __VLS_45 = {}.ADescriptions;
    /** @type {[typeof __VLS_components.ADescriptions, typeof __VLS_components.aDescriptions, typeof __VLS_components.ADescriptions, typeof __VLS_components.aDescriptions, ]} */ ;
    // @ts-ignore
    const __VLS_46 = __VLS_asFunctionalComponent(__VLS_45, new __VLS_45({
        column: (1),
        size: "small",
    }));
    const __VLS_47 = __VLS_46({
        column: (1),
        size: "small",
    }, ...__VLS_functionalComponentArgsRest(__VLS_46));
    __VLS_48.slots.default;
    const __VLS_49 = {}.ADescriptionsItem;
    /** @type {[typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_50 = __VLS_asFunctionalComponent(__VLS_49, new __VLS_49({
        label: "策略类别",
    }));
    const __VLS_51 = __VLS_50({
        label: "策略类别",
    }, ...__VLS_functionalComponentArgsRest(__VLS_50));
    __VLS_52.slots.default;
    (__VLS_ctx.strategy.category);
    var __VLS_52;
    const __VLS_53 = {}.ADescriptionsItem;
    /** @type {[typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_54 = __VLS_asFunctionalComponent(__VLS_53, new __VLS_53({
        label: "当前值",
    }));
    const __VLS_55 = __VLS_54({
        label: "当前值",
    }, ...__VLS_functionalComponentArgsRest(__VLS_54));
    __VLS_56.slots.default;
    (__VLS_ctx.strategy.current_value);
    var __VLS_56;
    const __VLS_57 = {}.ADescriptionsItem;
    /** @type {[typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_58 = __VLS_asFunctionalComponent(__VLS_57, new __VLS_57({
        label: "目标值",
    }));
    const __VLS_59 = __VLS_58({
        label: "目标值",
    }, ...__VLS_functionalComponentArgsRest(__VLS_58));
    __VLS_60.slots.default;
    (__VLS_ctx.strategy.target_value);
    var __VLS_60;
    const __VLS_61 = {}.ADescriptionsItem;
    /** @type {[typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_62 = __VLS_asFunctionalComponent(__VLS_61, new __VLS_61({
        label: "预期收益",
    }));
    const __VLS_63 = __VLS_62({
        label: "预期收益",
    }, ...__VLS_functionalComponentArgsRest(__VLS_62));
    __VLS_64.slots.default;
    (__VLS_ctx.strategy.expected_benefit);
    var __VLS_64;
    const __VLS_65 = {}.ADescriptionsItem;
    /** @type {[typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_66 = __VLS_asFunctionalComponent(__VLS_65, new __VLS_65({
        label: "风险说明",
    }));
    const __VLS_67 = __VLS_66({
        label: "风险说明",
    }, ...__VLS_functionalComponentArgsRest(__VLS_66));
    __VLS_68.slots.default;
    (__VLS_ctx.strategy.risk_note);
    var __VLS_68;
    var __VLS_48;
    var __VLS_44;
    var __VLS_40;
    const __VLS_69 = {}.ACol;
    /** @type {[typeof __VLS_components.ACol, typeof __VLS_components.aCol, typeof __VLS_components.ACol, typeof __VLS_components.aCol, ]} */ ;
    // @ts-ignore
    const __VLS_70 = __VLS_asFunctionalComponent(__VLS_69, new __VLS_69({
        span: (12),
    }));
    const __VLS_71 = __VLS_70({
        span: (12),
    }, ...__VLS_functionalComponentArgsRest(__VLS_70));
    __VLS_72.slots.default;
    const __VLS_73 = {}.ACard;
    /** @type {[typeof __VLS_components.ACard, typeof __VLS_components.aCard, typeof __VLS_components.ACard, typeof __VLS_components.aCard, ]} */ ;
    // @ts-ignore
    const __VLS_74 = __VLS_asFunctionalComponent(__VLS_73, new __VLS_73({
        title: "试运行信息",
        size: "small",
    }));
    const __VLS_75 = __VLS_74({
        title: "试运行信息",
        size: "small",
    }, ...__VLS_functionalComponentArgsRest(__VLS_74));
    __VLS_76.slots.default;
    if (__VLS_ctx.strategy.trial_run) {
        const __VLS_77 = {}.ADescriptions;
        /** @type {[typeof __VLS_components.ADescriptions, typeof __VLS_components.aDescriptions, typeof __VLS_components.ADescriptions, typeof __VLS_components.aDescriptions, ]} */ ;
        // @ts-ignore
        const __VLS_78 = __VLS_asFunctionalComponent(__VLS_77, new __VLS_77({
            column: (1),
            size: "small",
        }));
        const __VLS_79 = __VLS_78({
            column: (1),
            size: "small",
        }, ...__VLS_functionalComponentArgsRest(__VLS_78));
        __VLS_80.slots.default;
        const __VLS_81 = {}.ADescriptionsItem;
        /** @type {[typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, ]} */ ;
        // @ts-ignore
        const __VLS_82 = __VLS_asFunctionalComponent(__VLS_81, new __VLS_81({
            label: "状态",
        }));
        const __VLS_83 = __VLS_82({
            label: "状态",
        }, ...__VLS_functionalComponentArgsRest(__VLS_82));
        __VLS_84.slots.default;
        (__VLS_ctx.strategy.trial_run.status);
        var __VLS_84;
        const __VLS_85 = {}.ADescriptionsItem;
        /** @type {[typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, ]} */ ;
        // @ts-ignore
        const __VLS_86 = __VLS_asFunctionalComponent(__VLS_85, new __VLS_85({
            label: "天数",
        }));
        const __VLS_87 = __VLS_86({
            label: "天数",
        }, ...__VLS_functionalComponentArgsRest(__VLS_86));
        __VLS_88.slots.default;
        (__VLS_ctx.strategy.trial_run.trial_days);
        var __VLS_88;
        const __VLS_89 = {}.ADescriptionsItem;
        /** @type {[typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, ]} */ ;
        // @ts-ignore
        const __VLS_90 = __VLS_asFunctionalComponent(__VLS_89, new __VLS_89({
            label: "开始",
        }));
        const __VLS_91 = __VLS_90({
            label: "开始",
        }, ...__VLS_functionalComponentArgsRest(__VLS_90));
        __VLS_92.slots.default;
        (__VLS_ctx.strategy.trial_run.started_at);
        var __VLS_92;
        const __VLS_93 = {}.ADescriptionsItem;
        /** @type {[typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, typeof __VLS_components.ADescriptionsItem, typeof __VLS_components.aDescriptionsItem, ]} */ ;
        // @ts-ignore
        const __VLS_94 = __VLS_asFunctionalComponent(__VLS_93, new __VLS_93({
            label: "结束",
        }));
        const __VLS_95 = __VLS_94({
            label: "结束",
        }, ...__VLS_functionalComponentArgsRest(__VLS_94));
        __VLS_96.slots.default;
        (__VLS_ctx.strategy.trial_run.ends_at);
        var __VLS_96;
        var __VLS_80;
    }
    else {
        const __VLS_97 = {}.AEmpty;
        /** @type {[typeof __VLS_components.AEmpty, typeof __VLS_components.aEmpty, ]} */ ;
        // @ts-ignore
        const __VLS_98 = __VLS_asFunctionalComponent(__VLS_97, new __VLS_97({
            description: "尚未开始试运行",
            imageStyle: ({ height: '40px' }),
        }));
        const __VLS_99 = __VLS_98({
            description: "尚未开始试运行",
            imageStyle: ({ height: '40px' }),
        }, ...__VLS_functionalComponentArgsRest(__VLS_98));
    }
    var __VLS_76;
    var __VLS_72;
    var __VLS_36;
}
else if (!__VLS_ctx.loading) {
    const __VLS_101 = {}.AEmpty;
    /** @type {[typeof __VLS_components.AEmpty, typeof __VLS_components.aEmpty, ]} */ ;
    // @ts-ignore
    const __VLS_102 = __VLS_asFunctionalComponent(__VLS_101, new __VLS_101({
        description: "策略不存在",
    }));
    const __VLS_103 = __VLS_102({
        description: "策略不存在",
    }, ...__VLS_functionalComponentArgsRest(__VLS_102));
}
var __VLS_3;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            router: router,
            loading: loading,
            strategy: strategy,
            handleRollback: handleRollback,
            handlePromote: handlePromote,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
