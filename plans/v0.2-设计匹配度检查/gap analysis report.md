# 一站式全医技预约平台 — 实现差距分析与改造方案

> 版本 1.0 | 2026-03-23

---

## 目录

- [1 问题总览：三层不一致](#1-问题总览三层不一致)
- [2 详细差距清单](#2-详细差距清单)
  - [2.1 结构不一致（Migration vs PO vs 领域实体）](#21-结构不一致migration-vs-po-vs-领域实体)
  - [2.2 功能未实现](#22-功能未实现)
  - [2.3 功能部分实现](#23-功能部分实现)
- [3 改造方案](#3-改造方案)
  - [3.1 第一阶段 P0 — 结构补齐（3-5天）](#31-第一阶段-p0--结构补齐3-5天)
  - [3.2 第二阶段 P1 — 功能补全（7-10天）](#32-第二阶段-p1--功能补全7-10天)
  - [3.3 第三阶段 P2 — 增强功能（10-15天）](#33-第三阶段-p2--增强功能10-15天)
- [4 迁移脚本规范建议](#4-迁移脚本规范建议)
- [5 修正后的 Mock Data 说明](#5-修正后的-mock-data-说明)

---

## 1 问题总览：三层不一致

通过对比**规格书设计文档**、**领域层实体定义**（`domain/resource/entity.go`）、**GORM PO 持久化对象**（`po/*.go`）、以及 **Migration SQL**（`001_init_schema.sql`）四个层面，发现当前系统存在以下三类核心问题：

| 问题类型 | 数量 | 说明 |
|---------|------|------|
| 🔴 **结构不一致**（Migration vs PO vs 领域实体） | **6 处** | 数据库字段缺失或命名不匹配，导致 GORM 操作报错、功能无法持久化 |
| 🟠 **功能未实现**（设计有但代码无） | **7 处** | 整个子模块或关键功能在代码层面缺失 |
| 🔵 **功能部分实现** | **4 处** | 有框架代码但核心逻辑未完成 |

> **⚠️ 其中 P0 级别问题（3处）直接影响核心流程跑通，必须优先修复。**

---

## 2 详细差距清单

### 2.1 结构不一致（Migration vs PO vs 领域实体）

#### #1 🔴 devices 表缺少关键字段

| 维度 | 内容 |
|------|------|
| **级别** | **P0** |
| **模块** | 资源管理 |
| **设计** | Device 领域实体包含 `Model`, `Manufacturer`, `SupportedExamTypes []string`, `MaxDailySlots int` |
| **现状** | Migration SQL 和 DevicePO 仅有 `id, name, campus_id, department_id, status` |
| **影响** | 设备管理页面无法展示型号/厂商/支持检查类型；排班号源生成缺少 `max_daily_slots` 上限校验；预约引擎无法按检查类型筛选设备 |

**证据对比：**

```go
// domain/resource/entity.go — 领域实体（有字段）
type Device struct {
    ID                 string       
    Model              string       // ← Migration 无此列
    Manufacturer       string       // ← Migration 无此列
    SupportedExamTypes []string     // ← Migration 无此列
    MaxDailySlots      int          // ← Migration 无此列
    ...
}
```

```go
// persistence/po/resource.go — PO（无字段）
type DevicePO struct {
    ID           string  `gorm:"column:id"`
    Name         string  `gorm:"column:name"`
    CampusID     string  `gorm:"column:campus_id"`
    DepartmentID string  `gorm:"column:department_id"`
    Status       string  `gorm:"column:status"`
    // model, manufacturer, supported_exam_types, max_daily_slots 全部缺失
}
```

---

#### #2 🔴 SlotPool 表缺少配额与溢出字段

| 维度 | 内容 |
|------|------|
| **级别** | **P0** |
| **模块** | 资源管理 |
| **设计** | SlotPool 领域实体含 `AllocationRatio float64`, `OverflowEnabled bool`, `OverflowTargetPool string` |
| **现状** | Migration 和 SlotPoolPO 仅有 `id, name, type, status` |
| **影响** | 号源池配额比例配置、溢出规则无法持久化；门诊/住院号源分配比例功能缺失 |

---

#### #3 🔴 sorting_strategies 表 PO 与 Migration 字段不一致

| 维度 | 内容 |
|------|------|
| **级别** | **P1** |
| **模块** | 规则引擎 |
| **设计** | PO 定义 `scope_campuses, scope_depts, scope_devices` 三个独立字段 |
| **现状** | Migration SQL 定义为单一 `scope`(JSON) 字段 |
| **影响** | GORM AutoMigrate 创建三列但 Migration 只有一列；生产环境部署时 PO 读写会报 `Unknown column` |

**证据对比：**

```sql
-- 001_init_schema.sql（单字段）
CREATE TABLE sorting_strategies (
    ...
    scope JSON,                    -- ← 单一 JSON 字段
    ...
);
```

```go
// po/rule.go（三字段）
type SortingStrategyPO struct {
    ScopeCampuses string `gorm:"column:scope_campuses"` // ← Migration 无
    ScopeDepts    string `gorm:"column:scope_depts"`    // ← Migration 无
    ScopeDevices  string `gorm:"column:scope_devices"`  // ← Migration 无
}
```

---

#### #4 🟠 source_controls 表溢出目标池字段名不匹配

| 维度 | 内容 |
|------|------|
| **级别** | **P1** |
| **模块** | 规则引擎 |
| **PO 定义** | `OverflowTargetPoolID` → `gorm:"column:overflow_target_pool_id"` |
| **Migration** | 列名为 `overflow_target_pool` |
| **影响** | 查询/写入溢出目标池时列名不匹配导致 SQL 报错 |

---

#### #5 🟡 patient_adapt_rules 表 condition_value 长度不一致

| 维度 | 内容 |
|------|------|
| **级别** | P2 |
| **PO** | `size:50` |
| **Migration** | `VARCHAR(100)` |
| **影响** | 长条件值在 PO 层被截断 |

---

#### #6 🟡 roles 表缺少 status 字段

| 维度 | 内容 |
|------|------|
| **级别** | P2 |
| **设计** | 设计文档中角色有 status（active/inactive） |
| **现状** | RolePO 和 Migration 均无 status 字段，只有 `is_preset` |
| **影响** | 无法禁用角色 |

---

### 2.2 功能未实现

#### #7 doctors 表未创建

| 维度 | 内容 |
|------|------|
| **级别** | **P1** |
| **模块** | 资源管理 |
| **设计** | Doctor 领域实体完整定义（department_id, his_code, name, title, gender）；`DoctorRepository` 接口已在 `repository.go` 中声明 |
| **现状** | `001_init_schema.sql` 中无 `CREATE TABLE doctors`；无 `DoctorPO` 结构体 |
| **影响** | 医生排班、性别隔离匹配、医生专池等功能无法落地 |

---

#### #8 设备-检查项目关联关系缺失

| 维度 | 内容 |
|------|------|
| **级别** | **P0** |
| **设计** | `Device.SupportedExamTypes` 定义设备可执行的检查类型列表 |
| **现状** | `devices` 表无此字段，也无独立关联表 |
| **影响** | 号源生成无法校验设备是否支持该检查项目；预约引擎无法按检查类型筛选设备 |

---

#### #9 排班模板(ScheduleTemplate)未实现

| 维度 | 内容 |
|------|------|
| **级别** | **P1** |
| **设计** | `schedule_templates` 表、`ScheduleTemplate` 实体、`repeat_type` / `slot_pattern` |
| **现状** | 无 `ScheduleTemplatePO`，无模板 CRUD API，排班生成直接传参数 |
| **影响** | 无法保存常用排班模板，每次排班需重新配置 |

---

#### #10 检前通知(短信/微信推送)未实现

| 维度 | 内容 |
|------|------|
| **级别** | **P1** |
| **设计** | 预约成功后自动推送注意事项；停诊时通知患者改约 |
| **现状** | 无通知服务实现，无短信/微信渠道对接 |
| **影响** | 患者无法收到预约确认、检前提醒、停诊改约通知 |

---

#### #11 分诊大屏 WebSocket 推送未实现

| 维度 | 内容 |
|------|------|
| **级别** | P2 |
| **设计** | `/ws/v1/screen` 签到后实时更新分诊大屏队列 |
| **现状** | WebSocket handler 未找到实际实现 |

---

#### #12 效能优化子系统整体未实现

| 维度 | 内容 |
|------|------|
| **级别** | P2 |
| **设计** | 规格书 4.5 节：异常检测引擎、瓶颈归因分析、A/B/C 三类优化策略、试运行、评估报告 |
| **现状** | Migration 中无 `efficiency_metrics` 等表，无 optimization 相关 PO/Service/Handler |
| **影响** | 效能优化闭环功能全部缺失（属于高级功能可延后） |

---

#### #13 审计日志/操作日志表未创建

| 维度 | 内容 |
|------|------|
| **级别** | P2 |
| **设计** | `audit_logs`, `operation_logs` 表 |
| **现状** | Migration 中无对应 `CREATE TABLE` |
| **影响** | 操作行为无法追溯审计 |

---

### 2.3 功能部分实现

#### #14 缴费校验(HIS扣费接口)

| 维度 | 内容 |
|------|------|
| **级别** | **P1** |
| **设计** | 预约前验证缴费状态，未缴费 24h 自动释放号源 |
| **现状** | `AppointmentPO` 有 `payment_verified` 字段但无 HIS 扣费接口调用逻辑 |
| **影响** | 预约流程跳过缴费校验，无法自动释放未缴费号源 |

---

#### #15 自动预约引擎(智能最优方案计算)

| 维度 | 内容 |
|------|------|
| **级别** | **P1** |
| **设计** | 一键自动预约：缴费后触发智能引擎计算最优时间/地点组合 |
| **现状** | `auto` 预约接口存在但无智能排序/组合优化算法 |
| **影响** | 无法自动为多项目预约计算最优组合方案 |

---

#### #16 实时监控大屏数据采集

| 维度 | 内容 |
|------|------|
| **级别** | P2 |
| **设计** | 展示号源占用率、设备耗时、等待时长热力图 |
| **现状** | `DashboardSnapshotPO` 存在但采集逻辑简化 |

---

#### #17 角色 status 管理

| 维度 | 内容 |
|------|------|
| **级别** | P2 |
| **设计** | 角色有 active/inactive 状态 |
| **现状** | 只有 `is_preset` 区分预置/自定义，无法禁用角色 |

---

## 3 改造方案

按优先级分三个阶段推进，每个阶段内的改动相互独立可并行。

### 3.1 第一阶段 P0 — 结构补齐（3-5天）

**目标：补齐核心表字段，使基础数据 → 排班 → 号源 → 预约主流程可跑通。**

#### 3.1.1 devices 表字段补齐 + 设备-检查项目关联

1. 编写 `002_add_device_fields.sql` 增量迁移脚本：

```sql
ALTER TABLE devices 
  ADD COLUMN model VARCHAR(50) AFTER name,
  ADD COLUMN manufacturer VARCHAR(50) AFTER model,
  ADD COLUMN supported_exam_types JSON DEFAULT '[]' AFTER manufacturer,
  ADD COLUMN max_daily_slots INT NOT NULL DEFAULT 50 AFTER supported_exam_types;
```

2. 更新 `DevicePO` 结构体，增加四个 gorm tag 字段：

```go
type DevicePO struct {
    // ... 已有字段 ...
    Model              string `gorm:"column:model;size:50"`
    Manufacturer       string `gorm:"column:manufacturer;size:50"`
    SupportedExamTypes string `gorm:"column:supported_exam_types;type:json;default:'[]'"` 
    MaxDailySlots      int    `gorm:"column:max_daily_slots;not null;default:50"`
}
```

3. 更新 `DevicePO ↔ Device` 领域实体转换函数，`SupportedExamTypes` 做 JSON marshal/unmarshal。

4. 更新 `CreateDeviceReq` / `UpdateDeviceReq` DTO，增加对应字段。

5. 更新前端 `DeviceManager.vue`：表格增加型号、厂商、支持检查类型（Tag 展示）、每日最大号位列；新建/编辑表单增加对应字段。

6. 号源生成服务 `SlotGenerationService.Generate()` 中增加 `max_daily_slots` 上限校验逻辑。

#### 3.1.2 SlotPool 表字段补齐

1. 编写增量迁移：

```sql
ALTER TABLE slot_pools 
  ADD COLUMN allocation_ratio DECIMAL(5,2) NOT NULL DEFAULT 0.00,
  ADD COLUMN overflow_enabled TINYINT(1) NOT NULL DEFAULT 0,
  ADD COLUMN overflow_target_pool VARCHAR(36);
```

2. 更新 `SlotPoolPO` 增加 `AllocationRatio`, `OverflowEnabled`, `OverflowTargetPool` 字段。

3. 更新号源池管理 API 和前端 `SlotPoolView.vue`，支持配额比例配置和溢出规则设置。

#### 3.1.3 字段命名统一修复

1. **sorting_strategies 表**：Migration 中用 `scope`(JSON) 单字段存储，而 PO 用三个字段。**统一为 PO 方案（三字段）**，修改 Migration 增加 `scope_campuses`, `scope_depts`, `scope_devices` 列并迁移数据，删除 `scope` 列：

```sql
ALTER TABLE sorting_strategies
  ADD COLUMN scope_campuses TEXT,
  ADD COLUMN scope_depts TEXT,
  ADD COLUMN scope_devices TEXT;

-- 迁移数据（从 scope JSON 拆分到三列）
UPDATE sorting_strategies SET
  scope_campuses = JSON_EXTRACT(scope, '$.campuses'),
  scope_depts = JSON_EXTRACT(scope, '$.depts'),
  scope_devices = JSON_EXTRACT(scope, '$.devices')
WHERE scope IS NOT NULL;

ALTER TABLE sorting_strategies DROP COLUMN scope;
```

2. **source_controls 表**：PO 中 `overflow_target_pool_id` 与 Migration 中 `overflow_target_pool` 不一致。修改 PO gorm tag 为 `column:overflow_target_pool` 对齐 Migration：

```go
// 修改前
OverflowTargetPoolID string `gorm:"column:overflow_target_pool_id"`
// 修改后
OverflowTargetPoolID string `gorm:"column:overflow_target_pool"`
```

3. **patient_adapt_rules 表**：PO `condition_value size:50` 与 Migration `VARCHAR(100)` 不一致，修改 PO 为 `size:100` 对齐。

---

### 3.2 第二阶段 P1 — 功能补全（7-10天）

**目标：补全核心业务功能，使完整预约闭环可运转。**

#### 3.2.1 doctors 表创建与医生管理

1. 新增 Migration：

```sql
CREATE TABLE IF NOT EXISTS doctors (
    id            VARCHAR(36)  NOT NULL,
    department_id VARCHAR(36)  NOT NULL,
    his_code      VARCHAR(30),
    name          VARCHAR(30)  NOT NULL,
    title         VARCHAR(20),
    gender        VARCHAR(10)  NOT NULL DEFAULT 'unknown',
    status        VARCHAR(10)  NOT NULL DEFAULT 'active',
    synced_at     DATETIME(3),
    created_at    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_doctors_department (department_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

2. 新增 `DoctorPO` 结构体（`po/resource.go`）。
3. 实现 `DoctorRepository` 接口（已在 `repository.go` 中定义）。
4. 新增医生管理 CRUD API（Handler + Service + DTO）。
5. 前端新增医生管理页面或嵌入科室管理中。

#### 3.2.2 排班模板功能

1. 新增 Migration：

```sql
CREATE TABLE IF NOT EXISTS schedule_templates (
    id          VARCHAR(36)  NOT NULL,
    name        VARCHAR(50)  NOT NULL,
    repeat_type VARCHAR(20)  NOT NULL,
    slot_pattern JSON NOT NULL,
    created_at  DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

2. 新增 `ScheduleTemplatePO`、领域实体、仓储接口。
3. 排班生成接口增加 `template_id` 可选参数，支持从模板加载排班参数。
4. 前端排班日历页面增加「模板管理」入口。

#### 3.2.3 缴费校验与号源自动释放

1. 定义 HIS 缴费校验接口（`infrastructure/external/his_payment.go`），生产环境对接 HIS RESTful 接口，开发环境提供 Mock 实现。
2. 预约确认流程中调用缴费校验；未缴费状态设为 `pending`。
3. 新增定时任务（`cmd/scheduler`）：每小时扫描 `pending` 超 24 小时的预约，自动取消并释放号源。

#### 3.2.4 自动预约引擎优化

1. 实现多项目组合预约算法：输入多个 `exam_item_id`，调用规则引擎校验冲突/依赖，按**最短等待时间 + 空腹优先**排序生成推荐方案。
2. `combo` 预约接口返回最多 3 个候选方案供前端展示对比。

#### 3.2.5 检前通知服务

1. 新增 `notification-worker` 独立进程（`cmd/notification-worker`）。
2. 定义通知渠道接口（短信/微信模板消息），开发环境用日志 Mock。
3. 预约成功、改约、取消、停诊、检前提醒（T-1天）等事件触发通知。
4. 通知失败进入重试队列（每 10 分钟重试，最多 6 次）。

---

### 3.3 第三阶段 P2 — 增强功能（10-15天）

**目标：完善大屏推送、统计报表、审计日志、效能优化等高级功能。**

- **分诊大屏 WebSocket 推送**：实现 `/ws/v1/screen` 接口，签到/呼叫/完成事件实时广播到前端大屏。
- **审计日志 / 操作日志**：新增 `audit_logs`, `operation_logs` 表；通过 Gin 中间件自动记录关键操作。
- **实时监控大屏**：完善 DashboardSnapshot 采集逻辑，增加号源占用率、设备利用率、平均等待时长等指标计算。
- **效能优化子系统**：按规格书 4.5 节分期实现 ——
  - 第一期：指标体系 + 异常检测
  - 第二期：瓶颈归因 + A 类策略建议
  - 第三期：B/C 类策略 + 试运行 + 评估报告
- **角色 status 字段**：如需支持角色禁用，在 `roles` 表增加 `status` 字段并更新 RBAC 鉴权中间件。

---

## 4 迁移脚本规范建议

为防止再次出现三层不一致，建议建立以下规范：

1. **单一事实来源**：以 Migration SQL 为准，PO 结构体通过 gorm tag 严格映射 Migration 定义的列名和类型；禁止依赖 AutoMigrate 自动建表。

2. **增量迁移命名**：`002_add_device_fields.sql`、`003_create_doctors.sql` ... 按顺序编号，每次结构变更独立文件。

3. **CI 校验**：在 CI 流水线中增加「Migration 与 PO 一致性检查」步骤 — 可通过反射遍历所有 PO 的 gorm column tag 与 Migration DDL 做交叉比对。

4. **代码审查 Checklist**：任何涉及表结构变更的 PR 必须同步修改 **Migration SQL + PO + 领域实体 + DTO** 四层，缺一不合入。

---

## 5 修正后的 Mock Data 说明

同步输出两个版本的 mock data：

- **`mock_data_current.sql`** — 匹配当前实际数据库结构（devices 无 model 等字段、无 doctors 表），可直接导入现有环境运行。
- **`mock_data_target.sql`** — 匹配改造后的目标结构（包含所有补齐字段），用于改造完成后的集成测试。