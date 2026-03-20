/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { h, onMounted, ref } from 'vue';
import { message, Modal } from 'ant-design-vue';
import { ruleApi } from '@/api/rule';
import { usePagination } from '@/composables/usePagination';
const { loading, items, pagination, fetchData, onTableChange } = usePagination(params => ruleApi.listConflictRules(params));
onMounted(() => fetchData());
const showModal = ref(false);
const saving = ref(false);
const isEdit = ref(false);
const editId = ref('');
const editForm = ref({});
const LEVEL_MAP = { forbid: '禁止', warning: '警告' };
const LEVEL_COLOR = { forbid: 'red', warning: 'orange' };
const STATUS_MAP = {
    active: { color: 'success', label: '启用' },
    inactive: { color: 'default', label: '停用' },
};
const columns = [
    { title: '检查项A', dataIndex: 'item_a_name', key: 'item_a_name' },
    { title: '检查项B', dataIndex: 'item_b_name', key: 'item_b_name' },
    { title: '最短间隔', key: 'interval', customRender: ({ record }) => `${record.min_interval} ${record.interval_unit === 'hour' ? '小时' : '天'}` },
    { title: '冲突级别', dataIndex: 'level', key: 'level', customRender: ({ text }) => h('a-tag', { color: LEVEL_COLOR[text] ?? 'default' }, LEVEL_MAP[text] ?? text) },
    { title: '状态', dataIndex: 'status', key: 'status', customRender: ({ text }) => h('a-badge', { status: STATUS_MAP[text]?.color ?? 'default', text: STATUS_MAP[text]?.label ?? text }) },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at', customRender: ({ text }) => text ? text.slice(0, 10) : '' },
    { title: '操作', key: 'actions', width: 160 },
];
function openCreate() {
    isEdit.value = false;
    editId.value = '';
    editForm.value = { level: 'warning', interval_unit: 'hour', min_interval: 24 };
    showModal.value = true;
}
function openEdit(record) {
    isEdit.value = true;
    editId.value = record.id;
    editForm.value = { min_interval: record.min_interval, interval_unit: record.interval_unit, level: record.level, status: record.status };
    showModal.value = true;
}
async function handleSave() {
    saving.value = true;
    try {
        if (isEdit.value) {
            await ruleApi.updateConflictRule(editId.value, {
                min_interval: editForm.value.min_interval,
                level: editForm.value.level,
                status: editForm.value.status,
            });
            message.success('更新成功');
        }
        else {
            await ruleApi.createConflictRule(editForm.value);
            message.success('创建成功');
        }
        showModal.value = false;
        fetchData();
    }
    finally {
        saving.value = false;
    }
}
function handleToggleStatus(record) {
    const next = record.status === 'active' ? 'inactive' : 'active';
    Modal.confirm({
        title: `确认${next === 'inactive' ? '停用' : '启用'}`,
        content: `确定要${next === 'inactive' ? '停用' : '启用'}该冲突规则吗？`,
        onOk: async () => {
            await ruleApi.updateConflictRule(record.id, { status: next });
            message.success('操作成功');
            fetchData();
        },
    });
}
function handleDelete(record) {
    Modal.confirm({
        title: '确认删除',
        content: `确定删除「${record.item_a_name ?? record.item_a_id} ↔ ${record.item_b_name ?? record.item_b_id}」冲突规则吗？`,
        okType: 'danger',
        onOk: async () => {
            await ruleApi.deleteConflictRule(record.id);
            message.success('删除成功');
            fetchData();
        },
    });
}
debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
const __VLS_ctx = {};
let __VLS_components;
let __VLS_directives;
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "mb-4 flex items-center gap-2" },
});
const __VLS_0 = {}.AButton;
/** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
// @ts-ignore
const __VLS_1 = __VLS_asFunctionalComponent(__VLS_0, new __VLS_0({
    ...{ 'onClick': {} },
    type: "primary",
}));
const __VLS_2 = __VLS_1({
    ...{ 'onClick': {} },
    type: "primary",
}, ...__VLS_functionalComponentArgsRest(__VLS_1));
let __VLS_4;
let __VLS_5;
let __VLS_6;
const __VLS_7 = {
    onClick: (__VLS_ctx.openCreate)
};
__VLS_3.slots.default;
var __VLS_3;
const __VLS_8 = {}.AButton;
/** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
// @ts-ignore
const __VLS_9 = __VLS_asFunctionalComponent(__VLS_8, new __VLS_8({
    ...{ 'onClick': {} },
}));
const __VLS_10 = __VLS_9({
    ...{ 'onClick': {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_9));
let __VLS_12;
let __VLS_13;
let __VLS_14;
const __VLS_15 = {
    onClick: (__VLS_ctx.fetchData)
};
__VLS_11.slots.default;
var __VLS_11;
const __VLS_16 = {}.ATable;
/** @type {[typeof __VLS_components.ATable, typeof __VLS_components.aTable, typeof __VLS_components.ATable, typeof __VLS_components.aTable, ]} */ ;
// @ts-ignore
const __VLS_17 = __VLS_asFunctionalComponent(__VLS_16, new __VLS_16({
    ...{ 'onChange': {} },
    columns: (__VLS_ctx.columns),
    dataSource: (__VLS_ctx.items),
    loading: (__VLS_ctx.loading),
    pagination: (__VLS_ctx.pagination),
    rowKey: "id",
    size: "middle",
}));
const __VLS_18 = __VLS_17({
    ...{ 'onChange': {} },
    columns: (__VLS_ctx.columns),
    dataSource: (__VLS_ctx.items),
    loading: (__VLS_ctx.loading),
    pagination: (__VLS_ctx.pagination),
    rowKey: "id",
    size: "middle",
}, ...__VLS_functionalComponentArgsRest(__VLS_17));
let __VLS_20;
let __VLS_21;
let __VLS_22;
const __VLS_23 = {
    onChange: (__VLS_ctx.onTableChange)
};
__VLS_19.slots.default;
{
    const { bodyCell: __VLS_thisSlot } = __VLS_19.slots;
    const [{ column, record }] = __VLS_getSlotParams(__VLS_thisSlot);
    if (column.key === 'actions') {
        const __VLS_24 = {}.ASpace;
        /** @type {[typeof __VLS_components.ASpace, typeof __VLS_components.aSpace, typeof __VLS_components.ASpace, typeof __VLS_components.aSpace, ]} */ ;
        // @ts-ignore
        const __VLS_25 = __VLS_asFunctionalComponent(__VLS_24, new __VLS_24({}));
        const __VLS_26 = __VLS_25({}, ...__VLS_functionalComponentArgsRest(__VLS_25));
        __VLS_27.slots.default;
        const __VLS_28 = {}.AButton;
        /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
        // @ts-ignore
        const __VLS_29 = __VLS_asFunctionalComponent(__VLS_28, new __VLS_28({
            ...{ 'onClick': {} },
            type: "link",
            size: "small",
        }));
        const __VLS_30 = __VLS_29({
            ...{ 'onClick': {} },
            type: "link",
            size: "small",
        }, ...__VLS_functionalComponentArgsRest(__VLS_29));
        let __VLS_32;
        let __VLS_33;
        let __VLS_34;
        const __VLS_35 = {
            onClick: (...[$event]) => {
                if (!(column.key === 'actions'))
                    return;
                __VLS_ctx.openEdit(record);
            }
        };
        __VLS_31.slots.default;
        var __VLS_31;
        const __VLS_36 = {}.AButton;
        /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
        // @ts-ignore
        const __VLS_37 = __VLS_asFunctionalComponent(__VLS_36, new __VLS_36({
            ...{ 'onClick': {} },
            type: "link",
            size: "small",
        }));
        const __VLS_38 = __VLS_37({
            ...{ 'onClick': {} },
            type: "link",
            size: "small",
        }, ...__VLS_functionalComponentArgsRest(__VLS_37));
        let __VLS_40;
        let __VLS_41;
        let __VLS_42;
        const __VLS_43 = {
            onClick: (...[$event]) => {
                if (!(column.key === 'actions'))
                    return;
                __VLS_ctx.handleToggleStatus(record);
            }
        };
        __VLS_39.slots.default;
        (record.status === 'active' ? '停用' : '启用');
        var __VLS_39;
        const __VLS_44 = {}.AButton;
        /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
        // @ts-ignore
        const __VLS_45 = __VLS_asFunctionalComponent(__VLS_44, new __VLS_44({
            ...{ 'onClick': {} },
            type: "link",
            danger: true,
            size: "small",
        }));
        const __VLS_46 = __VLS_45({
            ...{ 'onClick': {} },
            type: "link",
            danger: true,
            size: "small",
        }, ...__VLS_functionalComponentArgsRest(__VLS_45));
        let __VLS_48;
        let __VLS_49;
        let __VLS_50;
        const __VLS_51 = {
            onClick: (...[$event]) => {
                if (!(column.key === 'actions'))
                    return;
                __VLS_ctx.handleDelete(record);
            }
        };
        __VLS_47.slots.default;
        var __VLS_47;
        var __VLS_27;
    }
}
var __VLS_19;
const __VLS_52 = {}.AModal;
/** @type {[typeof __VLS_components.AModal, typeof __VLS_components.aModal, typeof __VLS_components.AModal, typeof __VLS_components.aModal, ]} */ ;
// @ts-ignore
const __VLS_53 = __VLS_asFunctionalComponent(__VLS_52, new __VLS_52({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.showModal),
    title: (__VLS_ctx.isEdit ? '编辑冲突规则' : '新增冲突规则'),
    confirmLoading: (__VLS_ctx.saving),
}));
const __VLS_54 = __VLS_53({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.showModal),
    title: (__VLS_ctx.isEdit ? '编辑冲突规则' : '新增冲突规则'),
    confirmLoading: (__VLS_ctx.saving),
}, ...__VLS_functionalComponentArgsRest(__VLS_53));
let __VLS_56;
let __VLS_57;
let __VLS_58;
const __VLS_59 = {
    onOk: (__VLS_ctx.handleSave)
};
__VLS_55.slots.default;
const __VLS_60 = {}.AForm;
/** @type {[typeof __VLS_components.AForm, typeof __VLS_components.aForm, typeof __VLS_components.AForm, typeof __VLS_components.aForm, ]} */ ;
// @ts-ignore
const __VLS_61 = __VLS_asFunctionalComponent(__VLS_60, new __VLS_60({
    model: (__VLS_ctx.editForm),
    layout: "vertical",
}));
const __VLS_62 = __VLS_61({
    model: (__VLS_ctx.editForm),
    layout: "vertical",
}, ...__VLS_functionalComponentArgsRest(__VLS_61));
__VLS_63.slots.default;
if (!__VLS_ctx.isEdit) {
    const __VLS_64 = {}.AFormItem;
    /** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
    // @ts-ignore
    const __VLS_65 = __VLS_asFunctionalComponent(__VLS_64, new __VLS_64({
        label: "检查项A ID",
        required: true,
    }));
    const __VLS_66 = __VLS_65({
        label: "检查项A ID",
        required: true,
    }, ...__VLS_functionalComponentArgsRest(__VLS_65));
    __VLS_67.slots.default;
    const __VLS_68 = {}.AInput;
    /** @type {[typeof __VLS_components.AInput, typeof __VLS_components.aInput, ]} */ ;
    // @ts-ignore
    const __VLS_69 = __VLS_asFunctionalComponent(__VLS_68, new __VLS_68({
        value: (__VLS_ctx.editForm.item_a_id),
        placeholder: "请输入检查项A的ID",
    }));
    const __VLS_70 = __VLS_69({
        value: (__VLS_ctx.editForm.item_a_id),
        placeholder: "请输入检查项A的ID",
    }, ...__VLS_functionalComponentArgsRest(__VLS_69));
    var __VLS_67;
    const __VLS_72 = {}.AFormItem;
    /** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
    // @ts-ignore
    const __VLS_73 = __VLS_asFunctionalComponent(__VLS_72, new __VLS_72({
        label: "检查项B ID",
        required: true,
    }));
    const __VLS_74 = __VLS_73({
        label: "检查项B ID",
        required: true,
    }, ...__VLS_functionalComponentArgsRest(__VLS_73));
    __VLS_75.slots.default;
    const __VLS_76 = {}.AInput;
    /** @type {[typeof __VLS_components.AInput, typeof __VLS_components.aInput, ]} */ ;
    // @ts-ignore
    const __VLS_77 = __VLS_asFunctionalComponent(__VLS_76, new __VLS_76({
        value: (__VLS_ctx.editForm.item_b_id),
        placeholder: "请输入检查项B的ID",
    }));
    const __VLS_78 = __VLS_77({
        value: (__VLS_ctx.editForm.item_b_id),
        placeholder: "请输入检查项B的ID",
    }, ...__VLS_functionalComponentArgsRest(__VLS_77));
    var __VLS_75;
}
const __VLS_80 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_81 = __VLS_asFunctionalComponent(__VLS_80, new __VLS_80({
    label: "最短间隔",
}));
const __VLS_82 = __VLS_81({
    label: "最短间隔",
}, ...__VLS_functionalComponentArgsRest(__VLS_81));
__VLS_83.slots.default;
const __VLS_84 = {}.AInputNumber;
/** @type {[typeof __VLS_components.AInputNumber, typeof __VLS_components.aInputNumber, ]} */ ;
// @ts-ignore
const __VLS_85 = __VLS_asFunctionalComponent(__VLS_84, new __VLS_84({
    value: (__VLS_ctx.editForm.min_interval),
    min: (0),
    max: (720),
    ...{ style: {} },
}));
const __VLS_86 = __VLS_85({
    value: (__VLS_ctx.editForm.min_interval),
    min: (0),
    max: (720),
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_85));
const __VLS_88 = {}.ASelect;
/** @type {[typeof __VLS_components.ASelect, typeof __VLS_components.aSelect, typeof __VLS_components.ASelect, typeof __VLS_components.aSelect, ]} */ ;
// @ts-ignore
const __VLS_89 = __VLS_asFunctionalComponent(__VLS_88, new __VLS_88({
    value: (__VLS_ctx.editForm.interval_unit),
    ...{ style: {} },
}));
const __VLS_90 = __VLS_89({
    value: (__VLS_ctx.editForm.interval_unit),
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_89));
__VLS_91.slots.default;
const __VLS_92 = {}.ASelectOption;
/** @type {[typeof __VLS_components.ASelectOption, typeof __VLS_components.aSelectOption, typeof __VLS_components.ASelectOption, typeof __VLS_components.aSelectOption, ]} */ ;
// @ts-ignore
const __VLS_93 = __VLS_asFunctionalComponent(__VLS_92, new __VLS_92({
    value: "hour",
}));
const __VLS_94 = __VLS_93({
    value: "hour",
}, ...__VLS_functionalComponentArgsRest(__VLS_93));
__VLS_95.slots.default;
var __VLS_95;
const __VLS_96 = {}.ASelectOption;
/** @type {[typeof __VLS_components.ASelectOption, typeof __VLS_components.aSelectOption, typeof __VLS_components.ASelectOption, typeof __VLS_components.aSelectOption, ]} */ ;
// @ts-ignore
const __VLS_97 = __VLS_asFunctionalComponent(__VLS_96, new __VLS_96({
    value: "day",
}));
const __VLS_98 = __VLS_97({
    value: "day",
}, ...__VLS_functionalComponentArgsRest(__VLS_97));
__VLS_99.slots.default;
var __VLS_99;
var __VLS_91;
var __VLS_83;
const __VLS_100 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_101 = __VLS_asFunctionalComponent(__VLS_100, new __VLS_100({
    label: "冲突级别",
}));
const __VLS_102 = __VLS_101({
    label: "冲突级别",
}, ...__VLS_functionalComponentArgsRest(__VLS_101));
__VLS_103.slots.default;
const __VLS_104 = {}.ARadioGroup;
/** @type {[typeof __VLS_components.ARadioGroup, typeof __VLS_components.aRadioGroup, typeof __VLS_components.ARadioGroup, typeof __VLS_components.aRadioGroup, ]} */ ;
// @ts-ignore
const __VLS_105 = __VLS_asFunctionalComponent(__VLS_104, new __VLS_104({
    value: (__VLS_ctx.editForm.level),
}));
const __VLS_106 = __VLS_105({
    value: (__VLS_ctx.editForm.level),
}, ...__VLS_functionalComponentArgsRest(__VLS_105));
__VLS_107.slots.default;
const __VLS_108 = {}.ARadio;
/** @type {[typeof __VLS_components.ARadio, typeof __VLS_components.aRadio, typeof __VLS_components.ARadio, typeof __VLS_components.aRadio, ]} */ ;
// @ts-ignore
const __VLS_109 = __VLS_asFunctionalComponent(__VLS_108, new __VLS_108({
    value: "warning",
}));
const __VLS_110 = __VLS_109({
    value: "warning",
}, ...__VLS_functionalComponentArgsRest(__VLS_109));
__VLS_111.slots.default;
var __VLS_111;
const __VLS_112 = {}.ARadio;
/** @type {[typeof __VLS_components.ARadio, typeof __VLS_components.aRadio, typeof __VLS_components.ARadio, typeof __VLS_components.aRadio, ]} */ ;
// @ts-ignore
const __VLS_113 = __VLS_asFunctionalComponent(__VLS_112, new __VLS_112({
    value: "forbid",
}));
const __VLS_114 = __VLS_113({
    value: "forbid",
}, ...__VLS_functionalComponentArgsRest(__VLS_113));
__VLS_115.slots.default;
var __VLS_115;
var __VLS_107;
var __VLS_103;
if (__VLS_ctx.isEdit) {
    const __VLS_116 = {}.AFormItem;
    /** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
    // @ts-ignore
    const __VLS_117 = __VLS_asFunctionalComponent(__VLS_116, new __VLS_116({
        label: "状态",
    }));
    const __VLS_118 = __VLS_117({
        label: "状态",
    }, ...__VLS_functionalComponentArgsRest(__VLS_117));
    __VLS_119.slots.default;
    const __VLS_120 = {}.ARadioGroup;
    /** @type {[typeof __VLS_components.ARadioGroup, typeof __VLS_components.aRadioGroup, typeof __VLS_components.ARadioGroup, typeof __VLS_components.aRadioGroup, ]} */ ;
    // @ts-ignore
    const __VLS_121 = __VLS_asFunctionalComponent(__VLS_120, new __VLS_120({
        value: (__VLS_ctx.editForm.status),
    }));
    const __VLS_122 = __VLS_121({
        value: (__VLS_ctx.editForm.status),
    }, ...__VLS_functionalComponentArgsRest(__VLS_121));
    __VLS_123.slots.default;
    const __VLS_124 = {}.ARadio;
    /** @type {[typeof __VLS_components.ARadio, typeof __VLS_components.aRadio, typeof __VLS_components.ARadio, typeof __VLS_components.aRadio, ]} */ ;
    // @ts-ignore
    const __VLS_125 = __VLS_asFunctionalComponent(__VLS_124, new __VLS_124({
        value: "active",
    }));
    const __VLS_126 = __VLS_125({
        value: "active",
    }, ...__VLS_functionalComponentArgsRest(__VLS_125));
    __VLS_127.slots.default;
    var __VLS_127;
    const __VLS_128 = {}.ARadio;
    /** @type {[typeof __VLS_components.ARadio, typeof __VLS_components.aRadio, typeof __VLS_components.ARadio, typeof __VLS_components.aRadio, ]} */ ;
    // @ts-ignore
    const __VLS_129 = __VLS_asFunctionalComponent(__VLS_128, new __VLS_128({
        value: "inactive",
    }));
    const __VLS_130 = __VLS_129({
        value: "inactive",
    }, ...__VLS_functionalComponentArgsRest(__VLS_129));
    __VLS_131.slots.default;
    var __VLS_131;
    var __VLS_123;
    var __VLS_119;
}
var __VLS_63;
var __VLS_55;
/** @type {__VLS_StyleScopedClasses['mb-4']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['gap-2']} */ ;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            loading: loading,
            items: items,
            pagination: pagination,
            fetchData: fetchData,
            onTableChange: onTableChange,
            showModal: showModal,
            saving: saving,
            isEdit: isEdit,
            editForm: editForm,
            columns: columns,
            openCreate: openCreate,
            openEdit: openEdit,
            handleSave: handleSave,
            handleToggleStatus: handleToggleStatus,
            handleDelete: handleDelete,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
