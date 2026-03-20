-- MTAP 数据库初始化迁移脚本
-- 适用于 MySQL 8.0+ / PostgreSQL 13+
-- 说明：若使用 GORM AutoMigrate（开发阶段），此文件仅供生产参考
-- 执行顺序：按文件名顺序执行

-- ============================================================
-- 1. 权限与用户
-- ============================================================
CREATE TABLE IF NOT EXISTS roles (
    id           VARCHAR(36)  NOT NULL,
    name         VARCHAR(30)  NOT NULL,
    permissions  TEXT         NOT NULL DEFAULT '[]',
    is_preset    TINYINT(1)   NOT NULL DEFAULT 0,
    created_at   DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at   DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uk_roles_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

CREATE TABLE IF NOT EXISTS users (
    id            VARCHAR(36)  NOT NULL,
    username      VARCHAR(50)  NOT NULL,
    password_hash VARCHAR(100) NOT NULL,
    real_name     VARCHAR(30),
    role_id       VARCHAR(36)  NOT NULL,
    department_id VARCHAR(36),
    status        VARCHAR(10)  NOT NULL DEFAULT 'active',
    last_login_at DATETIME(3),
    created_at    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uk_users_username (username),
    KEY idx_users_role_id (role_id),
    KEY idx_users_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ============================================================
-- 2. 规则引擎
-- ============================================================
CREATE TABLE IF NOT EXISTS conflict_rules (
    id            VARCHAR(36) NOT NULL,
    item_a_id     VARCHAR(36) NOT NULL,
    item_b_id     VARCHAR(36) NOT NULL,
    min_interval  INT         NOT NULL DEFAULT 0,
    interval_unit VARCHAR(10) NOT NULL DEFAULT 'hour',
    level         VARCHAR(10) NOT NULL DEFAULT 'warning',
    status        VARCHAR(10) NOT NULL DEFAULT 'active',
    created_by    VARCHAR(36),
    created_at    DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at    DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_conflict_rules_item_a (item_a_id),
    KEY idx_conflict_rules_item_b (item_b_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='冲突规则表';

CREATE TABLE IF NOT EXISTS conflict_packages (
    id            VARCHAR(36)  NOT NULL,
    name          VARCHAR(60)  NOT NULL,
    min_interval  INT          NOT NULL DEFAULT 0,
    interval_unit VARCHAR(10)  NOT NULL DEFAULT 'day',
    level         VARCHAR(10)  NOT NULL DEFAULT 'warning',
    status        VARCHAR(10)  NOT NULL DEFAULT 'active',
    created_at    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uk_conflict_packages_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='冲突包表';

CREATE TABLE IF NOT EXISTS conflict_package_items (
    id             VARCHAR(36) NOT NULL,
    package_id     VARCHAR(36) NOT NULL,
    exam_item_id   VARCHAR(36) NOT NULL,
    created_at     DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_pkg_items_package (package_id),
    KEY idx_pkg_items_exam_item (exam_item_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='冲突包-项目关联表';

CREATE TABLE IF NOT EXISTS dependency_rules (
    id             VARCHAR(36) NOT NULL,
    pre_item_id    VARCHAR(36) NOT NULL,
    post_item_id   VARCHAR(36) NOT NULL,
    type           VARCHAR(15) NOT NULL DEFAULT 'mandatory',
    validity_hours INT         NOT NULL DEFAULT 72,
    status         VARCHAR(10) NOT NULL DEFAULT 'active',
    created_at     DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at     DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_dep_rules_pre (pre_item_id),
    KEY idx_dep_rules_post (post_item_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='依赖规则表';

CREATE TABLE IF NOT EXISTS priority_tags (
    id         VARCHAR(36) NOT NULL,
    name       VARCHAR(30) NOT NULL,
    weight     INT         NOT NULL DEFAULT 0,
    color      VARCHAR(20),
    is_preset  TINYINT(1)  NOT NULL DEFAULT 0,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uk_priority_tags_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='优先级标签表';

CREATE TABLE IF NOT EXISTS sorting_strategies (
    id         VARCHAR(36)  NOT NULL,
    type       VARCHAR(20)  NOT NULL,
    scope      JSON,
    start_date DATETIME(3),
    end_date   DATETIME(3),
    status     VARCHAR(10)  NOT NULL DEFAULT 'active',
    created_at DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='排序策略表';

CREATE TABLE IF NOT EXISTS patient_adapt_rules (
    id              VARCHAR(36)  NOT NULL,
    condition_type  VARCHAR(20)  NOT NULL,
    condition_value VARCHAR(100) NOT NULL,
    action          VARCHAR(30)  NOT NULL,
    action_params   JSON,
    priority        INT          NOT NULL DEFAULT 0,
    status          VARCHAR(10)  NOT NULL DEFAULT 'active',
    created_at      DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at      DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='患者属性适配规则表';

CREATE TABLE IF NOT EXISTS source_controls (
    id                    VARCHAR(36)  NOT NULL,
    source_type           VARCHAR(20)  NOT NULL,
    slot_pool_id          VARCHAR(36),
    allocation_ratio      DECIMAL(5,2) NOT NULL DEFAULT 1.00,
    overflow_enabled      TINYINT(1)   NOT NULL DEFAULT 0,
    overflow_target_pool  VARCHAR(36),
    status                VARCHAR(10)  NOT NULL DEFAULT 'active',
    created_at            DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at            DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='开单来源控制表';

-- ============================================================
-- 3. 资源管理
-- ============================================================
CREATE TABLE IF NOT EXISTS campuses (
    id         VARCHAR(36)  NOT NULL,
    name       VARCHAR(50)  NOT NULL,
    code       VARCHAR(20)  NOT NULL,
    address    VARCHAR(200),
    status     VARCHAR(10)  NOT NULL DEFAULT 'active',
    created_at DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uk_campuses_name (name),
    UNIQUE KEY uk_campuses_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='院区表';

CREATE TABLE IF NOT EXISTS departments (
    id         VARCHAR(36)  NOT NULL,
    campus_id  VARCHAR(36)  NOT NULL,
    name       VARCHAR(50)  NOT NULL,
    code       VARCHAR(20)  NOT NULL,
    floor      VARCHAR(20),
    status     VARCHAR(10)  NOT NULL DEFAULT 'active',
    synced_at  DATETIME(3),
    created_at DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uk_departments_code (code),
    KEY idx_departments_campus (campus_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='科室表';

CREATE TABLE IF NOT EXISTS devices (
    id            VARCHAR(36)  NOT NULL,
    name          VARCHAR(100) NOT NULL,
    campus_id     VARCHAR(36),
    department_id VARCHAR(36),
    status        VARCHAR(10)  NOT NULL DEFAULT 'active',
    created_at    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_devices_department (department_id),
    KEY idx_devices_campus (campus_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设备表';

CREATE TABLE IF NOT EXISTS exam_items (
    id           VARCHAR(36)  NOT NULL,
    name         VARCHAR(100) NOT NULL,
    duration_min INT          NOT NULL DEFAULT 30,
    is_fasting   TINYINT(1)   NOT NULL DEFAULT 0,
    fasting_desc VARCHAR(200),
    created_at   DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at   DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uk_exam_items_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='检查项目表';

CREATE TABLE IF NOT EXISTS item_aliases (
    id           VARCHAR(36)  NOT NULL,
    exam_item_id VARCHAR(36)  NOT NULL,
    alias        VARCHAR(50)  NOT NULL,
    created_at   DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uk_item_aliases_alias (alias),
    KEY idx_item_aliases_exam_item (exam_item_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目别名表';

CREATE TABLE IF NOT EXISTS slot_pools (
    id         VARCHAR(36)  NOT NULL,
    name       VARCHAR(60)  NOT NULL,
    type       VARCHAR(20)  NOT NULL,
    status     VARCHAR(10)  NOT NULL DEFAULT 'active',
    created_at DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uk_slot_pools_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='号源池表';

CREATE TABLE IF NOT EXISTS schedules (
    id             VARCHAR(36)  NOT NULL,
    device_id      VARCHAR(36)  NOT NULL,
    date           DATE         NOT NULL,
    start_time     VARCHAR(5)   NOT NULL COMMENT 'HH:mm',
    end_time       VARCHAR(5)   NOT NULL COMMENT 'HH:mm',
    status         VARCHAR(15)  NOT NULL DEFAULT 'normal',
    suspend_reason VARCHAR(200),
    created_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_schedules_device_date (device_id, date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='排班表';

CREATE TABLE IF NOT EXISTS time_slots (
    id                VARCHAR(36)  NOT NULL,
    device_id         VARCHAR(36)  NOT NULL,
    date              DATE         NOT NULL,
    exam_item_id      VARCHAR(36),
    pool_type         VARCHAR(15)  NOT NULL DEFAULT 'public',
    start_at          DATETIME(3)  NOT NULL,
    end_at            DATETIME(3)  NOT NULL,
    standard_duration INT          NOT NULL DEFAULT 0,
    adjusted_duration INT          NOT NULL DEFAULT 0,
    status            VARCHAR(15)  NOT NULL DEFAULT 'available',
    locked_by         VARCHAR(36),
    lock_until        DATETIME(3),
    created_at        DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at        DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_time_slots_device_date (device_id, date),
    KEY idx_time_slots_status (status),
    KEY idx_time_slots_exam_item (exam_item_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='号源时段表';

-- ============================================================
-- 4. 预约服务
-- ============================================================
CREATE TABLE IF NOT EXISTS appointments (
    id               VARCHAR(36)  NOT NULL,
    patient_id       VARCHAR(36)  NOT NULL,
    mode             VARCHAR(10)  NOT NULL,
    status           VARCHAR(15)  NOT NULL DEFAULT 'pending',
    override_by      VARCHAR(36),
    override_reason  VARCHAR(200),
    payment_verified TINYINT(1)   NOT NULL DEFAULT 0,
    change_count     INT          NOT NULL DEFAULT 0,
    cancel_reason    VARCHAR(200),
    confirmed_at     DATETIME(3),
    cancelled_at     DATETIME(3),
    created_at       DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at       DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_appointments_patient (patient_id),
    KEY idx_appointments_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='预约单表';

CREATE TABLE IF NOT EXISTS appointment_items (
    id             VARCHAR(36)  NOT NULL,
    appointment_id VARCHAR(36)  NOT NULL,
    exam_item_id   VARCHAR(36)  NOT NULL,
    slot_id        VARCHAR(36)  NOT NULL,
    device_id      VARCHAR(36)  NOT NULL,
    start_time     DATETIME(3)  NOT NULL,
    end_time       DATETIME(3)  NOT NULL,
    status         VARCHAR(15)  NOT NULL DEFAULT 'pending',
    created_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_appt_items_appointment (appointment_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='预约项目表';

CREATE TABLE IF NOT EXISTS appointment_credentials (
    id                 VARCHAR(36)  NOT NULL,
    appointment_id     VARCHAR(36)  NOT NULL,
    qr_code_data       TEXT         NOT NULL,
    patient_name_masked VARCHAR(30),
    exam_summary       TEXT,
    notice_content     TEXT,
    generated_at       DATETIME(3)  NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_credentials_appointment (appointment_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='预约凭证表';

CREATE TABLE IF NOT EXISTS appointment_change_logs (
    id             VARCHAR(36)  NOT NULL,
    appointment_id VARCHAR(36)  NOT NULL,
    change_type    VARCHAR(15)  NOT NULL,
    old_slot_id    VARCHAR(36),
    new_slot_id    VARCHAR(36),
    reason         VARCHAR(200),
    operator_id    VARCHAR(36)  NOT NULL,
    changed_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_change_logs_appointment (appointment_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='预约变更日志表';

CREATE TABLE IF NOT EXISTS blacklists (
    id             VARCHAR(36)  NOT NULL,
    patient_id     VARCHAR(36)  NOT NULL,
    trigger_time   DATETIME(3)  NOT NULL,
    expires_at     DATETIME(3)  NOT NULL,
    status         VARCHAR(10)  NOT NULL DEFAULT 'active',
    released_at    DATETIME(3),
    release_reason VARCHAR(200),
    created_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_blacklists_patient (patient_id),
    KEY idx_blacklists_expires (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='黑名单表';

CREATE TABLE IF NOT EXISTS no_show_records (
    id             VARCHAR(36)  NOT NULL,
    patient_id     VARCHAR(36)  NOT NULL,
    appointment_id VARCHAR(36)  NOT NULL,
    occurred_at    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_noshow_patient (patient_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='爽约记录表';

CREATE TABLE IF NOT EXISTS appeals (
    id           VARCHAR(36)  NOT NULL,
    blacklist_id VARCHAR(36)  NOT NULL,
    reason       VARCHAR(500) NOT NULL,
    status       VARCHAR(10)  NOT NULL DEFAULT 'pending',
    reviewed_by  VARCHAR(36),
    reviewed_at  DATETIME(3),
    created_at   DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_appeals_blacklist (blacklist_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='申诉表';

-- ============================================================
-- 5. 分诊执行
-- ============================================================
CREATE TABLE IF NOT EXISTS check_ins (
    id             VARCHAR(36)  NOT NULL,
    appointment_id VARCHAR(36)  NOT NULL,
    patient_id     VARCHAR(36)  NOT NULL,
    method         VARCHAR(10)  NOT NULL,
    check_in_time  DATETIME(3)  NOT NULL,
    is_late        TINYINT(1)   NOT NULL DEFAULT 0,
    remark         VARCHAR(100),
    created_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uk_checkins_appointment (appointment_id),
    KEY idx_checkins_patient (patient_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='签到记录表';

CREATE TABLE IF NOT EXISTS waiting_queues (
    id            VARCHAR(36)  NOT NULL,
    room_id       VARCHAR(36)  NOT NULL,
    device_id     VARCHAR(36)  NOT NULL,
    department_id VARCHAR(36)  NOT NULL,
    status        VARCHAR(10)  NOT NULL DEFAULT 'active',
    created_at    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uk_queues_room (room_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='候诊队列表';

CREATE TABLE IF NOT EXISTS queue_entries (
    id                  VARCHAR(36)  NOT NULL,
    queue_id            VARCHAR(36)  NOT NULL,
    patient_id          VARCHAR(36)  NOT NULL,
    patient_name_masked VARCHAR(30),
    appointment_id      VARCHAR(36)  NOT NULL,
    check_in_id         VARCHAR(36)  NOT NULL,
    queue_number        INT          NOT NULL,
    status              VARCHAR(15)  NOT NULL DEFAULT 'waiting',
    call_count          INT          NOT NULL DEFAULT 0,
    miss_count          INT          NOT NULL DEFAULT 0,
    entered_at          DATETIME(3)  NOT NULL,
    called_at           DATETIME(3),
    completed_at        DATETIME(3),
    created_at          DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at          DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_queue_entries_queue (queue_id),
    KEY idx_queue_entries_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='队列条目表';

CREATE TABLE IF NOT EXISTS exam_executions (
    id                  VARCHAR(36)  NOT NULL,
    appointment_item_id VARCHAR(36)  NOT NULL,
    patient_id          VARCHAR(36)  NOT NULL,
    device_id           VARCHAR(36)  NOT NULL,
    status              VARCHAR(15)  NOT NULL DEFAULT 'checked_in',
    started_at          DATETIME(3),
    completed_at        DATETIME(3),
    duration            INT          NOT NULL DEFAULT 0,
    operator_id         VARCHAR(36),
    undo_deadline       DATETIME(3),
    created_at          DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at          DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uk_executions_appt_item (appointment_item_id),
    KEY idx_executions_device (device_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='检查执行表';

-- ============================================================
-- 6. 统计分析
-- ============================================================
CREATE TABLE IF NOT EXISTS dashboard_snapshots (
    id           VARCHAR(36)  NOT NULL,
    campus_id    VARCHAR(36),
    snapshot     JSON         NOT NULL,
    snapshotted_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_snapshots_campus (campus_id),
    KEY idx_snapshots_at (snapshotted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='大屏快照表';

CREATE TABLE IF NOT EXISTS reports (
    id           VARCHAR(36)  NOT NULL,
    report_type  VARCHAR(20)  NOT NULL,
    dimensions   JSON,
    date_start   DATE         NOT NULL,
    date_end     DATE         NOT NULL,
    status       VARCHAR(15)  NOT NULL DEFAULT 'generating',
    file_path    VARCHAR(500),
    file_size    BIGINT       NOT NULL DEFAULT 0,
    format       VARCHAR(10)  NOT NULL DEFAULT 'xlsx',
    generated_at DATETIME(3),
    created_at   DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_reports_type (report_type),
    KEY idx_reports_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='报表表';

-- ============================================================
-- 7. 效能优化
-- ============================================================
CREATE TABLE IF NOT EXISTS efficiency_metrics (
    id           VARCHAR(36)  NOT NULL,
    name         VARCHAR(60)  NOT NULL,
    code         VARCHAR(30)  NOT NULL,
    calc_formula TEXT,
    unit         VARCHAR(20),
    normal_mean  DECIMAL(10,4),
    normal_stddev DECIMAL(10,4),
    normal_min   DECIMAL(10,4),
    normal_max   DECIMAL(10,4),
    is_custom    TINYINT(1)   NOT NULL DEFAULT 0,
    created_at   DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at   DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uk_metrics_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='效率指标定义表';

CREATE TABLE IF NOT EXISTS metric_snapshots (
    id         VARCHAR(36)   NOT NULL,
    metric_id  VARCHAR(36)   NOT NULL,
    value      DECIMAL(12,4) NOT NULL,
    dimensions JSON,
    sampled_at DATETIME(3)   NOT NULL,
    created_at DATETIME(3)   NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_snapshots_metric_at (metric_id, sampled_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='指标快照表';

CREATE TABLE IF NOT EXISTS bottleneck_alerts (
    id                   VARCHAR(36)   NOT NULL,
    metric_id            VARCHAR(36)   NOT NULL,
    alert_type           VARCHAR(20)   NOT NULL,
    severity             VARCHAR(10)   NOT NULL DEFAULT 'medium',
    deviation_pct        DECIMAL(8,2)  NOT NULL DEFAULT 0,
    consecutive_count    INT           NOT NULL DEFAULT 1,
    affected_scope       VARCHAR(200),
    root_cause_hypotheses JSON,
    suggested_category   VARCHAR(10),
    status               VARCHAR(15)   NOT NULL DEFAULT 'open',
    dismiss_reason       VARCHAR(200),
    created_at           DATETIME(3)   NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at           DATETIME(3)   NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_alerts_metric (metric_id),
    KEY idx_alerts_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='瓶颈告警表';

CREATE TABLE IF NOT EXISTS optimization_strategies (
    id               VARCHAR(36)  NOT NULL,
    title            VARCHAR(100) NOT NULL,
    category         VARCHAR(10)  NOT NULL,
    status           VARCHAR(20)  NOT NULL DEFAULT 'pending_review',
    alert_id         VARCHAR(36),
    current_value    VARCHAR(200),
    target_value     VARCHAR(200),
    expected_benefit VARCHAR(500),
    risk_note        VARCHAR(500),
    cost_estimate    JSON,
    approval_flow    JSON,
    created_at       DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at       DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_strategies_status (status),
    KEY idx_strategies_category (category)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='优化策略表';

CREATE TABLE IF NOT EXISTS trial_runs (
    id                          VARCHAR(36)   NOT NULL,
    strategy_id                 VARCHAR(36)   NOT NULL,
    gray_scope                  JSON,
    trial_days                  INT           NOT NULL DEFAULT 7,
    baseline                    JSON,
    started_at                  DATETIME(3)   NOT NULL,
    ends_at                     DATETIME(3)   NOT NULL,
    status                      VARCHAR(15)   NOT NULL DEFAULT 'running',
    emergency_rollback_threshold DECIMAL(8,2)  NOT NULL DEFAULT 20,
    created_at                  DATETIME(3)   NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at                  DATETIME(3)   NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uk_trial_runs_strategy (strategy_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='试运行记录表';

CREATE TABLE IF NOT EXISTS evaluation_reports (
    id              VARCHAR(36)  NOT NULL,
    strategy_id     VARCHAR(36)  NOT NULL,
    trial_run_id    VARCHAR(36),
    baseline_metrics JSON,
    trial_metrics    JSON,
    improvement_pct  DECIMAL(8,2),
    is_promoted      TINYINT(1)   NOT NULL DEFAULT 0,
    conclusion       TEXT,
    generated_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_eval_reports_strategy (strategy_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评估报告表';

CREATE TABLE IF NOT EXISTS roi_reports (
    id              VARCHAR(36)  NOT NULL,
    strategy_id     VARCHAR(36)  NOT NULL,
    cost_estimate   DECIMAL(12,2),
    benefit_estimate DECIMAL(12,2),
    payback_months  DECIMAL(8,2),
    detail          JSON,
    generated_at    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_roi_reports_strategy (strategy_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ROI报告表';

CREATE TABLE IF NOT EXISTS performance_scans (
    id           VARCHAR(36)  NOT NULL,
    scan_type    VARCHAR(20)  NOT NULL DEFAULT 'weekly',
    findings     JSON,
    opportunities JSON,
    scanned_at   DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_scans_at (scanned_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='效能扫描结果表';

CREATE TABLE IF NOT EXISTS strategy_decay_alerts (
    id              VARCHAR(36)   NOT NULL,
    strategy_id     VARCHAR(36)   NOT NULL,
    decay_pct       DECIMAL(8,2)  NOT NULL,
    original_benefit DECIMAL(8,2),
    current_benefit  DECIMAL(8,2),
    status          VARCHAR(15)   NOT NULL DEFAULT 'open',
    detected_at     DATETIME(3)   NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_decay_alerts_strategy (strategy_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='策略衰减告警表';

-- ============================================================
-- 8. 审计日志
-- ============================================================
CREATE TABLE IF NOT EXISTS audit_logs (
    id            VARCHAR(36)  NOT NULL,
    operator_id   VARCHAR(36)  NOT NULL,
    operator_name VARCHAR(30),
    action        VARCHAR(20)  NOT NULL,
    resource      VARCHAR(30)  NOT NULL,
    resource_id   VARCHAR(36)  NOT NULL,
    old_value     TEXT,
    new_value     TEXT,
    ip            VARCHAR(45),
    created_at    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_audit_resource (resource, resource_id),
    KEY idx_audit_operator (operator_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='审计日志表';

-- ============================================================
-- 9. 种子数据：预置角色 + 默认管理员
-- ============================================================
INSERT IGNORE INTO roles (id, name, permissions, is_preset) VALUES
('00000000-0000-0000-0000-000000000001', 'admin',           '["*"]',                                   1),
('00000000-0000-0000-0000-000000000002', 'scheduler_admin', '["rule:*","resource:*","appt:*"]',        1),
('00000000-0000-0000-0000-000000000003', 'operator',        '["appt:*","triage:*"]',                   1),
('00000000-0000-0000-0000-000000000004', 'nurse',           '["triage:*"]',                            1),
('00000000-0000-0000-0000-000000000005', 'viewer',          '["*:read"]',                              1);

-- 默认管理员密码: Admin@1234  (bcrypt hash, cost=10)
-- 生产环境请立即修改
INSERT IGNORE INTO users (id, username, password_hash, real_name, role_id, status)
VALUES (
    UUID(),
    'admin',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    '系统管理员',
    '00000000-0000-0000-0000-000000000001',
    'active'
);
