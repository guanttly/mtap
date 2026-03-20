/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { message, Modal } from 'ant-design-vue';
import { onMounted, ref } from 'vue';
import { adminApi } from '@/api/admin';
const loading = ref(false);
const roles = ref([]);
const showModal = ref(false);
const saving = ref(false);
const editingRole = ref(null);
const createForm = ref({ name: '', permissions: [] });
const editPermissions = ref([]);
// 预定义权限项（与后端角色权限对应）
const PERMISSION_OPTIONS = [
    { label: '规则引擎 - 读', value: 'rule:read' },
    { label: '规则引擎 - 写', value: 'rule:write' },
    { label: '规则引擎 - 全部', value: 'rule:*' },
    { label: '资源管理 - 读', value: 'resource:read' },
    { label: '资源管理 - 写', value: 'resource:write' },
    { label: '资源管理 - 全部', value: 'resource:*' },
    { label: '预约服务 - 读', value: 'appt:read' },
    { label: '预约服务 - 写', value: 'appt:write' },
    { label: '预约服务 - 全部', value: 'appt:*' },
    { label: '分诊执行 - 全部', value: 'triage:*' },
    { label: '统计分析 - 读', value: 'analytics:read' },
    { label: '效能优化 - 全部', value: 'optimization:*' },
    { label: '系统管理 - 全部', value: 'admin:*' },
    { label: '全部权限', value: '*' },
];
const columns = [
    { title: '角色名', dataIndex: 'name', key: 'name' },
    { title: '权限', dataIndex: 'permissions', key: 'permissions' },
    { title: '类型', dataIndex: 'is_preset', key: 'is_preset' },
    { title: '操作', key: 'actions' },
];
async function fetchData() {
    loading.value = true;
    try {
        const res = await adminApi.listRoles();
        roles.value = res.items;
    }
    finally {
        loading.value = false;
    }
}
onMounted(fetchData);
function openCreate() {
    editingRole.value = null;
    createForm.value = { name: '', permissions: [] };
    showModal.value = true;
}
function openEdit(record) {
    editingRole.value = record;
    editPermissions.value = [...record.permissions];
    showModal.value = true;
}
async function handleSave() {
    saving.value = true;
    try {
        if (editingRole.value) {
            await adminApi.updateRole(editingRole.value.id, editPermissions.value);
            message.success('权限已更新');
        }
        else {
            if (!createForm.value.name) {
                message.warning('请输入角色名称');
                return;
            }
            await adminApi.createRole(createForm.value);
            message.success('角色已创建');
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
        message.warning('预置角色不可删除');
        return;
    }
    Modal.confirm({
        title: '确认删除',
        content: `确定删除角色「${record.name}」吗？请确保该角色下无用户。`,
        okType: 'danger',
        onOk: async () => {
            await adminApi.deleteRole(record.id);
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
    dataSource: (__VLS_ctx.roles),
    loading: (__VLS_ctx.loading),
    rowKey: "id",
    size: "middle",
    pagination: (false),
}));
const __VLS_18 = __VLS_17({
    columns: (__VLS_ctx.columns),
    dataSource: (__VLS_ctx.roles),
    loading: (__VLS_ctx.loading),
    rowKey: "id",
    size: "middle",
    pagination: (false),
}, ...__VLS_functionalComponentArgsRest(__VLS_17));
__VLS_19.slots.default;
{
    const { bodyCell: __VLS_thisSlot } = __VLS_19.slots;
    const [{ column, record }] = __VLS_getSlotParams(__VLS_thisSlot);
    if (column.key === 'permissions') {
        const __VLS_20 = {}.ASpace;
        /** @type {[typeof __VLS_components.ASpace, typeof __VLS_components.aSpace, typeof __VLS_components.ASpace, typeof __VLS_components.aSpace, ]} */ ;
        // @ts-ignore
        const __VLS_21 = __VLS_asFunctionalComponent(__VLS_20, new __VLS_20({
            wrap: true,
        }));
        const __VLS_22 = __VLS_21({
            wrap: true,
        }, ...__VLS_functionalComponentArgsRest(__VLS_21));
        __VLS_23.slots.default;
        for (const [p] of __VLS_getVForSourceType((record.permissions))) {
            const __VLS_24 = {}.ATag;
            /** @type {[typeof __VLS_components.ATag, typeof __VLS_components.aTag, typeof __VLS_components.ATag, typeof __VLS_components.aTag, ]} */ ;
            // @ts-ignore
            const __VLS_25 = __VLS_asFunctionalComponent(__VLS_24, new __VLS_24({
                key: (p),
                color: "blue",
            }));
            const __VLS_26 = __VLS_25({
                key: (p),
                color: "blue",
            }, ...__VLS_functionalComponentArgsRest(__VLS_25));
            __VLS_27.slots.default;
            (p);
            var __VLS_27;
        }
        if (record.permissions.length === 0) {
            __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
                ...{ class: "text-gray-400" },
            });
        }
        var __VLS_23;
    }
    else if (column.key === 'is_preset') {
        const __VLS_28 = {}.ATag;
        /** @type {[typeof __VLS_components.ATag, typeof __VLS_components.aTag, typeof __VLS_components.ATag, typeof __VLS_components.aTag, ]} */ ;
        // @ts-ignore
        const __VLS_29 = __VLS_asFunctionalComponent(__VLS_28, new __VLS_28({
            color: (record.is_preset ? 'purple' : 'default'),
        }));
        const __VLS_30 = __VLS_29({
            color: (record.is_preset ? 'purple' : 'default'),
        }, ...__VLS_functionalComponentArgsRest(__VLS_29));
        __VLS_31.slots.default;
        (record.is_preset ? '系统预置' : '自定义');
        var __VLS_31;
    }
    else if (column.key === 'actions') {
        const __VLS_32 = {}.ASpace;
        /** @type {[typeof __VLS_components.ASpace, typeof __VLS_components.aSpace, typeof __VLS_components.ASpace, typeof __VLS_components.aSpace, ]} */ ;
        // @ts-ignore
        const __VLS_33 = __VLS_asFunctionalComponent(__VLS_32, new __VLS_32({}));
        const __VLS_34 = __VLS_33({}, ...__VLS_functionalComponentArgsRest(__VLS_33));
        __VLS_35.slots.default;
        const __VLS_36 = {}.AButton;
        /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
        // @ts-ignore
        const __VLS_37 = __VLS_asFunctionalComponent(__VLS_36, new __VLS_36({
            ...{ 'onClick': {} },
            type: "link",
            size: "small",
            disabled: (record.is_preset),
        }));
        const __VLS_38 = __VLS_37({
            ...{ 'onClick': {} },
            type: "link",
            size: "small",
            disabled: (record.is_preset),
        }, ...__VLS_functionalComponentArgsRest(__VLS_37));
        let __VLS_40;
        let __VLS_41;
        let __VLS_42;
        const __VLS_43 = {
            onClick: (...[$event]) => {
                if (!!(column.key === 'permissions'))
                    return;
                if (!!(column.key === 'is_preset'))
                    return;
                if (!(column.key === 'actions'))
                    return;
                __VLS_ctx.openEdit(record);
            }
        };
        __VLS_39.slots.default;
        var __VLS_39;
        const __VLS_44 = {}.AButton;
        /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
        // @ts-ignore
        const __VLS_45 = __VLS_asFunctionalComponent(__VLS_44, new __VLS_44({
            ...{ 'onClick': {} },
            type: "link",
            danger: true,
            size: "small",
            disabled: (record.is_preset),
        }));
        const __VLS_46 = __VLS_45({
            ...{ 'onClick': {} },
            type: "link",
            danger: true,
            size: "small",
            disabled: (record.is_preset),
        }, ...__VLS_functionalComponentArgsRest(__VLS_45));
        let __VLS_48;
        let __VLS_49;
        let __VLS_50;
        const __VLS_51 = {
            onClick: (...[$event]) => {
                if (!!(column.key === 'permissions'))
                    return;
                if (!!(column.key === 'is_preset'))
                    return;
                if (!(column.key === 'actions'))
                    return;
                __VLS_ctx.handleDelete(record);
            }
        };
        __VLS_47.slots.default;
        var __VLS_47;
        var __VLS_35;
    }
}
var __VLS_19;
const __VLS_52 = {}.AModal;
/** @type {[typeof __VLS_components.AModal, typeof __VLS_components.aModal, typeof __VLS_components.AModal, typeof __VLS_components.aModal, ]} */ ;
// @ts-ignore
const __VLS_53 = __VLS_asFunctionalComponent(__VLS_52, new __VLS_52({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.showModal),
    title: (__VLS_ctx.editingRole ? `编辑权限 - ${__VLS_ctx.editingRole.name}` : '新增角色'),
    confirmLoading: (__VLS_ctx.saving),
}));
const __VLS_54 = __VLS_53({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.showModal),
    title: (__VLS_ctx.editingRole ? `编辑权限 - ${__VLS_ctx.editingRole.name}` : '新增角色'),
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
    layout: "vertical",
}));
const __VLS_62 = __VLS_61({
    layout: "vertical",
}, ...__VLS_functionalComponentArgsRest(__VLS_61));
__VLS_63.slots.default;
if (!__VLS_ctx.editingRole) {
    const __VLS_64 = {}.AFormItem;
    /** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
    // @ts-ignore
    const __VLS_65 = __VLS_asFunctionalComponent(__VLS_64, new __VLS_64({
        label: "角色名称",
        required: true,
    }));
    const __VLS_66 = __VLS_65({
        label: "角色名称",
        required: true,
    }, ...__VLS_functionalComponentArgsRest(__VLS_65));
    __VLS_67.slots.default;
    const __VLS_68 = {}.AInput;
    /** @type {[typeof __VLS_components.AInput, typeof __VLS_components.aInput, ]} */ ;
    // @ts-ignore
    const __VLS_69 = __VLS_asFunctionalComponent(__VLS_68, new __VLS_68({
        value: (__VLS_ctx.createForm.name),
        placeholder: "2~30个字符",
    }));
    const __VLS_70 = __VLS_69({
        value: (__VLS_ctx.createForm.name),
        placeholder: "2~30个字符",
    }, ...__VLS_functionalComponentArgsRest(__VLS_69));
    var __VLS_67;
    const __VLS_72 = {}.AFormItem;
    /** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
    // @ts-ignore
    const __VLS_73 = __VLS_asFunctionalComponent(__VLS_72, new __VLS_72({
        label: "初始权限",
    }));
    const __VLS_74 = __VLS_73({
        label: "初始权限",
    }, ...__VLS_functionalComponentArgsRest(__VLS_73));
    __VLS_75.slots.default;
    const __VLS_76 = {}.ACheckboxGroup;
    /** @type {[typeof __VLS_components.ACheckboxGroup, typeof __VLS_components.aCheckboxGroup, ]} */ ;
    // @ts-ignore
    const __VLS_77 = __VLS_asFunctionalComponent(__VLS_76, new __VLS_76({
        value: (__VLS_ctx.createForm.permissions),
        options: (__VLS_ctx.PERMISSION_OPTIONS),
        ...{ style: {} },
    }));
    const __VLS_78 = __VLS_77({
        value: (__VLS_ctx.createForm.permissions),
        options: (__VLS_ctx.PERMISSION_OPTIONS),
        ...{ style: {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_77));
    var __VLS_75;
}
else {
    const __VLS_80 = {}.AFormItem;
    /** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
    // @ts-ignore
    const __VLS_81 = __VLS_asFunctionalComponent(__VLS_80, new __VLS_80({
        label: "权限配置",
    }));
    const __VLS_82 = __VLS_81({
        label: "权限配置",
    }, ...__VLS_functionalComponentArgsRest(__VLS_81));
    __VLS_83.slots.default;
    const __VLS_84 = {}.ACheckboxGroup;
    /** @type {[typeof __VLS_components.ACheckboxGroup, typeof __VLS_components.aCheckboxGroup, ]} */ ;
    // @ts-ignore
    const __VLS_85 = __VLS_asFunctionalComponent(__VLS_84, new __VLS_84({
        value: (__VLS_ctx.editPermissions),
        options: (__VLS_ctx.PERMISSION_OPTIONS),
        ...{ style: {} },
    }));
    const __VLS_86 = __VLS_85({
        value: (__VLS_ctx.editPermissions),
        options: (__VLS_ctx.PERMISSION_OPTIONS),
        ...{ style: {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_85));
    var __VLS_83;
}
var __VLS_63;
var __VLS_55;
/** @type {__VLS_StyleScopedClasses['mb-4']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['gap-2']} */ ;
/** @type {__VLS_StyleScopedClasses['text-gray-400']} */ ;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            loading: loading,
            roles: roles,
            showModal: showModal,
            saving: saving,
            editingRole: editingRole,
            createForm: createForm,
            editPermissions: editPermissions,
            PERMISSION_OPTIONS: PERMISSION_OPTIONS,
            columns: columns,
            fetchData: fetchData,
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
