package model

import "gorm.io/gorm"

// Role
// 角色表
type Role struct {
	gorm.Model
	Name        string         `gorm:"size:100;index;default:'';not null;comment:角色名称"` // 角色名称
	Description string         `gorm:"size:1024;default:'';not null;comment:备注"`        // 备注
	Menus       []ResourceMenu `gorm:"many2many:role_resource_menus;"`                  // 角色菜单
}
