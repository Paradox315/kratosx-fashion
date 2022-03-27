package model

import "gorm.io/gorm"

type RoleStatus uint8

const (
	RoleStatusActive RoleStatus = iota + 1
	RoleStatusInactive
)

func (s RoleStatus) String() string {
	switch s {
	case RoleStatusActive:
		return "active"
	case RoleStatusInactive:
		return "inactive"
	default:
		return "unknown"
	}
}

type Role struct {
	gorm.Model
	Name        string `gorm:"size:100;index;default:'';not null;comment:角色名称"`          // 角色名称
	Sort        uint32 `gorm:"index;default:0;not null;comment:排序值"`                     // 排序值
	Description string `gorm:"size:1024;default:'';not null;comment:备注"`                 // 备注
	Status      uint8  `gorm:"default:1;type:tinyint(1);not null;comment:状态(1:启用 2:禁用)"` // 状态(1:启用 2:禁用)
	CreatorID   uint64 `gorm:"default:0;not null;comment:创建者ID"`                         // 创建者
	Menus       []Menu `gorm:"many2many:role_menus;"`                                    // 角色菜单
}
