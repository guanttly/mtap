/// <reference types="../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { ref, onMounted, onUnmounted } from 'vue';
import { triageApi } from '@/api/triage';
import { useWebSocket } from '@/composables/useWebSocket';
const rooms = ['room-001', 'room-002', 'room-003'];
const statusMap = ref({});
const currentCall = ref(null);
const now = ref(new Date());
const { connect, on } = useWebSocket('/ws/triage');
async function fetchAll() {
    for (const roomId of rooms) {
        try {
            statusMap.value[roomId] = await triageApi.getQueueStatus(roomId);
        }
        catch { }
    }
}
let clockTimer;
let refreshTimer;
onMounted(() => {
    fetchAll();
    clockTimer = setInterval(() => { now.value = new Date(); }, 1000);
    refreshTimer = setInterval(fetchAll, 20000);
    connect();
    on('call', (payload) => { currentCall.value = payload; });
});
onUnmounted(() => {
    clearInterval(clockTimer);
    clearInterval(refreshTimer);
});
function timeStr() {
    return now.value.toLocaleTimeString('zh-CN', { hour12: false });
}
function dateStr() {
    return now.value.toLocaleDateString('zh-CN', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' });
}
debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
const __VLS_ctx = {};
let __VLS_components;
let __VLS_directives;
// CSS variable injection 
// CSS variable injection end 
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ style: {} },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ style: {} },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ style: {} },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ style: {} },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ style: {} },
});
(__VLS_ctx.timeStr());
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ style: {} },
});
(__VLS_ctx.dateStr());
if (__VLS_ctx.currentCall) {
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ style: {} },
    });
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ style: {} },
    });
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ style: {} },
    });
    (__VLS_ctx.currentCall.patient_name_masked);
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ style: {} },
    });
    (__VLS_ctx.currentCall.room_name);
    (__VLS_ctx.currentCall.queue_number);
}
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ style: {} },
});
for (const [roomId] of __VLS_getVForSourceType((__VLS_ctx.rooms))) {
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        key: (roomId),
        ...{ style: {} },
    });
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ style: {} },
    });
    (__VLS_ctx.statusMap[roomId]?.room_name ?? roomId);
    if (__VLS_ctx.statusMap[roomId]) {
        __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({});
        __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
            ...{ style: {} },
        });
        __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
            ...{ style: {} },
        });
        __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
            ...{ style: {} },
        });
        (__VLS_ctx.statusMap[roomId].waiting_count);
        __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
            ...{ style: {} },
        });
        __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
            ...{ style: {} },
        });
        __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
            ...{ style: {} },
        });
        (__VLS_ctx.statusMap[roomId].entries?.filter(e => e.status === 'completed').length ?? 0);
        __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
            ...{ style: {} },
        });
        for (const [entry, idx] of __VLS_getVForSourceType((__VLS_ctx.statusMap[roomId].entries?.slice(0, 5)))) {
            __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
                key: (entry.id),
                ...{ style: {} },
            });
            __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
                ...{ style: ({ fontWeight: '700', fontSize: '20px', color: idx === 0 ? '#40a9ff' : 'rgba(255,255,255,.5)', width: '24px' }) },
            });
            (idx + 1);
            __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
                ...{ style: ({ flex: 1, opacity: idx === 0 ? 1 : 0.6 }) },
            });
            (entry.patient_name_masked);
        }
    }
    else {
        const __VLS_0 = {}.ASpin;
        /** @type {[typeof __VLS_components.ASpin, typeof __VLS_components.aSpin, ]} */ ;
        // @ts-ignore
        const __VLS_1 = __VLS_asFunctionalComponent(__VLS_0, new __VLS_0({
            ...{ style: {} },
        }));
        const __VLS_2 = __VLS_1({
            ...{ style: {} },
        }, ...__VLS_functionalComponentArgsRest(__VLS_1));
    }
}
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            rooms: rooms,
            statusMap: statusMap,
            currentCall: currentCall,
            timeStr: timeStr,
            dateStr: dateStr,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
