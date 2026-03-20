/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { message, Modal } from 'ant-design-vue';
import { h, onMounted, ref } from 'vue';
import { adminApi } from '@/api/admin';
const loading = ref(false);
const users = ref([]);
const total = ref(0);
const roles = ref([]);
// 弹窗
const showCreateModal = ref(false);
const showPasswordModal = ref(false);
const saving = ref(false);
const selectedUserId = ref('');
const createForm = ref({
    username: '',
    password: '',
    real_name: '',
    role_id: '',
    department_id: '',
});
const passwordForm = ref({ new_password: '' });
const STATUS_MAP = {
    active: { label: '正常', color: 'green' },
    inactive: { label: '停用', color: 'red' },
};
const columns = [
    { title: '用户名', dataIndex: 'username', key: 'username' },
    { title: '真实姓名', dataIndex: 'real_name', key: 'real_name' },
    { title: '角色', dataIndex: 'role_name', key: 'role_name' },
    {
        title: '状态',
        dataIndex: 'status',
        key: 'status',
        customRender: ({ text }) => {
            const s = STATUS_MAP[text] ?? { label: text, color: 'default' };
            return h('a-tag', { color: s.color }, s.label);
        },
    },
    { title: '最后登录', dataIndex: 'last_login_at', key: 'last_login_at' },
    { title: '操作', key: 'actions' },
];
async function fetchData() {
    loading.value = true;
    try {
        const res = await adminApi.listUsers();
        users.value = res.items;
        total.value = res.total;
    }
    finally {
        loading.value = false;
    }
}
async function fetchRoles() {
    const res = await adminApi.listRoles();
    roles.value = res.items;
}
onMounted(() => {
    fetchData();
    fetchRoles();
});
function openCreate() {
    createForm.value = { username: '', password: '', real_name: '', role_id: '', department_id: '' };
    showCreateModal.value = true;
}
async function handleCreate() {
    if (!createForm.value.username || !createForm.value.password || !createForm.value.role_id) {
        message.warning('用户名、密码、角色为必填项');
        return;
    }
    saving.value = true;
    try {
        await adminApi.createUser(createForm.value);
        message.success('创建成功');
        showCreateModal.value = false;
        fetchData();
    }
    finally {
        saving.value = false;
    }
}
function openResetPassword(record) {
    selectedUserId.value = record.id;
    passwordForm.value = { new_password: '' };
    showPasswordModal.value = true;
}
async function handleResetPassword() {
    if (!passwordForm.value.new_password || passwordForm.value.new_password.length < 6) {
        message.warning('密码至少6位');
        return;
    }
    saving.value = true;
    try {
        await adminApi.resetPassword(selectedUserId.value, passwordForm.value.new_password);
        message.success('密码已重置');
        showPasswordModal.value = false;
    }
    finally {
        saving.value = false;
    }
}
function handleToggleStatus(record) {
    const next = record.status === 'active' ? 'inactive' : 'active';
    const label = next === 'inactive' ? '停用' : '启用';
    Modal.confirm({
        title: `确认${label}`,
        content: `确定要${label}用户「${record.username}」吗？`,
        okType: next === 'inactive' ? 'danger' : 'primary',
        onOk: async () => {
            await adminApi.updateUser(record.id, { status: next });
            message.success(`已${label}`);
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
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
    ...{ class: "ml-auto text-gray-500" },
});
(__VLS_ctx.total);
const __VLS_16 = {}.ATable;
/** @type {[typeof __VLS_components.ATable, typeof __VLS_components.aTable, typeof __VLS_components.ATable, typeof __VLS_components.aTable, ]} */ ;
// @ts-ignore
const __VLS_17 = __VLS_asFunctionalComponent(__VLS_16, new __VLS_16({
    columns: (__VLS_ctx.columns),
    dataSource: (__VLS_ctx.users),
    loading: (__VLS_ctx.loading),
    rowKey: "id",
    size: "middle",
    pagination: ({ pageSize: 20, total: __VLS_ctx.total, showTotal: (t) => `共 ${t} 条` }),
}));
const __VLS_18 = __VLS_17({
    columns: (__VLS_ctx.columns),
    dataSource: (__VLS_ctx.users),
    loading: (__VLS_ctx.loading),
    rowKey: "id",
    size: "middle",
    pagination: ({ pageSize: 20, total: __VLS_ctx.total, showTotal: (t) => `共 ${t} 条` }),
}, ...__VLS_functionalComponentArgsRest(__VLS_17));
__VLS_19.slots.default;
{
    const { bodyCell: __VLS_thisSlot } = __VLS_19.slots;
    const [{ column, record }] = __VLS_getSlotParams(__VLS_thisSlot);
    if (column.key === 'actions') {
        const __VLS_20 = {}.ASpace;
        /** @type {[typeof __VLS_components.ASpace, typeof __VLS_components.aSpace, typeof __VLS_components.ASpace, typeof __VLS_components.aSpace, ]} */ ;
        // @ts-ignore
        const __VLS_21 = __VLS_asFunctionalComponent(__VLS_20, new __VLS_20({}));
        const __VLS_22 = __VLS_21({}, ...__VLS_functionalComponentArgsRest(__VLS_21));
        __VLS_23.slots.default;
        const __VLS_24 = {}.AButton;
        /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
        // @ts-ignore
        const __VLS_25 = __VLS_asFunctionalComponent(__VLS_24, new __VLS_24({
            ...{ 'onClick': {} },
            type: "link",
            size: "small",
        }));
        const __VLS_26 = __VLS_25({
            ...{ 'onClick': {} },
            type: "link",
            size: "small",
        }, ...__VLS_functionalComponentArgsRest(__VLS_25));
        let __VLS_28;
        let __VLS_29;
        let __VLS_30;
        const __VLS_31 = {
            onClick: (...[$event]) => {
                if (!(column.key === 'actions'))
                    return;
                __VLS_ctx.openResetPassword(record);
            }
        };
        __VLS_27.slots.default;
        var __VLS_27;
        const __VLS_32 = {}.AButton;
        /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
        // @ts-ignore
        const __VLS_33 = __VLS_asFunctionalComponent(__VLS_32, new __VLS_32({
            ...{ 'onClick': {} },
            type: "link",
            size: "small",
            danger: (record.status === 'active'),
        }));
        const __VLS_34 = __VLS_33({
            ...{ 'onClick': {} },
            type: "link",
            size: "small",
            danger: (record.status === 'active'),
        }, ...__VLS_functionalComponentArgsRest(__VLS_33));
        let __VLS_36;
        let __VLS_37;
        let __VLS_38;
        const __VLS_39 = {
            onClick: (...[$event]) => {
                if (!(column.key === 'actions'))
                    return;
                __VLS_ctx.handleToggleStatus(record);
            }
        };
        __VLS_35.slots.default;
        (record.status === 'active' ? '停用' : '启用');
        var __VLS_35;
        var __VLS_23;
    }
}
var __VLS_19;
const __VLS_40 = {}.AModal;
/** @type {[typeof __VLS_components.AModal, typeof __VLS_components.aModal, typeof __VLS_components.AModal, typeof __VLS_components.aModal, ]} */ ;
// @ts-ignore
const __VLS_41 = __VLS_asFunctionalComponent(__VLS_40, new __VLS_40({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.showCreateModal),
    title: "新建用户",
    confirmLoading: (__VLS_ctx.saving),
}));
const __VLS_42 = __VLS_41({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.showCreateModal),
    title: "新建用户",
    confirmLoading: (__VLS_ctx.saving),
}, ...__VLS_functionalComponentArgsRest(__VLS_41));
let __VLS_44;
let __VLS_45;
let __VLS_46;
const __VLS_47 = {
    onOk: (__VLS_ctx.handleCreate)
};
__VLS_43.slots.default;
const __VLS_48 = {}.AForm;
/** @type {[typeof __VLS_components.AForm, typeof __VLS_components.aForm, typeof __VLS_components.AForm, typeof __VLS_components.aForm, ]} */ ;
// @ts-ignore
const __VLS_49 = __VLS_asFunctionalComponent(__VLS_48, new __VLS_48({
    model: (__VLS_ctx.createForm),
    layout: "vertical",
}));
const __VLS_50 = __VLS_49({
    model: (__VLS_ctx.createForm),
    layout: "vertical",
}, ...__VLS_functionalComponentArgsRest(__VLS_49));
__VLS_51.slots.default;
const __VLS_52 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_53 = __VLS_asFunctionalComponent(__VLS_52, new __VLS_52({
    label: "用户名",
    required: true,
}));
const __VLS_54 = __VLS_53({
    label: "用户名",
    required: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_53));
__VLS_55.slots.default;
const __VLS_56 = {}.AInput;
/** @type {[typeof __VLS_components.AInput, typeof __VLS_components.aInput, ]} */ ;
// @ts-ignore
const __VLS_57 = __VLS_asFunctionalComponent(__VLS_56, new __VLS_56({
    value: (__VLS_ctx.createForm.username),
    placeholder: "3~50个字符",
}));
const __VLS_58 = __VLS_57({
    value: (__VLS_ctx.createForm.username),
    placeholder: "3~50个字符",
}, ...__VLS_functionalComponentArgsRest(__VLS_57));
var __VLS_55;
const __VLS_60 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_61 = __VLS_asFunctionalComponent(__VLS_60, new __VLS_60({
    label: "初始密码",
    required: true,
}));
const __VLS_62 = __VLS_61({
    label: "初始密码",
    required: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_61));
__VLS_63.slots.default;
const __VLS_64 = {}.AInputPassword;
/** @type {[typeof __VLS_components.AInputPassword, typeof __VLS_components.aInputPassword, ]} */ ;
// @ts-ignore
const __VLS_65 = __VLS_asFunctionalComponent(__VLS_64, new __VLS_64({
    value: (__VLS_ctx.createForm.password),
    placeholder: "至少6位",
}));
const __VLS_66 = __VLS_65({
    value: (__VLS_ctx.createForm.password),
    placeholder: "至少6位",
}, ...__VLS_functionalComponentArgsRest(__VLS_65));
var __VLS_63;
const __VLS_68 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_69 = __VLS_asFunctionalComponent(__VLS_68, new __VLS_68({
    label: "真实姓名",
}));
const __VLS_70 = __VLS_69({
    label: "真实姓名",
}, ...__VLS_functionalComponentArgsRest(__VLS_69));
__VLS_71.slots.default;
const __VLS_72 = {}.AInput;
/** @type {[typeof __VLS_components.AInput, typeof __VLS_components.aInput, ]} */ ;
// @ts-ignore
const __VLS_73 = __VLS_asFunctionalComponent(__VLS_72, new __VLS_72({
    value: (__VLS_ctx.createForm.real_name),
}));
const __VLS_74 = __VLS_73({
    value: (__VLS_ctx.createForm.real_name),
}, ...__VLS_functionalComponentArgsRest(__VLS_73));
var __VLS_71;
const __VLS_76 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_77 = __VLS_asFunctionalComponent(__VLS_76, new __VLS_76({
    label: "角色",
    required: true,
}));
const __VLS_78 = __VLS_77({
    label: "角色",
    required: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_77));
__VLS_79.slots.default;
const __VLS_80 = {}.ASelect;
/** @type {[typeof __VLS_components.ASelect, typeof __VLS_components.aSelect, typeof __VLS_components.ASelect, typeof __VLS_components.aSelect, ]} */ ;
// @ts-ignore
const __VLS_81 = __VLS_asFunctionalComponent(__VLS_80, new __VLS_80({
    value: (__VLS_ctx.createForm.role_id),
    placeholder: "请选择角色",
}));
const __VLS_82 = __VLS_81({
    value: (__VLS_ctx.createForm.role_id),
    placeholder: "请选择角色",
}, ...__VLS_functionalComponentArgsRest(__VLS_81));
__VLS_83.slots.default;
for (const [r] of __VLS_getVForSourceType((__VLS_ctx.roles))) {
    const __VLS_84 = {}.ASelectOption;
    /** @type {[typeof __VLS_components.ASelectOption, typeof __VLS_components.aSelectOption, typeof __VLS_components.ASelectOption, typeof __VLS_components.aSelectOption, ]} */ ;
    // @ts-ignore
    const __VLS_85 = __VLS_asFunctionalComponent(__VLS_84, new __VLS_84({
        key: (r.id),
        value: (r.id),
    }));
    const __VLS_86 = __VLS_85({
        key: (r.id),
        value: (r.id),
    }, ...__VLS_functionalComponentArgsRest(__VLS_85));
    __VLS_87.slots.default;
    (r.name);
    var __VLS_87;
}
var __VLS_83;
var __VLS_79;
var __VLS_51;
var __VLS_43;
const __VLS_88 = {}.AModal;
/** @type {[typeof __VLS_components.AModal, typeof __VLS_components.aModal, typeof __VLS_components.AModal, typeof __VLS_components.aModal, ]} */ ;
// @ts-ignore
const __VLS_89 = __VLS_asFunctionalComponent(__VLS_88, new __VLS_88({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.showPasswordModal),
    title: "重置密码",
    confirmLoading: (__VLS_ctx.saving),
}));
const __VLS_90 = __VLS_89({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.showPasswordModal),
    title: "重置密码",
    confirmLoading: (__VLS_ctx.saving),
}, ...__VLS_functionalComponentArgsRest(__VLS_89));
let __VLS_92;
let __VLS_93;
let __VLS_94;
const __VLS_95 = {
    onOk: (__VLS_ctx.handleResetPassword)
};
__VLS_91.slots.default;
const __VLS_96 = {}.AForm;
/** @type {[typeof __VLS_components.AForm, typeof __VLS_components.aForm, typeof __VLS_components.AForm, typeof __VLS_components.aForm, ]} */ ;
// @ts-ignore
const __VLS_97 = __VLS_asFunctionalComponent(__VLS_96, new __VLS_96({
    model: (__VLS_ctx.passwordForm),
    layout: "vertical",
}));
const __VLS_98 = __VLS_97({
    model: (__VLS_ctx.passwordForm),
    layout: "vertical",
}, ...__VLS_functionalComponentArgsRest(__VLS_97));
__VLS_99.slots.default;
const __VLS_100 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_101 = __VLS_asFunctionalComponent(__VLS_100, new __VLS_100({
    label: "新密码",
    required: true,
}));
const __VLS_102 = __VLS_101({
    label: "新密码",
    required: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_101));
__VLS_103.slots.default;
const __VLS_104 = {}.AInputPassword;
/** @type {[typeof __VLS_components.AInputPassword, typeof __VLS_components.aInputPassword, ]} */ ;
// @ts-ignore
const __VLS_105 = __VLS_asFunctionalComponent(__VLS_104, new __VLS_104({
    value: (__VLS_ctx.passwordForm.new_password),
    placeholder: "至少6位",
}));
const __VLS_106 = __VLS_105({
    value: (__VLS_ctx.passwordForm.new_password),
    placeholder: "至少6位",
}, ...__VLS_functionalComponentArgsRest(__VLS_105));
var __VLS_103;
var __VLS_99;
var __VLS_91;
/** @type {__VLS_StyleScopedClasses['mb-4']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['gap-2']} */ ;
/** @type {__VLS_StyleScopedClasses['ml-auto']} */ ;
/** @type {__VLS_StyleScopedClasses['text-gray-500']} */ ;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            loading: loading,
            users: users,
            total: total,
            roles: roles,
            showCreateModal: showCreateModal,
            showPasswordModal: showPasswordModal,
            saving: saving,
            createForm: createForm,
            passwordForm: passwordForm,
            columns: columns,
            fetchData: fetchData,
            openCreate: openCreate,
            handleCreate: handleCreate,
            openResetPassword: openResetPassword,
            handleResetPassword: handleResetPassword,
            handleToggleStatus: handleToggleStatus,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
