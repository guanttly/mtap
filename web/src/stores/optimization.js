import { defineStore } from 'pinia';
import { ref } from 'vue';
import { optimizationApi } from '@/api/optimization';
export const useOptimizationStore = defineStore('optimization', () => {
    const strategies = ref([]);
    const alerts = ref([]);
    const metrics = ref([]);
    const total = ref(0);
    const loading = ref(false);
    async function fetchStrategies(params = {}) {
        loading.value = true;
        try {
            const res = await optimizationApi.listStrategies(params);
            strategies.value = res.items;
            total.value = res.total;
            return res;
        }
        finally {
            loading.value = false;
        }
    }
    async function fetchAlerts(params = {}) {
        loading.value = true;
        try {
            const res = await optimizationApi.listAlerts(params);
            alerts.value = res.items;
            return res;
        }
        finally {
            loading.value = false;
        }
    }
    async function fetchMetrics() {
        loading.value = true;
        try {
            const res = await optimizationApi.listMetrics();
            metrics.value = res.items;
            return res;
        }
        finally {
            loading.value = false;
        }
    }
    return { strategies, alerts, metrics, total, loading, fetchStrategies, fetchAlerts, fetchMetrics };
});
