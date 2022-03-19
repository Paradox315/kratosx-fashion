package model

import "gorm.io/gorm"

type Menu struct {
	gorm.Model
	Name        string       `gorm:"size:50;index;default:'';not null;comment:菜单名称"`  // 菜单名称
	Icon        string       `gorm:"size:255;default:'';not null;comment:菜单图标"`       // 菜单图标
	Router      string       `gorm:"size:255;default:'';not null;comment:服务端路由路径URL"` // 访问路由组(后端)
	Component   string       `gorm:"size:255;default:'';not null;comment:对应前端的路径"`    //对应前端的路径
	ParentID    uint64       `gorm:"index;default:0;not null;comment:父级内码"`           // 父级内码
	Hidden      uint8        `gorm:"index;default:0;type:tinyint(1);not null;comment:是否显示"`           // 是否显示
	KeepAlive   uint8        `gorm:"index;default:0;type:tinyint(1);not null;comment:是否缓存"`           //是否缓存
	Status      uint8        `gorm:"index;default:0;type:tinyint(1);not null;comment:状态(1:启用 2:禁用)"`  // 状态(1:启用 2:禁用)
	Sort        uint32       `gorm:"index;default:0;not null;comment:排序值"`            // 排序值
	Description string       `gorm:"size:1024;default:'';not null;comment:备注"`        // 备注
	CreatorID   uint64       `gorm:"default:0;not null;comment:创建人ID"`                // 创建人ID
	MenuActions []MenuAction `gorm:"foreignkey:MenuID"`
}
