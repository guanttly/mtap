// Package main WebSocket 推送服务器
// 管理分诊大屏和监控大屏的实时推送连接
// 同时订阅 Redis Pub/Sub，将后端事件广播至对应频道
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	infraCache "github.com/euler/mtap/internal/infrastructure/cache"
	infraConfig "github.com/euler/mtap/internal/infrastructure/config"
	wsInterface "github.com/euler/mtap/internal/interfaces/ws"
	"github.com/gin-gonic/gin"
)

func main() {
	cfgPath := os.Getenv("MTAP_CONFIG")
	if cfgPath == "" {
		cfgPath = "configs/config.yaml"
	}
	cfg, err := infraConfig.Load(cfgPath)
	if err != nil {
		log.Printf("load config failed (using defaults): %v", err)
		cfg, _ = infraConfig.Load("")
	}

	// 初始化 Redis 客户端（用于 Pub/Sub 转发）
	redisAddr := cfg.Redis.Addr
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisClient, err := infraCache.NewClient(redisAddr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Printf("[WsServer] 连接 Redis 失败（Pub/Sub 不可用）: %v", err)
		redisClient = nil
	}

	// 初始化 WebSocket Hub
	hub := wsInterface.NewHub()
	go hub.Run()

	// 启动 Redis Pub/Sub 转发 goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if redisClient != nil {
		go subscribeAndForward(ctx, redisClient, hub)
	}

	// 初始化 Gin + WebSocket 路由
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())

	wsHandler := wsInterface.NewHandler(hub)
	wsHandler.RegisterRoutes(r)

	// 健康检查
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	port := cfg.Server.Port
	if port == 0 {
		port = 8081
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 0, // WebSocket 连接不限写超时
	}

	go func() {
		log.Printf("[WsServer] WebSocket 服务器已启动，监听 :%d", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[WsServer] 服务器错误: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[WsServer] 收到退出信号，正在停止...")
	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	srv.Shutdown(shutdownCtx)
	log.Println("[WsServer] 已停止")
}

// subscribeAndForward 订阅 Redis Pub/Sub 频道，将消息转发至对应 WebSocket 频道
func subscribeAndForward(ctx context.Context, redisClient *infraCache.Client, hub *wsInterface.Hub) {
	msgCh := redisClient.Subscribe(ctx, wsInterface.ChannelTriage, wsInterface.ChannelDashboard)
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-msgCh:
			if !ok {
				return
			}
			// msg.Channel 即 Redis 频道名，与 WS 频道名一致
			var payload map[string]any
			if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
				hub.Broadcast(msg.Channel, []byte(msg.Payload))
			} else {
				hub.Broadcast(msg.Channel, []byte(msg.Payload))
			}
		}
	}
}
