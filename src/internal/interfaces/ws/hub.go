// Package ws 接口层 - WebSocket连接管理
// 核心目的：管理WebSocket连接的生命周期
// 模块功能：
//   - Hub: 连接池管理（注册/注销/广播）
//   - 按频道分组（分诊大屏频道、监控大屏频道）
//   - 心跳检测与自动重连
package ws

import (
	"log"
	"sync"
	"time"
)

const (
	// ChannelTriage 分诊大屏频道
	ChannelTriage = "triage"
	// ChannelDashboard 监控大屏频道
	ChannelDashboard = "dashboard"

	// 心跳配置
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512 * 1024 // 512KB
)

// Client 单个 WebSocket 连接
type Client struct {
	hub     *Hub
	conn    WSConn // 连接接口
	send    chan []byte
	channel string
	id      string
}

// WSConn WebSocket 连接抽象（便于测试 mock）
type WSConn interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	SetReadLimit(limit int64)
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
	SetPongHandler(h func(appData string) error)
	Close() error
}

// Hub WebSocket 连接集线器
type Hub struct {
	mu         sync.RWMutex
	clients    map[string]map[*Client]struct{} // channel -> clients
	register   chan *Client
	unregister chan *Client
	broadcast  chan broadcastMsg
}

type broadcastMsg struct {
	channel string
	data    []byte
}

// NewHub 创建 Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]struct{}),
		register:   make(chan *Client, 64),
		unregister: make(chan *Client, 64),
		broadcast:  make(chan broadcastMsg, 256),
	}
}

// Run 启动 Hub 事件循环（需在 goroutine 中运行）
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.channel] == nil {
				h.clients[client.channel] = make(map[*Client]struct{})
			}
			h.clients[client.channel][client] = struct{}{}
			h.mu.Unlock()
			log.Printf("[WS] client %s joined channel=%s, total=%d",
				client.id, client.channel, h.countChannel(client.channel))

		case client := <-h.unregister:
			h.mu.Lock()
			if ch, ok := h.clients[client.channel]; ok {
				if _, ok := ch[client]; ok {
					delete(ch, client)
					close(client.send)
				}
			}
			h.mu.Unlock()
			log.Printf("[WS] client %s left channel=%s", client.id, client.channel)

		case msg := <-h.broadcast:
			h.mu.RLock()
			clients := h.clients[msg.channel]
			h.mu.RUnlock()
			for c := range clients {
				select {
				case c.send <- msg.data:
				default:
					// 缓冲满，断开连接
					h.unregister <- c
				}
			}
		}
	}
}

// Broadcast 向指定频道广播消息
func (h *Hub) Broadcast(channel string, data []byte) {
	h.broadcast <- broadcastMsg{channel: channel, data: data}
}

// CountChannel 统计频道在线人数
func (h *Hub) CountChannel(channel string) int {
	return h.countChannel(channel)
}

func (h *Hub) countChannel(channel string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients[channel])
}

// writePump 向客户端写消息（每个连接独立 goroutine）
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(8, nil) // CloseMessage
				return
			}
			if err := c.conn.WriteMessage(1, msg); err != nil { // TextMessage
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(9, nil); err != nil { // PingMessage
				return
			}
		}
	}
}

// readPump 读取客户端消息（保持连接活跃）
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
