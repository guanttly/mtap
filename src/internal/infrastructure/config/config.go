// Package config 基础设施层 - 配置加载
// 核心目的：加载和管理应用配置
// 模块功能：
//   - 从YAML文件加载配置
//   - 环境变量覆盖
//   - 数据库/Redis/MQ/外部系统连接配置
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 应用总配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Kafka    KafkaConfig    `yaml:"kafka"`
	JWT      JWTConfig      `yaml:"jwt"`
	HIS      HISConfig      `yaml:"his"`
	Payment  PaymentConfig  `yaml:"payment"`
	Message  MessageConfig  `yaml:"message"`
	Log      LogConfig      `yaml:"log"`
}

// ServerConfig HTTP 服务配置
type ServerConfig struct {
	Port         int           `yaml:"port"`
	WSPort       int           `yaml:"ws_port"` // WebSocket 服务端口，默认 8081
	Mode         string        `yaml:"mode"`    // debug / release
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string        `yaml:"driver"` // 只支持 mysql
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	Name            string        `yaml:"name"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	Charset         string        `yaml:"charset"`    // 默认 utf8mb4
	ParseTime       bool          `yaml:"parse_time"` // 默认 true
}

// DSNString 生成 MySQL GORM DSN
func (c DatabaseConfig) DSNString() string {
	// MySQL DSN: user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	charset := c.Charset
	if charset == "" {
		charset = "utf8mb4"
	}
	parseTime := "True"
	if !c.ParseTime {
		parseTime = "False"
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Name, charset, parseTime)
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"`
}

// KafkaConfig Kafka 配置
type KafkaConfig struct {
	Brokers       []string `yaml:"brokers"`        // broker 地址列表，如 ["localhost:9092"]
	GroupID       string   `yaml:"group_id"`       // 消费者组 ID 前缀
	NumPartitions int      `yaml:"num_partitions"` // Topic 分区数（自动创建时使用）
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret        string        `yaml:"secret"`
	AccessExpire  time.Duration `yaml:"access_expire"`
	RefreshExpire time.Duration `yaml:"refresh_expire"`
}

// HISConfig HIS 系统配置
type HISConfig struct {
	BaseURL string        `yaml:"base_url"`
	Timeout time.Duration `yaml:"timeout"`
	Retry   int           `yaml:"retry"`
	APIKey  string        `yaml:"api_key"`
}

// PaymentConfig 收费系统配置
type PaymentConfig struct {
	BaseURL string        `yaml:"base_url"`
	Timeout time.Duration `yaml:"timeout"`
	APIKey  string        `yaml:"api_key"`
}

// MessageConfig 消息平台配置
type MessageConfig struct {
	SMSProvider string `yaml:"sms_provider"` // aliyun / tencent
	SMSKey      string `yaml:"sms_key"`
	SMSSecret   string `yaml:"sms_secret"`
	SMSSign     string `yaml:"sms_sign"`
	WechatAppID string `yaml:"wechat_app_id"`
	WechatToken string `yaml:"wechat_token"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level     string `yaml:"level"` // debug / info / warn / error
	AuditPath string `yaml:"audit_path"`
}

// Load 从 YAML 文件加载配置，并用环境变量覆盖敏感项
func Load(path string) (*Config, error) {
	cfg := defaultConfig()

	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read config file: %w", err)
		}
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("parse config file: %w", err)
		}
	}

	// 环境变量覆盖
	overrideFromEnv(cfg)
	return cfg, nil
}

// MustLoad 加载配置，失败则 panic
func MustLoad(path string) *Config {
	cfg, err := Load(path)
	if err != nil {
		panic("load config failed: " + err.Error())
	}
	return cfg
}

func defaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         8080,
			Mode:         "release",
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		Database: DatabaseConfig{
			Driver:          "mysql",
			Host:            "127.0.0.1",
			Port:            3306,
			Name:            "mtap",
			User:            "root",
			Password:        "",
			MaxOpenConns:    50,
			MaxIdleConns:    10,
			ConnMaxLifetime: time.Hour,
			Charset:         "utf8mb4",
			ParseTime:       true,
		},
		Redis: RedisConfig{
			Addr:     "localhost:6379",
			DB:       0,
			PoolSize: 20,
		},
		Kafka: KafkaConfig{
			Brokers:       []string{"localhost:9092"},
			GroupID:       "mtap",
			NumPartitions: 3,
		},
		JWT: JWTConfig{
			Secret:        "dev-secret-change-me-32-bytes-len!",
			AccessExpire:  2 * time.Hour,
			RefreshExpire: 7 * 24 * time.Hour,
		},
		HIS: HISConfig{
			BaseURL: "http://localhost:8081",
			Timeout: 5 * time.Second,
			Retry:   3,
		},
		Log: LogConfig{
			Level:     "info",
			AuditPath: "audit.log",
		},
	}
}

// splitNonEmpty 按分隔符切分字符串并过滤空项
func splitNonEmpty(s, sep string) []string {
	parts := []string{}
	for _, p := range append(parts, s) {
		if p != "" {
			for _, item := range splitString(p, sep) {
				if item != "" {
					parts = append(parts, item)
				}
			}
		}
	}
	return parts
}

func splitString(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
		}
	}
	result = append(result, s[start:])
	return result
}

func overrideFromEnv(cfg *Config) {
	if v := os.Getenv("MTAP_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Server.Port = port
		}
	}
	if v := os.Getenv("MTAP_DB_HOST"); v != "" {
		cfg.Database.Host = v
	}
	if v := os.Getenv("MTAP_DB_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Database.Port = port
		}
	}
	if v := os.Getenv("MTAP_DB_NAME"); v != "" {
		cfg.Database.Name = v
	}
	if v := os.Getenv("MTAP_DB_USER"); v != "" {
		cfg.Database.User = v
	}
	if v := os.Getenv("MTAP_DB_PASSWORD"); v != "" {
		cfg.Database.Password = v
	}
	if v := os.Getenv("MTAP_REDIS_ADDR"); v != "" {
		cfg.Redis.Addr = v
	}
	if v := os.Getenv("MTAP_REDIS_PASSWORD"); v != "" {
		cfg.Redis.Password = v
	}
	if v := os.Getenv("MTAP_KAFKA_BROKERS"); v != "" {
		// 逗号分隔，如 "broker1:9092,broker2:9092"
		cfg.Kafka.Brokers = splitNonEmpty(v, ",")
	}
	if v := os.Getenv("MTAP_KAFKA_GROUP"); v != "" {
		cfg.Kafka.GroupID = v
	}
	if v := os.Getenv("MTAP_JWT_SECRET"); v != "" {
		cfg.JWT.Secret = v
	}
	if v := os.Getenv("MTAP_AUDIT_LOG"); v != "" {
		cfg.Log.AuditPath = v
	}
	if v := os.Getenv("MTAP_HIS_URL"); v != "" {
		cfg.HIS.BaseURL = v
	}
}
