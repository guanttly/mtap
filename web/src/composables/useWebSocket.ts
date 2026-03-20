// 核心目的：WebSocket组合式函数
// 模块功能：WebSocket连接管理、自动重连、心跳检测、消息处理
import { ref, onUnmounted } from 'vue'
import { useUserStore } from '@/stores/user'

export type WSStatus = 'connecting' | 'connected' | 'disconnected' | 'error'

export interface WSMessage<T = unknown> {
  type: string
  payload: T
}

export function useWebSocket<T = unknown>(path: string) {
  const status = ref<WSStatus>('disconnected')
  const lastMessage = ref<WSMessage<T> | null>(null)
  const handlers = new Map<string, ((payload: T) => void)[]>()

  let ws: WebSocket | null = null
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null
  let shouldReconnect = true
  let reconnectCount = 0
  const MAX_RECONNECT = 5  // 最多重连 5 次，防止服务不可用时刷屏控制台

  function connect() {
    const userStore = useUserStore()
    const token = userStore.accessToken
    const base = (import.meta.env.VITE_WS_URL as string) ?? `ws://${location.host}`
    const url = `${base}${path}?token=${token}`

    status.value = 'connecting'
    ws = new WebSocket(url)

    ws.onopen = () => {
      status.value = 'connected'
      reconnectCount = 0  // 连接成功后重置计数器
    }

    ws.onmessage = (event) => {
      try {
        const msg: WSMessage<T> = JSON.parse(event.data as string)
        lastMessage.value = msg
        const list = handlers.get(msg.type) ?? []
        list.forEach(fn => fn(msg.payload))
      }
      catch {}
    }

    ws.onerror = () => { status.value = 'error' }

    ws.onclose = () => {
      status.value = 'disconnected'
      if (shouldReconnect && reconnectCount < MAX_RECONNECT) {
        reconnectCount++
        reconnectTimer = setTimeout(connect, 3000 * reconnectCount) // 退刯重连
      }
    }
  }

  function disconnect() {
    shouldReconnect = false
    if (reconnectTimer) clearTimeout(reconnectTimer)
    ws?.close()
  }

  function on(type: string, handler: (payload: T) => void) {
    const list = handlers.get(type) ?? []
    handlers.set(type, [...list, handler])
  }

  function send(type: string, payload: unknown) {
    if (ws?.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ type, payload }))
    }
  }

  onUnmounted(disconnect)

  return { status, lastMessage, connect, disconnect, on, send }
}
