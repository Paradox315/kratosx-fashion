package model

import "gorm.io/gorm"

type RoleMenu struct {
	gorm.Model
	RoleID   uint64 `gorm:"index;default:0;not null;"` // 角色ID
	MenuID   uint64 `gorm:"index;default:0;not null;"` // 菜单ID
	ActionID uint64 `gorm:"index;default:0;not null;"` // 动作ID
}
