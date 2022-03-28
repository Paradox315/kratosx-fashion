package model

import "gorm.io/gorm"

type ResourceAction struct {
	gorm.Model
	MenuID uint64 `gorm:"index;not null;comment:菜单ID"`               // 菜单ID
	Name   string `gorm:"size:100;default:'';not null;comment:动作名称"` // 动作名称
	Code   string `gorm:"size:100;default:'';not null;comment:动作编码"` // 动作编码
}
