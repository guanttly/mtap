-- MTAP 增量迁移脚本 002
-- 日期: 2026-03-23
-- 描述: 补齐缺失字段、修复命名不一致、新增 doctors / schedule_templates 表
-- 执行方式: 在 001_init_schema.sql 执行后顺序执行

-- ============================================================
-- 1. devices 表补齐关键字段（#1 P0）
-- ============================================================
ALTER TABLE devices
  ADD COLUMN IF NOT EXISTS model              VARCHAR(50)  AFTER name,
  ADD COLUMN IF NOT EXISTS manufacturer      VARCHAR(50)  AFTER model,
  ADD COLUMN IF NOT EXISTS supported_exam_types JSON       AFTER manufacturer,
  ADD COLUMN IF NOT EXISTS max_daily_slots   INT NOT NULL DEFAULT 50 AFTER supported_exam_types;

-- ============================================================
-- 2. slot_pools 表补齐配额与溢出字段（#2 P0）
-- ============================================================
ALTER TABLE slot_pools
  ADD COLUMN IF NOT EXISTS allocation_ratio  DECIMAL(5,2) NOT NULL DEFAULT 0.00 AFTER status,
  ADD COLUMN IF NOT EXISTS overflow_enabled  TINYINT(1)   NOT NULL DEFAULT 0    AFTER allocation_ratio,
  ADD COLUMN IF NOT EXISTS overflow_target_pool VARCHAR(36)                     AFTER overflow_enabled;

-- ============================================================
-- 3. sorting_strategies 表：单 JSON scope → 三独立列（#3 P1）
-- ============================================================
ALTER TABLE sorting_strategies
  ADD COLUMN IF NOT EXISTS scope_campuses TEXT AFTER scope,
  ADD COLUMN IF NOT EXISTS scope_depts    TEXT AFTER scope_campuses,
  ADD COLUMN IF NOT EXISTS scope_devices  TEXT AFTER scope_depts;

-- 迁移已有数据：从 scope JSON 拆解到三列
UPDATE sorting_strategies
SET
  scope_campuses = JSON_UNQUOTE(JSON_EXTRACT(scope, '$.campuses')),
  scope_depts    = JSON_UNQUOTE(JSON_EXTRACT(scope, '$.depts')),
  scope_devices  = JSON_UNQUOTE(JSON_EXTRACT(scope, '$.devices'))
WHERE scope IS NOT NULL;

-- 删除旧 scope 列（若数据库支持 IF EXISTS）
-- MySQL 8.0 支持：ALTER TABLE ... DROP COLUMN IF EXISTS
ALTER TABLE sorting_strategies
  DROP COLUMN IF EXISTS scope;

-- ============================================================
-- 4. source_controls 表：列名对齐（#4 P1）
--    PO 使用 overflow_target_pool，Migration 已是此名，无需改动
--    若旧环境存在 overflow_target_pool_id 列，则执行以下重命名：
-- ============================================================
-- (幂等处理：若列已存在正确名称则跳过)
-- ALTER TABLE source_controls RENAME COLUMN overflow_target_pool_id TO overflow_target_pool;
-- 注：此行为注释形式保留，由 DBA 在存在旧列的环境手动执行

-- ============================================================
-- 5. patient_adapt_rules condition_value 长度对齐（#5 P2）
-- ============================================================
ALTER TABLE patient_adapt_rules
  MODIFY COLUMN condition_value VARCHAR(100) NOT NULL;

-- ============================================================
-- 6. roles 表补充 status 字段（#6 P2）
-- ============================================================
ALTER TABLE roles
  ADD COLUMN IF NOT EXISTS status VARCHAR(10) NOT NULL DEFAULT 'active' AFTER is_preset;

-- ============================================================
-- 7. doctors 表创建（#7 P1）
-- ============================================================
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
    KEY idx_doctors_department (department_id),
    KEY idx_doctors_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='医生表';

-- ============================================================
-- 8. schedule_templates 表创建（#9 P1）
-- ============================================================
CREATE TABLE IF NOT EXISTS schedule_templates (
    id           VARCHAR(36)  NOT NULL,
    name         VARCHAR(50)  NOT NULL,
    repeat_type  VARCHAR(20)  NOT NULL COMMENT 'once/daily/weekly',
    slot_pattern JSON         NOT NULL COMMENT '{"start_time":"08:00","end_time":"12:00","slot_minutes":30,"exam_item_id":"...","pool_type":"public"}',
    skip_weekends TINYINT(1)  NOT NULL DEFAULT 0,
    created_at   DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at   DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uk_schedule_templates_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='排班模板表';
