package model

import "gorm.io/gorm"

// UserRole
// 用户角色关联表
type UserRole struct {
	gorm.Model
	UserID uint64 `gorm:"index;default:0;not null;comment:用户ID"` // 用户内码
	RoleID uint64 `gorm:"index;default:0;not null;comment:角色ID"` // 角色内码
}
