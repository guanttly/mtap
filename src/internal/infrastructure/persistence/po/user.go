// Package po 基础设施层 - 用户/角色持久化对象
package po

import "time"

// RolePO 角色表
type RolePO struct {
	ID          string    `gorm:"column:id;primaryKey;size:36"`
	Name        string    `gorm:"column:name;not null;size:30;uniqueIndex"`
	Permissions string    `gorm:"column:permissions;type:text;not null"` // JSON数组，写入时由应用层保证不为空（最小值 "[]"）
	IsPreset    bool      `gorm:"column:is_preset;not null;default:false"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (RolePO) TableName() string { return "roles" }

// UserPO 用户表
type UserPO struct {
	ID           string     `gorm:"column:id;primaryKey;size:36"`
	Username     string     `gorm:"column:username;not null;size:50;uniqueIndex"`
	PasswordHash string     `gorm:"column:password_hash;not null;size:100"`
	RealName     string     `gorm:"column:real_name;size:30"`
	RoleID       string     `gorm:"column:role_id;not null;size:36;index"`
	RoleName     string     `gorm:"-"` // 关联查询时填充
	DepartmentID string     `gorm:"column:department_id;size:36"`
	Status       string     `gorm:"column:status;not null;size:10;default:active;index"`
	LastLoginAt  *time.Time `gorm:"column:last_login_at"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (UserPO) TableName() string { return "users" }
