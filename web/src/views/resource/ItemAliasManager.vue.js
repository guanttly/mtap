/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { onMounted, ref } from 'vue';
import { message, Modal } from 'ant-design-vue';
import { resourceApi } from '@/api/resource';
import { usePagination } from '@/composables/usePagination';
const { loading, items, pagination, fetchData, onTableChange } = usePagination(params => resourceApi.listExamItems(params));
onMounted(() => fetchData());
// 展开行管理
const expandedRowKeys = ref([]);
const newAlias = ref({});
const addingAlias = ref({});
function toggleExpand(record) {
    const idx = expandedRowKeys.value.indexOf(record.id);
    if (idx >= 0) {
        expandedRowKeys.value.splice(idx, 1);
    }
    else {
        expandedRowKeys.value.push(record.id);
    }
}
async function addAlias(item) {
    const alias = newAlias.value[item.id]?.trim();
    if (!alias) {
        message.warning('别名不能为空');
        return;
    }
    addingAlias.value[item.id] = true;
    try {
        await resourceApi.addAlias(item.id, alias);
        message.success('别名添加成功');
        newAlias.value[item.id] = '';
        fetchData();
    }
    finally {
        addingAlias.value[item.id] = false;
    }
}
function removeAlias(item, aliasName) {
    Modal.confirm({
        title: '确认删除',
        content: `确定删除别名「${aliasName}」吗？`,
        okType: 'danger',
        onOk: async () => {
            await resourceApi.deleteAlias(item.id, aliasName);
            message.success('删除成功');
            fetchData();
        },
    });
}
const columns = [
    { title: '检查项目名称', dataIndex: 'name', key: 'name' },
    { title: '标准时长(分钟)', dataIndex: 'duration_min', key: 'duration_min' },
    { title: '别名数量', key: 'alias_count', customRender: ({ record }) => record.aliases?.length ?? 0 },
    { title: '操作', key: 'actions', width: 100 },
];
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
    loading: (__VLS_ctx.loading),
}));
const __VLS_2 = __VLS_1({
    ...{ 'onClick': {} },
    loading: (__VLS_ctx.loading),
}, ...__VLS_functionalComponentArgsRest(__VLS_1));
let __VLS_4;
let __VLS_5;
let __VLS_6;
const __VLS_7 = {
    onClick: (__VLS_ctx.fetchData)
};
__VLS_3.slots.default;
var __VLS_3;
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
    ...{ class: "text-gray-400 text-sm" },
});
const __VLS_8 = {}.ATable;
/** @type {[typeof __VLS_components.ATable, typeof __VLS_components.aTable, typeof __VLS_components.ATable, typeof __VLS_components.aTable, ]} */ ;
// @ts-ignore
const __VLS_9 = __VLS_asFunctionalComponent(__VLS_8, new __VLS_8({
    ...{ 'onChange': {} },
    columns: (__VLS_ctx.columns),
    dataSource: (__VLS_ctx.items),
    loading: (__VLS_ctx.loading),
    pagination: (__VLS_ctx.pagination),
    expandedRowKeys: (__VLS_ctx.expandedRowKeys),
    rowKey: "id",
    size: "middle",
}));
const __VLS_10 = __VLS_9({
    ...{ 'onChange': {} },
    columns: (__VLS_ctx.columns),
    dataSource: (__VLS_ctx.items),
    loading: (__VLS_ctx.loading),
    pagination: (__VLS_ctx.pagination),
    expandedRowKeys: (__VLS_ctx.expandedRowKeys),
    rowKey: "id",
    size: "middle",
}, ...__VLS_functionalComponentArgsRest(__VLS_9));
let __VLS_12;
let __VLS_13;
let __VLS_14;
const __VLS_15 = {
    onChange: (__VLS_ctx.onTableChange)
};
__VLS_11.slots.default;
{
    const { bodyCell: __VLS_thisSlot } = __VLS_11.slots;
    const [{ column, record }] = __VLS_getSlotParams(__VLS_thisSlot);
    if (column.key === 'alias_count') {
        if (record.aliases?.length) {
            const __VLS_16 = {}.ATag;
            /** @type {[typeof __VLS_components.ATag, typeof __VLS_components.aTag, typeof __VLS_components.ATag, typeof __VLS_components.aTag, ]} */ ;
            // @ts-ignore
            const __VLS_17 = __VLS_asFunctionalComponent(__VLS_16, new __VLS_16({
                color: "blue",
            }));
            const __VLS_18 = __VLS_17({
                color: "blue",
            }, ...__VLS_functionalComponentArgsRest(__VLS_17));
            __VLS_19.slots.default;
            (record.aliases.length);
            var __VLS_19;
        }
        else {
            __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
                ...{ class: "text-gray-400" },
            });
        }
    }
    else if (column.key === 'actions') {
        const __VLS_20 = {}.AButton;
        /** @type {[typeof __VLS_components.AButton, typeof __VLS_components.aButton, typeof __VLS_components.AButton, typeof __VLS_components.aButton, ]} */ ;
        // @ts-ignore
        const __VLS_21 = __VLS_asFunctionalComponent(__VLS_20, new __VLS_20({
            ...{ 'onClick': {} },
            type: "link",
            size: "small",
        }));
        const __VLS_22 = __VLS_21({
            ...{ 'onClick': {} },
            type: "link",
            size: "small",
        }, ...__VLS_functionalComponentArgsRest(__VLS_21));
        let __VLS_24;
        let __VLS_25;
        let __VLS_26;
        const __VLS_27 = {
            onClick: (...[$event]) => {
                if (!!(column.key === 'alias_count'))
                    return;
                if (!(column.key === 'actions'))
                    return;
                __VLS_ctx.toggleExpand(record);
            }
        };
        __VLS_23.slots.default;
        (__VLS_ctx.expandedRowKeys.includes(record.id) ? '收起' : '管理别名');
        var __VLS_23;
    }
}
{
    const { expandedRowRender: __VLS_thisSlot } = __VLS_11.slots;
    const [{ record }] = __VLS_getSlotParams(__VLS_thisSlot);
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ style: {} },
    });
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ style: {} },
    });
    (record.name);
    __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
        ...{ style: {} },
    });
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ style: {} },
    });
    for (const [alias] of __VLS_getVForSourceType((record.aliases))) {
        const __VLS_28 = {}.ATag;
        /** @type {[typeof __VLS_components.ATag, typeof __VLS_components.aTag, typeof __VLS_components.ATag, typeof __VLS_components.aTag, ]} */ ;
        // @ts-ignore
        const __VLS_29 = __VLS_asFunctionalComponent(__VLS_28, new __VLS_28({
            ...{ 'onClose': {} },
            key: (alias.alias),
            closable: true,
            color: "blue",
            ...{ style: {} },
        }));
        const __VLS_30 = __VLS_29({
            ...{ 'onClose': {} },
            key: (alias.alias),
            closable: true,
            color: "blue",
            ...{ style: {} },
        }, ...__VLS_functionalComponentArgsRest(__VLS_29));
        let __VLS_32;
        let __VLS_33;
        let __VLS_34;
        const __VLS_35 = {
            onClose: (...[$event]) => {
                __VLS_ctx.removeAlias(record, alias.alias);
            }
        };
        __VLS_31.slots.default;
        (alias.alias);
        var __VLS_31;
    }
    if (!record.aliases?.length) {
        __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
            ...{ style: {} },
        });
    }
    const __VLS_36 = {}.AInputSearch;
    /** @type {[typeof __VLS_components.AInputSearch, typeof __VLS_components.aInputSearch, ]} */ ;
    // @ts-ignore
    const __VLS_37 = __VLS_asFunctionalComponent(__VLS_36, new __VLS_36({
        ...{ 'onSearch': {} },
        value: (__VLS_ctx.newAlias[record.id]),
        placeholder: "输入新别名，回车或点击添加",
        enterButton: "添加",
        ...{ style: {} },
        loading: (__VLS_ctx.addingAlias[record.id]),
    }));
    const __VLS_38 = __VLS_37({
        ...{ 'onSearch': {} },
        value: (__VLS_ctx.newAlias[record.id]),
        placeholder: "输入新别名，回车或点击添加",
        enterButton: "添加",
        ...{ style: {} },
        loading: (__VLS_ctx.addingAlias[record.id]),
    }, ...__VLS_functionalComponentArgsRest(__VLS_37));
    let __VLS_40;
    let __VLS_41;
    let __VLS_42;
    const __VLS_43 = {
        onSearch: (...[$event]) => {
            __VLS_ctx.addAlias(record);
        }
    };
    var __VLS_39;
}
var __VLS_11;
/** @type {__VLS_StyleScopedClasses['mb-4']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['gap-2']} */ ;
/** @type {__VLS_StyleScopedClasses['text-gray-400']} */ ;
/** @type {__VLS_StyleScopedClasses['text-sm']} */ ;
/** @type {__VLS_StyleScopedClasses['text-gray-400']} */ ;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            loading: loading,
            items: items,
            pagination: pagination,
            fetchData: fetchData,
            onTableChange: onTableChange,
            expandedRowKeys: expandedRowKeys,
            newAlias: newAlias,
            addingAlias: addingAlias,
            toggleExpand: toggleExpand,
            addAlias: addAlias,
            removeAlias: removeAlias,
            columns: columns,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
