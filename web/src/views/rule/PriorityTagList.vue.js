/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { onMounted, ref } from 'vue';
import { message, Modal } from 'ant-design-vue';
import { ruleApi } from '@/api/rule';
const loading = ref(false);
const items = ref([]);
async function fetchData() {
    loading.value = true;
    try {
        const res = await ruleApi.listPriorityTags();
        items.value = res.items;
    }
    finally {
        loading.value = false;
    }
}
onMounted(fetchData);
const showModal = ref(false);
const saving = ref(false);
const isEdit = ref(false);
const editId = ref('');
const editForm = ref({ color: '#1890ff', weight: 10 });
const columns = [
    { title: '标签名', dataIndex: 'name', key: 'name' },
    { title: '权重', dataIndex: 'weight', key: 'weight' },
    { title: '颜色', key: 'color' },
    { title: '类型', key: 'is_preset' },
    { title: '操作', key: 'actions', width: 120 },
];
function openCreate() {
    isEdit.value = false;
    editId.value = '';
    editForm.value = { color: '#1890ff', weight: 10 };
    showModal.value = true;
}
function openEdit(record) {
    if (record.is_preset) {
        message.warning('预置标签仅允许调整权重和颜色');
    }
    isEdit.value = true;
    editId.value = record.id;
    editForm.value = { name: record.name, weight: record.weight, color: record.color };
    showModal.value = true;
}
async function handleSave() {
    if (!editForm.value.name?.trim() && !isEdit.value) {
        message.warning('请填写标签名称');
        return;
    }
    saving.value = true;
    try {
        if (isEdit.value) {
            await ruleApi.updatePriorityTag(editId.value, {
                name: editForm.value.name,
                weight: editForm.value.weight,
                color: editForm.value.color,
            });
            message.success('更新成功');
        }
        else {
            await ruleApi.createPriorityTag(editForm.value);
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
    if (record.is_preset) {
        message.warning('预置标签不可删除');
        return;
    }
    Modal.confirm({
        title: '确认删除',
        content: `确定删除优先级标签「${record.name}」吗？`,
        okType: 'danger',
        onOk: async () => {
            await ruleApi.deletePriorityTag(record.id);
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
    columns: (__VLS_ctx.columns),
    dataSource: (__VLS_ctx.items),
    loading: (__VLS_ctx.loading),
    pagination: (false),
    rowKey: "id",
    size: "middle",
}));
const __VLS_18 = __VLS_17({
    columns: (__VLS_ctx.columns),
    dataSource: (__VLS_ctx.items),
    loading: (__VLS_ctx.loading),
    pagination: (false),
    rowKey: "id",
    size: "middle",
}, ...__VLS_functionalComponentArgsRest(__VLS_17));
__VLS_19.slots.default;
{
    const { bodyCell: __VLS_thisSlot } = __VLS_19.slots;
    const [{ column, record }] = __VLS_getSlotParams(__VLS_thisSlot);
    if (column.key === 'color') {
        __VLS_asFunctionalElement(__VLS_intrinsicElements.span)({
            ...{ style: ({ display: 'inline-block', width: '16px', height: '16px', borderRadius: '4px', background: record.color, verticalAlign: 'middle', marginRight: '6px' }) },
        });
        __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
            ...{ style: {} },
        });
        (record.color);
    }
    else if (column.key === 'is_preset') {
        const __VLS_20 = {}.ATag;
        /** @type {[typeof __VLS_components.ATag, typeof __VLS_components.aTag, typeof __VLS_components.ATag, typeof __VLS_components.aTag, ]} */ ;
        // @ts-ignore
        const __VLS_21 = __VLS_asFunctionalComponent(__VLS_20, new __VLS_20({
            color: (record.is_preset ? 'purple' : 'default'),
        }));
        const __VLS_22 = __VLS_21({
            color: (record.is_preset ? 'purple' : 'default'),
        }, ...__VLS_functionalComponentArgsRest(__VLS_21));
        __VLS_23.slots.default;
        (record.is_preset ? '系统预置' : '自定义');
        var __VLS_23;
    }
    else if (column.key === 'actions') {
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
                if (!!(column.key === 'color'))
                    return;
                if (!!(column.key === 'is_preset'))
                    return;
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
            disabled: (record.is_preset),
        }));
        const __VLS_38 = __VLS_37({
            ...{ 'onClick': {} },
            type: "link",
            danger: true,
            size: "small",
            disabled: (record.is_preset),
        }, ...__VLS_functionalComponentArgsRest(__VLS_37));
        let __VLS_40;
        let __VLS_41;
        let __VLS_42;
        const __VLS_43 = {
            onClick: (...[$event]) => {
                if (!!(column.key === 'color'))
                    return;
                if (!!(column.key === 'is_preset'))
                    return;
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
    title: (__VLS_ctx.isEdit ? '编辑优先级标签' : '新增优先级标签'),
    confirmLoading: (__VLS_ctx.saving),
}));
const __VLS_46 = __VLS_45({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.showModal),
    title: (__VLS_ctx.isEdit ? '编辑优先级标签' : '新增优先级标签'),
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
    label: "标签名称",
    required: (!__VLS_ctx.isEdit),
}));
const __VLS_58 = __VLS_57({
    label: "标签名称",
    required: (!__VLS_ctx.isEdit),
}, ...__VLS_functionalComponentArgsRest(__VLS_57));
__VLS_59.slots.default;
const __VLS_60 = {}.AInput;
/** @type {[typeof __VLS_components.AInput, typeof __VLS_components.aInput, ]} */ ;
// @ts-ignore
const __VLS_61 = __VLS_asFunctionalComponent(__VLS_60, new __VLS_60({
    value: (__VLS_ctx.editForm.name),
    placeholder: "例如：急诊优先、老年患者",
    maxlength: (20),
    showCount: true,
}));
const __VLS_62 = __VLS_61({
    value: (__VLS_ctx.editForm.name),
    placeholder: "例如：急诊优先、老年患者",
    maxlength: (20),
    showCount: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_61));
var __VLS_59;
const __VLS_64 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_65 = __VLS_asFunctionalComponent(__VLS_64, new __VLS_64({
    label: "权重（1~100，越高越优先）",
}));
const __VLS_66 = __VLS_65({
    label: "权重（1~100，越高越优先）",
}, ...__VLS_functionalComponentArgsRest(__VLS_65));
__VLS_67.slots.default;
const __VLS_68 = {}.AInputNumber;
/** @type {[typeof __VLS_components.AInputNumber, typeof __VLS_components.aInputNumber, ]} */ ;
// @ts-ignore
const __VLS_69 = __VLS_asFunctionalComponent(__VLS_68, new __VLS_68({
    value: (__VLS_ctx.editForm.weight),
    min: (1),
    max: (100),
    ...{ style: {} },
}));
const __VLS_70 = __VLS_69({
    value: (__VLS_ctx.editForm.weight),
    min: (1),
    max: (100),
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_69));
var __VLS_67;
const __VLS_72 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_73 = __VLS_asFunctionalComponent(__VLS_72, new __VLS_72({
    label: "标签颜色",
}));
const __VLS_74 = __VLS_73({
    label: "标签颜色",
}, ...__VLS_functionalComponentArgsRest(__VLS_73));
__VLS_75.slots.default;
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ style: {} },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.input)({
    type: "color",
    ...{ style: {} },
});
(__VLS_ctx.editForm.color);
const __VLS_76 = {}.AInput;
/** @type {[typeof __VLS_components.AInput, typeof __VLS_components.aInput, ]} */ ;
// @ts-ignore
const __VLS_77 = __VLS_asFunctionalComponent(__VLS_76, new __VLS_76({
    value: (__VLS_ctx.editForm.color),
    ...{ style: {} },
    placeholder: "#1890ff",
}));
const __VLS_78 = __VLS_77({
    value: (__VLS_ctx.editForm.color),
    ...{ style: {} },
    placeholder: "#1890ff",
}, ...__VLS_functionalComponentArgsRest(__VLS_77));
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
    ...{ style: ({ display: 'inline-block', padding: '2px 12px', borderRadius: '12px', background: __VLS_ctx.editForm.color, color: '#fff', fontSize: '12px' }) },
});
var __VLS_75;
var __VLS_55;
var __VLS_47;
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
            fetchData: fetchData,
            showModal: showModal,
            saving: saving,
            isEdit: isEdit,
            editForm: editForm,
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
