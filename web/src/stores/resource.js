import { defineStore } from 'pinia';
import { ref } from 'vue';
import { resourceApi } from '@/api/resource';
export const useResourceStore = defineStore('resource', () => {
    const devices = ref([]);
    const examItems = ref([]);
    const slotPools = ref([]);
    const loading = ref(false);
    async function fetchDevices(params = {}) {
        loading.value = true;
        try {
            const res = await resourceApi.listDevices(params);
            devices.value = res.items;
            return res;
        }
        finally {
            loading.value = false;
        }
    }
    async function fetchExamItems(params = {}) {
        loading.value = true;
        try {
            const res = await resourceApi.listExamItems(params);
            examItems.value = res.items;
            return res;
        }
        finally {
            loading.value = false;
        }
    }
    async function fetchSlotPools() {
        loading.value = true;
        try {
            const res = await resourceApi.listSlotPools();
            slotPools.value = res.items;
            return res;
        }
        finally {
            loading.value = false;
        }
    }
    return { devices, examItems, slotPools, loading, fetchDevices, fetchExamItems, fetchSlotPools };
});
