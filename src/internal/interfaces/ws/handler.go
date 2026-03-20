// Package ws 接口层 - WebSocket消息处理
// 核心目的：处理WebSocket连接与消息推送
// 模块功能：
//   - ScreenHandler: 分诊大屏推送（候诊队列 + 呼叫信息，10秒刷新）
//   - DashboardHandler: 监控大屏推送（号源占用率 + 设备状态 + 等待时长）
//   - 连接断开重连机制（每5秒重连）
package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// upgrader WebSocket 升级器
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境需校验 Origin
	},
}

// Handler WebSocket HTTP 处理器
type Handler struct {
	hub *Hub
}

// NewHandler 创建 WebSocket 处理器
func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

// RegisterRoutes 注册 WebSocket 路由
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	ws := r.Group("/ws")
	ws.GET("/triage", h.TriageScreenAll)      // 分诊大屏全局频道
	ws.GET("/triage/:roomId", h.TriageScreen) // 指定房间
	ws.GET("/dashboard", h.Dashboard)
}

// TriageScreenAll 分诊大屏全局 WebSocket 端点（订阅所有分诊事件）
// GET /ws/triage
func (h *Handler) TriageScreenAll(c *gin.Context) {
	h.serveWS(c.Writer, c.Request, ChannelTriage)
}

// TriageScreen 分诊大屏 WebSocket 端点
// GET /ws/triage/:roomId
func (h *Handler) TriageScreen(c *gin.Context) {
	roomID := c.Param("roomId")
	channel := ChannelTriage + ":" + roomID
	h.serveWS(c.Writer, c.Request, channel)
}

// Dashboard 监控大屏 WebSocket 端点
// GET /ws/dashboard
func (h *Handler) Dashboard(c *gin.Context) {
	h.serveWS(c.Writer, c.Request, ChannelDashboard)
}

// serveWS 升级连接并注册到 Hub
func (h *Handler) serveWS(w http.ResponseWriter, r *http.Request, channel string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS] upgrade error: %v", err)
		return
	}
	client := &Client{
		hub:     h.hub,
		conn:    conn,
		send:    make(chan []byte, 256),
		channel: channel,
		id:      uuid.New().String(),
	}
	h.hub.register <- client

	go client.writePump()
	go client.readPump()
}

// BroadcastTriageUpdate 广播分诊大屏更新（由分诊模块调用）
func (h *Handler) BroadcastTriageUpdate(roomID string, payload any) {
	data, err := json.Marshal(map[string]any{
		"type":      "triage_update",
		"room_id":   roomID,
		"timestamp": time.Now().Unix(),
		"data":      payload,
	})
	if err != nil {
		return
	}
	h.hub.Broadcast(ChannelTriage+":"+roomID, data)
}

// BroadcastDashboardUpdate 广播监控大屏更新（由 analytics 模块调用）
func (h *Handler) BroadcastDashboardUpdate(payload any) {
	data, err := json.Marshal(map[string]any{
		"type":      "dashboard_update",
		"timestamp": time.Now().Unix(),
		"data":      payload,
	})
	if err != nil {
		return
	}
	h.hub.Broadcast(ChannelDashboard, data)
}
