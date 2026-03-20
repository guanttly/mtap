import { defineStore } from 'pinia';
import { ref } from 'vue';
import { analyticsApi } from '@/api/analytics';
export const useAnalyticsStore = defineStore('analytics', () => {
    const snapshot = ref(null);
    const reports = ref([]);
    const totalReports = ref(0);
    const loading = ref(false);
    async function fetchDashboard(campusId) {
        loading.value = true;
        try {
            snapshot.value = await analyticsApi.getDashboard(campusId);
            return snapshot.value;
        }
        finally {
            loading.value = false;
        }
    }
    async function fetchReports(params = {}) {
        loading.value = true;
        try {
            const res = await analyticsApi.listReports(params);
            reports.value = res.items;
            totalReports.value = res.total;
            return res;
        }
        finally {
            loading.value = false;
        }
    }
    async function generateReport(data) {
        return analyticsApi.generateReport(data);
    }
    function clearSnapshot() {
        snapshot.value = null;
    }
    return {
        snapshot,
        reports,
        totalReports,
        loading,
        fetchDashboard,
        fetchReports,
        generateReport,
        clearSnapshot,
    };
});
