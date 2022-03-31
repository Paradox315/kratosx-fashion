package model

import "gorm.io/gorm"

type ResourceType uint8

// RoleResource
// 角色资源关联表
type RoleResource struct {
	gorm.Model
	RoleID     uint64       `gorm:"index;default:0;not null;comment:角色ID"`                 // 角色ID
	ResourceID uint64       `gorm:"index;default:0;not null;comment:资源ID"`                 // 资源ID
	Type       ResourceType `gorm:"index;type:tinyint(1);default:0;not null;comment:资源类型"` // 资源类型
}
