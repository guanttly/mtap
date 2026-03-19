// Package logger 提供结构化日志和审计日志能力
package logger

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultLogger *zap.SugaredLogger

func init() {
	l, _ := NewLogger("info")
	defaultLogger = l
}

// NewLogger 创建结构化日志实例
func NewLogger(level string) (*zap.SugaredLogger, error) {
	lvl, err := zapcore.ParseLevel(level)
	if err != nil {
		lvl = zapcore.InfoLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		lvl,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return logger.Sugar(), nil
}

// L 获取默认的全局日志实例
func L() *zap.SugaredLogger {
	return defaultLogger
}

// SetDefault 设置全局日志实例
func SetDefault(l *zap.SugaredLogger) {
	defaultLogger = l
}

// Info 结构化信息日志
func Info(msg string, keysAndValues ...interface{}) {
	defaultLogger.Infow(msg, keysAndValues...)
}

// Warn 结构化警告日志
func Warn(msg string, keysAndValues ...interface{}) {
	defaultLogger.Warnw(msg, keysAndValues...)
}

// Error 结构化错误日志
func Error(msg string, keysAndValues ...interface{}) {
	defaultLogger.Errorw(msg, keysAndValues...)
}

// AuditEntry 审计日志条目
type AuditEntry struct {
	OperatorID   string    `json:"operator_id"`
	OperatorName string    `json:"operator_name"`
	Action       string    `json:"action"`
	Resource     string    `json:"resource"`
	ResourceID   string    `json:"resource_id"`
	OldValue     string    `json:"old_value,omitempty"`
	NewValue     string    `json:"new_value,omitempty"`
	IP           string    `json:"ip,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}

// AuditLogger 审计日志接口
type AuditLogger interface {
	Log(ctx context.Context, entry AuditEntry) error
}

var defaultAuditLogger AuditLogger

// SetAuditLogger 设置全局审计日志实现
func SetAuditLogger(l AuditLogger) {
	defaultAuditLogger = l
}

// Audit 写入一条审计日志（若未配置实现则忽略）
func Audit(ctx context.Context, entry AuditEntry) error {
	if defaultAuditLogger == nil {
		return nil
	}
	return defaultAuditLogger.Log(ctx, entry)
}
