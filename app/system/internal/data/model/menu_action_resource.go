package model

import "gorm.io/gorm"

type MenuActionResource struct {
	gorm.Model
	ActionID uint64 `gorm:"index;default:0;not null;comment:菜单对应的action_id"`      // 菜单动作ID
	Method   string `gorm:"size:50;default:'';not null;comment:资源请求方式(REST|RPC)"` // 资源请求方式(支持正则)
	Path     string `gorm:"size:255;default:'';not null;comment:资源请求路径（支持/:id匹配"` // 资源请求路径（支持/:id匹配）
}
