package model

// Role
// 角色表
type Role struct {
	Model
	Name        string `gorm:"size:100;index;default:'';not null;comment:角色名称"` // 角色名称
	Description string `gorm:"size:1024;default:'';not null;comment:备注"`        // 备注
}
