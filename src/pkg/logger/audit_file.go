package logger

import (
	"context"
	"encoding/json"
	"os"
	"sync"
)

// FileAuditLogger 以 JSONL 方式追加写入审计日志
type FileAuditLogger struct {
	mu   sync.Mutex
	file *os.File
}

func NewFileAuditLogger(path string) (*FileAuditLogger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return nil, err
	}
	return &FileAuditLogger{file: f}, nil
}

func (l *FileAuditLogger) Log(_ context.Context, entry AuditEntry) error {
	b, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, err := l.file.Write(append(b, '\n')); err != nil {
		return err
	}
	return nil
}

func (l *FileAuditLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file == nil {
		return nil
	}
	return l.file.Close()
}

