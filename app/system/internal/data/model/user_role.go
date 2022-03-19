package model

import "gorm.io/gorm"

type UserRole struct {
	gorm.Model
	UserID uint64 `gorm:"index;default:0;not null;"` // 用户内码
	RoleID uint64 `gorm:"index;default:0;not null;"` // 角色内码
}
