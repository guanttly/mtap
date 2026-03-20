/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { message } from 'ant-design-vue';
import { onMounted, ref } from 'vue';
import { resourceApi } from '@/api/resource';
import { usePagination } from '@/composables/usePagination';
// ---- 排班列表 ----
const filterDate = ref('');
const { loading, items, pagination, fetchData, onTableChange } = usePagination(params => resourceApi.listSchedules({ ...params, date: filterDate.value || undefined }));
onMounted(fetchData);
const STATUS_COLOR = { normal: 'success', suspended: 'warning', substitute: 'processing' };
const STATUS_LABEL = { normal: '正常', suspended: '已暂停', substitute: '替班' };
const columns = [
    { title: '设备', dataIndex: 'device_name', key: 'device_name' },
    { title: '工作日期', dataIndex: 'work_date', key: 'work_date' },
    { title: '时间段', key: 'time' },
    { title: '时隙时长(分)', dataIndex: 'slot_minutes', key: 'slot_minutes' },
    { title: '状态', dataIndex: 'status', key: 'status' },
    { title: '操作', key: 'actions' },
];
async function suspend(record) {
    await resourceApi.suspendSchedule(record.device_id, record.work_date, record.start_time, record.end_time, '手动暂停');
    message.success('已暂停该排班');
    fetchData();
}
// ---- 设备列表（供下拉选择） ----
const devices = ref([]);
async function loadDevices() {
    const res = await resourceApi.listDevices({ page: 1, page_size: 200 });
    devices.value = res.items ?? [];
}
onMounted(loadDevices);
// ---- 新增单条排班 ----
const showAddModal = ref(false);
const addSaving = ref(false);
const addForm = ref({ status: 'normal' });
function openAdd() {
    addForm.value = { status: 'normal' };
    showAddModal.value = true;
}
async function handleAdd() {
    if (!addForm.value.device_id || !addForm.value.work_date) {
        message.warning('请选择设备和工作日期');
        return;
    }
    addSaving.value = true;
    try {
        await resourceApi.createSchedule(addForm.value);
        message.success('排班创建成功');
        showAddModal.value = false;
        fetchData();
    }
    finally {
        addSaving.value = false;
    }
}
// ---- 批量生成排班 ----
const showGenModal = ref(false);
const genSaving = ref(false);
const genForm = ref({
    device_id: '',
    start_date: '',
    end_date: '',
    start_time: '08:00',
    end_time: '17:00',
    slot_minutes: 30,
    skip_weekends: false,
});
function openGenerate() {
    genForm.value = { device_id: '', start_date: '', end_date: '', start_time: '08:00', end_time: '17:00', slot_minutes: 30, skip_weekends: false };
    showGenModal.value = true;
}
async function handleGenerate() {
    if (!genForm.value.device_id || !genForm.value.start_date || !genForm.value.end_date) {
        message.warning('请填写设备、开始日期和结束日期');
        return;
    }
    genSaving.value = true;
    try {
        await resourceApi.generateSchedule(genForm.value);
        message.success('批量生成成功');
        showGenModal.value = false;
        fetchData();
    }
    finally {
        genSaving.value = false;
    }
}
// ---- 替班 ----
const showSubModal = ref(false);
const subSaving = ref(false);
const subForm = ref({ source_device_id: '', target_device_id: '', date: '' });
function openSubstitute() {
    subForm.value = { source_device_id: '', target_device_id: '', date: '' };
    showSubModal.value = true;
}
async function handleSubstitute() {
    if (!subForm.value.source_device_id || !subForm.value.target_device_id || !subForm.value.date) {
        message.warning('请填写所有替班信息');
        return;
    }
    subSaving.value = true;
    try {
        await resourceApi.substituteSchedule(subForm.value.source_device_id, subForm.value.target_device_id, subForm.value.date);
        message.success('替班操作成功');
        showSubModal.value = false;
        fetchData();
    }
    finally {
        subSaving.value = false;
    }
}
debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
const __VLS_ctx = {};
let __VLS_components;
let __VLS_directives;
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "mb-4 flex flex-wrap items-center gap-2" },
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
    onClick: (__VLS_ctx.openAdd)
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
    onClick: (__VLS_ctx.openGenerate)
};
__VLS_11.slots.default;
var __VLS_11;
const __VLS_16 = {}.AButton;
/** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
// @ts-ignore
const __VLS_17 = __VLS_asFunctionalComponent(__VLS_16, new __VLS_16({
    ...{ 'onClick': {} },
}));
const __VLS_18 = __VLS_17({
    ...{ 'onClick': {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_17));
let __VLS_20;
let __VLS_21;
let __VLS_22;
const __VLS_23 = {
    onClick: (__VLS_ctx.openSubstitute)
};
__VLS_19.slots.default;
var __VLS_19;
const __VLS_24 = {}.ADatePicker;
/** @type {[typeof __VLS_components.ADatePicker, typeof __VLS_components.aDatePicker, ]} */ ;
// @ts-ignore
const __VLS_25 = __VLS_asFunctionalComponent(__VLS_24, new __VLS_24({
    ...{ 'onChange': {} },
    value: (__VLS_ctx.filterDate),
    placeholder: "按日期筛选",
    valueFormat: "YYYY-MM-DD",
    ...{ style: {} },
    allowClear: true,
}));
const __VLS_26 = __VLS_25({
    ...{ 'onChange': {} },
    value: (__VLS_ctx.filterDate),
    placeholder: "按日期筛选",
    valueFormat: "YYYY-MM-DD",
    ...{ style: {} },
    allowClear: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_25));
let __VLS_28;
let __VLS_29;
let __VLS_30;
const __VLS_31 = {
    onChange: (__VLS_ctx.fetchData)
};
var __VLS_27;
const __VLS_32 = {}.AButton;
/** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
// @ts-ignore
const __VLS_33 = __VLS_asFunctionalComponent(__VLS_32, new __VLS_32({
    ...{ 'onClick': {} },
}));
const __VLS_34 = __VLS_33({
    ...{ 'onClick': {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_33));
let __VLS_36;
let __VLS_37;
let __VLS_38;
const __VLS_39 = {
    onClick: (__VLS_ctx.fetchData)
};
__VLS_35.slots.default;
var __VLS_35;
const __VLS_40 = {}.ATable;
/** @type {[typeof __VLS_components.ATable, typeof __VLS_components.aTable, typeof __VLS_components.ATable, typeof __VLS_components.aTable, ]} */ ;
// @ts-ignore
const __VLS_41 = __VLS_asFunctionalComponent(__VLS_40, new __VLS_40({
    ...{ 'onChange': {} },
    columns: (__VLS_ctx.columns),
    dataSource: (__VLS_ctx.items),
    loading: (__VLS_ctx.loading),
    pagination: (__VLS_ctx.pagination),
    rowKey: "id",
    size: "middle",
}));
const __VLS_42 = __VLS_41({
    ...{ 'onChange': {} },
    columns: (__VLS_ctx.columns),
    dataSource: (__VLS_ctx.items),
    loading: (__VLS_ctx.loading),
    pagination: (__VLS_ctx.pagination),
    rowKey: "id",
    size: "middle",
}, ...__VLS_functionalComponentArgsRest(__VLS_41));
let __VLS_44;
let __VLS_45;
let __VLS_46;
const __VLS_47 = {
    onChange: (__VLS_ctx.onTableChange)
};
__VLS_43.slots.default;
{
    const { bodyCell: __VLS_thisSlot } = __VLS_43.slots;
    const [{ column, record }] = __VLS_getSlotParams(__VLS_thisSlot);
    if (column.key === 'time') {
        (record.start_time);
        (record.end_time);
    }
    else if (column.key === 'status') {
        const __VLS_48 = {}.ABadge;
        /** @type {[typeof __VLS_components.ABadge, typeof __VLS_components.aBadge, ]} */ ;
        // @ts-ignore
        const __VLS_49 = __VLS_asFunctionalComponent(__VLS_48, new __VLS_48({
            status: __VLS_ctx.STATUS_COLOR[record.status],
            text: (__VLS_ctx.STATUS_LABEL[record.status]),
        }));
        const __VLS_50 = __VLS_49({
            status: __VLS_ctx.STATUS_COLOR[record.status],
            text: (__VLS_ctx.STATUS_LABEL[record.status]),
        }, ...__VLS_functionalComponentArgsRest(__VLS_49));
    }
    else if (column.key === 'actions') {
        if (record.status === 'normal') {
            const __VLS_52 = {}.AButton;
            /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
            // @ts-ignore
            const __VLS_53 = __VLS_asFunctionalComponent(__VLS_52, new __VLS_52({
                ...{ 'onClick': {} },
                type: "link",
                danger: true,
                size: "small",
            }));
            const __VLS_54 = __VLS_53({
                ...{ 'onClick': {} },
                type: "link",
                danger: true,
                size: "small",
            }, ...__VLS_functionalComponentArgsRest(__VLS_53));
            let __VLS_56;
            let __VLS_57;
            let __VLS_58;
            const __VLS_59 = {
                onClick: (...[$event]) => {
                    if (!!(column.key === 'time'))
                        return;
                    if (!!(column.key === 'status'))
                        return;
                    if (!(column.key === 'actions'))
                        return;
                    if (!(record.status === 'normal'))
                        return;
                    __VLS_ctx.suspend(record);
                }
            };
            __VLS_55.slots.default;
            var __VLS_55;
        }
    }
}
var __VLS_43;
const __VLS_60 = {}.AModal;
/** @type {[typeof __VLS_components.AModal, typeof __VLS_components.aModal, typeof __VLS_components.AModal, typeof __VLS_components.aModal, ]} */ ;
// @ts-ignore
const __VLS_61 = __VLS_asFunctionalComponent(__VLS_60, new __VLS_60({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.showAddModal),
    title: "新增单条排班",
    confirmLoading: (__VLS_ctx.addSaving),
}));
const __VLS_62 = __VLS_61({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.showAddModal),
    title: "新增单条排班",
    confirmLoading: (__VLS_ctx.addSaving),
}, ...__VLS_functionalComponentArgsRest(__VLS_61));
let __VLS_64;
let __VLS_65;
let __VLS_66;
const __VLS_67 = {
    onOk: (__VLS_ctx.handleAdd)
};
__VLS_63.slots.default;
const __VLS_68 = {}.AForm;
/** @type {[typeof __VLS_components.AForm, typeof __VLS_components.aForm, typeof __VLS_components.AForm, typeof __VLS_components.aForm, ]} */ ;
// @ts-ignore
const __VLS_69 = __VLS_asFunctionalComponent(__VLS_68, new __VLS_68({
    model: (__VLS_ctx.addForm),
    layout: "vertical",
}));
const __VLS_70 = __VLS_69({
    model: (__VLS_ctx.addForm),
    layout: "vertical",
}, ...__VLS_functionalComponentArgsRest(__VLS_69));
__VLS_71.slots.default;
const __VLS_72 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_73 = __VLS_asFunctionalComponent(__VLS_72, new __VLS_72({
    label: "设备",
    required: true,
}));
const __VLS_74 = __VLS_73({
    label: "设备",
    required: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_73));
__VLS_75.slots.default;
const __VLS_76 = {}.ASelect;
/** @type {[typeof __VLS_components.ASelect, typeof __VLS_components.aSelect, typeof __VLS_components.ASelect, typeof __VLS_components.aSelect, ]} */ ;
// @ts-ignore
const __VLS_77 = __VLS_asFunctionalComponent(__VLS_76, new __VLS_76({
    value: (__VLS_ctx.addForm.device_id),
    placeholder: "请选择设备",
    showSearch: true,
    optionFilterProp: "label",
}));
const __VLS_78 = __VLS_77({
    value: (__VLS_ctx.addForm.device_id),
    placeholder: "请选择设备",
    showSearch: true,
    optionFilterProp: "label",
}, ...__VLS_functionalComponentArgsRest(__VLS_77));
__VLS_79.slots.default;
for (const [d] of __VLS_getVForSourceType((__VLS_ctx.devices))) {
    const __VLS_80 = {}.ASelectOption;
    /** @type {[typeof __VLS_components.ASelectOption, typeof __VLS_components.aSelectOption, typeof __VLS_components.ASelectOption, typeof __VLS_components.aSelectOption, ]} */ ;
    // @ts-ignore
    const __VLS_81 = __VLS_asFunctionalComponent(__VLS_80, new __VLS_80({
        key: (d.id),
        value: (d.id),
        label: (d.name),
    }));
    const __VLS_82 = __VLS_81({
        key: (d.id),
        value: (d.id),
        label: (d.name),
    }, ...__VLS_functionalComponentArgsRest(__VLS_81));
    __VLS_83.slots.default;
    (d.name);
    var __VLS_83;
}
var __VLS_79;
var __VLS_75;
const __VLS_84 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_85 = __VLS_asFunctionalComponent(__VLS_84, new __VLS_84({
    label: "工作日期",
    required: true,
}));
const __VLS_86 = __VLS_85({
    label: "工作日期",
    required: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_85));
__VLS_87.slots.default;
const __VLS_88 = {}.ADatePicker;
/** @type {[typeof __VLS_components.ADatePicker, typeof __VLS_components.aDatePicker, ]} */ ;
// @ts-ignore
const __VLS_89 = __VLS_asFunctionalComponent(__VLS_88, new __VLS_88({
    value: (__VLS_ctx.addForm.work_date),
    valueFormat: "YYYY-MM-DD",
    ...{ style: {} },
}));
const __VLS_90 = __VLS_89({
    value: (__VLS_ctx.addForm.work_date),
    valueFormat: "YYYY-MM-DD",
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_89));
var __VLS_87;
const __VLS_92 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_93 = __VLS_asFunctionalComponent(__VLS_92, new __VLS_92({
    label: "工作时间",
}));
const __VLS_94 = __VLS_93({
    label: "工作时间",
}, ...__VLS_functionalComponentArgsRest(__VLS_93));
__VLS_95.slots.default;
const __VLS_96 = {}.ATimePicker;
/** @type {[typeof __VLS_components.ATimePicker, typeof __VLS_components.aTimePicker, ]} */ ;
// @ts-ignore
const __VLS_97 = __VLS_asFunctionalComponent(__VLS_96, new __VLS_96({
    value: (__VLS_ctx.addForm.start_time),
    format: "HH:mm",
    valueFormat: "HH:mm",
    placeholder: "开始",
    ...{ style: {} },
}));
const __VLS_98 = __VLS_97({
    value: (__VLS_ctx.addForm.start_time),
    format: "HH:mm",
    valueFormat: "HH:mm",
    placeholder: "开始",
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_97));
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
    ...{ style: {} },
});
const __VLS_100 = {}.ATimePicker;
/** @type {[typeof __VLS_components.ATimePicker, typeof __VLS_components.aTimePicker, ]} */ ;
// @ts-ignore
const __VLS_101 = __VLS_asFunctionalComponent(__VLS_100, new __VLS_100({
    value: (__VLS_ctx.addForm.end_time),
    format: "HH:mm",
    valueFormat: "HH:mm",
    placeholder: "结束",
    ...{ style: {} },
}));
const __VLS_102 = __VLS_101({
    value: (__VLS_ctx.addForm.end_time),
    format: "HH:mm",
    valueFormat: "HH:mm",
    placeholder: "结束",
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_101));
var __VLS_95;
const __VLS_104 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_105 = __VLS_asFunctionalComponent(__VLS_104, new __VLS_104({
    label: "时隙时长（分钟）",
}));
const __VLS_106 = __VLS_105({
    label: "时隙时长（分钟）",
}, ...__VLS_functionalComponentArgsRest(__VLS_105));
__VLS_107.slots.default;
const __VLS_108 = {}.AInputNumber;
/** @type {[typeof __VLS_components.AInputNumber, typeof __VLS_components.aInputNumber, ]} */ ;
// @ts-ignore
const __VLS_109 = __VLS_asFunctionalComponent(__VLS_108, new __VLS_108({
    value: (__VLS_ctx.addForm.slot_minutes),
    min: (5),
    max: (120),
    ...{ style: {} },
}));
const __VLS_110 = __VLS_109({
    value: (__VLS_ctx.addForm.slot_minutes),
    min: (5),
    max: (120),
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_109));
var __VLS_107;
var __VLS_71;
var __VLS_63;
const __VLS_112 = {}.AModal;
/** @type {[typeof __VLS_components.AModal, typeof __VLS_components.aModal, typeof __VLS_components.AModal, typeof __VLS_components.aModal, ]} */ ;
// @ts-ignore
const __VLS_113 = __VLS_asFunctionalComponent(__VLS_112, new __VLS_112({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.showGenModal),
    title: "批量生成排班",
    confirmLoading: (__VLS_ctx.genSaving),
    okText: "立即生成",
}));
const __VLS_114 = __VLS_113({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.showGenModal),
    title: "批量生成排班",
    confirmLoading: (__VLS_ctx.genSaving),
    okText: "立即生成",
}, ...__VLS_functionalComponentArgsRest(__VLS_113));
let __VLS_116;
let __VLS_117;
let __VLS_118;
const __VLS_119 = {
    onOk: (__VLS_ctx.handleGenerate)
};
__VLS_115.slots.default;
const __VLS_120 = {}.AForm;
/** @type {[typeof __VLS_components.AForm, typeof __VLS_components.aForm, typeof __VLS_components.AForm, typeof __VLS_components.aForm, ]} */ ;
// @ts-ignore
const __VLS_121 = __VLS_asFunctionalComponent(__VLS_120, new __VLS_120({
    model: (__VLS_ctx.genForm),
    layout: "vertical",
}));
const __VLS_122 = __VLS_121({
    model: (__VLS_ctx.genForm),
    layout: "vertical",
}, ...__VLS_functionalComponentArgsRest(__VLS_121));
__VLS_123.slots.default;
const __VLS_124 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_125 = __VLS_asFunctionalComponent(__VLS_124, new __VLS_124({
    label: "设备",
    required: true,
}));
const __VLS_126 = __VLS_125({
    label: "设备",
    required: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_125));
__VLS_127.slots.default;
const __VLS_128 = {}.ASelect;
/** @type {[typeof __VLS_components.ASelect, typeof __VLS_components.aSelect, typeof __VLS_components.ASelect, typeof __VLS_components.aSelect, ]} */ ;
// @ts-ignore
const __VLS_129 = __VLS_asFunctionalComponent(__VLS_128, new __VLS_128({
    value: (__VLS_ctx.genForm.device_id),
    placeholder: "请选择设备",
    showSearch: true,
    optionFilterProp: "label",
}));
const __VLS_130 = __VLS_129({
    value: (__VLS_ctx.genForm.device_id),
    placeholder: "请选择设备",
    showSearch: true,
    optionFilterProp: "label",
}, ...__VLS_functionalComponentArgsRest(__VLS_129));
__VLS_131.slots.default;
for (const [d] of __VLS_getVForSourceType((__VLS_ctx.devices))) {
    const __VLS_132 = {}.ASelectOption;
    /** @type {[typeof __VLS_components.ASelectOption, typeof __VLS_components.aSelectOption, typeof __VLS_components.ASelectOption, typeof __VLS_components.aSelectOption, ]} */ ;
    // @ts-ignore
    const __VLS_133 = __VLS_asFunctionalComponent(__VLS_132, new __VLS_132({
        key: (d.id),
        value: (d.id),
        label: (d.name),
    }));
    const __VLS_134 = __VLS_133({
        key: (d.id),
        value: (d.id),
        label: (d.name),
    }, ...__VLS_functionalComponentArgsRest(__VLS_133));
    __VLS_135.slots.default;
    (d.name);
    var __VLS_135;
}
var __VLS_131;
var __VLS_127;
const __VLS_136 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_137 = __VLS_asFunctionalComponent(__VLS_136, new __VLS_136({
    label: "日期范围",
    required: true,
}));
const __VLS_138 = __VLS_137({
    label: "日期范围",
    required: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_137));
__VLS_139.slots.default;
const __VLS_140 = {}.ADatePicker;
/** @type {[typeof __VLS_components.ADatePicker, typeof __VLS_components.aDatePicker, ]} */ ;
// @ts-ignore
const __VLS_141 = __VLS_asFunctionalComponent(__VLS_140, new __VLS_140({
    value: (__VLS_ctx.genForm.start_date),
    valueFormat: "YYYY-MM-DD",
    placeholder: "开始日期",
    ...{ style: {} },
}));
const __VLS_142 = __VLS_141({
    value: (__VLS_ctx.genForm.start_date),
    valueFormat: "YYYY-MM-DD",
    placeholder: "开始日期",
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_141));
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
    ...{ style: {} },
});
const __VLS_144 = {}.ADatePicker;
/** @type {[typeof __VLS_components.ADatePicker, typeof __VLS_components.aDatePicker, ]} */ ;
// @ts-ignore
const __VLS_145 = __VLS_asFunctionalComponent(__VLS_144, new __VLS_144({
    value: (__VLS_ctx.genForm.end_date),
    valueFormat: "YYYY-MM-DD",
    placeholder: "结束日期",
    ...{ style: {} },
}));
const __VLS_146 = __VLS_145({
    value: (__VLS_ctx.genForm.end_date),
    valueFormat: "YYYY-MM-DD",
    placeholder: "结束日期",
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_145));
var __VLS_139;
const __VLS_148 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_149 = __VLS_asFunctionalComponent(__VLS_148, new __VLS_148({
    label: "每日工作时间",
}));
const __VLS_150 = __VLS_149({
    label: "每日工作时间",
}, ...__VLS_functionalComponentArgsRest(__VLS_149));
__VLS_151.slots.default;
const __VLS_152 = {}.ATimePicker;
/** @type {[typeof __VLS_components.ATimePicker, typeof __VLS_components.aTimePicker, ]} */ ;
// @ts-ignore
const __VLS_153 = __VLS_asFunctionalComponent(__VLS_152, new __VLS_152({
    value: (__VLS_ctx.genForm.start_time),
    format: "HH:mm",
    valueFormat: "HH:mm",
    placeholder: "开始",
    ...{ style: {} },
}));
const __VLS_154 = __VLS_153({
    value: (__VLS_ctx.genForm.start_time),
    format: "HH:mm",
    valueFormat: "HH:mm",
    placeholder: "开始",
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_153));
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
    ...{ style: {} },
});
const __VLS_156 = {}.ATimePicker;
/** @type {[typeof __VLS_components.ATimePicker, typeof __VLS_components.aTimePicker, ]} */ ;
// @ts-ignore
const __VLS_157 = __VLS_asFunctionalComponent(__VLS_156, new __VLS_156({
    value: (__VLS_ctx.genForm.end_time),
    format: "HH:mm",
    valueFormat: "HH:mm",
    placeholder: "结束",
    ...{ style: {} },
}));
const __VLS_158 = __VLS_157({
    value: (__VLS_ctx.genForm.end_time),
    format: "HH:mm",
    valueFormat: "HH:mm",
    placeholder: "结束",
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_157));
var __VLS_151;
const __VLS_160 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_161 = __VLS_asFunctionalComponent(__VLS_160, new __VLS_160({
    label: "时隙时长（分钟）",
}));
const __VLS_162 = __VLS_161({
    label: "时隙时长（分钟）",
}, ...__VLS_functionalComponentArgsRest(__VLS_161));
__VLS_163.slots.default;
const __VLS_164 = {}.AInputNumber;
/** @type {[typeof __VLS_components.AInputNumber, typeof __VLS_components.aInputNumber, ]} */ ;
// @ts-ignore
const __VLS_165 = __VLS_asFunctionalComponent(__VLS_164, new __VLS_164({
    value: (__VLS_ctx.genForm.slot_minutes),
    min: (5),
    max: (120),
    ...{ style: {} },
}));
const __VLS_166 = __VLS_165({
    value: (__VLS_ctx.genForm.slot_minutes),
    min: (5),
    max: (120),
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_165));
var __VLS_163;
const __VLS_168 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_169 = __VLS_asFunctionalComponent(__VLS_168, new __VLS_168({
    label: "跳过周末",
}));
const __VLS_170 = __VLS_169({
    label: "跳过周末",
}, ...__VLS_functionalComponentArgsRest(__VLS_169));
__VLS_171.slots.default;
const __VLS_172 = {}.ASwitch;
/** @type {[typeof __VLS_components.ASwitch, typeof __VLS_components.aSwitch, ]} */ ;
// @ts-ignore
const __VLS_173 = __VLS_asFunctionalComponent(__VLS_172, new __VLS_172({
    checked: (__VLS_ctx.genForm.skip_weekends),
    checkedChildren: "是",
    unCheckedChildren: "否",
}));
const __VLS_174 = __VLS_173({
    checked: (__VLS_ctx.genForm.skip_weekends),
    checkedChildren: "是",
    unCheckedChildren: "否",
}, ...__VLS_functionalComponentArgsRest(__VLS_173));
var __VLS_171;
var __VLS_123;
var __VLS_115;
const __VLS_176 = {}.AModal;
/** @type {[typeof __VLS_components.AModal, typeof __VLS_components.aModal, typeof __VLS_components.AModal, typeof __VLS_components.aModal, ]} */ ;
// @ts-ignore
const __VLS_177 = __VLS_asFunctionalComponent(__VLS_176, new __VLS_176({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.showSubModal),
    title: "替班管理",
    confirmLoading: (__VLS_ctx.subSaving),
    okText: "确认替班",
}));
const __VLS_178 = __VLS_177({
    ...{ 'onOk': {} },
    open: (__VLS_ctx.showSubModal),
    title: "替班管理",
    confirmLoading: (__VLS_ctx.subSaving),
    okText: "确认替班",
}, ...__VLS_functionalComponentArgsRest(__VLS_177));
let __VLS_180;
let __VLS_181;
let __VLS_182;
const __VLS_183 = {
    onOk: (__VLS_ctx.handleSubstitute)
};
__VLS_179.slots.default;
const __VLS_184 = {}.AAlert;
/** @type {[typeof __VLS_components.AAlert, typeof __VLS_components.aAlert, ]} */ ;
// @ts-ignore
const __VLS_185 = __VLS_asFunctionalComponent(__VLS_184, new __VLS_184({
    message: "替班将把指定日期原设备的所有预约转移给替班设备",
    type: "info",
    showIcon: true,
    ...{ class: "mb-4" },
}));
const __VLS_186 = __VLS_185({
    message: "替班将把指定日期原设备的所有预约转移给替班设备",
    type: "info",
    showIcon: true,
    ...{ class: "mb-4" },
}, ...__VLS_functionalComponentArgsRest(__VLS_185));
const __VLS_188 = {}.AForm;
/** @type {[typeof __VLS_components.AForm, typeof __VLS_components.aForm, typeof __VLS_components.AForm, typeof __VLS_components.aForm, ]} */ ;
// @ts-ignore
const __VLS_189 = __VLS_asFunctionalComponent(__VLS_188, new __VLS_188({
    model: (__VLS_ctx.subForm),
    layout: "vertical",
}));
const __VLS_190 = __VLS_189({
    model: (__VLS_ctx.subForm),
    layout: "vertical",
}, ...__VLS_functionalComponentArgsRest(__VLS_189));
__VLS_191.slots.default;
const __VLS_192 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_193 = __VLS_asFunctionalComponent(__VLS_192, new __VLS_192({
    label: "原设备",
    required: true,
}));
const __VLS_194 = __VLS_193({
    label: "原设备",
    required: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_193));
__VLS_195.slots.default;
const __VLS_196 = {}.ASelect;
/** @type {[typeof __VLS_components.ASelect, typeof __VLS_components.aSelect, typeof __VLS_components.ASelect, typeof __VLS_components.aSelect, ]} */ ;
// @ts-ignore
const __VLS_197 = __VLS_asFunctionalComponent(__VLS_196, new __VLS_196({
    value: (__VLS_ctx.subForm.source_device_id),
    placeholder: "选择原设备",
    showSearch: true,
    optionFilterProp: "label",
}));
const __VLS_198 = __VLS_197({
    value: (__VLS_ctx.subForm.source_device_id),
    placeholder: "选择原设备",
    showSearch: true,
    optionFilterProp: "label",
}, ...__VLS_functionalComponentArgsRest(__VLS_197));
__VLS_199.slots.default;
for (const [d] of __VLS_getVForSourceType((__VLS_ctx.devices))) {
    const __VLS_200 = {}.ASelectOption;
    /** @type {[typeof __VLS_components.ASelectOption, typeof __VLS_components.aSelectOption, typeof __VLS_components.ASelectOption, typeof __VLS_components.aSelectOption, ]} */ ;
    // @ts-ignore
    const __VLS_201 = __VLS_asFunctionalComponent(__VLS_200, new __VLS_200({
        key: (d.id),
        value: (d.id),
        label: (d.name),
    }));
    const __VLS_202 = __VLS_201({
        key: (d.id),
        value: (d.id),
        label: (d.name),
    }, ...__VLS_functionalComponentArgsRest(__VLS_201));
    __VLS_203.slots.default;
    (d.name);
    var __VLS_203;
}
var __VLS_199;
var __VLS_195;
const __VLS_204 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_205 = __VLS_asFunctionalComponent(__VLS_204, new __VLS_204({
    label: "替班设备",
    required: true,
}));
const __VLS_206 = __VLS_205({
    label: "替班设备",
    required: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_205));
__VLS_207.slots.default;
const __VLS_208 = {}.ASelect;
/** @type {[typeof __VLS_components.ASelect, typeof __VLS_components.aSelect, typeof __VLS_components.ASelect, typeof __VLS_components.aSelect, ]} */ ;
// @ts-ignore
const __VLS_209 = __VLS_asFunctionalComponent(__VLS_208, new __VLS_208({
    value: (__VLS_ctx.subForm.target_device_id),
    placeholder: "选择替班设备",
    showSearch: true,
    optionFilterProp: "label",
}));
const __VLS_210 = __VLS_209({
    value: (__VLS_ctx.subForm.target_device_id),
    placeholder: "选择替班设备",
    showSearch: true,
    optionFilterProp: "label",
}, ...__VLS_functionalComponentArgsRest(__VLS_209));
__VLS_211.slots.default;
for (const [d] of __VLS_getVForSourceType((__VLS_ctx.devices))) {
    const __VLS_212 = {}.ASelectOption;
    /** @type {[typeof __VLS_components.ASelectOption, typeof __VLS_components.aSelectOption, typeof __VLS_components.ASelectOption, typeof __VLS_components.aSelectOption, ]} */ ;
    // @ts-ignore
    const __VLS_213 = __VLS_asFunctionalComponent(__VLS_212, new __VLS_212({
        key: (d.id),
        value: (d.id),
        label: (d.name),
    }));
    const __VLS_214 = __VLS_213({
        key: (d.id),
        value: (d.id),
        label: (d.name),
    }, ...__VLS_functionalComponentArgsRest(__VLS_213));
    __VLS_215.slots.default;
    (d.name);
    var __VLS_215;
}
var __VLS_211;
var __VLS_207;
const __VLS_216 = {}.AFormItem;
/** @type {[typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, typeof __VLS_components.AFormItem, typeof __VLS_components.aFormItem, ]} */ ;
// @ts-ignore
const __VLS_217 = __VLS_asFunctionalComponent(__VLS_216, new __VLS_216({
    label: "替班日期",
    required: true,
}));
const __VLS_218 = __VLS_217({
    label: "替班日期",
    required: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_217));
__VLS_219.slots.default;
const __VLS_220 = {}.ADatePicker;
/** @type {[typeof __VLS_components.ADatePicker, typeof __VLS_components.aDatePicker, ]} */ ;
// @ts-ignore
const __VLS_221 = __VLS_asFunctionalComponent(__VLS_220, new __VLS_220({
    value: (__VLS_ctx.subForm.date),
    valueFormat: "YYYY-MM-DD",
    ...{ style: {} },
}));
const __VLS_222 = __VLS_221({
    value: (__VLS_ctx.subForm.date),
    valueFormat: "YYYY-MM-DD",
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_221));
var __VLS_219;
var __VLS_191;
var __VLS_179;
/** @type {__VLS_StyleScopedClasses['mb-4']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['flex-wrap']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['gap-2']} */ ;
/** @type {__VLS_StyleScopedClasses['mb-4']} */ ;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            filterDate: filterDate,
            loading: loading,
            items: items,
            pagination: pagination,
            fetchData: fetchData,
            onTableChange: onTableChange,
            STATUS_COLOR: STATUS_COLOR,
            STATUS_LABEL: STATUS_LABEL,
            columns: columns,
            suspend: suspend,
            devices: devices,
            showAddModal: showAddModal,
            addSaving: addSaving,
            addForm: addForm,
            openAdd: openAdd,
            handleAdd: handleAdd,
            showGenModal: showGenModal,
            genSaving: genSaving,
            genForm: genForm,
            openGenerate: openGenerate,
            handleGenerate: handleGenerate,
            showSubModal: showSubModal,
            subSaving: subSaving,
            subForm: subForm,
            openSubstitute: openSubstitute,
            handleSubstitute: handleSubstitute,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
