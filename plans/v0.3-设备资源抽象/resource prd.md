# v0.2 诊位资源模型重构 — PRD

> **版本** v0.2-draft &nbsp;|&nbsp; **日期** 2026-03-23 &nbsp;|&nbsp; **状态** 待评审

---

## 1 变更背景

### 1.1 问题陈述

当前系统的排班→号源→预约全链路以 `Device`（物理设备）作为唯一资源主体。该模型能覆盖放射科（CT/MRI/DR）、超声科、内镜中心等**设备驱动型**科室，但无法覆盖医院中大量**非设备驱动型**的医技预约场景：

| 科室/场景 | 实际排班单元 | 用 Device 硬套的问题 |
|----------|------------|-------------------|
| 检验科（抽血窗口） | 窗口/工位 | 无型号/厂商；多窗口并发不是"一台设备按时段切分" |
| 康复科/推拿 | 治疗师+治疗床 | 资源是人+空间组合；治疗师可换床位 |
| 门诊手术室 | 手术间 | 一个手术间内含多台设备但排班单元是手术间整体 |
| 输液室 | 座位/床位 | 几十个座位并发；按容量管理而非按设备时段切分 |
| 皮肤科光疗 | 光疗舱(共享) | 同一设备可同时服务多人；号源不等于独占设备 |
| 麻醉评估 | 诊间+医生 | 按"医生+时段"排班，类似门诊 |

### 1.2 招标依据

多份项目招标文件明确要求：

> "支持资源维护以**诊位**作为资源，同一诊室内多台检查设备支持每台设备作为资源单独排班生成号源" — 永春县医院、贵州省人民医院、海军安庆医院招标文件

> "资源可以根据配置对应到诊室、设备或者设备组等，可以合理自定义时段资源区间" — 海军安庆医院

这说明行业标准是**诊位是资源的统一抽象**，设备只是诊位的一种特例。

### 1.3 改造目标

引入 `Resource`（可预约资源/诊位）作为排班→号源→预约链路的统一主体，使 `Device` 成为 `Resource` 的一种子类型，同时兼容窗口、诊间、床位等非设备场景。

**核心原则：向上兼容，不破坏已有设备类资源的全部功能。**

---

## 2 术语定义

| 术语 | 定义 |
|------|------|
| **诊位 / 资源（Resource）** | 可被排班和产生号源的最小预约单元。是系统中"排班→号源→预约"链路的统一主体 |
| **资源类型（resource_type）** | 诊位的物理形态分类，决定号源生成策略和页面展示形态 |
| **并发容量（capacity）** | 该资源在同一号源时段内可同时服务的患者数。设备类通常=1，抽血窗口=1，输液座位区可能=30 |
| **占号模式（slot_mode）** | 号源分配策略。不同资源类型使用不同占号逻辑 |
| **诊室（Room）** | 物理空间。一个诊室内可包含一个或多个资源（如 CT 室内有 1 台 CT，或手术间内有多台设备） |

---

## 3 资源类型体系

### 3.1 类型枚举与特征

| resource_type | 中文名 | 典型科室 | 号源生成策略 | capacity 典型值 | 关联设备? | 关联医生? |
|---------------|-------|---------|-----------|---------------|---------|---------|
| `device` | 设备 | 放射、超声、内镜、核医学、功能检查 | 按检查项目标准耗时切分时段 | 1 | 是（1:1 绑定） | 可选 |
| `window` | 窗口/工位 | 检验科抽血、采样 | 按固定时长切分（如 5 分钟/号） | 1 | 否 | 可选（绑定操作护士） |
| `room` | 诊间/手术间 | 门诊手术、麻醉评估、诊间预约 | 按医生排班+时段切分 | 1 | 可选（间内有设备） | 是（必须绑定出诊医生） |
| `bed` | 床位/座位 | 输液室、康复治疗 | 按容量×时段生成（如 30 床×半小时=15 号/段） | N（多座位并发） | 可选（绑定治疗设备） | 可选（绑定治疗师） |
| `device_group` | 设备组 | 多台同类设备共享队列 | 号源绑定到组，执行时系统自动分配具体设备 | N（=组内设备数） | 是（1:N 绑定） | 否 |

### 3.2 类型详解与场景

#### 3.2.1 device（设备）— 现有逻辑不变

沿用当前模型。一台物理设备 = 一个资源，按检查项目耗时切分时段，每个号源独占设备。capacity=1。

**覆盖科室：** 放射科（CT/MRI/DR）、超声科、内镜中心、核医学（PET-CT）、功能检查科（心电图/脑电图/肺功能）。

#### 3.2.2 window（窗口）— 检验科抽血

一个窗口 = 一个资源。检验科有 5 个抽血窗口，则创建 5 个 window 类型资源。

**号源生成：** 管理员配置每个窗口的工作时段和每号耗时（如 5 分钟），系统按固定间隔切分。

**场景示例：** 检验科 1 号窗口，08:00-12:00 工作，5 分钟/号 → 生成 48 个号源。患者预约"血常规"检查项目 → 系统查询关联该项目的所有 window 资源 → 返回最早可用号源。

#### 3.2.3 room（诊间）— 门诊手术/麻醉评估

一个诊间 = 一个资源，但必须绑定出诊医生才能排班。

**号源生成：** 以"医生+诊间+时段"为组合排班。同一诊间上午由 A 医生出诊、下午由 B 医生出诊，产生两段排班。

**场景示例：** 麻醉评估诊间，08:00-12:00 张医生出诊，15 分钟/号 → 生成 16 个号源。患者预约"麻醉评估" → 号源带有诊间位置 + 出诊医生信息。

#### 3.2.4 bed（床位）— 输液室/康复

多个床位/座位组成一个资源，通过 capacity 表达并发量。

**号源生成：** 按容量并发 × 时段切分。输液室 30 个座位，每半小时一档 → 每个时段号源数 = 30。患者预约不指定具体座位号，签到时现场分配。

**场景示例：** 输液室资源，capacity=30，08:00-08:30 / 08:30-09:00 / ... 每段 30 个号源。

**康复场景变体：** 康复科治疗床 8 张，capacity=8，绑定治疗师。排班以治疗师为维度（李治疗师管 3 张床 → capacity=3 的 bed 资源）。

#### 3.2.5 device_group（设备组）— 同类设备共享队列

将多台同类设备编组为一个虚拟资源。号源绑定到组而非具体设备，执行时由系统自动分配空闲设备。

**场景示例：** CT 室有 3 台 CT，组成"CT 设备组"资源。排班对组整体排，号源总量 = 3 台设备的合并产能。患者预约"CT 平扫" → 锁定组内号源 → 签到时系统分配当前空闲的具体 CT 设备。

**适用场景：** 医院不关心具体用哪台设备，只关心总产能和等待时间。

---

## 4 数据模型设计

### 4.1 核心实体关系

```
Campus 1──N Department 1──N Room（诊室，物理空间）
                                  │
                           Room 1──N Resource（诊位，排班主体）
                                  │
                           Resource 1──N Schedule 1──N TimeSlot
                                  │
                           Resource N──N ExamItem（资源-项目关联）
                                  │
                           Resource 0..1── Device（设备类资源的扩展属性）
                           Resource 0..N── Doctor（诊间类资源绑定医生）
```

### 4.2 新增/变更实体

#### 4.2.1 Room（诊室 — 新增）

```sql
CREATE TABLE rooms (
    id            VARCHAR(36)  NOT NULL PRIMARY KEY,
    department_id VARCHAR(36)  NOT NULL,
    name          VARCHAR(50)  NOT NULL,           -- "CT 1号室" / "抽血窗口区" / "输液大厅"
    floor         VARCHAR(20),
    location_desc VARCHAR(200),                    -- 导航描述 "门诊楼B1层放射科左转"
    status        VARCHAR(10)  NOT NULL DEFAULT 'active',
    created_at    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
);
```

#### 4.2.2 Resource（诊位/可预约资源 — 核心新增）

```sql
CREATE TABLE resources (
    id              VARCHAR(36)  NOT NULL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,            -- 管理员可见名称
    resource_type   VARCHAR(20)  NOT NULL,            -- device / window / room / bed / device_group
    department_id   VARCHAR(36)  NOT NULL,
    campus_id       VARCHAR(36)  NOT NULL,
    room_id         VARCHAR(36),                      -- 所在诊室（可选）
    capacity        INT          NOT NULL DEFAULT 1,  -- 并发容量
    slot_mode       VARCHAR(20)  NOT NULL DEFAULT 'dynamic', -- dynamic(按耗时) / fixed(固定时长) / capacity(按容量)
    default_slot_minutes INT     NOT NULL DEFAULT 15, -- 固定模式下的默认号源时长
    max_daily_slots INT          NOT NULL DEFAULT 100,
    -- 设备类扩展属性（仅 resource_type=device/device_group 时使用）
    device_model         VARCHAR(50),
    device_manufacturer  VARCHAR(50),
    -- 通用属性
    status          VARCHAR(10)  NOT NULL DEFAULT 'active',
    created_at      DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at      DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    
    KEY idx_resources_dept (department_id),
    KEY idx_resources_type (resource_type),
    KEY idx_resources_room (room_id)
);
```

#### 4.2.3 resource_exam_items（资源-检查项目关联 — 新增）

替代原来 `Device.SupportedExamTypes` JSON 字段，改为标准多对多关联。

```sql
CREATE TABLE resource_exam_items (
    resource_id  VARCHAR(36) NOT NULL,
    exam_item_id VARCHAR(36) NOT NULL,
    PRIMARY KEY (resource_id, exam_item_id),
    KEY idx_rei_exam (exam_item_id)
);
```

#### 4.2.4 resource_doctors（资源-医生绑定 — 新增）

room 类型资源必须绑定出诊医生；其他类型可选。

```sql
CREATE TABLE resource_doctors (
    resource_id VARCHAR(36) NOT NULL,
    doctor_id   VARCHAR(36) NOT NULL,
    is_primary  TINYINT(1)  NOT NULL DEFAULT 0,  -- 主要负责人
    PRIMARY KEY (resource_id, doctor_id)
);
```

#### 4.2.5 device_group_members（设备组成员 — 新增）

```sql
CREATE TABLE device_group_members (
    group_resource_id  VARCHAR(36) NOT NULL,  -- device_group 类型的 resource.id
    member_resource_id VARCHAR(36) NOT NULL,  -- device 类型的 resource.id
    PRIMARY KEY (group_resource_id, member_resource_id)
);
```

#### 4.2.6 现有表变更

| 表 | 变更 |
|---|------|
| `schedules` | `device_id` → `resource_id`（重命名或新增列+迁移） |
| `time_slots` | `device_id` → `resource_id` |
| `appointment_items` | `device_id` → `resource_id` |
| `waiting_queues` | `device_id` → `resource_id` |
| `exam_executions` | `device_id` → `resource_id` |
| `devices` | **保留但降级为扩展表**，或直接废弃（字段合入 resources） |

---

## 5 号源生成策略

### 5.1 策略选择矩阵

| slot_mode | 适用 resource_type | 号源生成逻辑 |
|-----------|-------------------|-------------|
| `dynamic` | device | 按检查项目标准耗时逐个切分（现有逻辑不变）。支持混合排程模板。支持年龄折算。 |
| `fixed` | window, room | 按 `default_slot_minutes` 等间隔切分。每个号源时长固定。 |
| `capacity` | bed | 按时段分档（如半小时一档），每档号源数 = capacity。患者不指定具体座位。 |
| `dynamic` | device_group | 按组内设备合并时段计算，号源总量 = Σ各成员设备号源。执行时自动路由到空闲成员。 |

### 5.2 号源生成接口变更

```
POST /api/v1/resources/schedules/generate
```

```json
{
    "resource_id": "RES_CT_FORCE",      // 替代原来的 device_id
    "start_date": "2026-04-01",
    "end_date": "2026-04-30",
    "start_time": "08:00",
    "end_time": "12:00",
    "slot_minutes": 15,                  // fixed 模式使用；dynamic 模式忽略（由项目耗时决定）
    "exam_item_id": "EXAM_CT_PLAIN",     // dynamic 模式必填
    "pool_type": "public",
    "skip_weekends": true
}
```

系统根据 `resource.slot_mode` 自动选择对应的生成策略，调用方无需关心具体算法。

---

## 6 占号模式

### 6.1 模式定义

招标文件要求支持 6 种占号模式。在 Resource 模型下重新定义：

| 占号模式 | 编码 | 说明 | 典型场景 |
|---------|------|------|---------|
| 单系数模式 | `single_coefficient` | 号源数 = 工作时长 ÷ 标准耗时 × 系数 | 超声科，系数=0.9 预留空隙 |
| 对半分模式 | `half_split` | 号源总量按门诊/住院 50:50 拆分 | 通用 |
| 明细占号模式 | `detail` | 每个号源绑定具体检查项目和时段，不可混用 | MRI（平扫/增强时长差异大） |
| 双系数模式 | `dual_coefficient` | 门诊和住院分别使用不同系数计算号源 | 门诊系数0.8，住院系数1.0 |
| 全部模式 | `full` | 所有号源放入公共池，先到先得 | 心电图（无需区分来源） |
| 平均模式 | `average` | 号源按渠道数量平均分配 | 多渠道预约场景 |

### 6.2 配置方式

占号模式配置挂在 Resource 维度：

```sql
ALTER TABLE resources ADD COLUMN slot_allocation_mode VARCHAR(30) DEFAULT 'full';
ALTER TABLE resources ADD COLUMN slot_allocation_params JSON;  -- 系数、比例等参数
```

---

## 7 迁移策略：从 Device 到 Resource

### 7.1 核心原则

- **渐进式迁移**，不一次性全切。
- **现有 device 数据自动转为 resource**，业务无感。
- **API 双版本并行**，旧接口（`device_id`）短期内继续可用。

### 7.2 迁移步骤

**Step 1 — 建新表，不动旧表（1-2天）**

创建 `rooms`, `resources`, `resource_exam_items`, `resource_doctors`, `device_group_members` 五张新表。

**Step 2 — 数据迁移脚本（1天）**

```sql
-- 为每台 device 自动创建对应 resource
INSERT INTO resources (id, name, resource_type, department_id, campus_id, capacity, 
                       slot_mode, device_model, device_manufacturer, status)
SELECT CONCAT('RES_', id), name, 'device', department_id, campus_id, 1, 
       'dynamic', model, manufacturer, status
FROM devices;
```

**Step 3 — 排班/号源表增加 resource_id（1天）**

```sql
ALTER TABLE schedules ADD COLUMN resource_id VARCHAR(36);
ALTER TABLE time_slots ADD COLUMN resource_id VARCHAR(36);
ALTER TABLE appointment_items ADD COLUMN resource_id VARCHAR(36);

-- 回填
UPDATE schedules SET resource_id = CONCAT('RES_', device_id);
UPDATE time_slots SET resource_id = CONCAT('RES_', device_id);
UPDATE appointment_items SET resource_id = CONCAT('RES_', device_id);
```

**Step 4 — 业务代码切换到 resource_id（3-5天）**

领域实体、PO、Service、Handler 逐层替换。保留旧 `device_id` 列但标记为 deprecated。

**Step 5 — 清理旧列（后续版本）**

确认全部稳定后删除 `device_id` 列和 `devices` 表。

---

## 8 对现有子系统的影响

| 子系统 | 影响范围 | 改造内容 |
|--------|---------|---------|
| **资源管理** | 高 | Device CRUD → Resource CRUD；新增 Room CRUD；排班日历纵轴从"设备列表"改为"资源列表"（按类型分组） |
| **规则引擎** | 低 | 冲突/依赖规则绑定的是 exam_item_id，不涉及 device/resource；患者属性适配中 `filter_device` 改为 `filter_resource` |
| **预约服务** | 中 | 号源查询从 `device_id` 改为 `resource_id`；预约引擎按 resource_type 选择号源匹配策略 |
| **分诊执行** | 中 | 候诊队列从 `device_id` 改为 `resource_id`；检查执行记录同步改 |
| **前端页面** | 高 | 设备管理页面重构为"资源管理"；排班日历支持多类型资源切换视图 |

---

## 9 前端页面变更

### 9.1 资源管理（替代原设备管理）

**路由：** `/resource/manage`（替代 `/resource/device`）

**页面结构：**
- 顶部 Tab 按资源类型切换：全部 / 设备 / 窗口 / 诊间 / 床位 / 设备组
- 表格列：资源名称 / 类型(Tag) / 所属科室 / 所在诊室 / 关联检查项目(Tags) / 并发容量 / 占号模式 / 状态 / 操作
- 新建资源表单根据选择的类型动态展示不同字段

### 9.2 排班日历变更

纵轴从"设备列表"改为"资源列表"，支持按资源类型筛选。其余交互（拖拽、停诊、替班、追加）逻辑不变，仅主体 ID 从 device_id 切换为 resource_id。

---

## 10 约束与验证标准

| 约束类型 | 约束内容 |
|---------|---------|
| 数据约束 | resource_type 为枚举值，不可为空；capacity ≥ 1 |
| 数据约束 | device 类型 capacity 必须 = 1 |
| 数据约束 | room 类型资源必须绑定至少一位医生，否则不可排班 |
| 数据约束 | device_group 必须包含 ≥ 2 个 device 类型成员 |
| 功能约束 | 资源删除前须检查是否存在未来排班或已确认预约，有则拒绝删除 |
| 兼容约束 | 迁移后所有历史预约数据的 resource_id 必须有效且可回溯到原 device |
| 性能约束 | 资源列表查询响应 ≤ 200ms；排班生成性能不低于当前 Device 模式 |
| 验证标准 | 迁移完成后，对同一台 CT 设备执行排班→号源→预约→签到→完成全流程，结果与迁移前一致 |

---

## 11 里程碑计划

| 阶段 | 工期 | 内容 | 交付物 |
|------|------|------|--------|
| M1 数据模型 | 3天 | 新表 DDL、迁移脚本、Resource 领域实体/PO/Repository | Migration SQL + Go 代码 |
| M2 后端核心 | 5天 | Resource CRUD API、号源生成策略重构、排班接口切换 | 通过单元测试 |
| M3 前端适配 | 4天 | 资源管理页面重构、排班日历适配、预约流程适配 | 前端可操作 |
| M4 集成测试 | 3天 | 端到端测试（设备/窗口/诊间/床位四种类型全流程） | 测试报告 |
| M5 旧数据清理 | 2天 | 废弃 devices 表、删除 device_id 列 | 清理脚本 |
| **合计** | **17天** | | |

---

## 12 风险与对策

| 风险 | 可能性 | 影响 | 对策 |
|------|--------|------|------|
| 迁移脚本数据不一致 | 中 | 历史预约找不到资源 | 迁移前全量备份；迁移后自动校验脚本比对记录数 |
| 前端改动量超预期 | 中 | 工期延长 | M2/M3 并行开发；资源管理页面优先适配 device 类型，其他类型增量添加 |
| device_group 号源路由复杂 | 低 | 执行时分配逻辑 bug | device_group 作为 Phase 2 实现，Phase 1 先不开放 |
| 第三方系统（HIS/PACS）仍用 device_id | 高 | 接口不兼容 | 对外接口保留 device_id 字段（内部映射为 resource_id），实现透传兼容 |

---

## 13 评审检查清单

请评审人对以下问题逐项确认：

- [ ] 五种资源类型是否覆盖了目标医院的全部科室场景？是否有遗漏？
- [ ] Resource 表的字段设计是否满足所有类型的配置需求？是否有类型需要额外的扩展表？
- [ ] 号源生成三种策略（dynamic/fixed/capacity）是否能覆盖六种占号模式？
- [ ] 迁移策略（渐进式，保留 device_id 兼容期）是否可接受？还是应该一步到位？
- [ ] device_group 是否确实需要在 v0.2 实现，还是可以推迟到 v0.3？
- [ ] Room（诊室）作为独立实体是否有必要？还是直接作为 Resource 的一个属性即可？
- [ ] 17 天工期评估是否合理？
- [ ] 是否需要同步修改 HIS 数据同步接口以支持非设备类资源？