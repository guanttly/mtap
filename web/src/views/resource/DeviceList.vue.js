/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { message, Modal } from 'ant-design-vue';
import { onMounted, ref } from 'vue';
import { resourceApi } from '@/api/resource';
import { usePagination } from '@/composables/usePagination';
const { loading, items, pagination, fetchData, onTableChange } = usePagination(params => resourceApi.listDevices(params));
onMounted(() => fetchData());
const STATUS_COLOR = { online: 'success', offline: 'default', maintenance: 'warning' };
const STATUS_LABEL = { online: '在线', offline: '离线', maintenance: '维护中' };
const columns = [
    { title: '设备名称', dataIndex: 'name', key: 'name' },
    { title: '型号', dataIndex: 'model', key: 'model' },
    { title: '厂商', dataIndex: 'manufacturer', key: 'manufacturer' },
    { title: '支持检查类型', key: 'exam_types' },
    { title: '每日最大号位', dataIndex: 'max_daily_slots', key: 'max_daily_slots' },
    { title: '状态', dataIndex: 'status', key: 'status' },
    { title: '操作', key: 'actions' },
];
const showModal = ref(false);
const saving = ref(false);
const editForm = ref({ status: 'online', supported_exam_types: [] });
const examTypesText = ref('');
function openCreate() {
    editForm.value = { status: 'online', supported_exam_types: [], max_daily_slots: 20 };
    examTypesText.value = '';
    showModal.value = true;
}
function openEdit(record) {
    editForm.value = { ...record };
    examTypesText.value = (record.supported_exam_types ?? []).join(',');
    showModal.value = true;
}
async function handleSave() {
    if (!editForm.value.name?.trim()) {
        message.warning('请填写设备名称');
        return;
    }
    editForm.value.supported_exam_types = examTypesText.value
        .split(',')
        .map(s => s.trim())
        .filter(Boolean);
    saving.value = true;
    try {
        if (editForm.value.id) {
            await resourceApi.updateDevice(editForm.value.id, editForm.value);
            message.success('更新成功');
        }
        else {
            await resourceApi.createDevice(editForm.value);
            message.success('创建成功');
        }
        showModal.value = false;
        fetchData();
    }
    finally {
        saving.value = false;
    }
}
function handleDelete(record) {
    Modal.confirm({
        title: '确认删除',
        content: `确定删除设备「${record.name}」吗？该操作将同时影响关联排班。`,
        okType: 'danger',
        onOk: async () => {
            await resourceApi.deleteDevice(record.id);
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
    if (column.key === 'exam_types') {
        const __VLS_24 = {}.ASpace;
        /** @type {[typeof __VLS_components.ASpace, typeof __VLS_components.aSpace, typeof __VLS_components.ASpace, typeof __VLS_components.aSpace, ]} */ ;
        // @ts-ignore
        const __VLS_25 = __VLS_asFunctionalComponent(__VLS_24, new __VLS_24({
            wrap: true,
        }));
        const __VLS_26 = __VLS_25({
            wrap: true,
        }, ...__VLS_functionalComponentArgsRest(__VLS_25));
        __VLS_27.slots.default;
        for (const [t] of __VLS_getVForSourceType((record.supported_exam_types))) {
            const __VLS_28 = {}.ATag;
            /** @type {[typeof __VLS_components.ATag, typeof __VLS_components.aTag, typeof __VLS_components.ATag, typeof __VLS_components.aTag, ]} */ ;
            // @ts-ignore
            const __VLS_29 = __VLS_asFunctionalComponent(__VLS_28, new __VLS_28({
                key: (t),
                color: "blue",
            }));
            const __VLS_30 = __VLS_29({
                key: (t),
                color: "blue",
            }, ...__VLS_functionalComponentArgsRest(__VLS_29));
            __VLS_31.slots.default;
            (t);
            var __VLS_31;
        }
        if (!record.supported_exam_types?.length) {
            __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
                ...{ style: {} },
            });
        }
        var __VLS_27;
    }
    else if (column.key === 'status') {
        const __VLS_32 = {}.ABadge;
        /** @type {[typeof __VLS_components.ABadge, typeof __VLS_components.aBadge, ]} */ ;
        // @ts-ignore
        const __VLS_33 = __VLS_asFunctionalComponent(__VLS_32, new __VLS_32({
            status: __VLS_ctx.STATUS_COLOR[record.status],
            text: (__VLS_ctx.STATUS_LABEL[record.status]),
        }));
        const __VLS_34 = __VLS_33({
            status: __VLS_ctx.STATUS_COLOR[record.status],
            text: (__VLS_ctx.STATUS_LABEL[record.status]),
        }, ...__VLS_functionalComponentArgsRest(__VLS_33));
    }
    else if (column.key === 'actions') {
        const __VLS_36 = {}.ASpace;
        /** @type {[typeof __VLS_components.ASpace, typeof __VLS_components.aSpace, typeof __VLS_components.ASpace, typeof __VLS_components.aSpace, ]} */ ;
        // @ts-ignore
        const __VLS_37 = __VLS_asFunctionalComponent(__VLS_36, new __VLS_36({}));
        const __VLS_38 = __VLS_37({}, ...__VLS_functionalComponentArgsRest(__VLS_37));
        __VLS_39.slots.default;
        const __VLS_40 = {}.AButton;
        /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
        // @ts-ignore
        const __VLS_41 = __VLS_asFunctionalComponent(__VLS_40, new __VLS_40({
            ...{ 'onClick': {} },
            type: "link",
            size: "small",
        }));
        const __VLS_42 = __VLS_41({
            ...{ 'onClick': {} },
            type: "link",
            size: "small",
        }, ...__VLS_functionalComponentArgsRest(__VLS_41));
        let __VLS_44;
        let __VLS_45;
        let __VLS_46;
        const __VLS_47 = {
            onClick: (...[$event]) => {
                if (!!(column.key === 'exam_types'))
                    return;
                if (!!(column.key === 'status'))
                    return;
                if (!(column.key === 'actions'))
                    return;
                __VLS_ctx.openEdit(record);
            }
        };
        __VLS_43.slots.default;
        var __VLS_43;
        const __VLS_48 = {}.AButton;
        /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
        // @ts-ignore
        const __VLS_49 = __VLS_asFunctionalComponent(__VLS_48, new __VLS_48({
            ...{ 'onClick': {} },
            type: "link",
            danger: true,
            size: "small",
        }));
        const __VLS_50 = __VLS_49({
            ...{ 'onClick': {} },
            type: "link",
            danger: true,
            size: "small",
        }, ...__VLS_functionalComponentArgsRest(__VLS_49));
        let __VLS_52;
        let __VLS_53;
        let __VLS_54;
        const __VLS_55 = {
            onClick: (...[$event]) => {
                if (!!(column.key === 'exam_types'))
                    return;
                if (!!(column.key === 'status'))
                    return;
                if (!(column.key === 'actions'))
                    return;
                __VLS_ctx.handleDelete(record);
            }
        };
        __VLS_51.slots.default;
        var __VLS_51;
        var __VLS_39;
    }
}
var __VLS_19;
const __VLS_56 = {}.AModal;
/** @type {[typeof __VLS_components.AModal, typeof __VLS_components.aModal, typeof __VLS_components.AModal, typeof __VLS_components.aModal, ]} */ ;
// @ts-ignore
const __VLS_57 = __VLS_asFunctionalComponent(__VLS_56, new __VLS_56({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.showModal),
    title: (__VLS_ctx.editForm.id ? '编辑设备' : '新增设备'),
    confirmLoading: (__VLS_ctx.saving),
}));
const __VLS_58 = __VLS_57({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.showModal),
    title: (__VLS_ctx.editForm.id ? '编辑设备' : '新增设备'),
    confirmLoading: (__VLS_ctx.saving),
}, ...__VLS_functionalComponentArgsRest(__VLS_57));
let __VLS_60;
let __VLS_61;
let __VLS_62;
const __VLS_63 = {
    onOk: (__VLS_ctx.handleSave)
};
__VLS_59.slots.default;
const __VLS_64 = {}.AForm;
/** @type {[typeof __VLS_components.AForm, typeof __VLS_components.aForm, typeof __VLS_components.AForm, typeof __VLS_components.aForm, ]} */ ;
// @ts-ignore
const __VLS_65 = __VLS_asFunctionalComponent(__VLS_64, new __VLS_64({
    model: (__VLS_ctx.editForm),
    layout: "vertical",
}));
const __VLS_66 = __VLS_65({
    model: (__VLS_ctx.editForm),
    layout: "vertical",
}, ...__VLS_functionalComponentArgsRest(__VLS_65));
__VLS_67.slots.default;
const __VLS_68 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_69 = __VLS_asFunctionalComponent(__VLS_68, new __VLS_68({
    label: "设备名称",
    required: true,
}));
const __VLS_70 = __VLS_69({
    label: "设备名称",
    required: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_69));
__VLS_71.slots.default;
const __VLS_72 = {}.AInput;
/** @type {[typeof __VLS_components.AInput, typeof __VLS_components.aInput, ]} */ ;
// @ts-ignore
const __VLS_73 = __VLS_asFunctionalComponent(__VLS_72, new __VLS_72({
    value: (__VLS_ctx.editForm.name),
    placeholder: "如：MRI-西门子Vida 3T",
}));
const __VLS_74 = __VLS_73({
    value: (__VLS_ctx.editForm.name),
    placeholder: "如：MRI-西门子Vida 3T",
}, ...__VLS_functionalComponentArgsRest(__VLS_73));
var __VLS_71;
const __VLS_76 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_77 = __VLS_asFunctionalComponent(__VLS_76, new __VLS_76({
    label: "型号",
}));
const __VLS_78 = __VLS_77({
    label: "型号",
}, ...__VLS_functionalComponentArgsRest(__VLS_77));
__VLS_79.slots.default;
const __VLS_80 = {}.AInput;
/** @type {[typeof __VLS_components.AInput, typeof __VLS_components.aInput, ]} */ ;
// @ts-ignore
const __VLS_81 = __VLS_asFunctionalComponent(__VLS_80, new __VLS_80({
    value: (__VLS_ctx.editForm.model),
}));
const __VLS_82 = __VLS_81({
    value: (__VLS_ctx.editForm.model),
}, ...__VLS_functionalComponentArgsRest(__VLS_81));
var __VLS_79;
const __VLS_84 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_85 = __VLS_asFunctionalComponent(__VLS_84, new __VLS_84({
    label: "厂商",
}));
const __VLS_86 = __VLS_85({
    label: "厂商",
}, ...__VLS_functionalComponentArgsRest(__VLS_85));
__VLS_87.slots.default;
const __VLS_88 = {}.AInput;
/** @type {[typeof __VLS_components.AInput, typeof __VLS_components.aInput, ]} */ ;
// @ts-ignore
const __VLS_89 = __VLS_asFunctionalComponent(__VLS_88, new __VLS_88({
    value: (__VLS_ctx.editForm.manufacturer),
}));
const __VLS_90 = __VLS_89({
    value: (__VLS_ctx.editForm.manufacturer),
}, ...__VLS_functionalComponentArgsRest(__VLS_89));
var __VLS_87;
const __VLS_92 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_93 = __VLS_asFunctionalComponent(__VLS_92, new __VLS_92({
    label: "每日最大号位数",
}));
const __VLS_94 = __VLS_93({
    label: "每日最大号位数",
}, ...__VLS_functionalComponentArgsRest(__VLS_93));
__VLS_95.slots.default;
const __VLS_96 = {}.AInputNumber;
/** @type {[typeof __VLS_components.AInputNumber, typeof __VLS_components.aInputNumber, ]} */ ;
// @ts-ignore
const __VLS_97 = __VLS_asFunctionalComponent(__VLS_96, new __VLS_96({
    value: (__VLS_ctx.editForm.max_daily_slots),
    min: (1),
    max: (500),
    ...{ style: {} },
}));
const __VLS_98 = __VLS_97({
    value: (__VLS_ctx.editForm.max_daily_slots),
    min: (1),
    max: (500),
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_97));
var __VLS_95;
const __VLS_100 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_101 = __VLS_asFunctionalComponent(__VLS_100, new __VLS_100({
    label: "支持检查类型（逗号分隔）",
}));
const __VLS_102 = __VLS_101({
    label: "支持检查类型（逗号分隔）",
}, ...__VLS_functionalComponentArgsRest(__VLS_101));
__VLS_103.slots.default;
const __VLS_104 = {}.ATextarea;
/** @type {[typeof __VLS_components.ATextarea, typeof __VLS_components.aTextarea, ]} */ ;
// @ts-ignore
const __VLS_105 = __VLS_asFunctionalComponent(__VLS_104, new __VLS_104({
    value: (__VLS_ctx.examTypesText),
    rows: (2),
    placeholder: "如：MRI平扫,MRI增强,MRI功能成像",
}));
const __VLS_106 = __VLS_105({
    value: (__VLS_ctx.examTypesText),
    rows: (2),
    placeholder: "如：MRI平扫,MRI增强,MRI功能成像",
}, ...__VLS_functionalComponentArgsRest(__VLS_105));
var __VLS_103;
const __VLS_108 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_109 = __VLS_asFunctionalComponent(__VLS_108, new __VLS_108({
    label: "状态",
}));
const __VLS_110 = __VLS_109({
    label: "状态",
}, ...__VLS_functionalComponentArgsRest(__VLS_109));
__VLS_111.slots.default;
const __VLS_112 = {}.ASelect;
/** @type {[typeof __VLS_components.ASelect, typeof __VLS_components.aSelect, typeof __VLS_components.ASelect, typeof __VLS_components.aSelect, ]} */ ;
// @ts-ignore
const __VLS_113 = __VLS_asFunctionalComponent(__VLS_112, new __VLS_112({
    value: (__VLS_ctx.editForm.status),
}));
const __VLS_114 = __VLS_113({
    value: (__VLS_ctx.editForm.status),
}, ...__VLS_functionalComponentArgsRest(__VLS_113));
__VLS_115.slots.default;
for (const [label, val] of __VLS_getVForSourceType((__VLS_ctx.STATUS_LABEL))) {
    const __VLS_116 = {}.ASelectOption;
    /** @type {[typeof __VLS_components.ASelectOption, typeof __VLS_components.aSelectOption, typeof __VLS_components.ASelectOption, typeof __VLS_components.aSelectOption, ]} */ ;
    // @ts-ignore
    const __VLS_117 = __VLS_asFunctionalComponent(__VLS_116, new __VLS_116({
        key: (val),
        value: (val),
    }));
    const __VLS_118 = __VLS_117({
        key: (val),
        value: (val),
    }, ...__VLS_functionalComponentArgsRest(__VLS_117));
    __VLS_119.slots.default;
    (label);
    var __VLS_119;
}
var __VLS_115;
var __VLS_111;
var __VLS_67;
var __VLS_59;
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
            STATUS_COLOR: STATUS_COLOR,
            STATUS_LABEL: STATUS_LABEL,
            columns: columns,
            showModal: showModal,
            saving: saving,
            editForm: editForm,
            examTypesText: examTypesText,
            openCreate: openCreate,
            openEdit: openEdit,
            handleSave: handleSave,
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
