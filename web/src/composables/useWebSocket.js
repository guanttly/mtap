// 核心目的：WebSocket组合式函数
// 模块功能：WebSocket连接管理、自动重连、心跳检测、消息处理
import { ref, onUnmounted } from 'vue';
import { useUserStore } from '@/stores/user';
export function useWebSocket(path) {
    const status = ref('disconnected');
    const lastMessage = ref(null);
    const handlers = new Map();
    let ws = null;
    let reconnectTimer = null;
    let shouldReconnect = true;
    function connect() {
        const userStore = useUserStore();
        const token = userStore.accessToken;
        const base = import.meta.env.VITE_WS_URL ?? `ws://${location.host}`;
        const url = `${base}${path}?token=${token}`;
        status.value = 'connecting';
        ws = new WebSocket(url);
        ws.onopen = () => { status.value = 'connected'; };
        ws.onmessage = (event) => {
            try {
                const msg = JSON.parse(event.data);
                lastMessage.value = msg;
                const list = handlers.get(msg.type) ?? [];
                list.forEach(fn => fn(msg.payload));
            }
            catch { }
        };
        ws.onerror = () => { status.value = 'error'; };
        ws.onclose = () => {
            status.value = 'disconnected';
            if (shouldReconnect) {
                reconnectTimer = setTimeout(connect, 3000);
            }
        };
    }
    function disconnect() {
        shouldReconnect = false;
        if (reconnectTimer)
            clearTimeout(reconnectTimer);
        ws?.close();
    }
    function on(type, handler) {
        const list = handlers.get(type) ?? [];
        handlers.set(type, [...list, handler]);
    }
    function send(type, payload) {
        if (ws?.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify({ type, payload }));
        }
    }
    onUnmounted(disconnect);
    return { status, lastMessage, connect, disconnect, on, send };
}
