/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { message, Modal } from 'ant-design-vue';
import { h, onMounted, ref } from 'vue';
import { resourceApi } from '@/api/resource';
import { usePagination } from '@/composables/usePagination';
const { loading, items, pagination, fetchData, onTableChange } = usePagination(params => resourceApi.listExamItems(params));
onMounted(() => fetchData());
const showModal = ref(false);
const saving = ref(false);
const editForm = ref({});
const isEdit = ref(false);
const columns = [
    { title: '项目名称', dataIndex: 'name', key: 'name' },
    {
        title: '标准时长(分钟)',
        dataIndex: 'duration_min',
        key: 'duration_min',
    },
    {
        title: '空腹要求',
        dataIndex: 'is_fasting',
        key: 'is_fasting',
        customRender: ({ text }) => h('a-tag', { color: text ? 'orange' : 'default' }, text ? '需空腹' : '无要求'),
    },
    { title: '空腹说明', dataIndex: 'fasting_desc', key: 'fasting_desc' },
    { title: '操作', key: 'actions' },
];
function openCreate() {
    editForm.value = { is_fasting: false, duration_min: 30, fasting_desc: '' };
    isEdit.value = false;
    showModal.value = true;
}
function openEdit(record) {
    editForm.value = { ...record };
    isEdit.value = true;
    showModal.value = true;
}
async function handleSave() {
    if (!editForm.value.name?.trim()) {
        message.warning('请填写项目名称');
        return;
    }
    saving.value = true;
    try {
        if (isEdit.value && editForm.value.id) {
            await resourceApi.updateExamItem(editForm.value.id, editForm.value);
            message.success('更新成功');
        }
        else {
            await resourceApi.createExamItem(editForm.value);
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
        content: `确定删除检查项目「${record.name}」吗？此操作不可恢复。`,
        okType: 'danger',
        onOk: async () => {
            await resourceApi.deleteExamItem(record.id);
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
    ...{ class: "action-bar mb-4 flex items-center gap-2" },
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
            danger: true,
            size: "small",
        }));
        const __VLS_38 = __VLS_37({
            ...{ 'onClick': {} },
            type: "link",
            danger: true,
            size: "small",
        }, ...__VLS_functionalComponentArgsRest(__VLS_37));
        let __VLS_40;
        let __VLS_41;
        let __VLS_42;
        const __VLS_43 = {
            onClick: (...[$event]) => {
                if (!(column.key === 'actions'))
                    return;
                __VLS_ctx.handleDelete(record);
            }
        };
        __VLS_39.slots.default;
        var __VLS_39;
        var __VLS_27;
    }
}
var __VLS_19;
const __VLS_44 = {}.AModal;
/** @type {[typeof __VLS_components.AModal, typeof __VLS_components.aModal, typeof __VLS_components.AModal, typeof __VLS_components.aModal, ]} */ ;
// @ts-ignore
const __VLS_45 = __VLS_asFunctionalComponent(__VLS_44, new __VLS_44({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.showModal),
    title: (__VLS_ctx.isEdit ? '编辑检查项目' : '新增检查项目'),
    confirmLoading: (__VLS_ctx.saving),
}));
const __VLS_46 = __VLS_45({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.showModal),
    title: (__VLS_ctx.isEdit ? '编辑检查项目' : '新增检查项目'),
    confirmLoading: (__VLS_ctx.saving),
}, ...__VLS_functionalComponentArgsRest(__VLS_45));
let __VLS_48;
let __VLS_49;
let __VLS_50;
const __VLS_51 = {
    onOk: (__VLS_ctx.handleSave)
};
__VLS_47.slots.default;
const __VLS_52 = {}.AForm;
/** @type {[typeof __VLS_components.AForm, typeof __VLS_components.aForm, typeof __VLS_components.AForm, typeof __VLS_components.aForm, ]} */ ;
// @ts-ignore
const __VLS_53 = __VLS_asFunctionalComponent(__VLS_52, new __VLS_52({
    model: (__VLS_ctx.editForm),
    layout: "vertical",
}));
const __VLS_54 = __VLS_53({
    model: (__VLS_ctx.editForm),
    layout: "vertical",
}, ...__VLS_functionalComponentArgsRest(__VLS_53));
__VLS_55.slots.default;
const __VLS_56 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_57 = __VLS_asFunctionalComponent(__VLS_56, new __VLS_56({
    label: "项目名称",
    required: true,
}));
const __VLS_58 = __VLS_57({
    label: "项目名称",
    required: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_57));
__VLS_59.slots.default;
const __VLS_60 = {}.AInput;
/** @type {[typeof __VLS_components.AInput, typeof __VLS_components.aInput, ]} */ ;
// @ts-ignore
const __VLS_61 = __VLS_asFunctionalComponent(__VLS_60, new __VLS_60({
    value: (__VLS_ctx.editForm.name),
    placeholder: "请输入检查项目名称",
}));
const __VLS_62 = __VLS_61({
    value: (__VLS_ctx.editForm.name),
    placeholder: "请输入检查项目名称",
}, ...__VLS_functionalComponentArgsRest(__VLS_61));
var __VLS_59;
const __VLS_64 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_65 = __VLS_asFunctionalComponent(__VLS_64, new __VLS_64({
    label: "标准时长（分钟）",
}));
const __VLS_66 = __VLS_65({
    label: "标准时长（分钟）",
}, ...__VLS_functionalComponentArgsRest(__VLS_65));
__VLS_67.slots.default;
const __VLS_68 = {}.AInputNumber;
/** @type {[typeof __VLS_components.AInputNumber, typeof __VLS_components.aInputNumber, ]} */ ;
// @ts-ignore
const __VLS_69 = __VLS_asFunctionalComponent(__VLS_68, new __VLS_68({
    value: (__VLS_ctx.editForm.duration_min),
    min: (1),
    max: (480),
    ...{ style: {} },
}));
const __VLS_70 = __VLS_69({
    value: (__VLS_ctx.editForm.duration_min),
    min: (1),
    max: (480),
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_69));
var __VLS_67;
const __VLS_72 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_73 = __VLS_asFunctionalComponent(__VLS_72, new __VLS_72({
    label: "空腹要求",
}));
const __VLS_74 = __VLS_73({
    label: "空腹要求",
}, ...__VLS_functionalComponentArgsRest(__VLS_73));
__VLS_75.slots.default;
const __VLS_76 = {}.ASwitch;
/** @type {[typeof __VLS_components.ASwitch, typeof __VLS_components.aSwitch, ]} */ ;
// @ts-ignore
const __VLS_77 = __VLS_asFunctionalComponent(__VLS_76, new __VLS_76({
    checked: (__VLS_ctx.editForm.is_fasting),
    checkedChildren: "需空腹",
    unCheckedChildren: "无要求",
}));
const __VLS_78 = __VLS_77({
    checked: (__VLS_ctx.editForm.is_fasting),
    checkedChildren: "需空腹",
    unCheckedChildren: "无要求",
}, ...__VLS_functionalComponentArgsRest(__VLS_77));
var __VLS_75;
if (__VLS_ctx.editForm.is_fasting) {
    const __VLS_80 = {}.AFormItem;
    /** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
    // @ts-ignore
    const __VLS_81 = __VLS_asFunctionalComponent(__VLS_80, new __VLS_80({
        label: "空腹说明",
    }));
    const __VLS_82 = __VLS_81({
        label: "空腹说明",
    }, ...__VLS_functionalComponentArgsRest(__VLS_81));
    __VLS_83.slots.default;
    const __VLS_84 = {}.ATextarea;
    /** @type {[typeof __VLS_components.ATextarea, typeof __VLS_components.aTextarea, ]} */ ;
    // @ts-ignore
    const __VLS_85 = __VLS_asFunctionalComponent(__VLS_84, new __VLS_84({
        value: (__VLS_ctx.editForm.fasting_desc),
        rows: (3),
        placeholder: "请描述空腹要求（如：检查前8小时禁食禁水）",
        maxlength: (200),
        showCount: true,
    }));
    const __VLS_86 = __VLS_85({
        value: (__VLS_ctx.editForm.fasting_desc),
        rows: (3),
        placeholder: "请描述空腹要求（如：检查前8小时禁食禁水）",
        maxlength: (200),
        showCount: true,
    }, ...__VLS_functionalComponentArgsRest(__VLS_85));
    var __VLS_83;
}
var __VLS_55;
var __VLS_47;
/** @type {__VLS_StyleScopedClasses['action-bar']} */ ;
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
            editForm: editForm,
            isEdit: isEdit,
            columns: columns,
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
