# 公共基础设施与安全模块详细设计

| 项目 | 内容 |
|------|------|
| 模块编号 | MOD-00 |
| 对应规格书 | 5.x 安全与性能需求 |
| 涉及目录 | pkg/, internal/infrastructure/, configs/ |

---

## 1 模块定位

公共基础设施模块为所有业务子系统提供**横切关注点**：认证鉴权、统一错误处理、加解密、日志审计、消息队列、缓存、外部系统对接、配置管理等能力。

---

## 2 认证鉴权（pkg/auth + middleware）

### 2.1 JWT Token 流程

```mermaid
sequenceDiagram
    actor 用户
    participant 前端
    participant API as api-server
    participant Redis

    用户->>前端: 输入用户名/密码
    前端->>API: POST /api/v1/auth/login
    API->>API: 校验密码(bcrypt)
    API->>Redis: 存储Session信息
    API-->>前端: { access_token(2h), refresh_token(7d) }
    前端->>API: GET /api/v1/xxx (Header: Bearer token)
    API->>API: JWT验证 + 提取UserInfo
    API-->>前端: 业务数据

    Note over 前端,API: Token即将过期
    前端->>API: POST /api/v1/auth/refresh
    API->>Redis: 校验refresh_token有效性
    API-->>前端: 新access_token
```

### 2.2 RBAC 权限模型

| 角色 | 权限范围 | 说明 |
|------|----------|------|
| `admin` | 全部接口 | 系统管理员 |
| `scheduler_admin` | 排班/规则/策略管理 | 预约管理员 |
| `operator` | 预约/分诊操作 | 操作员 |
| `nurse` | 签到/呼叫/检查状态 | 护士 |
| `viewer` | 只读查询 | 查看者 |

```sql
CREATE TABLE users (
    id            VARCHAR(36)  PRIMARY KEY,
    username      VARCHAR(50)  NOT NULL UNIQUE,
    password_hash VARCHAR(100) NOT NULL,
    real_name     VARCHAR(30),
    role_id       VARCHAR(36)  NOT NULL REFERENCES roles(id),
    department_id VARCHAR(36)  REFERENCES departments(id),
    status        VARCHAR(10)  NOT NULL DEFAULT 'active',
    last_login_at TIMESTAMP,
    created_at    TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE roles (
    id          VARCHAR(36)  PRIMARY KEY,
    name        VARCHAR(30)  NOT NULL UNIQUE,
    permissions JSONB        NOT NULL DEFAULT '[]',  -- ["rule:read","rule:write","appt:*"]
    is_preset   BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);
```

### 2.3 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/auth/login` | 登录 |
| POST | `/api/v1/auth/refresh` | 刷新Token |
| POST | `/api/v1/auth/logout` | 登出 |
| GET | `/api/v1/auth/profile` | 当前用户信息 |

---

## 3 统一响应格式（pkg/response）

```go
type Response struct {
    Code    int         `json:"code"`    // 0=成功，非0=错误
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

type PageResult[T any] struct {
    Items      []T   `json:"items"`
    Total      int64 `json:"total"`
    Page       int   `json:"page"`
    PageSize   int   `json:"page_size"`
    TotalPages int   `json:"total_pages"`
}
```

---

## 4 统一错误码（pkg/errors）

```go
// 系统级错误码 1xxx
const (
    ErrUnauthorized     = 1001 // 未认证
    ErrForbidden        = 1002 // 权限不足
    ErrNotFound         = 1003 // 资源不存在
    ErrInvalidParam     = 1004 // 参数校验失败
    ErrInternalServer   = 1005 // 服务器内部错误
    ErrRateLimit        = 1006 // 请求频率过高
    ErrDuplicate        = 1007 // 重复操作
)

// 业务错误码 各模块前缀
// RULE_0xx  规则引擎
// RES_0xx   资源管理
// APPT_0xx  预约服务
// TRIAGE_0xx 分诊执行
// STATS_0xx  统计分析
// OPT_0xx   效能优化
```

---

## 5 数据安全（pkg/encrypt）

| 场景 | 方案 |
|------|------|
| 传输加密 | HTTPS/TLS 1.2+ |
| 存储加密 | 患者姓名/手机号 AES-256加密后存储 |
| 展示脱敏 | 姓名中间字符替换为* |
| 审计日志 | 不可篡改，保留2年 |

---

## 6 日志审计（pkg/logger）

```go
// StructuredLogger 结构化日志
type StructuredLogger interface {
    Info(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    WithContext(ctx context.Context) StructuredLogger
}

// AuditLogger 审计日志（追加写入，不可修改）
type AuditLogger interface {
    Log(ctx context.Context, entry AuditEntry) error
}

type AuditEntry struct {
    OperatorID   string
    OperatorName string
    Action       string    // create/update/delete/override/approve/reject
    Resource     string    // conflict_rule/appointment/strategy
    ResourceID   string
    OldValue     string    // JSON
    NewValue     string    // JSON
    IP           string
    Timestamp    time.Time
}
```

```sql
CREATE TABLE audit_logs (
    id           VARCHAR(36) PRIMARY KEY,
    operator_id  VARCHAR(36) NOT NULL,
    operator_name VARCHAR(30),
    action       VARCHAR(20) NOT NULL,
    resource     VARCHAR(30) NOT NULL,
    resource_id  VARCHAR(36) NOT NULL,
    old_value    TEXT,
    new_value    TEXT,
    ip           VARCHAR(45),
    created_at   TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_audit_resource ON audit_logs(resource, resource_id);
CREATE INDEX idx_audit_operator ON audit_logs(operator_id, created_at);

CREATE TABLE operation_logs (
    id           VARCHAR(36) PRIMARY KEY,
    level        VARCHAR(5) NOT NULL,
    message      TEXT NOT NULL,
    fields       JSONB,
    caller       VARCHAR(200),
    created_at   TIMESTAMP NOT NULL DEFAULT NOW()
);
```

---

## 7 消息队列（internal/infrastructure/mq）

### 7.1 Exchange & Queue 规划

| Exchange | Queue | 路由键 | 消费者 |
|----------|-------|--------|--------|
| `mtap.events` | `q.notification` | `appointment.confirmed` / `appointment.cancelled` / `schedule.changed` | notify-worker |
| `mtap.events` | `q.analytics` | `#` (全部事件) | ws-server (数据聚合) |
| `mtap.events` | `q.optimization` | `bottleneck.detected` / `trial.*` / `strategy.*` | analyzer-worker |

---

## 8 配置管理（internal/infrastructure/config）

```yaml
# configs/config.yaml 结构示例
server:
  port: 8080
  mode: release          # debug / release
  read_timeout: 10s
  write_timeout: 10s

database:
  host: localhost
  port: 5432
  name: mtap
  user: mtap_app
  password: ""           # 生产环境从环境变量读取
  max_open_conns: 50
  max_idle_conns: 10
  conn_max_lifetime: 1h

redis:
  addr: localhost:6379
  password: ""
  db: 0
  pool_size: 20

rabbitmq:
  url: amqp://guest:guest@localhost:5672/

jwt:
  secret: ""            # 从环境变量读取
  access_expire: 2h
  refresh_expire: 168h

his:
  base_url: http://his-api:8081
  timeout: 5s
  retry: 3

ratelimit:
  max_requests: 60
  window: 1m
```

---

## 9 前端公共设计

### 9.1 API 请求封装（web/src/api/request.ts）

- Axios 实例：baseURL、超时30秒
- 请求拦截器：自动附带 `Authorization: Bearer <token>`
- 响应拦截器：code≠0自动弹错误提示；401自动跳转登录；429 显示限流提示

### 9.2 路由守卫

- 登录校验：无Token重定向到 `/login`
- 权限校验：按角色判断路由可见性
- 页面标题自动更新

### 9.3 全局样式设计系统

- 主色调：医疗蓝 `#1890FF`
- 辅助色：成功绿`#52C41A`、警告橙`#FAAD14`、危险红`#FF4D4F`
- 字体：系统默认 + 备用 Inter
- 卡片圆角：8px
- 阴影：`0 2px 8px rgba(0,0,0,0.08)`
