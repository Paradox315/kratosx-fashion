package model

import "gorm.io/gorm"

type ResourceMenu struct {
	gorm.Model
	Name        string `gorm:"size:50;index;default:'';not null;comment:菜单名称"`         // 菜单名称
	Path        string `gorm:"size:255;default:'';not null;comment:服务端路由路径URL"`        // 访问路由组(后端)
	Component   string `gorm:"size:255;default:'';not null;comment:对应前端的路径"`           //对应前端的路径
	ParentID    uint   `json:"parent_id" gorm:"index;default:0;not null;comment:父级内码"` // 父级内码
	Locale      string `gorm:"size:50;default:'';not null;comment:语言"`                 // 语言
	RequireAuth uint8  `gorm:"default:1;type:tinyint(1);not null;comment:是否需要权限"`      // 是否需要权限
	Icon        string `gorm:"size:50;default:'';not null;comment:图标"`                 // 图标
	Order       uint32 `gorm:"default:0;type:tinyint(1);not null;comment:排序"`          // 排序
	HideInMenu  uint8  `gorm:"default:0;type:tinyint(1);not null;comment:是否在菜单中隐藏"`    // 是否在菜单中隐藏
	NoAffix     uint8  `gorm:"default:0;type:tinyint(1);not null;comment:是否不附加到菜单"`    // 是否不附加到菜单
	IgnoreCache uint8  `gorm:"default:0;type:tinyint(1);not null;comment:是否忽略缓存"`      // 是否忽略缓存
	Actions     string `gorm:"size:1024;default:'';not null;comment:前端动作按钮列表，json存储"`  //动作按钮
}
