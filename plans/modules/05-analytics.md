# 统计分析与监控子系统详细设计

| 项目 | 内容 |
|------|------|
| 模块编号 | MOD-05 |
| 对应规格书 | 4.4.3 统计分析大屏 |
| 对应限界上下文 | analytics |
| 上游依赖 | 资源管理、预约服务、分诊执行（数据采集） |
| 下游消费者 | 效能优化子系统（运营数据） |

---

## 1 模块定位

统计分析子系统从各业务模块采集运行数据，提供**实时监控大屏**和**多维度报表导出**两大能力，辅助管理层进行资源调配决策。同时作为效能优化子系统的数据基座。

---

## 2 领域模型

### 2.1 聚合根

```go
// DashboardSnapshot 大屏快照数据
type DashboardSnapshot struct {
    ID              string
    Timestamp       time.Time
    SlotUsage       SlotUsageData       // 号源占用率
    DeviceStatus    []DeviceStatusData  // 设备状态
    WaitTrend       []WaitTrendPoint    // 等待时长趋势
    Alerts          []AlertItem         // 异常告警
}

type SlotUsageData struct {
    TotalSlots    int
    UsedSlots     int
    ExpiredSlots  int
    AvailableSlots int
    UsageRate     float64
}

type DeviceStatusData struct {
    DeviceID   string
    DeviceName string
    Status     string  // idle/in_use/maintenance
    QueueCount int
}

type WaitTrendPoint struct {
    Time        time.Time
    AvgWaitMin  float64
}

type AlertItem struct {
    Type     string // device_queue_overflow / slot_exhausted
    Message  string
    DeviceID string
    Value    int
}

// Report 报表
type Report struct {
    ID            string
    ReportType    string       // daily / weekly / monthly
    Dimensions    []string     // 筛选维度
    DateRange     DateRange
    Status        string       // generating / ready / failed
    FilePath      string       // 导出文件路径
    FileSize      int64
    GeneratedAt   *time.Time
}
```

---

## 3 领域服务

### 3.1 DashboardService（大屏数据服务）

```go
type DashboardService interface {
    // GetSnapshot 获取实时大屏数据
    GetSnapshot(ctx context.Context, campusID string) (*DashboardSnapshot, error)

    // GetDeviceDetail 点击设备卡片展开详情
    GetDeviceDetail(ctx context.Context, deviceID string, date time.Time) (*DeviceDetail, error)
}
```

### 3.2 ReportService（报表服务）

```go
type ReportService interface {
    // Generate 生成报表（同步<30秒 / 异步）
    Generate(ctx context.Context, input ReportInput) (*Report, error)

    // Export 导出报表文件
    Export(ctx context.Context, reportID string, format string) ([]byte, error)
}

type ReportInput struct {
    ReportType   string   // daily/weekly/monthly
    Metrics      []string // slot_usage/device_usage/avg_wait/no_show_rate/override_rate
    CampusID     string
    DepartmentID string
    DeviceID     string
    DateRange    DateRange
    Format       string   // xlsx / pdf
}
```

---

## 4 接口设计

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/v1/analytics/dashboard` | 实时大屏数据 | 管理员 |
| GET | `/api/v1/analytics/dashboard/device/:id` | 设备详情 | 管理员 |
| POST | `/api/v1/analytics/reports` | 生成报表 | 管理员 |
| GET | `/api/v1/analytics/reports/:id` | 报表状态/详情 | 管理员 |
| GET | `/api/v1/analytics/reports/:id/export` | 导出报表文件 | 管理员 |

**WebSocket**：`/ws/v1/dashboard` — 实时推送大屏数据（每10秒刷新）

---

## 5 数据库设计

```sql
-- 大屏快照按时间分区，保留最近7天的高频数据
CREATE TABLE dashboard_snapshots (
    id          VARCHAR(36) PRIMARY KEY,
    campus_id   VARCHAR(36),
    snapshot    JSONB NOT NULL,          -- 完整快照JSON
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_dashboard_campus_time ON dashboard_snapshots(campus_id, created_at);

CREATE TABLE reports (
    id            VARCHAR(36) PRIMARY KEY,
    report_type   VARCHAR(10) NOT NULL,
    dimensions    JSONB,
    date_start    DATE NOT NULL,
    date_end      DATE NOT NULL,
    status        VARCHAR(15) NOT NULL DEFAULT 'generating',
    file_path     VARCHAR(500),
    file_size     BIGINT,
    format        VARCHAR(5) NOT NULL DEFAULT 'xlsx',
    generated_at  TIMESTAMP,
    created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by    VARCHAR(36) NOT NULL
);
```

---

## 6 前端页面设计

| 页面 | 路由 | 核心交互 |
|------|------|----------|
| 实时监控大屏 | `/analytics/dashboard` | 号源饼图 + 设备状态卡片 + 等待趋势折线图 + 告警面板 |
| 报表导出 | `/analytics/report` | 维度选择 + 日期范围 + 生成/下载 |

**大屏技术要点**：
- ECharts 图表，10秒 WebSocket 刷新
- 支持按院区/科室切换
- 设备卡片可点击展开当日详情
- 告警阈值：等待队列 > 20人高亮红色

---

## 7 错误码定义

| 错误码 | 说明 |
|--------|------|
| `STATS_001` | 报表查询超时（已切换异步生成） |
| `STATS_002` | 查询时间范围超过24个月 |
| `STATS_003` | 异步报表生成失败 |
| `STATS_004` | 导出文件超过50MB（已压缩为ZIP） |
